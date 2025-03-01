package processor

import (
	"bytes"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessLogContent(t *testing.T) {
	// 准备测试数据
	filePath := "./testdata/hmiserver.log"
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("无法打开测试文件: %v", err)
	}
	defer file.Close()

	// 准备输出buffer
	var output bytes.Buffer

	// 设置时间范围（根据日志文件内容调整）
	timePattern := regexp.MustCompile(`\[(.*?)\]`) // 假设日志格式为 [2024-01-01 12:00:00]
	timeFormat := "2006-01-02 15:04:05"
	startTime, _ := time.ParseInLocation(timeFormat, "2024-01-01 00:00:00", time.Local)
	endTime, _ := time.ParseInLocation(timeFormat, "2024-01-02 00:00:00", time.Local)

	// 写入文件头
	err = WriteFileHeader(&output, filePath, startTime, endTime)
	assert.NoError(t, err, "写入文件头应该成功")

	// 处理日志内容
	lineCount, matchCount, totalSize, err := ProcessLogContent(file, &output, timePattern, timeFormat, startTime, endTime)

	// 验证处理结果
	assert.NoError(t, err, "处理日志内容应该成功")
	assert.Greater(t, lineCount, 0, "应该处理了至少一行")
	assert.GreaterOrEqual(t, lineCount, matchCount, "匹配行数不应该超过总行数")
	assert.Greater(t, totalSize, int64(0), "文件大小应该大于0")

	// 验证输出内容
	outputStr := output.String()
	assert.Contains(t, outputStr, "原始日志文件:", "输出应该包含文件头")
	assert.Contains(t, outputStr, "时间范围:", "输出应该包含时间范围信息")

	// 验证时间过滤是否正确
	lines := bytes.Split(output.Bytes(), []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue // 跳过空行和注释行
		}
		matches := timePattern.FindSubmatch(line)
		if len(matches) >= 2 {
			timeStr := string(matches[1])
			timestamp, err := time.ParseInLocation(timeFormat, timeStr, time.Local)
			assert.NoError(t, err, "日志时间格式应该正确")
			assert.True(t,
				(timestamp.After(startTime) || timestamp.Equal(startTime)) &&
					(timestamp.Before(endTime) || timestamp.Equal(endTime)),
				"过滤后的日志行时间应该在指定范围内: %s", timeStr)
		}
	}
}

func TestFilterLogLineByTime(t *testing.T) {
	timePattern := regexp.MustCompile(`\[(.*?)\]`)
	timeFormat := "2006-01-02 15:04:05"
	startTime, _ := time.ParseInLocation(timeFormat, "2024-01-01 00:00:00", time.Local)
	endTime, _ := time.ParseInLocation(timeFormat, "2024-01-02 00:00:00", time.Local)

	tests := []struct {
		name     string
		line     string
		expected bool
		hasError bool
	}{
		{
			name:     "valid time in range",
			line:     "[2024-01-01 12:00:00] test log",
			expected: true,
			hasError: false,
		},
		{
			name:     "valid time before range",
			line:     "[2023-12-31 23:59:59] test log",
			expected: false,
			hasError: false,
		},
		{
			name:     "valid time after range",
			line:     "[2024-01-02 00:00:01] test log",
			expected: false,
			hasError: false,
		},
		{
			name:     "invalid time format",
			line:     "[invalid-time] test log",
			expected: false,
			hasError: true,
		},
		{
			name:     "no timestamp",
			line:     "test log without timestamp",
			expected: false,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := FilterLogLineByTime(tt.line, timePattern, timeFormat, startTime, endTime)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, match)
		})
	}
}
