package collector

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 创建模拟的LogProcessor
type MockLogProcessor struct {
	mock.Mock
}

func (m *MockLogProcessor) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLogProcessor) GetLogPath() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockLogProcessor) GetOutputDir() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLogProcessor) Collect(startTime, endTime time.Time, outputDir string) (string, []FileProcessResult, error) {
	args := m.Called(startTime, endTime, outputDir)
	return args.String(0), args.Get(1).([]FileProcessResult), args.Error(2)
}

func TestNewCollector(t *testing.T) {
	// 创建模拟的LogProcessor
	processor := new(MockLogProcessor)
	processor.On("GetName").Return("测试处理器")
	
	// 创建处理器列表
	processors := []LogProcessor{processor}
	
	// 创建收集器
	outputDir := "/test/output"
	collector := NewCollector(processors, outputDir)
	
	// 验证收集器
	assert.NotNil(t, collector, "收集器不应为空")
	assert.Equal(t, processors, collector.logProcessors, "处理器列表应该正确")
	assert.Equal(t, outputDir, collector.outputDir, "输出目录应该正确")
}

func TestAddProcessor(t *testing.T) {
	// 创建收集器
	collector := NewCollector([]LogProcessor{}, "/test/output")
	
	// 创建模拟的LogProcessor
	processor := new(MockLogProcessor)
	processor.On("GetName").Return("测试处理器")
	
	// 添加处理器
	collector.AddProcessor(processor)
	
	// 验证处理器是否被添加
	assert.Len(t, collector.logProcessors, 1, "应该有1个处理器")
	assert.Equal(t, processor, collector.logProcessors[0], "处理器应该被正确添加")
}

func TestRemoveProcessor(t *testing.T) {
	// 创建模拟的LogProcessor
	processor1 := new(MockLogProcessor)
	processor1.On("GetName").Return("处理器1")
	
	processor2 := new(MockLogProcessor)
	processor2.On("GetName").Return("处理器2")
	
	// 创建收集器
	collector := NewCollector([]LogProcessor{processor1, processor2}, "/test/output")
	
	// 移除处理器
	result := collector.RemoveProcessor("处理器1")
	
	// 验证结果
	assert.True(t, result, "移除应该成功")
	assert.Len(t, collector.logProcessors, 1, "应该剩余1个处理器")
	assert.Equal(t, processor2, collector.logProcessors[0], "处理器2应该被保留")
	
	// 测试移除不存在的处理器
	result = collector.RemoveProcessor("不存在的处理器")
	assert.False(t, result, "移除不存在的处理器应该失败")
}

func TestGetProcessors(t *testing.T) {
	// 创建模拟的LogProcessor
	processor1 := new(MockLogProcessor)
	processor1.On("GetName").Return("处理器1")
	
	processor2 := new(MockLogProcessor)
	processor2.On("GetName").Return("处理器2")
	
	// 创建处理器列表
	processors := []LogProcessor{processor1, processor2}
	
	// 创建收集器
	collector := NewCollector(processors, "/test/output")
	
	// 获取处理器列表
	result := collector.GetProcessors()
	
	// 验证结果
	assert.Equal(t, processors, result, "返回的处理器列表应该正确")
}

func TestGetProcessorByName(t *testing.T) {
	// 创建模拟的LogProcessor
	processor1 := new(MockLogProcessor)
	processor1.On("GetName").Return("处理器1")
	
	processor2 := new(MockLogProcessor)
	processor2.On("GetName").Return("处理器2")
	
	// 创建收集器
	collector := NewCollector([]LogProcessor{processor1, processor2}, "/test/output")
	
	// 获取处理器
	result, found := collector.GetProcessorByName("处理器1")
	
	// 验证结果
	assert.True(t, found, "应该找到处理器")
	assert.Equal(t, processor1, result, "返回的处理器应该正确")
	
	// 测试获取不存在的处理器
	result, found = collector.GetProcessorByName("不存在的处理器")
	assert.False(t, found, "不应该找到不存在的处理器")
	assert.Nil(t, result, "不存在的处理器应该返回nil")
}

func TestHasProcessor(t *testing.T) {
	// 创建模拟的LogProcessor
	processor := new(MockLogProcessor)
	processor.On("GetName").Return("测试处理器")
	
	// 创建收集器
	collector := NewCollector([]LogProcessor{processor}, "/test/output")
	
	// 测试存在的处理器
	result := collector.HasProcessor("测试处理器")
	assert.True(t, result, "应该有测试处理器")
	
	// 测试不存在的处理器
	result = collector.HasProcessor("不存在的处理器")
	assert.False(t, result, "不应该有不存在的处理器")
}

func TestClearProcessors(t *testing.T) {
	// 创建模拟的LogProcessor
	processor := new(MockLogProcessor)
	processor.On("GetName").Return("测试处理器")
	
	// 创建收集器
	collector := NewCollector([]LogProcessor{processor}, "/test/output")
	
	// 清空处理器
	collector.ClearProcessors()
	
	// 验证结果
	assert.Empty(t, collector.logProcessors, "处理器列表应该为空")
}

func TestSetOutputDir(t *testing.T) {
	// 创建收集器
	collector := NewCollector([]LogProcessor{}, "/test/output")
	
	// 设置新的输出目录
	newOutputDir := "/new/output"
	collector.SetOutputDir(newOutputDir)
	
	// 验证结果
	assert.Equal(t, newOutputDir, collector.outputDir, "输出目录应该被更新")
}

func TestGetOutputDir(t *testing.T) {
	// 创建收集器
	outputDir := "/test/output"
	collector := NewCollector([]LogProcessor{}, outputDir)
	
	// 获取输出目录
	result := collector.GetOutputDir()
	
	// 验证结果
	assert.Equal(t, outputDir, result, "返回的输出目录应该正确")
}

func TestCollect(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "logsnap-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 创建模拟的LogProcessor
	processor1 := new(MockLogProcessor)
	processor1.On("GetName").Return("处理器1")
	processor1.On("GetLogPath").Return("/logs/test1.log", nil)
	processor1.On("Collect", mock.Anything, mock.Anything, mock.Anything).Return(
		filepath.Join(tempDir, "output1.log"),
		[]FileProcessResult{
			{FilePath: "test1.log", TotalLines: 100, MatchLines: 50},
		},
		nil,
	)
	
	processor2 := new(MockLogProcessor)
	processor2.On("GetName").Return("处理器2")
	processor2.On("GetLogPath").Return("/logs/test2.log", nil)
	processor2.On("Collect", mock.Anything, mock.Anything, mock.Anything).Return(
		filepath.Join(tempDir, "output2.log"),
		[]FileProcessResult{
			{FilePath: "test2.log", TotalLines: 200, MatchLines: 100},
		},
		nil,
	)
	
	// 创建收集器
	collector := NewCollector([]LogProcessor{processor1, processor2}, tempDir)
	
	// 创建测试文件
	testFile1 := filepath.Join(tempDir, "output1.log")
	if err := os.WriteFile(testFile1, []byte("测试内容1"), 0644); err != nil {
		t.Fatalf("创建测试文件1失败: %v", err)
	}
	
	testFile2 := filepath.Join(tempDir, "output2.log")
	if err := os.WriteFile(testFile2, []byte("测试内容2"), 0644); err != nil {
		t.Fatalf("创建测试文件2失败: %v", err)
	}
	
	// 设置测试时间范围
	startTime := time.Now().Add(-time.Hour)
	endTime := time.Now()
	
	// 执行收集
	zipPath, err := collector.Collect(startTime, endTime)
	
	// 验证结果
	assert.NoError(t, err, "收集不应返回错误")
	assert.NotEmpty(t, zipPath, "ZIP路径不应为空")
	assert.FileExists(t, zipPath, "ZIP文件应该存在")
	
	// 验证模拟对象的调用
	processor1.AssertExpectations(t)
	processor2.AssertExpectations(t)
}
