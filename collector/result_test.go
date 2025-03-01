package collector

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorResult_GetTotalLines(t *testing.T) {
	// 创建测试数据
	result := &ProcessorResult{
		processorName: "测试处理器",
		outputPath:    "/test/output.log",
		results: []FileProcessResult{
			{
				FilePath:   "/test/file1.log",
				TotalLines: 100,
				MatchLines: 50,
			},
			{
				FilePath:   "/test/file2.log",
				TotalLines: 200,
				MatchLines: 100,
			},
		},
		err: nil,
	}

	// 测试获取总行数
	totalLines := result.GetTotalLines()

	// 验证结果
	assert.Equal(t, 300, totalLines, "总行数应该为300")
}

func TestProcessorResult_GetMatchLines(t *testing.T) {
	// 创建测试数据
	result := &ProcessorResult{
		processorName: "测试处理器",
		outputPath:    "/test/output.log",
		results: []FileProcessResult{
			{
				FilePath:   "/test/file1.log",
				TotalLines: 100,
				MatchLines: 50,
			},
			{
				FilePath:   "/test/file2.log",
				TotalLines: 200,
				MatchLines: 100,
			},
		},
		err: nil,
	}

	// 测试获取匹配行数
	matchLines := result.GetMatchLines()

	// 验证结果
	assert.Equal(t, 150, matchLines, "匹配行数应该为150")
}

func TestProcessorResult_GetFileCount(t *testing.T) {
	// 创建测试数据
	result := &ProcessorResult{
		processorName: "测试处理器",
		outputPath:    "/test/output.log",
		results: []FileProcessResult{
			{
				FilePath:   "/test/file1.log",
				FileCount:  1,
				TotalLines: 100,
				MatchLines: 50,
			},
			{
				FilePath:   "/test/file2.log",
				FileCount:  2,
				TotalLines: 200,
				MatchLines: 100,
			},
		},
		err: nil,
	}

	// 测试获取文件数量
	fileCount := result.GetFileCount()

	// 验证结果
	assert.Equal(t, 3, fileCount, "文件数量应该为3")
}

func TestProcessorResult_GetFileSize(t *testing.T) {
	// 创建测试数据
	result := &ProcessorResult{
		processorName: "测试处理器",
		outputPath:    "/test/output.log",
		results: []FileProcessResult{
			{
				FilePath:   "/test/file1.log",
				FileSize:   1024,
				TotalLines: 100,
				MatchLines: 50,
			},
			{
				FilePath:   "/test/file2.log",
				FileSize:   2048,
				TotalLines: 200,
				MatchLines: 100,
			},
		},
		err: nil,
	}

	// 测试获取文件大小
	fileSize := result.GetFileSize()

	// 验证结果
	assert.Equal(t, int64(3072), fileSize, "文件大小应该为3072")
}

func TestFileProcessResult(t *testing.T) {
	// 测试创建FileProcessResult
	result := FileProcessResult{
		FilePath:   "/test/file.log",
		Err:        errors.New("测试错误"),
		FileCount:  5,
		FileSize:   1024,
		TotalLines: 100,
		MatchLines: 50,
		MatchFiles: 3,
	}

	// 验证字段
	assert.Equal(t, "/test/file.log", result.FilePath, "文件路径应该正确")
	assert.Error(t, result.Err, "错误应该被设置")
	assert.Equal(t, "测试错误", result.Err.Error(), "错误消息应该正确")
	assert.Equal(t, 5, result.FileCount, "文件数量应该正确")
	assert.Equal(t, int64(1024), result.FileSize, "文件大小应该正确")
	assert.Equal(t, 100, result.TotalLines, "总行数应该正确")
	assert.Equal(t, 50, result.MatchLines, "匹配行数应该正确")
	assert.Equal(t, 3, result.MatchFiles, "匹配文件数应该正确")
}

func TestProcessorResult_EmptyResults(t *testing.T) {
	// 创建没有结果的ProcessorResult
	result := &ProcessorResult{
		processorName: "空处理器",
		outputPath:    "",
		results:       []FileProcessResult{},
		err:           errors.New("处理失败"),
	}

	// 测试各种获取方法
	assert.Equal(t, 0, result.GetTotalLines(), "没有结果时总行数应该为0")
	assert.Equal(t, 0, result.GetMatchLines(), "没有结果时匹配行数应该为0")
	assert.Equal(t, 0, result.GetFileCount(), "没有结果时文件数量应该为0")
	assert.Equal(t, int64(0), result.GetFileSize(), "没有结果时文件大小应该为0")
}
