package processor

import (
	"fmt"
	"io"
	collector "logsnap/collector"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// DefaultFindLogFiles 查找指定目录下的所有日志文件
// suffixes 为可选参数，如果为空，则查找所有日志文件
func DefaultFindLogFiles(dirPath string, suffixes ...string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// 如果目录不存在，返回空切片
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 如果suffixes为空，则将所有文件添加到files中
		if len(suffixes) == 0 {
			files = append(files, filepath.Join(dirPath, entry.Name()))
			continue
		}

		name := entry.Name()
		for _, suffix := range suffixes {
			if strings.HasSuffix(name, suffix) {
				files = append(files, filepath.Join(dirPath, name))
				break
			}
		}
	}

	return files, nil
}


// DefaultProcessDir 处理目录中的日志文件
func DefaultProcessDir(provider FileProcessorProvider, dirPath, outputDir string, startTime, endTime time.Time) ([]collector.FileProcessResult, error) {
	// 查找子目录下所有日志文件
	logFiles, err := provider.FindFiles(dirPath, provider.GetFileSuffixes()...)
	if err != nil {
		return nil, fmt.Errorf("查找目录 %s 下的日志文件失败: %w", dirPath, err)
	}

	// 如果没有日志文件，跳过此目录
	if len(logFiles) == 0 {
		return nil, nil
	}

	// 解析文件信息并筛选时间范围内的文件
	fileInfos, err := provider.FilterFiles(logFiles, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("分析目录 %s 下的日志文件信息失败: %w", dirPath, err)
	}

	// 如果没有符合条件的文件，直接返回
	if len(fileInfos) == 0 {
		return nil, nil
	}

	// 确定工作线程数量
	workerCount := Min(len(fileInfos), runtime.NumCPU())
	if workerCount < 1 {
		workerCount = 1
	}

	logrus.Infof("使用 %d 个工作线程处理 %d 个文件", workerCount, len(fileInfos))

	// 创建任务通道和结果通道
	taskChan := make(chan LogFileInfo, len(fileInfos))
	resultChan := make(chan collector.FileProcessResult, len(fileInfos))

	// 创建等待组
	var wg sync.WaitGroup

	// 启动工作线程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			// 工作线程从任务通道获取文件并处理
			for fileInfo := range taskChan {
				// 获取文件名（假设T类型有某种方式获取名称）
				fileName := fileInfo.FileName
				logrus.Debugf("工作线程 %d 开始处理文件: %s", workerId, fileName)

				result, err := provider.ProcessFile(fileInfo, startTime, endTime, outputDir)
				if err != nil {
					logrus.Errorf("处理文件 %s 失败: %v", fileName, err)
					result.Err = err
				}
				// 发送处理结果到结果通道
				resultChan <- result

				logrus.Debugf("工作线程 %d 完成处理文件: %s", workerId, fileName)
			}
		}(i)
	}

	// 将所有文件发送到任务通道
	for _, fileInfo := range fileInfos {
		taskChan <- fileInfo
	}
	close(taskChan)

	// 等待所有工作线程完成并关闭结果通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集处理结果
	var lastError error
	results := make([]collector.FileProcessResult, 0)

	for result := range resultChan {
		if result.Err != nil {
			logrus.Errorf("处理文件 %s 失败: %v", result.FilePath, result.Err)
			lastError = result.Err
			continue
		}
		results = append(results, result)
	}

	if lastError != nil {
		return nil, lastError
	}

	return results, nil
}

// FileNameProcessor 是一个函数类型，用于处理文件名
type FileNameProcessor func(fileInfo LogFileInfo) string

// DefaultGenerateOutputFileName 根据文件类型生成默认的输出文件名
func DefaultGenerateOutputFileName(fileInfo LogFileInfo) string {
	return fileInfo.FileName
}


// ReaderCreator 是一个函数类型，用于创建文件读取器
type ReaderCreator func(fileInfo LogFileInfo) (io.ReadCloser, error)

// DefaultCreateReaderForFile 创建默认的文件读取器
// 默认按文本文件读取
func DefaultCreateReaderForFile(fileInfo LogFileInfo) (io.ReadCloser, error) {
	file, err := os.Open(fileInfo.Path)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}
	return file, nil
}