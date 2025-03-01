package hmi_server

import (
	"fmt"
	"logsnap/collector"
	processor "logsnap/collector/processor"
	"logsnap/collector/utils"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 用于从归档日志文件名解析时间戳的正则表达式
// 例如：2025-02-28_13-35-27_015111
var archiveTimePattern = regexp.MustCompile(`(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}_\d{6})`)

type HMIServerArchiveLogFileInfoFilter struct {}

// parseArchiveFileInfo 解析归档日志文件信息
func parseArchiveFileInfo(filePath string) (processor.LogFileInfo, error) {
	fileName := filepath.Base(filePath)
	matches := archiveTimePattern.FindStringSubmatch(fileName)
	if len(matches) <= 1 {
		return processor.LogFileInfo{}, fmt.Errorf("无法从文件名解析时间戳: %s", fileName)
	}

	timeStr := matches[1]
	fileTime, err := utils.ParseArchiveTimeStamp(timeStr)
	if err != nil {
		return processor.LogFileInfo{}, fmt.Errorf("解析文件时间戳失败 %s: %w", fileName, err)
	}

	return processor.LogFileInfo{
		Path:      filePath,
		StartTime: fileTime,
		FileName:  fileName,
		FileType: "zip",
	}, nil
}


func (l *HMIServerArchiveLogFileInfoFilter) IsMatch(fileName string) bool {
	return strings.HasSuffix(fileName, ".log.zip")
}

func (l *HMIServerArchiveLogFileInfoFilter) ParseFileInfos(files []string) ([]processor.LogFileInfo, error) {
	var fileInfos []processor.LogFileInfo

	for _, file := range files {
		if !l.IsMatch(file) {
			continue
		}
		fileInfo, err := parseArchiveFileInfo(file)
		if err != nil {
			return nil, err
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}


// ==========================

type HMIServerArchiveLogFileProcessorProvider struct {
	processor.BaseProcessorProvider
}

func NewHMIServerArchiveLogFileProcessorProvider() *HMIServerArchiveLogFileProcessorProvider {
	return &HMIServerArchiveLogFileProcessorProvider{
		BaseProcessorProvider: *processor.NewBaseProcessorProvider(
			&HMIServerArchiveLogFileInfoFilter{},
			[]string{".log.zip"},
		),
	}
}

func (p *HMIServerArchiveLogFileProcessorProvider) FindFiles(dirPath string, suffixes ...string) ([]string, error) {
	return processor.DefaultFindLogFiles(dirPath, suffixes...)
}

func (p *HMIServerArchiveLogFileProcessorProvider) FilterFiles(files []string, startTime, endTime time.Time) ([]processor.LogFileInfo, error) {
	return processor.FilterFiles(files, p.FileInfoFilter, startTime, endTime, nil)
}

func (p *HMIServerArchiveLogFileProcessorProvider) ProcessDir(dirPath, outputDir string, startTime, endTime time.Time) ([]collector.FileProcessResult, error) {
	return processor.DefaultProcessDir(p, dirPath, outputDir, startTime, endTime)
}

func (p *HMIServerArchiveLogFileProcessorProvider) ProcessFile(fileInfo processor.LogFileInfo, startTime, endTime time.Time, outputDir string) (collector.FileProcessResult, error) {
	var timePattern *regexp.Regexp
	var timeFormat string

	// 根据日志类型选择不同的时间模式和格式
	timePattern = logTimePattern
	timeFormat = "2006-01-02 15:04:05.000" // 程序日志的时间格式
	return processor.ProcessLogWithStrategy(fileInfo, outputDir, &processor.FilterLogProcessor{
		TimePattern: timePattern,
		TimeFormat: timeFormat,
		StartTime: startTime,
		EndTime: endTime,
		FileNameProcessor: HMIServerFileNameProcessor,
		ReaderCreator: HMIServerLogFileReaderCreator,
	})
}

func (p *HMIServerArchiveLogFileProcessorProvider) GetFileSuffixes() []string {
	return p.Suffixes
}