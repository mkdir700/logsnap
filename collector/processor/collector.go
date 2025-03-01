package processor

import (
	"fmt"
	collector "logsnap/collector"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// CollectWithProcessor 通用的日志收集方法
// 使用提供的文件处理器处理日志文件
// 参数:
//   - p: 基础处理器
//   - fileProcessor: 文件处理器
//   - startTime: 开始时间
//   - endTime: 结束时间
//   - outputDir: 输出目录
//
// 返回:
//   - outputPath: 收集的日志文件路径
//   - lineCount: 收集的日志行数
//   - matchCount: 匹配的日志行数
//   - error: 错误信息
func CollectWithProcessor(
	p collector.LogProcessor,
	fileProcessors []FileProcessorProvider,
	startTime, endTime time.Time,
	outputDir string,
) (string, []collector.FileProcessResult, error) {
	logrus.Infof("创建日志输出目录: %s", outputDir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", []collector.FileProcessResult{}, fmt.Errorf("创建日志输出目录失败: %w", err)
	}

	logPath, err := p.GetLogPath()
	if err != nil {
		return "", []collector.FileProcessResult{}, fmt.Errorf("获取日志路径失败: %w", err)
	}

	var results []collector.FileProcessResult
	// 递归处理目录
	for _, fileProcessor := range fileProcessors {
		err = processDirectory(logPath, outputDir, fileProcessor, startTime, endTime, &results)
		if err != nil {
			logrus.Errorf("处理目录时发生 %v", err)
		}
	}

	return outputDir, results, nil
}

// processDirectory 递归处理目录中的文件
// 参数:
//   - dirPath: 要处理的目录路径
//   - outputDir: 输出目录路径
//   - fileProcessor: 文件处理器
//   - startTime: 开始时间
//   - endTime: 结束时间
//   - totalLineCount: 总行数计数器的指针
//   - totalMatchCount: 匹配行数计数器的指针
//
// 返回:
//   - error: 错误信息
func processDirectory(
	dirPath, outputDir string,
	fileProcessor FileProcessorProvider,
	startTime, endTime time.Time,
	results *[]collector.FileProcessResult,
) error {
	// 处理当前目录下的文件
	_results, err := fileProcessor.ProcessDir(dirPath, outputDir, startTime, endTime)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("目录 %s 不存在", dirPath)
			return nil
		}
		logrus.Errorf("处理目录 %s 失败: %v\n", dirPath, err)
	} else {
		*results = append(*results, _results...)
	}

	// 获取当前目录下的所有条目
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("目录 %s 不存在", dirPath)
			return nil
		}
		return fmt.Errorf("读取目录 %s 失败: %w", dirPath, err)
	}

	// 递归处理每个子目录
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}

		subDirName := entry.Name()
		subDirPath := filepath.Join(dirPath, subDirName)

		// 创建对应的输出子目录
		subOutputDir := filepath.Join(outputDir, subDirName)
		if err := os.MkdirAll(subOutputDir, 0755); err != nil {
			logrus.Errorf("创建输出子目录 %s 失败: %v\n", subOutputDir, err)
			continue
		}

		// 递归处理子目录
		if err := processDirectory(subDirPath, subOutputDir, fileProcessor, startTime, endTime, results); err != nil {
			logrus.Errorf("处理子目录 %s 失败: %v\n", subDirPath, err)
		}
	}

	return nil
}
