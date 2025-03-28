package processor

import (
	"fmt"
	collector "logsnap/collector"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 日志文件后缀是动态的数字（例如：2966778）
// 日志文件名是：xyz_studio_max_bin.xyz-Workstation.xyz.log.ERROR.20250228-094825.2966778
// 例如：xyz_studio_max_bin.xyz-Workstation.xyz.log.ERROR.20250228-094825.2966778
var timePatternForProgramLogFile = regexp.MustCompile(`\.(\d{8}-\d{6})\.`)

// 用于从程序日志行提取时间戳的正则表达式
// 例如：E20250228 09:48:25.654057 2966778 station_status_label.cpp:53] 0:  0
// [IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg
var timePatternForProgramLogLine = regexp.MustCompile(`[IWEF](\d{4}\d{2}\d{2} \d{2}:\d{2}:\d{2}\.\d{6})`)

type CppLogFileInfoFilter struct{}

func (f *CppLogFileInfoFilter) parseLogFileInfo(filePath string) (LogFileInfo, error) {
	fileName := filepath.Base(filePath)

	// 从文件名中提取时间戳
	matches := timePatternForProgramLogFile.FindStringSubmatch(fileName)
	if len(matches) <= 1 {
		return LogFileInfo{}, fmt.Errorf("无法从日志文件名中解析出时间戳: %s", fileName)
	}

	timeStr := matches[1]
	// 解析时间戳，格式为 20250228-094825
	fileTime, err := time.ParseInLocation("20060102-150405", timeStr, time.Local)
	if err != nil {
		return LogFileInfo{}, fmt.Errorf("解析程序日志文件时间戳失败 %s: %w", fileName, err)
	}

	return LogFileInfo{
		Path:      filePath,
		StartTime: fileTime,
		FileName:  fileName,
	}, nil

}

// ParseFileInfos 解析文件信息列表
// 1. 解析文件信息
// 2. 如果解析失败，则跳过该文件
// 3. 如果解析成功，则将文件信息添加到列表中
func (l *CppLogFileInfoFilter) ParseFileInfos(files []string) ([]LogFileInfo, error) {
	var fileInfos []LogFileInfo

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

func (l *CppLogFileInfoFilter) IsMatch(fileName string) bool {
	return timePatternForProgramLogFile.MatchString(fileName)
}

type CppLogFileProcessorProvider struct {
	fileInfoFilter *CppLogFileInfoFilter
}

func NewCppLogFileProcessorProvider() *CppLogFileProcessorProvider {
	return &CppLogFileProcessorProvider{
		fileInfoFilter: &CppLogFileInfoFilter{},
	}
}

func (p *CppLogFileProcessorProvider) FindFiles(dirPath string, suffixes ...string) ([]string, error) {
	return DefaultFindLogFiles(dirPath, suffixes...)
}

func (p *CppLogFileProcessorProvider) FilterFiles(files []string, startTime, endTime time.Time) ([]LogFileInfo, error) {
	// 在同一个目录下，文件名会出现相同创建时间的多个文件，包括 INFO、WARNING、ERROR
	// xyz_studio_max_bin.xyz-Workstation.xyz.log.INFO.20250301-171208
	// xyz_studio_max_bin.xyz-Workstation.xyz.log.WARNING.20250301-171208
	// xyz_studio_max_bin.xyz-Workstation.xyz.log.ERROR.20250301-171208
	// 由于 FilterFiles 的排序规则，是按照文件名的时间作为开始时间，该文件的下一个时间作为结束时间
	// 因此需要在 FilterFiles 中将这些文件按日志等级进行分组
	var infoFiles, warnFiles, errorFiles, otherFiles []string
	var err error
	for _, file := range files {
		if strings.Contains(file, ".INFO.") {
			infoFiles = append(infoFiles, file)
		} else if strings.Contains(file, ".WARNING.") {
			warnFiles = append(warnFiles, file)
		} else if strings.Contains(file, ".ERROR.") {
			errorFiles = append(errorFiles, file)
		} else {
			otherFiles = append(otherFiles, file)
		}
	}

	infoFileInfos, err := FilterFiles(infoFiles, p.fileInfoFilter, startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	warnFileInfos, err := FilterFiles(warnFiles, p.fileInfoFilter, startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	errorFileInfos, err := FilterFiles(errorFiles, p.fileInfoFilter, startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	otherFileInfos, err := FilterFiles(otherFiles, p.fileInfoFilter, startTime, endTime, nil)
	if err != nil {
		return nil, err
	}

	allFileInfos := append(infoFileInfos, warnFileInfos...)
	allFileInfos = append(allFileInfos, errorFileInfos...)
	allFileInfos = append(allFileInfos, otherFileInfos...)
	SortByTime(allFileInfos)

	return allFileInfos, nil
}

func (p *CppLogFileProcessorProvider) ProcessDir(dirPath, outputDir string, startTime, endTime time.Time) ([]collector.FileProcessResult, error) {
	return DefaultProcessDir(p, dirPath, outputDir, startTime, endTime)
}

func (p *CppLogFileProcessorProvider) ProcessFile(fileInfo LogFileInfo, startTime, endTime time.Time, outputDir string) (collector.FileProcessResult, error) {
	var timePattern *regexp.Regexp
	var timeFormat string

	// 根据日志类型选择不同的时间模式和格式
	timePattern = timePatternForProgramLogLine
	timeFormat = "20060102 15:04:05.000000" // 程序日志的时间格式
	return ProcessLogWithStrategy(fileInfo, outputDir, &FilterLogProcessor{
		TimePattern: timePattern,
		TimeFormat: timeFormat,
		StartTime: startTime,
		EndTime: endTime,
		FileNameProcessor: nil,
		ReaderCreator: nil,
	})
}

func (p *CppLogFileProcessorProvider) GetFileSuffixes() []string {
	return []string{}
}