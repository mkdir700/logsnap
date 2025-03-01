package collector

import (
	"archive/zip"
	"fmt"
	"logsnap/collector/utils"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ProcessorType 定义处理器类型
type ProcessorType string

const (
	// HMIProcessorType HMI日志处理器类型
	HMIProcessorType ProcessorType = "xyz-hmi"
	// HMIServerProcessorType HMI服务器日志处理器类型
	HMIServerProcessorType ProcessorType = "xyz-max-hmi-server"
	// StudioMaxProcessorType StudioMax日志处理器类型
	StudioMaxProcessorType ProcessorType = "xyz-studio-max"
	// BinPackingProcessorType BinPacking日志处理器类型
	BinPackingProcessorType ProcessorType = "xyz-bin-packing"
	// VisionLogViewerProcessorType VisionLogViewer日志处理器类型
	VisionLogViewerProcessorType ProcessorType = "vision-log-viewer"
	// RobotDriverNodeProcessorType RobotDriverNode日志处理器类型
	RobotDriverNodeProcessorType ProcessorType = "robot-driver-node"
)

// LogProcessor 定义日志处理接口
type LogProcessor interface {
	// Name 返回日志处理器的名称
	GetName() string

	// GetLogPath 返回日志文件的路径
	GetLogPath() (string, error)

	// GetOutputDir 返回日志文件的输出目录
	GetOutputDir() string

	// Collect 处理日志文件，提取指定时间范围内的日志
	// 参数:
	//   - startTime: 开始时间
	//   - outputDir: 输出目录
	// 返回:
	//   - outputPath: 处理后的日志文件路径
	//   - lineCount: 处理的总行数
	//   - matchCount: 匹配的日志行数
	//   - error: 错误信息
	Collect(startTime, endTime time.Time, outputDir string) (outputPath string, results []FileProcessResult, err error)
}

// Collector 负责收集和打包日志
type Collector struct {
	logProcessors []LogProcessor
	outputDir     string // 最终ZIP文件的输出目录
}

// NewCollector 创建新的收集器
func NewCollector(logProcessors []LogProcessor, outputDir string) *Collector {
	return &Collector{
		logProcessors: logProcessors,
		outputDir:     outputDir,
	}
}

// AddProcessor 添加日志处理器
func (c *Collector) AddProcessor(processor LogProcessor) {
	c.logProcessors = append(c.logProcessors, processor)
	logrus.Infof("已添加日志处理器: %s", processor.GetName())
}

// RemoveProcessor 移除指定名称的日志处理器
func (c *Collector) RemoveProcessor(name string) bool {
	for i, p := range c.logProcessors {
		if p.GetName() == name {
			// 从切片中移除元素
			c.logProcessors = append(c.logProcessors[:i], c.logProcessors[i+1:]...)
			logrus.Infof("已移除日志处理器: %s", name)
			return true
		}
	}
	logrus.Warnf("未找到名为 %s 的日志处理器", name)
	return false
}

// GetProcessors 获取所有日志处理器
func (c *Collector) GetProcessors() []LogProcessor {
	return c.logProcessors
}

// GetProcessorByName 根据名称获取日志处理器
func (c *Collector) GetProcessorByName(name string) (LogProcessor, bool) {
	for _, p := range c.logProcessors {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}

// HasProcessor 检查是否包含指定名称的日志处理器
func (c *Collector) HasProcessor(name string) bool {
	_, found := c.GetProcessorByName(name)
	return found
}

// ClearProcessors 清空所有日志处理器
func (c *Collector) ClearProcessors() {
	c.logProcessors = []LogProcessor{}
	logrus.Info("已清空所有日志处理器")
}

// SetOutputDir 设置输出目录
func (c *Collector) SetOutputDir(outputDir string) {
	c.outputDir = outputDir
	logrus.Infof("已设置输出目录: %s", outputDir)
}

// GetOutputDir 获取输出目录
func (c *Collector) GetOutputDir() string {
	return c.outputDir
}

// Collect 收集指定时间范围内的日志（多线程版本）
func (c *Collector) Collect(startTime, endTime time.Time) (string, error) {
	// 验证时间范围
	if endTime.Before(startTime) {
		return "", fmt.Errorf("结束时间不能早于开始时间")
	}

	// 创建临时目录用于处理日志文件
	tempDir, err := os.MkdirTemp("", "logsnap_*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}
	// 函数结束时删除临时目录
	defer os.RemoveAll(tempDir)

	logrus.Infof("使用临时处理目录: %s", tempDir)

	// 查看有多少个处理器
	processorCount := len(c.logProcessors)
	logrus.Infof("有 %d 个解析器", processorCount)

	if processorCount == 0 {
		return "", fmt.Errorf("没有配置日志处理器")
	}

	// 创建等待组和结果通道
	var wg sync.WaitGroup
	resultChan := make(chan ProcessorResult, processorCount)

	// 启动多个goroutine并行处理日志
	for i, processor := range c.logProcessors {
		wg.Add(1)
		processorName := processor.GetName()
		logrus.Infof("准备启动协程 #%d 处理器: %s", i+1, processorName)

		go func(index int, p LogProcessor, name string) {
			defer func() {
				wg.Done()
			}()

			outputPath, results, err := p.Collect(startTime, endTime, tempDir)
			logrus.Debugf("协程 #%d Collect 方法返回，处理器: %s, 输出路径: %s, 结果数: %d, 错误: %v",
				index+1, name, outputPath, len(results), err)

			// 获取日志路径用于日志记录
			logPath, pathErr := p.GetLogPath()
			if pathErr != nil {
				logrus.Warnf("获取日志路径失败: %v", pathErr)
				logPath = "未知路径"
			}
			logrus.Debugf("协程 #%d 获取日志路径: %s", index+1, logPath)

			// 发送处理结果到通道
			logrus.Debugf("协程 #%d 准备发送结果到通道，处理器: %s", index+1, name)
			resultChan <- ProcessorResult{
				processorName: name,
				outputPath:    outputPath,
				results:       results,
				err:           err,
			}
			logrus.Debugf("协程 #%d 已发送结果到通道，处理器: %s", index+1, name)

			if err != nil {
				logrus.Warnf("解析器 %s 收集失败: %v", name, err)
			} else if outputPath != "" {
				logrus.Debugf("解析器 %s 已处理完成: 路径=%s, 结果数量=%d",
					name, outputPath, len(results))

				// 添加详细的处理结果日志
				totalLines := 0
				matchLines := 0
				for _, result := range results {
					totalLines += result.TotalLines
					matchLines += result.MatchLines
				}
				logrus.Infof("已处理 %s: 共 %d 行, 匹配 %d 条日志",
					logPath, totalLines, matchLines)
			} else {
				logrus.Warnf("日志文件不存在或为空: %s", logPath)
			}

			logrus.Debugf("===== 协程 #%d 处理器 %s 已完成工作 =====", index+1, name)
		}(i, processor, processorName)
	}

	// 启动一个goroutine等待所有处理器完成并关闭结果通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集处理结果
	totalLineCount := 0
	totalMatchCount := 0

	// 从通道读取结果
	for result := range resultChan {
		if result.err == nil && result.outputPath != "" {
			totalLineCount += result.GetTotalLines()
			totalMatchCount += result.GetMatchLines()
		}
	}

	hasFiles := totalLineCount > 0 && totalMatchCount > 0

	// 检查临时目录中的文件
	files, err := os.ReadDir(tempDir)
	if err != nil {
		logrus.Warnf("无法读取临时目录: %v", err)
	} else {
		logrus.Infof("临时目录中有 %d 个文件/目录:", len(files))
		for _, f := range files {
			info, _ := f.Info()
			if info != nil {
				logrus.Infof("  - %s (大小: %d 字节, 目录: %v)", f.Name(), info.Size(), f.IsDir())
			} else {
				logrus.Infof("  - %s", f.Name())
			}
		}
	}

	// 如果没有收集到任何文件，返回错误
	if !hasFiles {
		logrus.Infof("没有找到任何匹配的日志文件")
		return "", fmt.Errorf("指定时间范围内没有找到任何日志")
	}

	// 创建快照文件名
	snapFile := fmt.Sprintf("logsnap_%s_%s.zip",
		startTime.Format("20060102_150405"),
		endTime.Format("20060102_150405"))

	// 确定最终ZIP文件的路径
	var snapPath string
	if c.outputDir != "" {
		// 确保输出目录存在
		if err := os.MkdirAll(c.outputDir, 0755); err != nil {
			return "", fmt.Errorf("创建输出目录失败: %w", err)
		}
		snapPath = filepath.Join(c.outputDir, snapFile)
	} else {
		// 使用当前目录
		snapPath = snapFile
	}

	logrus.Infof("开始创建ZIP文件: %s", snapPath)

	// 压缩收集的日志
	err = utils.ZipDirectory(tempDir, snapPath)
	if err != nil {
		logrus.Errorf("创建ZIP文件失败: %v", err)
		return "", fmt.Errorf("创建日志快照失败: %w", err)
	}

	// 验证生成的ZIP文件
	fileInfo, err := os.Stat(snapPath)
	if err != nil {
		logrus.Errorf("无法访问生成的ZIP文件: %v", err)
		return "", fmt.Errorf("无法访问生成的ZIP文件: %w", err)
	}

	logrus.Infof("已创建ZIP文件: %s, 大小: %d 字节", snapPath, fileInfo.Size())

	if fileInfo.Size() == 0 {
		logrus.Warnf("生成的ZIP文件为空")
		os.Remove(snapPath)
		return "", fmt.Errorf("生成的ZIP文件为空")
	}

	// 尝试验证ZIP文件
	_, err = zip.OpenReader(snapPath)
	if err != nil {
		logrus.Warnf("无法打开生成的ZIP文件进行验证: %v", err)
	} else {
		logrus.Infof("ZIP文件验证成功")
	}

	return snapPath, nil
}
