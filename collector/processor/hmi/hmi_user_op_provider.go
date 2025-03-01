package hmi

import (
	"fmt"
	"logsnap/collector"
	processor "logsnap/collector/processor"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 用于从 HMI 用户操作日志文件名解析时间戳的正则表达式
// 例如：20250302-084053.269
var timePatternForUserOpLogFile = regexp.MustCompile(`^(\d{8}-\d{6})\.`)
// 用于从用户操作日志行提取时间戳的正则表达式
// 例如：20250302 08:41:13.163] User clicked [StartTask].
var timePatternForUserOpLogLine = regexp.MustCompile(`^(\d{8} \d{2}:\d{2}:\d{2}\.\d{3})]`)


type UserOpFileInfoFilter struct {}

func (p *UserOpFileInfoFilter) parseLogFileInfo(filePath string) (processor.LogFileInfo, error) {
	fileName := filepath.Base(filePath)
	// 从文件名中提取时间戳
	matches := timePatternForUserOpLogFile.FindStringSubmatch(fileName)
	if len(matches) <= 1 {
		return processor.LogFileInfo{}, fmt.Errorf("无法从用户操作日志文件名解析时间戳: %s", fileName)
	}

	timeStr := matches[1]
	// 解析时间戳，格式为 20250307-135338.555
	timeStr = strings.Split(timeStr, ".")[0]
	fileTime, err := time.ParseInLocation("20060102-150405", timeStr, time.Local)
	fileTime = fileTime.Local()
	if err != nil {
		return processor.LogFileInfo{}, fmt.Errorf("解析用户操作日志文件时间戳失败 %s: %w", fileName, err)
	}

	return processor.LogFileInfo{
		Path:      filePath,
		StartTime: fileTime,
		FileName:  fileName,
		FileType: "user_op",
	}, nil
}

// ParseFileInfos 解析文件信息列表
func (p *UserOpFileInfoFilter) ParseFileInfos(files []string) ([]processor.LogFileInfo, error) {
	var fileInfos []processor.LogFileInfo

	for _, file := range files {
		fileName := filepath.Base(file)

		// 检查是否为HMI程序日志文件或用户操作日志文件
		if p.IsMatch(fileName) {
			fileInfo, err := p.parseLogFileInfo(file)
			if err != nil {
				return nil, err
			}
			fileInfos = append(fileInfos, fileInfo)
		}
	}

	return fileInfos, nil
}

func (f *UserOpFileInfoFilter) IsMatch(fileName string) bool {
	return timePatternForUserOpLogFile.MatchString(fileName)
}

type UserOpFileProcessorProvider struct {
	processor.BaseProcessorProvider
}

func NewUserOpFileProcessorProvider() *UserOpFileProcessorProvider {
	return &UserOpFileProcessorProvider{
		BaseProcessorProvider: *processor.NewBaseProcessorProvider(
			&UserOpFileInfoFilter{},
			[]string{},
		),
	}
}

func (p *UserOpFileProcessorProvider) FindFiles(dirPath string, suffixes ...string) ([]string, error) {
	return processor.DefaultFindLogFiles(dirPath, suffixes...)
}

func (p *UserOpFileProcessorProvider) FilterFiles(files []string, startTime, endTime time.Time) ([]processor.LogFileInfo, error) {
	return processor.FilterFiles(files, p.FileInfoFilter, startTime, endTime, nil)
}

func (p *UserOpFileProcessorProvider) ProcessDir(dirPath, outputDir string, startTime, endTime time.Time) ([]collector.FileProcessResult, error) {
	return processor.DefaultProcessDir(p, dirPath, outputDir, startTime, endTime)
}

func (p *UserOpFileProcessorProvider) ProcessFile(fileInfo processor.LogFileInfo, startTime, endTime time.Time, outputDir string) (collector.FileProcessResult, error) {
	var timePattern *regexp.Regexp
	var timeFormat string

	timePattern = timePatternForUserOpLogLine
	timeFormat = "20060102 15:04:05.000" // 用户操作日志的时间格式

	// 使用通用的日志处理函数，提供自定义的文件名处理器
	return processor.ProcessLogWithStrategy(
		fileInfo,
		outputDir,
		&processor.FilterLogProcessor{
			TimePattern: timePattern,
			TimeFormat: timeFormat,
			StartTime: startTime,
			EndTime: endTime,
			FileNameProcessor: nil,
			ReaderCreator: nil,
		},
	)
}

func (p *UserOpFileProcessorProvider) GetFileSuffixes() []string {
	return []string{}
}