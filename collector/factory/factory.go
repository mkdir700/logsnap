package factory

import (
	"fmt"
	"logsnap/collector"
	binPackingProcessor "logsnap/collector/processor/bin_packing"
	cppLogProcessor "logsnap/collector/processor/cpp_log"
	hmiProcessor "logsnap/collector/processor/hmi"
	hmiServerProcessor "logsnap/collector/processor/hmi_server"
	studioMaxProcessor "logsnap/collector/processor/studio_max"
)

// ProcessorFactory 处理器工厂接口
type ProcessorFactory interface {
	// CreateProcessor 创建处理器
	CreateProcessor(logDir, outputDir string) (collector.LogProcessor, error)
}

// HMIProcessorFactory HMI处理器工厂
type HMIProcessorFactory struct{}

// CreateProcessor 实现 ProcessorFactory 接口
func (f *HMIProcessorFactory) CreateProcessor(logDir, outputDir string) (collector.LogProcessor, error) {
	return hmiProcessor.NewHMILogProcessor(logDir+"/xyz_hmi", outputDir+"/xyz_hmi"), nil
}

type HMIServerProcessorFactory struct{}

func (f *HMIServerProcessorFactory) CreateProcessor(logDir, outputDir string) (collector.LogProcessor, error) {
	return hmiServerProcessor.NewHMIServerLogProcessor(logDir+"/xyz_max_hmi/server", outputDir+"/xyz_max_hmi/server"), nil
}

type StudioMaxProcessorFactory struct{}

func (f *StudioMaxProcessorFactory) CreateProcessor(logDir, outputDir string) (collector.LogProcessor, error) {
	return studioMaxProcessor.NewStudioMaxLogProcessor(logDir+"/xyz_studio_max", outputDir+"/xyz_studio_max"), nil
}

type BinPackingProcessorFactory struct{}

func (f *BinPackingProcessorFactory) CreateProcessor(logDir, outputDir string) (collector.LogProcessor, error) {
	return binPackingProcessor.NewBinPackingLogProcessor(logDir+"/xyz_bin_packing", outputDir+"/xyz_bin_packing"), nil
}

type GenericLogProcessorFactory struct {
	Path string
}

func (f *GenericLogProcessorFactory) CreateProcessor(logDir, outputDir string) (collector.LogProcessor, error) {
	return cppLogProcessor.NewCppLogProcessor(logDir+f.Path, outputDir+f.Path), nil
}

// ProcessorFactoryRegistry 处理器工厂注册表
var ProcessorFactoryRegistry = map[collector.ProcessorType]ProcessorFactory{
	collector.HMIProcessorType:             &HMIProcessorFactory{},
	collector.HMIServerProcessorType:       &HMIServerProcessorFactory{},
	collector.StudioMaxProcessorType:       &StudioMaxProcessorFactory{},
	collector.BinPackingProcessorType:      &BinPackingProcessorFactory{},
	collector.VisionLogViewerProcessorType: &GenericLogProcessorFactory{Path: "/vision_log_viewer"},
	collector.RobotDriverNodeProcessorType: &GenericLogProcessorFactory{Path: "/xyz_robot_driver_node"},
}

// CreateProcessor 创建指定类型的日志处理器
// 参数:
//   - processorType: 处理器类型
//   - logDir: 日志目录
//   - outputDir: 输出目录
//
// 返回:
//   - processor: 日志处理器
//   - err: 错误信息
func CreateProcessor(processorType collector.ProcessorType, logDir, outputDir string) (collector.LogProcessor, error) {
	factory, exists := ProcessorFactoryRegistry[processorType]
	if !exists {
		return nil, fmt.Errorf("不支持的处理器类型: %s", processorType)
	}

	return factory.CreateProcessor(logDir, outputDir)
}

// CreateCollector 创建收集器并添加指定类型的处理器
// 参数:
//   - processorTypes: 处理器类型列表
//   - logDirs: 日志目录列表，与处理器类型一一对应
//   - outputDir: 输出目录
//
// 返回:
//   - collector: 收集器
//   - err: 错误信息
func CreateCollector(processorTypes []collector.ProcessorType, logDirs []string, outputDir string) (*collector.Collector, error) {
	if len(processorTypes) != len(logDirs) {
		return nil, fmt.Errorf("处理器类型数量与日志目录数量不匹配")
	}

	// 创建处理器列表
	processors := make([]collector.LogProcessor, 0, len(processorTypes))
	for i, processorType := range processorTypes {
		processor, err := CreateProcessor(processorType, logDirs[i], outputDir)
		if err != nil {
			return nil, fmt.Errorf("创建处理器失败: %w", err)
		}
		processors = append(processors, processor)
	}

	// 创建收集器
	return collector.NewCollector(processors, outputDir), nil
}

// GetSupportedProcessorTypes 获取支持的处理器类型列表
func GetSupportedProcessorTypes() []collector.ProcessorType {
	return []collector.ProcessorType{
		collector.HMIProcessorType,
		collector.HMIServerProcessorType,
		collector.StudioMaxProcessorType,
		collector.BinPackingProcessorType,
		collector.VisionLogViewerProcessorType,
		collector.RobotDriverNodeProcessorType,
	}
}
