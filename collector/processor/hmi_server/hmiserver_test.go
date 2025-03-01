package hmi_server

import (
	"archive/zip"
	"fmt"
	"io"
	"logsnap/collector/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestHMIServerLogProcessorWithRealSamples 使用真实日志样本测试处理器
func TestHMIServerLogProcessorWithRealSamples(t *testing.T) {
	// 确保测试开始时清空 testresult 目录
	clearTestResultDir(t, "testresult")

	// 使用testdata目录下的样本
	testRootDir := "testdata"

	// 创建输出目录
	outputDir, err := os.MkdirTemp("", "real_samples_output_")
	if err != nil {
		t.Fatalf("创建输出目录失败: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// 创建解析器
	processor := NewHMIServerLogProcessor(testRootDir, "hmiserver")

	// 设置测试时间范围（根据样本数据调整）
	// 修改时间范围，确保能够覆盖测试数据中的日志文件
	// 根据目录中的文件名，日志时间为 2025-02-28
	startTime := time.Date(2025, 2, 28, 12, 0, 0, 0, time.Local)         // 设置为当天开始
	endTime := time.Date(2025, 2, 28, 23, 59, 59, 999999999, time.Local) // 设置为当天结束

	// 执行解析
	outputPath, _, err := processor.Collect(startTime, endTime, outputDir)
	if err != nil {
		t.Fatalf("处理日志失败: %v", err)
	}

	// 将输出的目录拷贝至 result 目录
	copyOutputToResult(t, outputPath, "testresult")

	// 验证结果
	// 遍历 testresult 目录下的所有文件，找出最早时间然后与 startTime 比较
	earliestTime := time.Now() // 初始化为当前时间，确保任何过去的时间都会更早

	err = filepath.Walk("testresult", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 尝试从文件名解析时间
		filename := info.Name()
		var fileTime time.Time
		var parseErr error

		if strings.Contains(filename, "_") && strings.HasSuffix(filename, ".log") {
			// 格式1: name_YYYYMMDD_HHMMSS.log
			parts := strings.Split(filename, "_")
			if len(parts) >= 3 {
				dateStr := parts[len(parts)-2]
				timeStr := strings.TrimSuffix(parts[len(parts)-1], ".log")
				if len(dateStr) == 8 && len(timeStr) == 6 {
					// YYYYMMDD_HHMMSS
					timeString := fmt.Sprintf("%s %s", dateStr, timeStr)
					fileTime, parseErr = time.ParseInLocation("20060102 150405", timeString, time.Local)
				}
			}
		} else if strings.Contains(filename, ".") && strings.Contains(filename, "-") {
			// 格式2: name.YYYY-MM-DD_HH-MM-SS_SSSSSS.log - 使用工具函数解析
			parts := strings.SplitN(filename, ".", 2)
			if len(parts) == 2 {
				timeStr := strings.TrimSuffix(parts[1], ".log")
				fileTime, parseErr = utils.ParseArchiveTimeStamp(timeStr)
			}
		}

		if parseErr == nil && !fileTime.IsZero() {
			// 找到有效的时间
			if fileTime.Before(earliestTime) {
				earliestTime = fileTime
				t.Logf("找到更早的时间戳: %s 来自文件 %s", earliestTime.Format(time.RFC3339), path)
			}
		} else if parseErr != nil {
			t.Logf("无法从文件 %s 解析时间: %v", path, parseErr)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("遍历testresult目录失败: %v", err)
	}

	// 验证最早时间是否等于或晚于startTime
	if earliestTime.IsZero() {
		t.Error("未找到任何带有有效时间戳的文件")
	} else {
		if earliestTime.Before(startTime) {
			t.Errorf("找到的最早时间 %v 早于开始时间 %v",
				earliestTime.Format(time.RFC3339),
				startTime.Format(time.RFC3339))
		} else {
			t.Logf("验证通过: 找到的最早时间 %v 不早于开始时间 %v",
				earliestTime.Format(time.RFC3339),
				startTime.Format(time.RFC3339))
		}
	}

	// 验证 all 和 sqlalchemy 目录中是否有文件
	verifyDirectoryHasFiles(t, filepath.Join("testresult", "all"), "all 目录")
	verifyDirectoryHasFiles(t, filepath.Join("testresult", "sqlalchemy"), "sqlalchemy 目录")
}

// setupTestEnvironment 设置测试环境，创建必要的目录和测试日志文件
func setupTestEnvironment(t *testing.T, rootDir string) {
	// 创建子目录结构
	dirs := []string{"all", "error", "debug", "wcs"}
	for _, dir := range dirs {
		dirPath := filepath.Join(rootDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("创建测试目录 %s 失败: %v", dirPath, err)
		}

		// 为每个目录创建当前日志文件和归档日志文件
		path := filepath.Join(dirPath, dir+".log")
		createTestLogFile(t, path, 100, time.Now().Add(-30*time.Minute))
		t.Logf("创建测试日志文件: %s", path)

		// 创建一个归档文件
		archiveTime := time.Now().Add(-2 * time.Hour)
		archiveName := filepath.Join(dirPath,
			dir+"."+archiveTime.Format("2006-01-02_15-04-05")+"_000000.log.zip")
		createTestArchiveFile(t, archiveName, dir+".log", 50, archiveTime)
		t.Logf("创建测试归档文件: %s", archiveName)
	}
}

// createTestLogFile 创建测试日志文件
func createTestLogFile(t *testing.T, path string, lineCount int, startTime time.Time) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("创建测试日志文件 %s 失败: %v", path, err)
	}
	defer file.Close()

	// 写入测试日志行
	timeIncrement := time.Minute
	currentTime := startTime
	for i := 0; i < lineCount; i++ {
		// 确保日志格式与 logTimePattern 正则表达式完全匹配
		// 原正则: ^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) \|
		logLine := currentTime.Format("2006-01-02 15:04:05.000") + " | " +
			"INFO | HMIServer | Module | 这是测试日志行 " + strconv.Itoa(i) +
			" | SessionID=12345 | UserID=test_user\n"
		file.WriteString(logLine)
		currentTime = currentTime.Add(timeIncrement)
	}
}

// createTestArchiveFile 创建测试归档文件
func createTestArchiveFile(t *testing.T, archivePath, logFileName string, lineCount int, startTime time.Time) {
	// 首先创建一个临时日志文件
	tempDir, err := os.MkdirTemp("", "archive_temp_")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFilePath := filepath.Join(tempDir, logFileName)
	createTestLogFile(t, logFilePath, lineCount, startTime)

	// 创建zip文件
	zipFile, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("创建ZIP文件 %s 失败: %v", archivePath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 添加日志文件到zip
	logFileContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("读取日志文件 %s 失败: %v", logFilePath, err)
	}

	zipFileWriter, err := zipWriter.Create(logFileName)
	if err != nil {
		t.Fatalf("创建ZIP条目失败: %v", err)
	}

	_, err = zipFileWriter.Write(logFileContent)
	if err != nil {
		t.Fatalf("写入ZIP内容失败: %v", err)
	}
}

// verifyOutputStructure 验证输出的目录结构和文件内容
func verifyOutputStructure(t *testing.T, outputPath string) {
	// 验证输出目录是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("输出目录不存在: %s", outputPath)
	}

	// 检查是否创建了对应的子目录
	expectedDirs := []string{"all", "error", "debug", "wcs"}
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(outputPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("期望的输出子目录不存在: %s", dirPath)
			continue
		}

		// 检查目录中是否有输出文件 - 如果没有匹配项，这里可能没有文件
		files, err := os.ReadDir(dirPath)
		if err != nil {
			t.Errorf("读取输出目录 %s 失败: %v", dirPath, err)
			continue
		}

		// 只有当目录中有文件时才检查文件内容
		if len(files) > 0 {
			for _, file := range files {
				t.Logf("发现输出文件: %s", filepath.Join(dirPath, file.Name()))

				// 读取文件内容并验证
				content, err := os.ReadFile(filepath.Join(dirPath, file.Name()))
				if err != nil {
					t.Errorf("读取输出文件失败: %v", err)
					continue
				}

				// 验证文件内容包含预期的头信息
				if !strings.Contains(string(content), "# 原始日志文件:") {
					t.Errorf("输出文件缺少预期的头信息")
				}

				t.Logf("文件 %s 的内容验证通过", file.Name())
			}
		} else {
			t.Logf("输出目录 %s 中没有文件，可能是因为没有匹配的日志行", dirPath)
		}
	}
}

// copyOutputToResult 将输出目录下的内容拷贝到指定的结果目录
func copyOutputToResult(t *testing.T, outputPath, resultDirName string) {
	// 如果输出路径为空或不存在，则直接返回
	if outputPath == "" {
		t.Logf("输出路径为空，不执行拷贝")
		return
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Logf("输出路径不存在，不执行拷贝: %s", outputPath)
		return
	}

	// 确定目标目录路径（在当前包目录下）
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前工作目录失败: %v", err)
	}

	resultDir := filepath.Join(currentDir, resultDirName)

	// 清理之前的结果目录（如果存在）
	if _, err := os.Stat(resultDir); err == nil {
		if err := os.RemoveAll(resultDir); err != nil {
			t.Fatalf("清理旧的结果目录失败: %v", err)
		}
	}

	// 创建结果目录
	if err := os.MkdirAll(resultDir, 0755); err != nil {
		t.Fatalf("创建结果目录失败: %v", err)
	}

	// 递归拷贝文件
	err = filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算目标路径（相对路径部分）
		relPath, err := filepath.Rel(outputPath, path)
		if err != nil {
			return fmt.Errorf("计算相对路径失败: %w", err)
		}

		// 目标完整路径
		destPath := filepath.Join(resultDir, relPath)

		// 如果是目录，创建对应的目录
		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// 如果是文件，拷贝文件
		return copyFile(path, destPath)
	})

	if err != nil {
		t.Fatalf("拷贝输出文件到结果目录失败: %v", err)
	}

	t.Logf("成功将输出文件拷贝到结果目录: %s", resultDir)
}

// copyFile 拷贝单个文件
func copyFile(src, dest string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	// 拷贝文件内容
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("拷贝文件内容失败: %w", err)
	}

	// 确保写入到磁盘
	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("同步文件到磁盘失败: %w", err)
	}

	return nil
}

// clearTestResultDir 确保测试结果目录是空的
func clearTestResultDir(t *testing.T, resultDirName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前工作目录失败: %v", err)
	}

	resultDir := filepath.Join(currentDir, resultDirName)

	// 检查目录是否存在
	if _, err := os.Stat(resultDir); err == nil {
		// 目录存在，清空它
		t.Logf("清空测试结果目录: %s", resultDir)
		if err := os.RemoveAll(resultDir); err != nil {
			t.Fatalf("清理测试结果目录失败: %v", err)
		}
	}

	// 创建空目录
	if err := os.MkdirAll(resultDir, 0755); err != nil {
		t.Fatalf("创建测试结果目录失败: %v", err)
	}
}

// verifyDirectoryHasFiles 验证目录中是否有文件
func verifyDirectoryHasFiles(t *testing.T, dir string, dirName string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Errorf("读取 %s 目录失败: %v", dirName, err)
		return
	}

	if len(files) == 0 {
		t.Errorf("%s 目录中没有文件", dirName)
	} else {
		t.Logf("%s 目录中有 %d 个文件", dirName, len(files))
		for _, file := range files {
			t.Logf("  - %s", file.Name())
		}
	}
}
