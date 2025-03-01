package hmi

import (
	"fmt"
	"io"
	processor "logsnap/collector/processor"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// 测试常量定义
const (
	testResultDir = "testresult"
)

// setupTestEnvironment 设置测试环境，创建必要的目录和测试日志文件
func setupTestEnvironment(t *testing.T, rootDir string) {
	// 创建测试目录
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		t.Fatalf("创建测试目录 %s 失败: %v", rootDir, err)
	}

	// 创建程序日志文件
	programLogName := "xyz_hmi_bin.xyz-Workstation.xyz.log.ERROR.20250228-094825.2966778"
	programLogPath := filepath.Join(rootDir, programLogName)
	createProgramLogFile(t, programLogPath, 100)

	// 创建用户操作日志文件
	userOpDirPath := filepath.Join(rootDir, "useroperations")
	if err := os.MkdirAll(userOpDirPath, 0755); err != nil {
		t.Fatalf("创建用户操作日志目录失败: %v", err)
	}

	userOpLogName := "20250302-084053.269"
	userOpLogPath := filepath.Join(userOpDirPath, userOpLogName)
	createUserOpLogFile(t, userOpLogPath, 50)
}

// createProgramLogFile 创建测试程序日志文件
func createProgramLogFile(t *testing.T, path string, lineCount int) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("创建测试程序日志文件 %s 失败: %v", path, err)
	}
	defer file.Close()

	// 写入日志头部
	file.WriteString("Log file created at: 2025/02/28 09:48:25\n")
	file.WriteString("Running on machine: xyz-Workstation\n")
	file.WriteString("Running duration (h:mm:ss): 0:00:00\n")
	file.WriteString("Log line format: [IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line] msg\n")

	// 写入测试日志行
	startTime := time.Date(2025, 2, 28, 9, 48, 25, 0, time.Local)
	timeIncrement := time.Second
	currentTime := startTime

	for i := 0; i < lineCount; i++ {
		// 格式：E20250228 09:48:25.654057 2966778 station_status_label.cpp:53] 0:  0
		logLine := fmt.Sprintf("E%s %06d station_status_label.cpp:%d] 测试日志行 %d\n",
			currentTime.Format("20060102 15:04:05.000000"),
			2966778+i,
			53+i,
			i)
		file.WriteString(logLine)
		currentTime = currentTime.Add(timeIncrement)
	}
}

// createUserOpLogFile 创建测试用户操作日志文件
func createUserOpLogFile(t *testing.T, path string, lineCount int) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("创建测试用户操作日志文件 %s 失败: %v", path, err)
	}
	defer file.Close()

	// 写入测试日志行
	startTime := time.Date(2025, 3, 2, 8, 40, 53, 0, time.Local)
	timeIncrement := time.Second
	currentTime := startTime

	for i := 0; i < lineCount; i++ {
		// 格式：20250302 08:41:13.163] User clicked [StartTask].
		logLine := fmt.Sprintf("%s] 用户点击了 [操作%d].\n",
			currentTime.Format("20060102 15:04:05.000"),
			i)
		file.WriteString(logLine)
		currentTime = currentTime.Add(timeIncrement)
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

// verifyOutputFiles 验证输出文件
func verifyOutputFiles(t *testing.T, resultDir string) {
	files, err := os.ReadDir(resultDir)
	if err != nil {
		t.Errorf("读取结果目录 %s 失败: %v", resultDir, err)
		return
	}

	if len(files) == 0 {
		t.Errorf("结果目录中没有文件")
		return
	}

	// 检查是否有程序日志和用户操作日志
	var hasProgramLog, hasUserOpLog bool

	for _, file := range files {
		fileName := file.Name()
		t.Logf("发现结果文件: %s", fileName)

		// useroperations 目录下的文件
		if strings.Contains(fileName, "useroperations") {
			opfiles, err := os.ReadDir(filepath.Join(resultDir, fileName))
			if err != nil {
				t.Errorf("读取用户操作日志目录 %s 失败: %v", fileName, err)
			}
			if len(opfiles) > 0 {
				hasUserOpLog = true
			}
		}

		if strings.Contains(fileName, ".ERROR.") {
			hasProgramLog = true
		}
	}

	if !hasProgramLog {
		t.Errorf("未找到程序日志文件")
	}

	if !hasUserOpLog {
		t.Errorf("未找到用户操作日志文件")
	}
}

// TestHMILogProcessor_WithRealSamples 使用真实日志样本测试处理器
func TestHMILogProcessor_WithRealSamples(t *testing.T) {
	// 确保测试开始时清空 testresult 目录
	clearTestResultDir(t, testResultDir)

	// 使用testdata目录下的样本
	testRootDir := "testdata"

	// 创建输出目录
	outputDir, err := os.MkdirTemp("", "hmi_real_samples_output_")
	if err != nil {
		t.Fatalf("创建输出目录失败: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// 创建处理器
	hmiProcessor := &HMILogProcessor{
		BaseProcessor: processor.NewBaseProcessor("HMI日志", testRootDir, "hmi"),
	}

	// 设置测试时间范围（根据样本数据调整）
	startTime := time.Date(2025, 2, 28, 12, 0, 0, 0, time.Local)
	endTime := time.Date(2025, 3, 2, 23, 59, 59, 999999999, time.Local)

	// 执行解析
	outputPath, _, err := hmiProcessor.Collect(startTime, endTime, outputDir)
	if err != nil {
		t.Fatalf("处理日志失败: %v", err)
	}

	// 将输出的目录拷贝至 result 目录
	copyOutputToResult(t, outputPath, testResultDir)

	// 验证结果
	verifyOutputFiles(t, testResultDir)
}

// TestHMILogProcessor_FindFiles 测试查找文件功能
func TestHMILogProcessor_FindFiles(t *testing.T) {
	// 创建临时测试目录
	testDir, err := os.MkdirTemp("", "hmi_find_files_test_")
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 设置测试环境
	setupTestEnvironment(t, testDir)

	// 创建处理器
	hmiProcessor := &HMILogProcessor{
		BaseProcessor: processor.NewBaseProcessor("HMI日志", testDir, "hmi"),
	}

	// 创建文件处理器
	fileProcessors := hmiProcessor.CreateFileProcessor()
	if len(fileProcessors) == 0 {
		t.Fatalf("未创建文件处理器")
	}

	// 记录是否至少有一个处理器找到了文件
	foundAnyFiles := false

	// 测试每个文件处理器的FindFiles方法
	for i, fileProcessor := range fileProcessors {
		files, err := fileProcessor.FindFiles(testDir)
		if err != nil {
			t.Fatalf("处理器%d查找文件失败: %v", i, err)
		}

		t.Logf("处理器%d找到的文件: %v", i, files)
		if len(files) > 0 {
			foundAnyFiles = true
			t.Logf("处理器%d找到了%d个文件", i, len(files))
		}
	}

	// 验证总体结果
	if !foundAnyFiles {
		t.Errorf("所有处理器都未找到文件")
	}
}

// TestHMILogProcessor_FilterFiles 测试过滤文件功能
func TestHMILogProcessor_FilterFiles(t *testing.T) {
	// 创建临时测试目录
	testDir, err := os.MkdirTemp("", "hmi_filter_files_test_")
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 设置测试环境
	setupTestEnvironment(t, testDir)

	// 创建处理器
	hmiProcessor := &HMILogProcessor{
		BaseProcessor: processor.NewBaseProcessor("HMI日志", testDir, "hmi"),
	}

	// 创建文件处理器
	fileProcessors := hmiProcessor.CreateFileProcessor()
	if len(fileProcessors) == 0 {
		t.Fatalf("未创建文件处理器")
	}

	// 设置时间范围
	startTime := time.Date(2025, 2, 27, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2025, 3, 3, 0, 0, 0, 0, time.Local)

	// 记录至少有一个处理器找到了文件
	foundAnyFiles := false
	// 记录至少有一个处理器过滤出了文件
	filteredAnyFiles := false

	// 测试每个文件处理器的FindFiles和FilterFiles方法
	for i, fileProcessor := range fileProcessors {
		files, err := fileProcessor.FindFiles(testDir)
		if err != nil {
			t.Fatalf("处理器%d查找文件失败: %v", i, err)
		}

		if len(files) > 0 {
			foundAnyFiles = true
			t.Logf("处理器%d找到了%d个文件", i, len(files))
			logrus.Infof("处理器%d过滤文件: %v", i, files)

			// 过滤文件
			fileInfos, err := fileProcessor.FilterFiles(files, startTime, endTime)
			if err != nil {
				t.Fatalf("处理器%d过滤文件失败: %v", i, err)
			}

			t.Logf("处理器%d过滤后的文件: %v", i, fileInfos)

			if len(fileInfos) > 0 {
				filteredAnyFiles = true
				t.Logf("处理器%d成功过滤出%d个文件", i, len(fileInfos))
			}
		} else {
			t.Logf("处理器%d未找到文件，跳过过滤测试", i)
		}
	}

	// 验证总体结果
	if !foundAnyFiles {
		t.Errorf("所有处理器都未找到文件")
	}

	if foundAnyFiles && !filteredAnyFiles {
		t.Errorf("找到了文件但所有处理器过滤后都没有匹配的文件")
	}
}

// TestHMILogProcessor_ProcessFile 测试处理单个文件功能
func TestHMILogProcessor_ProcessFile(t *testing.T) {
	// 创建临时测试目录
	testDir, err := os.MkdirTemp("", "hmi_process_file_test_")
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 设置测试环境
	setupTestEnvironment(t, testDir)

	// 创建输出目录
	outputDir, err := os.MkdirTemp("", "hmi_process_file_output_")
	if err != nil {
		t.Fatalf("创建输出目录失败: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// 创建处理器
	hmiProcessor := &HMILogProcessor{
		BaseProcessor: processor.NewBaseProcessor("HMI日志", testDir, "hmi"),
	}

	// 创建文件处理器
	fileProcessors := hmiProcessor.CreateFileProcessor()
	if len(fileProcessors) == 0 {
		t.Fatalf("未创建文件处理器")
	}

	// 设置时间范围
	startTime := time.Date(2025, 2, 27, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2025, 3, 3, 0, 0, 0, 0, time.Local)

	// 记录是否至少有一个处理器成功处理了文件
	processedAnyFile := false

	// 测试每个文件处理器的处理文件功能
	for i, fileProcessor := range fileProcessors {
		files, err := fileProcessor.FindFiles(testDir)
		if err != nil {
			t.Fatalf("处理器%d查找文件失败: %v", i, err)
		}

		if len(files) == 0 {
			t.Logf("处理器%d未找到文件，跳过处理测试", i)
			continue
		}

		// 过滤文件
		fileInfos, err := fileProcessor.FilterFiles(files, startTime, endTime)
		if err != nil {
			t.Fatalf("处理器%d过滤文件失败: %v", i, err)
		}

		if len(fileInfos) == 0 {
			t.Logf("处理器%d过滤后没有文件，跳过处理测试", i)
			continue
		}

		// 处理第一个文件
		result, err := fileProcessor.ProcessFile(fileInfos[0], startTime, endTime, outputDir)
		if err != nil {
			t.Fatalf("处理器%d处理文件失败: %v", i, err)
		}

		processedAnyFile = true
		t.Logf("处理器%d处理文件结果: 输出文件=%s, 总行数=%d, 匹配行数=%d",
			i, result.FilePath, result.TotalLines, result.MatchLines)

		// 只有当有匹配的行时，才验证输出文件是否存在
		if result.MatchLines > 0 {
			if _, err := os.Stat(result.FilePath); os.IsNotExist(err) {
				t.Errorf("处理器%d的输出文件不存在: %s", i, result.FilePath)
			}
		} else {
			t.Logf("处理器%d没有匹配的行，跳过验证输出文件", i)
		}
	}

	// 验证总体结果
	if !processedAnyFile {
		t.Errorf("所有处理器都未能成功处理文件")
	}
}

// TestHMILogProcessor_ProcessDir 测试处理目录功能
func TestHMILogProcessor_ProcessDir(t *testing.T) {
	// 创建临时测试目录
	testDir, err := os.MkdirTemp("", "hmi_process_dir_test_")
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 设置测试环境
	setupTestEnvironment(t, testDir)

	// 创建输出目录
	outputDir, err := os.MkdirTemp("", "hmi_process_dir_output_")
	if err != nil {
		t.Fatalf("创建输出目录失败: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// 创建处理器
	hmiProcessor := &HMILogProcessor{
		BaseProcessor: processor.NewBaseProcessor("HMI日志", testDir, "hmi"),
	}

	// 创建文件处理器
	fileProcessors := hmiProcessor.CreateFileProcessor()
	if len(fileProcessors) == 0 {
		t.Fatalf("未创建文件处理器")
	}

	// 设置时间范围
	startTime := time.Date(2025, 2, 27, 0, 0, 0, 0, time.Local)
	endTime := time.Date(2025, 3, 3, 0, 0, 0, 0, time.Local)

	// 记录是否至少有一个处理器成功处理了目录
	processedAnyDir := false

	// 测试每个文件处理器的处理目录功能
	for i, fileProcessor := range fileProcessors {
		// 处理目录
		results, err := fileProcessor.ProcessDir(testDir, outputDir, startTime, endTime)
		if err != nil {
			t.Fatalf("处理器%d处理目录失败: %v", i, err)
		}

		t.Logf("处理器%d处理目录结果: 处理了%d个文件", i, len(results))

		if len(results) > 0 {
			processedAnyDir = true
			// 验证每个处理结果
			for j, result := range results {
				t.Logf("处理器%d的结果%d: 输出文件=%s, 总行数=%d, 匹配行数=%d",
					i, j, result.FilePath, result.TotalLines, result.MatchLines)

				// 只有当有匹配的行时，才验证输出文件是否存在
				if result.MatchLines > 0 {
					if _, err := os.Stat(result.FilePath); os.IsNotExist(err) {
						t.Errorf("处理器%d的结果%d的输出文件不存在: %s", i, j, result.FilePath)
					}
				} else {
					t.Logf("处理器%d的结果%d没有匹配的行，跳过验证输出文件", i, j)
				}
			}
		}
	}

	// 验证总体结果
	if !processedAnyDir {
		t.Errorf("所有处理器都未能成功处理目录")
	}
}
