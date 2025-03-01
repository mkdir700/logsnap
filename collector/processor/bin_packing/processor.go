package bin_packing

import (
	"fmt"
	"logsnap/collector"
	processor "logsnap/collector/processor"
	"path/filepath"
	"time"
)

type BinPackingLogProcessor struct {
	*processor.BaseProcessor
}

func NewBinPackingLogProcessor(logDir string, outputDir string) *BinPackingLogProcessor {
	return &BinPackingLogProcessor{
		BaseProcessor: processor.NewBaseProcessor("BinPacking日志", logDir, outputDir),
	}
}

// Collect 处理日志文件的通用方法，子类可以覆盖
func (p *BinPackingLogProcessor) Collect(startTime, endTime time.Time, rootOutputDir string) (string, []collector.FileProcessResult, error) {
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
func (p *BinPackingLogProcessor) CreateFileProcessor() []processor.FileProcessorProvider {
	return []processor.FileProcessorProvider{
		NewJsonFileProcessorProvider(),
		processor.NewCppLogFileProcessorProvider(),
	}
} 