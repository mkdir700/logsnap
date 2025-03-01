package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ZipDirectory 将指定目录下的所有文件打包成zip文件
// 使用并发处理提高性能
func ZipDirectory(sourceDir, destZip string, concurrency ...int) error {
	startTime := time.Now()
	
	// 设置默认并发数为CPU核心数
	workerCount := runtime.NumCPU()
	if len(concurrency) > 0 && concurrency[0] > 0 {
		workerCount = concurrency[0]
	}

	// 确保源目录存在
	info, err := os.Stat(sourceDir)
	if err != nil {
		return fmt.Errorf("无法访问源目录: %w", err)
	}
	
	if !info.IsDir() {
		return fmt.Errorf("源路径不是目录: %s", sourceDir)
	}
	
	// 获取源目录的绝对路径
	absSourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("获取源目录绝对路径失败: %w", err)
	}

	logrus.Infof("开始压缩目录: %s (并发数: %d)", absSourceDir, workerCount)

	// 确保目标文件不存在
	if _, err := os.Stat(destZip); err == nil {
		if err := os.Remove(destZip); err != nil {
			return fmt.Errorf("无法删除已存在的目标文件: %w", err)
		}
	}
	
	// 收集目录中的所有文件
	var filesToZip []string
	err = filepath.Walk(absSourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 跳过目录本身
		if info.IsDir() {
			return nil
		}
		
		filesToZip = append(filesToZip, path)
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("遍历目录失败: %w", err)
	}

	fileCount := len(filesToZip)
	logrus.Infof("找到 %d 个文件需要压缩", fileCount)
	
	if fileCount == 0 {
		logrus.Warnf("没有找到需要压缩的文件")
		// 创建一个空的ZIP文件
		emptyZip, err := os.Create(destZip)
		if err != nil {
			return fmt.Errorf("创建空ZIP文件失败: %w", err)
		}
		defer emptyZip.Close()
		
		zipWriter := zip.NewWriter(emptyZip)
		err = zipWriter.Close()
		if err != nil {
			return fmt.Errorf("关闭空ZIP writer失败: %w", err)
		}
		
		return nil
	}

	// 调整并发数，避免过多的并发
	maxOpenFiles := 100 // 设置一个合理的同时打开文件数上限
	if workerCount > maxOpenFiles {
		workerCount = maxOpenFiles
	}
	
	if fileCount < workerCount {
		workerCount = fileCount
	}
	
	logrus.Infof("使用 %d 个工作协程进行压缩", workerCount)

	// 创建临时目录存放中间文件
	tempDir, err := os.MkdirTemp("", "zip_temp_*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// 将文件分成多个批次
	batchSize := (fileCount + workerCount - 1) / workerCount
	batches := make([][]string, 0, workerCount)
	
	for i := 0; i < fileCount; i += batchSize {
		end := i + batchSize
		if end > fileCount {
			end = fileCount
		}
		batches = append(batches, filesToZip[i:end])
	}
	
	// 为每个批次创建一个临时ZIP文件
	var wg sync.WaitGroup
	tempZips := make([]string, len(batches))
	errors := make([]error, len(batches))
	
	for i, batch := range batches {
		wg.Add(1)
		go func(batchIndex int, fileBatch []string) {
			defer wg.Done()
			
			// 创建临时ZIP文件
			tempZipPath := filepath.Join(tempDir, fmt.Sprintf("temp_%d.zip", batchIndex))
			tempZips[batchIndex] = tempZipPath
			
			// 创建ZIP文件
			tempZipFile, err := os.Create(tempZipPath)
			if err != nil {
				errors[batchIndex] = fmt.Errorf("创建临时ZIP文件失败: %w", err)
				return
			}
			
			zipWriter := zip.NewWriter(tempZipFile)
			
			// 处理批次中的每个文件
			for _, filePath := range fileBatch {
				// 获取相对路径
				relPath, err := filepath.Rel(absSourceDir, filePath)
				if err != nil {
					logrus.Warnf("获取相对路径失败 %s: %v", filePath, err)
					continue
				}
				
				// 打开源文件
				file, err := os.Open(filePath)
				if err != nil {
					logrus.Warnf("无法打开文件 %s: %v", filePath, err)
					continue
				}
				
				// 使用匿名函数确保文件被关闭
				func() {
					defer file.Close()
					
					// 获取文件信息
					info, err := file.Stat()
					if err != nil {
						logrus.Warnf("无法获取文件信息 %s: %v", filePath, err)
						return
					}
					
					// 创建zip文件头
					header, err := zip.FileInfoHeader(info)
					if err != nil {
						logrus.Warnf("无法创建文件头 %s: %v", filePath, err)
						return
					}
					
					// 设置文件名为相对路径，并统一使用斜杠作为分隔符
					header.Name = filepath.ToSlash(relPath)
					header.Method = zip.Deflate
					
					// 创建文件条目
					writer, err := zipWriter.CreateHeader(header)
					if err != nil {
						logrus.Warnf("无法创建ZIP条目 %s: %v", filePath, err)
						return
					}
					
					// 复制文件内容到zip
					_, err = io.Copy(writer, file)
					if err != nil {
						logrus.Warnf("无法写入文件内容 %s: %v", filePath, err)
						return
					}
					
					logrus.Debugf("工作协程 %d: 已添加文件: %s", batchIndex, relPath)
				}()
			}
			
			// 关闭ZIP writer
			if err := zipWriter.Close(); err != nil {
				errors[batchIndex] = fmt.Errorf("关闭临时ZIP writer失败: %w", err)
				tempZipFile.Close()
				return
			}
			
			// 关闭临时文件
			if err := tempZipFile.Close(); err != nil {
				errors[batchIndex] = fmt.Errorf("关闭临时ZIP文件失败: %w", err)
				return
			}
			
			logrus.Infof("工作协程 %d: 已完成 %d 个文件的压缩", batchIndex, len(fileBatch))
		}(i, batch)
	}
	
	// 等待所有工作协程完成
	wg.Wait()
	
	// 检查是否有错误
	for i, err := range errors {
		if err != nil {
			return fmt.Errorf("批次 %d 处理失败: %w", i, err)
		}
	}
	
	// 合并所有临时ZIP文件
	err = mergeZipFiles(tempZips, destZip)
	if err != nil {
		return fmt.Errorf("合并ZIP文件失败: %w", err)
	}
	
	// 验证生成的文件
	fileInfo, err := os.Stat(destZip)
	if err != nil {
		return fmt.Errorf("无法访问生成的ZIP文件: %w", err)
	}
	
	elapsedTime := time.Since(startTime)
	logrus.Infof("ZIP文件创建成功: %s (大小: %d 字节, 耗时: %v)", destZip, fileInfo.Size(), elapsedTime)
	
	return nil
}

// mergeZipFiles 合并多个ZIP文件到一个目标文件
func mergeZipFiles(sourceZips []string, destZip string) error {
	// 创建目标文件
	destFile, err := os.Create(destZip)
	if err != nil {
		return fmt.Errorf("创建目标ZIP文件失败: %w", err)
	}
	defer destFile.Close()
	
	// 创建新的zip writer
	destZipWriter := zip.NewWriter(destFile)
	defer destZipWriter.Close()
	
	// 处理每个源ZIP文件
	for _, sourceZip := range sourceZips {
		// 打开源ZIP文件
		reader, err := zip.OpenReader(sourceZip)
		if err != nil {
			return fmt.Errorf("打开源ZIP文件失败 %s: %w", sourceZip, err)
		}
		
		// 复制每个文件到目标ZIP
		for _, file := range reader.File {
			err = copyZipFile(file, destZipWriter)
			if err != nil {
				reader.Close()
				return fmt.Errorf("复制ZIP文件条目失败 %s: %w", file.Name, err)
			}
		}
		
		reader.Close()
	}
	
	return nil
}

// copyZipFile 从源ZIP文件复制一个文件到目标ZIP writer
func copyZipFile(file *zip.File, destZipWriter *zip.Writer) error {
	// 打开源文件
	sourceFile, err := file.Open()
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	// 创建目标文件条目
	destFile, err := destZipWriter.CreateHeader(&file.FileHeader)
	if err != nil {
		return err
	}
	
	// 复制内容
	_, err = io.Copy(destFile, sourceFile)
	return err
}
