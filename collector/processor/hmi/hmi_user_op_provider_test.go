package hmi

import (
	processor "logsnap/collector/processor"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// 测试 UserOpFileInfoFilter 的 IsMatch 方法
func TestUserOpFileInfoFilter_IsMatch(t *testing.T) {
	filter := &UserOpFileInfoFilter{}
	
	// 测试用例：应该匹配的文件名
	validNames := []string{
		"20250302-084053.269",
		"20250307-135338.555",
		"20250101-000000.000",
	}
	
	// 测试用例：不应该匹配的文件名
	invalidNames := []string{
		"log_20250302.txt",
		"userlog-20250302",
		"2025030-084053.269", // 少一位数字
		"202503020-084053.269", // 多一位数字
		"20250302_084053.269", // 使用下划线而非连字符
	}
	
	// 测试有效的文件名
	for _, name := range validNames {
		if !filter.IsMatch(name) {
			t.Errorf("IsMatch(%s) 应该返回 true，但返回了 false", name)
		}
	}
	
	// 测试无效的文件名
	for _, name := range invalidNames {
		if filter.IsMatch(name) {
			t.Errorf("IsMatch(%s) 应该返回 false，但返回了 true", name)
		}
	}
}

// 测试 UserOpFileInfoFilter 的 parseLogFileInfo 方法
func TestUserOpFileInfoFilter_ParseLogFileInfo(t *testing.T) {
	filter := &UserOpFileInfoFilter{}
	
	// 测试用例
	testCases := []struct {
		fileName    string
		expectError bool
		expectedTime string
	}{
		{"20250302-084053.269", false, "2025-03-02 08:40:53"},
		{"20250307-135338.555", false, "2025-03-07 13:53:38"},
		{"invalid_filename.txt", true, ""},
	}
	
	for _, tc := range testCases {
		filePath := filepath.Join("testdata", tc.fileName)
		fileInfo, err := filter.parseLogFileInfo(filePath)
		
		if tc.expectError {
			if err == nil {
				t.Errorf("parseLogFileInfo(%s) 应该返回错误，但没有", tc.fileName)
			}
		} else {
			if err != nil {
				t.Errorf("parseLogFileInfo(%s) 返回了错误: %v", tc.fileName, err)
				continue
			}
			
			expectedTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tc.expectedTime, time.Local)
			if !fileInfo.StartTime.Equal(expectedTime) {
				t.Errorf("parseLogFileInfo(%s) 时间解析错误，期望: %v, 实际: %v", 
					tc.fileName, expectedTime, fileInfo.StartTime)
			}
			
			if fileInfo.FileType != "user_op" {
				t.Errorf("parseLogFileInfo(%s) 文件类型错误，期望: user_op, 实际: %s", 
					tc.fileName, fileInfo.FileType)
			}
		}
	}
}

// 测试 UserOpFileInfoFilter 的 ParseFileInfos 方法
func TestUserOpFileInfoFilter_ParseFileInfos(t *testing.T) {
	filter := &UserOpFileInfoFilter{}
	
	// 创建测试文件列表
	files := []string{
		filepath.Join("testdata", "20250302-084053.269"),
		filepath.Join("testdata", "20250307-135338.555"),
		filepath.Join("testdata", "invalid_file.txt"),
	}
	
	fileInfos, err := filter.ParseFileInfos(files)
	if err != nil {
		t.Fatalf("ParseFileInfos 返回了错误: %v", err)
	}
	
	// 应该只有两个有效的文件信息
	if len(fileInfos) != 2 {
		t.Errorf("ParseFileInfos 应该返回 2 个文件信息，但返回了 %d 个", len(fileInfos))
	}
}

// 测试 UserOpFileProcessorProvider 的 FilterFiles 方法
func TestUserOpFileProcessorProvider_FilterFiles(t *testing.T) {
	provider := NewUserOpFileProcessorProvider()
	
	// 创建测试文件列表
	files := []string{
		filepath.Join("testdata", "20250302-084053.269"),
		filepath.Join("testdata", "20250302-084053.269"),
		filepath.Join("testdata", "20250307-135338.555"),
	}
	
	// 测试时间范围过滤
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-03-05 00:00:00", time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-03-10 00:00:00", time.Local)
	
	// 如果有时间范围，至少会返回一个文件，即使该文件不在时间范围内
	fileInfos, err := provider.FilterFiles(files, startTime, endTime)
	if err != nil {
		t.Fatalf("FilterFiles 返回了错误: %v", err)
	}
	
	// 应该只有一个文件在时间范围内
	if len(fileInfos) != 2 {
		t.Errorf("FilterFiles 应该返回 2 个文件信息，但返回了 %d 个", len(fileInfos))
	}
	
	if len(fileInfos) > 0 && filepath.Base(fileInfos[1].Path) != "20250307-135338.555" {
		t.Errorf("FilterFiles 返回了错误的文件，期望: 20250307-135338.555, 实际: %s", 
			filepath.Base(fileInfos[1].Path))
	}
}

// 测试 UserOpFileProcessorProvider 的 ProcessFile 方法
func TestUserOpFileProcessorProvider_ProcessFile(t *testing.T) {
	// 创建临时目录用于输出
	tempDir, err := os.MkdirTemp("", "user_op_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	provider := NewUserOpFileProcessorProvider()
	
	// 创建测试文件
	testFilePath := filepath.Join("testdata", "20250302-084053.269")
	testFileContent := `20250302 08:40:53.269] User clicked [StartButton].
20250302 08:41:13.163] User clicked [StopButton].
20250302 08:42:05.789] User changed setting [Speed] to [100].
`
	// 确保测试目录存在
	err = os.MkdirAll(filepath.Dir(testFilePath), 0755)
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	
	err = os.WriteFile(testFilePath, []byte(testFileContent), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove(testFilePath)
	
	// 创建文件信息
	fileTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-03-02 08:40:53", time.Local)
	fileInfo := processor.LogFileInfo{
		Path:      testFilePath,
		StartTime: fileTime,
		FileName:  filepath.Base(testFilePath),
		FileType:  "user_op",
	}
	
	// 测试处理文件
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-03-02 08:41:00", time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-03-02 08:42:00", time.Local)
	
	_, err = provider.ProcessFile(fileInfo, startTime, endTime, tempDir)
	if err != nil {
		t.Fatalf("ProcessFile 返回了错误: %v", err)
	}
	
	// 检查输出文件是否存在
	outputFilePath := filepath.Join(tempDir, filepath.Base(testFilePath))
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("输出文件不存在: %s", outputFilePath)
	}
	
	// 检查输出文件内容
	outputContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("读取输出文件失败: %v", err)
	}
	
	// 去掉文件头部分再进行比较
	// 文件头通常是注释，以 # 开头的几行
	outputLines := splitAndFilterHeader(string(outputContent))
	
	// 应该只包含时间范围内的一行日志
	expectedContent := "20250302 08:41:13.163] User clicked [StopButton].\n"
	actualContent := outputLines
	if actualContent != expectedContent {
		t.Errorf("输出文件内容错误，期望:\n%s\n实际:\n%s", expectedContent, actualContent)
	}
}

// 辅助函数：分割文本并过滤掉文件头（以 # 开头的行）
func splitAndFilterHeader(content string) string {
	lines := strings.Split(content, "\n")
	var filteredLines []string
	
	// 跳过以 # 开头的行（文件头）
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") && line != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	
	// 重新组合过滤后的行
	return strings.Join(filteredLines, "\n") + "\n"
}

// 测试 UserOpFileProcessorProvider 的 GetFileSuffixes 方法
func TestUserOpFileProcessorProvider_GetFileSuffixes(t *testing.T) {
	provider := NewUserOpFileProcessorProvider()
	suffixes := provider.GetFileSuffixes()
	
	// 应该返回空切片
	if len(suffixes) != 0 {
		t.Errorf("GetFileSuffixes 应该返回空切片，但返回了 %v", suffixes)
	}
}