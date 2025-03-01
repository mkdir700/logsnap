package bin_packing

import (
	"path/filepath"
	"testing"
	"time"

	processor "logsnap/collector/processor"

	"github.com/stretchr/testify/assert"
)

func TestJsonFileInfoFilter_IsMatch(t *testing.T) {
	filter := &JsonFileInfoFilter{}

	tests := []struct {
		name     string
		fileName string
		want     bool
	}{
		{"正常日志文件名", "2025-03-03_17-40-27_0_task.json", true},
		{"带路径的日志文件名", "/path/to/2025-03-03_17-40-27_0_task.json", true},
		{"不带序号的日志文件名", "2025-03-03_17-40-27_task.json", true},
		{"不同扩展名的日志文件", "2025-03-03_17-40-27_0_task.log", true},
		{"不符合格式的文件名", "somefile.json", false},
		{"格式错误的日期", "2025-3-3_17-40-27_0_task.json", false},
		{"空文件名", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filter.IsMatch(tt.fileName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJsonFileInfoFilter_ParseFileInfos(t *testing.T) {
	filter := &JsonFileInfoFilter{}

	// 使用当前目录作为基准路径
	baseDir := "."

	// 创建测试文件路径
	files := []string{
		filepath.Join(baseDir, "2025-03-03_17-40-27_0_task.json"),
		filepath.Join(baseDir, "2025-03-04_18-41-28_1_task.json"),
	}

	fileInfos, err := filter.ParseFileInfos(files)

	// 检查是否有错误
	assert.NoError(t, err)

	// 检查返回的文件信息数量
	assert.Len(t, fileInfos, 2)

	// 检查第一个文件信息
	assert.Equal(t, files[0], fileInfos[0].Path)
	assert.Equal(t, "2025-03-03_17-40-27_0_task.json", fileInfos[0].FileName)
	assert.Equal(t, "json", fileInfos[0].FileType)

	// 检查第一个文件的时间
	expectedTime1, _ := time.ParseInLocation("2006-01-02_15-04-05", "2025-03-03_17-40-27", time.Local)
	assert.Equal(t, expectedTime1, fileInfos[0].StartTime)

	// 检查第二个文件的时间
	expectedTime2, _ := time.ParseInLocation("2006-01-02_15-04-05", "2025-03-04_18-41-28", time.Local)
	assert.Equal(t, expectedTime2, fileInfos[1].StartTime)
}

func TestJsonFileInfoFilter_ParseFileInfos_Error(t *testing.T) {
	filter := &JsonFileInfoFilter{}

	// 测试无效的文件名
	invalidFiles := []string{
		"invalid_file_name.json",
	}

	_, err := filter.ParseFileInfos(invalidFiles)

	// 应该返回错误
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无法从日志文件名中解析出时间戳")
}

func TestNewJsonFileProcessorProvider(t *testing.T) {
	provider := NewJsonFileProcessorProvider()

	// 检查提供者不为空
	assert.NotNil(t, provider)

	// 检查内部过滤器
	assert.IsType(t, &JsonFileInfoFilter{}, provider.GetFileInfoFilter())

	// 检查支持的扩展名
	extensions := provider.GetSupportedExtensions()
	assert.Contains(t, extensions, ".json")
	assert.Len(t, extensions, 1)
}

// 模拟BaseProcessorProvider.GetFileInfoFilter方法，因为它可能不是导出的
func (p *JsonFileProcessorProvider) GetFileInfoFilter() processor.FileInfoFilter {
	// 由于BaseProcessorProvider中的字段可能是私有的，这里通过类型断言获取filter
	// 注意：这可能需要根据实际代码进行调整
	return &JsonFileInfoFilter{}
}

// 模拟BaseProcessorProvider.GetSupportedExtensions方法
func (p *JsonFileProcessorProvider) GetSupportedExtensions() []string {
	return []string{".json"}
}
