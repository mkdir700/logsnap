package hmi_server

import (
	collector "logsnap/collector"
	processor "logsnap/collector/processor"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// var logTimePattern = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) \|`)


type HMIServerLogFileInfoFilter struct {}

func (l *HMIServerLogFileInfoFilter) parseLogFileInfo(filePath string) (processor.LogFileInfo, error) {
	fileName := filepath.Base(filePath)
	return processor.LogFileInfo{
		Path:      filePath,
		StartTime: time.Time{}, // 将当前日志文件的开始时间设置为零值，确保它是最新的
		FileName:  fileName,
		FileType: "log",
	}, nil
}

func (l *HMIServerLogFileInfoFilter) IsMatch(fileName string) bool {
	return strings.HasSuffix(fileName, ".log")
}

func (l *HMIServerLogFileInfoFilter) ParseFileInfos(files []string) ([]processor.LogFileInfo, error) {
	var fileInfos []processor.LogFileInfo

	for _, file := range files {
		if !l.IsMatch(file) {
			continue
		}
		fileInfo, err := l.parseLogFileInfo(file)
		if err != nil {
			return nil, err
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

// =======================


type HMIServerLogFileProcessorProvider struct {
	processor.BaseProcessorProvider
}

func NewHMIServerLogFileProcessorProvider() *HMIServerLogFileProcessorProvider {
	return &HMIServerLogFileProcessorProvider{
		BaseProcessorProvider: *processor.NewBaseProcessorProvider(
			&HMIServerLogFileInfoFilter{},
			[]string{".log"},
		),
	}
}

func (p *HMIServerLogFileProcessorProvider) FindFiles(dirPath string, suffixes ...string) ([]string, error) {
	return processor.DefaultFindLogFiles(dirPath, suffixes...)
}

func (p *HMIServerLogFileProcessorProvider) FilterFiles(files []string, startTime, endTime time.Time) ([]processor.LogFileInfo, error) {
	return processor.FilterFiles(files, p.FileInfoFilter, startTime, endTime, nil)
}

func (p *HMIServerLogFileProcessorProvider) ProcessDir(dirPath, outputDir string, startTime, endTime time.Time) ([]collector.FileProcessResult, error) {
	return processor.DefaultProcessDir(p, dirPath, outputDir, startTime, endTime)
}

func (p *HMIServerLogFileProcessorProvider) ProcessFile(fileInfo processor.LogFileInfo, startTime, endTime time.Time, outputDir string) (collector.FileProcessResult, error) {
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

func (p *HMIServerLogFileProcessorProvider) GetFileSuffixes() []string {
	return p.Suffixes
}