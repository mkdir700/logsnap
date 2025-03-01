package hmi_server

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"logsnap/collector"
	processor "logsnap/collector/processor"
)

// HMIServer 的日志文件
// 日志文件后缀：
//   - .log, 当前日志文件
//   - .log.zip, 归档日志文件
// 日志文件格式：
//   - 2025-02-28 13:35:27.015 |
//   - 2025-02-28 13:35:27.015 |
//   - 2025-02-28 13:35:27.015 |
//   - 2025-02-28 13:35:27.015 |
//   - 2025-02-28 13:35:27.015 |

// 用于从日志行提取时间戳的正则表达式
// 注意：这个正则表达式需要根据实际日志格式调整
// 例如：2025-02-28 13:35:27.015 |
var logTimePattern = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) \|`)

// HMIServerFileNameProcessor 为 HMI 服务器日志生成输出文件名
func HMIServerFileNameProcessor(fileInfo processor.LogFileInfo) string {
	if fileInfo.FileType == "zip" {
		trimedFileName := strings.TrimSuffix(fileInfo.FileName, ".zip")
		return strings.ReplaceAll(trimedFileName, ".log", ".archive.log")
	} else {
		// 对于普通文件，加上时间戳以避免覆盖
		baseName := fileInfo.FileName
		ext := filepath.Ext(baseName)
		nameWithoutExt := strings.TrimSuffix(baseName, ext)
		return fmt.Sprintf("%s_%s%s",
			nameWithoutExt,
			time.Now().Format("20060102_150405"),
			ext)
	}
}

// HMIServerLogFileReaderCreator 为 HMI 服务器日志生成读取器
func HMIServerLogFileReaderCreator(fileInfo processor.LogFileInfo) (io.ReadCloser, error) {
	// 根据文件类型处理
	if fileInfo.FileType == "zip" || strings.HasSuffix(fileInfo.FileName, ".zip") {
		// 打开ZIP文件
		zipReader, err := zip.OpenReader(fileInfo.Path)
		if err != nil {
			return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
		}

		// 在ZIP中查找日志文件
		baseFileName := strings.TrimSuffix(fileInfo.FileName, ".zip")
		var logFileInZip *zip.File
		for _, f := range zipReader.File {
			if filepath.Base(f.Name) == baseFileName {
				logFileInZip = f
				break
			}
		}

		if logFileInZip == nil {
			zipReader.Close()
			return nil, fmt.Errorf("在ZIP文件中未找到日志文件: %s", fileInfo.Path)
		}

		// 打开ZIP中的日志文件
		reader, err := logFileInZip.Open()
		if err != nil {
			zipReader.Close()
			return nil, fmt.Errorf("打开ZIP中的日志文件失败: %w", err)
		}

		// 返回一个复合的ReadCloser，它会同时关闭内部reader和zipReader
		return processor.NewCompositeReadCloser(reader, zipReader), nil
	}
	return processor.DefaultCreateReaderForFile(fileInfo)
}


type HMIServerLogProcessor struct {
	processor.BaseProcessor
}

func NewHMIServerLogProcessor(logDir string, outputDir string) *HMIServerLogProcessor {
	return &HMIServerLogProcessor{
		BaseProcessor: *processor.NewBaseProcessor("HMI服务器日志", logDir, outputDir),
	}
}

// Collect 处理日志文件的通用方法，子类可以覆盖
func (p *HMIServerLogProcessor) Collect(startTime, endTime time.Time, rootOutputDir string) (string, []collector.FileProcessResult, error) {
	// 创建文件处理器
	fileProcessors := p.CreateFileProcessor()

	if len(fileProcessors) == 0 {
		return "", nil, fmt.Errorf("没有找到文件处理器")
	}
	
	// 拼接成完整的输出目录
	outputDir := filepath.Join(rootOutputDir, p.OutputDir)

	// 使用通用的收集方法
	return processor.CollectWithProcessor(p, fileProcessors, startTime, endTime, outputDir)
}


// CreateFileProcessor 创建文件处理器，子类应该覆盖此方法
func (p *HMIServerLogProcessor) CreateFileProcessor() []processor.FileProcessorProvider {
	return []processor.FileProcessorProvider{
		NewHMIServerLogFileProcessorProvider(),
		NewHMIServerArchiveLogFileProcessorProvider(),
	}
}