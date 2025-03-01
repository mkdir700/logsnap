package processor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	collector "logsnap/collector"

	"github.com/sirupsen/logrus"
)

// LogProcessStrategy 定义日志处理策略接口
type LogProcessStrategy interface {
	// Process 处理日志文件
	// 返回：处理结果，错误信息
	Process(fileInfo LogFileInfo, outputDir string) (collector.FileProcessResult, error)
}

// FilterLogProcessor 按时间过滤日志的处理器
type FilterLogProcessor struct {
	TimePattern       *regexp.Regexp
	TimeFormat        string
	StartTime         time.Time
	EndTime           time.Time
	FileNameProcessor FileNameProcessor
	ReaderCreator     ReaderCreator
}

func NewFilterLogProcessor(timePattern *regexp.Regexp, timeFormat string, startTime time.Time, endTime time.Time, fileNameProcessor FileNameProcessor, readerCreator ReaderCreator) *FilterLogProcessor {
	return &FilterLogProcessor{
		TimePattern:       timePattern,
		TimeFormat:        timeFormat,
		StartTime:         startTime,
		EndTime:           endTime,
		FileNameProcessor: fileNameProcessor,
		ReaderCreator:     readerCreator,
	}
}

// Process 实现 LogProcessStrategy 接口
func (p *FilterLogProcessor) Process(fileInfo LogFileInfo, outputDir string) (collector.FileProcessResult, error) {
	lines, matches, totalSize, outputPath, err := ProcessLogFileByLine(
		fileInfo,
		p.TimePattern,
		p.TimeFormat,
		p.StartTime,
		p.EndTime,
		outputDir,
		p.FileNameProcessor,
		p.ReaderCreator,
	)

	result := collector.FileProcessResult{
		FilePath:   outputPath, // 设置输出文件路径
		FileCount:  1,          // 处理了一个文件
		FileSize:   totalSize,  // 处理的文件大小
		TotalLines: lines,      // 处理的总行数
		MatchLines: matches,    // 匹配的行数
	}

	return result, err
}

// CopyLogProcessor 直接复制日志的处理器
type CopyLogProcessor struct{}

func NewCopyLogProcessor() *CopyLogProcessor {
	return &CopyLogProcessor{}
}

// Process 实现 LogProcessStrategy 接口
func (p *CopyLogProcessor) Process(fileInfo LogFileInfo, outputDir string) (collector.FileProcessResult, error) {
	bytes, files, outputPath, err := ProcessLogFileByCopying(fileInfo, outputDir)

	result := collector.FileProcessResult{
		FilePath:   outputPath,   // 设置输出文件路径
		FileCount:  1,            // 处理了一个文件
		FileSize:   int64(bytes), // 处理的文件大小
		MatchFiles: files,        // 复制的文件数
	}

	return result, err
}

// ProcessLogWithStrategy 使用指定策略处理日志文件
func ProcessLogWithStrategy(fileInfo LogFileInfo, outputDir string, strategy LogProcessStrategy) (collector.FileProcessResult, error) {
	return strategy.Process(fileInfo, outputDir)
}

// ProcessLogFileByLine 按行处理日志文件，筛选出符合时间范围的日志行，并写入到输出文件中
// 参数:
//   - fileInfo: 日志文件信息
//   - timePattern: 时间正则表达式
//   - timeFormat: 时间格式
//   - startTime: 开始时间
//   - endTime: 结束时间
//   - outputDir: 输出目录
//   - fileNameProcessor: 可选的文件名处理函数，如果为nil则使用默认函数， 默认直接使用文件名
//   - readerCreator: 可选的读取器创建函数，如果为nil则使用默认函数， 默认将文件作为普通文本文件读取
//
// 返回:
//   - 匹配的行数
//   - 匹配的文件数
//   - 文件总大小
//   - 输出文件路径（如果有匹配的行）
//   - 错误信息
func ProcessLogFileByLine(
	fileInfo LogFileInfo,
	timePattern *regexp.Regexp,
	timeFormat string,
	startTime, endTime time.Time,
	outputDir string,
	fileNameProcessor FileNameProcessor,
	readerCreator ReaderCreator,
) (int, int, int64, string, error) {
	// 使用提供的处理函数或默认函数
	if fileNameProcessor == nil {
		fileNameProcessor = DefaultGenerateOutputFileName
	}
	if readerCreator == nil {
		readerCreator = DefaultCreateReaderForFile
	}

	// 生成输出文件名
	outputFileName := fileNameProcessor(fileInfo)
	outputPath := filepath.Join(outputDir, outputFileName)

	// 打开输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outputFile.Close()

	// 创建读取器
	reader, err := readerCreator(fileInfo)
	if err != nil {
		return 0, 0, 0, "", err
	}
	defer reader.Close()

	// 处理日志内容
	lineCount, matchCount, totalSize, err := ProcessLogContent(reader, outputFile, timePattern, timeFormat, startTime, endTime)
	if err != nil {
		return lineCount, matchCount, 0, "", err
	}

	// 如果有匹配的内容，写入文件头
	if lineCount > 0 && matchCount > 0 {
		// 重新打开文件以在开头写入头信息
		outputFile.Seek(0, 0)
		tempContent, err := io.ReadAll(outputFile)
		if err != nil {
			return lineCount, matchCount, 0, "", fmt.Errorf("读取临时内容失败: %w", err)
		}

		outputFile.Seek(0, 0)
		outputFile.Truncate(0)

		// 写入文件头
		err = WriteFileHeader(outputFile, fileInfo.Path, startTime, endTime)
		if err != nil {
			return lineCount, matchCount, 0, "", fmt.Errorf("写入文件头失败: %w", err)
		}

		// 写回原内容
		_, err = outputFile.Write(tempContent)
		if err != nil {
			return lineCount, matchCount, 0, "", fmt.Errorf("写回原内容失败: %w", err)
		}

		return lineCount, matchCount, totalSize, outputPath, nil
	} else {
		// 删除空文件
		os.Remove(outputPath)
		return lineCount, matchCount, totalSize, "", nil
	}
}

// ProcessLogFileByCopying 处理单个日志文件，处理逻辑是直接将原文件内容复制一份到目标目录
// 参数:
//   - fileInfo: 日志文件信息
//   - outputDir: 输出目录
//
// 返回:
//   - 处理的文件大小(字节)
//   - 复制的文件数量(成功为1，失败为0)
//   - 输出文件路径
//   - 错误信息
func ProcessLogFileByCopying(
	fileInfo LogFileInfo,
	outputDir string,
) (int64, int, string, error) {
	// 生成输出文件名
	outputFileName := filepath.Base(fileInfo.Path)
	outputPath := filepath.Join(outputDir, outputFileName)

	// 打开源文件
	sourceFile, err := os.Open(fileInfo.Path)
	if err != nil {
		return 0, 0, "", fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(outputPath)
	if err != nil {
		return 0, 0, "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	// 复制文件内容
	bytesCopied, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return int64(bytesCopied), 0, "", fmt.Errorf("复制文件内容失败: %w", err)
	}

	return int64(bytesCopied), 1, outputPath, nil
}

// FilterLogLineByTime 根据时间范围过滤日志行
// 参数:
//   - line: 日志行
//   - timePattern: 时间正则表达式
//   - timeFormat: 时间格式
//   - startTime: 开始时间
//   - endTime: 结束时间
//
// 返回:
//   - 是否匹配
//   - 错误信息
func FilterLogLineByTime(line string, timePattern *regexp.Regexp, timeFormat string, startTime, endTime time.Time) (bool, error) {
	// 从日志行中提取时间戳
	matches := timePattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		// 如果行没有时间戳，返回false
		return false, nil
	}

	// 解析时间戳，使用本地时区
	timeStr := matches[1]
	timestamp, err := time.ParseInLocation(timeFormat, timeStr, time.Local)
	if err != nil {
		return false, err
	}

	// 检查是否在时间范围内
	return (timestamp.After(startTime) || timestamp.Equal(startTime)) &&
		(timestamp.Before(endTime) || timestamp.Equal(endTime)), nil
}

// ProcessLogContent 处理日志内容，根据时间范围过滤
// 参数:
//   - reader: 日志内容读取器
//   - writer: 日志内容写入器
//   - timePattern: 时间正则表达式
//   - timeFormat: 时间格式
//   - startTime: 开始时间
//   - endTime: 结束时间
//
// 返回:
//   - 处理的总行数
//   - 匹配的行数
//   - 文件总大小
//   - 错误信息
func ProcessLogContent(
	reader io.Reader,
	writer io.Writer,
	timePattern *regexp.Regexp,
	timeFormat string,
	startTime, endTime time.Time,
) (int, int, int64, error) {
	scanner := bufio.NewScanner(reader)
	// 设置 1MB 的缓冲区，对于普通日志应该足够了
	buf := make([]byte, 1024*1024) // 1MB
	scanner.Buffer(buf, 1024*1024)

	bufWriter := bufio.NewWriter(writer)
	defer bufWriter.Flush()

	totalSize := int64(0)
	lineCount := 0
	matchCount := 0
	
	var currentLogEntry []string    // 当前日志条目的所有行
	var currentEntryMatched bool    // 当前条目是否匹配时间范围
	var hasCurrentEntry bool = false // 是否有正在处理的日志条目

	// 处理一个完整的日志条目
	processLogEntry := func() {
		if len(currentLogEntry) > 0 && currentEntryMatched {
			for _, line := range currentLogEntry {
				if _, err := fmt.Fprintln(bufWriter, line); err != nil {
					logrus.Errorf("写入输出失败: %v", err)
				}
			}
			matchCount++
		}
		// 重置当前日志条目
		currentLogEntry = nil
		currentEntryMatched = false
		hasCurrentEntry = false
	}

	for scanner.Scan() {
		line := scanner.Text()
		totalSize += int64(len(line))
		lineCount++

		// 检查这一行是否包含时间戳（是否是新日志条目的开始）
		if timePattern.MatchString(line) {
			// 如果已经有一个日志条目在处理中，先处理完它
			if hasCurrentEntry {
				processLogEntry()
			}
			
			// 开始一个新的日志条目
			hasCurrentEntry = true
			match, err := FilterLogLineByTime(line, timePattern, timeFormat, startTime, endTime)
			if err != nil {
				// 记录解析错误但继续处理
				logrus.Errorf("第 %d 行解析失败: %v", lineCount, err)
				continue
			}
			
			currentEntryMatched = match
			currentLogEntry = append(currentLogEntry, line)
		} else if hasCurrentEntry {
			// 这一行不包含时间戳，属于当前日志条目的一部分
			currentLogEntry = append(currentLogEntry, line)
		} else {
			// 这一行不包含时间戳，也不属于任何日志条目
			// 可能是文件开头的注释或者其他非日志内容，忽略它
			logrus.Debugf("忽略不属于任何日志条目的行: %s", line)
		}
	}

	// 处理最后一个日志条目
	if hasCurrentEntry {
		processLogEntry()
	}

	if err := scanner.Err(); err != nil {
		return lineCount, matchCount, totalSize, fmt.Errorf("读取日志内容失败 (行数: %d): %w", lineCount, err)
	}

	return lineCount, matchCount, totalSize, nil
}

// WriteFileHeader 写入文件头信息
func WriteFileHeader(writer io.Writer, originalPath string, startTime, endTime time.Time) error {
	_, err := fmt.Fprintf(writer, "# 原始日志文件: %s\n", originalPath)
	if err != nil {
		return err
	}
	// 使用本地时区格式化时间
	_, err = fmt.Fprintf(writer, "# 时间范围: %s 到 %s\n\n",
		startTime.In(time.Local).Format("2006-01-02 15:04:05"),
		endTime.In(time.Local).Format("2006-01-02 15:04:05"))
	return err
}
