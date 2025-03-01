package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigEnsureDefaultValues(t *testing.T) {
	// 测试空配置
	config := &Config{}
	config.EnsureDefaultValues()

	// 验证默认值
	assert.NotNil(t, config.StartTime, "StartTime应该有默认值")
	assert.NotNil(t, config.EndTime, "EndTime应该有默认值")
	
	// 验证时间范围
	duration := config.EndTime.Sub(*config.StartTime)
	assert.InDelta(t, 30*time.Minute.Minutes(), duration.Minutes(), 1, "默认时间范围应该是30分钟")

	// 验证配置目录
	homeDir, _ := os.UserHomeDir()
	expectedConfigDir := filepath.Join(homeDir, ".logsnap")
	assert.Equal(t, expectedConfigDir, config.ConfigDir, "配置目录应该设置为默认值")

	// 验证日志目录
	expectedLogDir := filepath.Join(homeDir, "xyz_log")
	assert.Equal(t, expectedLogDir, config.LogRootDir, "日志目录应该设置为默认值")

	// 测试自定义配置
	customTime := time.Now().Add(-1 * time.Hour)
	customConfig := &Config{
		StartTime:  &customTime,
		ConfigDir:  "/custom/config",
		LogRootDir: "/custom/logs",
	}
	customConfig.EnsureDefaultValues()

	// 验证自定义值被保留
	assert.Equal(t, customTime, *customConfig.StartTime, "自定义StartTime应该被保留")
	assert.Equal(t, "/custom/config", customConfig.ConfigDir, "自定义ConfigDir应该被保留")
	assert.Equal(t, "/custom/logs", customConfig.LogRootDir, "自定义LogRootDir应该被保留")
}
