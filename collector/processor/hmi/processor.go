package hmi

import (
	"fmt"
	"logsnap/collector"
	processor "logsnap/collector/processor"
	"path/filepath"
	"time"
)

type HMILogProcessor struct {
	*processor.BaseProcessor
}

// NewHMILogProcessor 创建一个新的HMI日志处理器
func NewHMILogProcessor(logDir string, outputDir string) *HMILogProcessor {
	return &HMILogProcessor{
		BaseProcessor: processor.NewBaseProcessor("HMI日志", logDir, outputDir),
	}
}

// Collect 处理日志文件的通用方法，子类可以覆盖
func (p *HMILogProcessor) Collect(startTime, endTime time.Time, rootOutputDir string) (string, []collector.FileProcessResult, error) {
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
func (p *HMILogProcessor) CreateFileProcessor() []processor.FileProcessorProvider {
	return []processor.FileProcessorProvider{
		NewUserOpFileProcessorProvider(),
		processor.NewCppLogFileProcessorProvider(),
	}
}
