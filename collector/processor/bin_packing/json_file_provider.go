package bin_packing

import (
	"fmt"
	processor "logsnap/collector/processor"
	"path/filepath"
	"regexp"
	"time"
)

// 2025-03-03_17-40-27_0_task.json
var timePatternForProgramJSONLogFile = regexp.MustCompile(`(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})`)


type JsonFileInfoFilter struct {}


func (f *JsonFileInfoFilter) ParseFileInfos(files []string) ([]processor.LogFileInfo, error) {
	var fileInfos []processor.LogFileInfo
	for _, filePath := range files {
		fileName := filepath.Base(filePath)

		// 从文件名中提取时间戳
		matches := timePatternForProgramJSONLogFile.FindStringSubmatch(fileName)
		if len(matches) <= 1 {
			return nil, fmt.Errorf("无法从日志文件名中解析出时间戳: %s", fileName)
		}
	
		timeStr := matches[1]
		fileTime, err := time.ParseInLocation("2006-01-02_15-04-05", timeStr, time.Local)
		if err != nil {
			return nil, fmt.Errorf("解析程序日志文件时间戳失败 %s: %w", fileName, err)
		}
		
		fileInfos = append(fileInfos, processor.LogFileInfo{
			Path: filePath,
			StartTime: fileTime,
			FileName: fileName,
			FileType: "json",
			Extra: nil,
		})
	}
	return fileInfos, nil
}

func (f *JsonFileInfoFilter) IsMatch(fileName string) bool {
	return timePatternForProgramJSONLogFile.MatchString(fileName)
}

type JsonFileProcessorProvider struct {
	processor.BaseProcessorProvider
}

func NewJsonFileProcessorProvider() *JsonFileProcessorProvider {
	return &JsonFileProcessorProvider{
		BaseProcessorProvider: *processor.NewBaseProcessorProvider(
			&JsonFileInfoFilter{},
			[]string{".json"},
		),
	}
}