package service

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// Config 包含服务配置
type Config struct {
	StartTime        *time.Time
	EndTime          *time.Time
	OutputDir        string           // 输出目录
	ShouldUpload     bool             // 是否上传
	KeepLocalSnap    bool             // 是否保留本地日志快照
	ConfigDir        string           // 配置目录
	LogRootDir       string           // 日志目录
	SkipVersionCheck bool             // 跳过版本检查
	ProgressCallback ProgressCallback // 进度回调函数
	Programs         []string         // 日志类型过滤（可选）
}

// EnsureDefaultValues 确保配置具有默认值
func (c *Config) EnsureDefaultValues() {
	// 如果没有设置配置目录，设置默认值
	if c.ConfigDir == "" {
		// 根据操作系统设置默认配置目录
		homeDir, err := os.UserHomeDir()
		if err == nil {
			c.ConfigDir = filepath.Join(homeDir, ".logsnap")
		} else {
			c.ConfigDir = "/etc/logsnap" // 备选目录
		}
	}

	// 如果没有设置日志目录，设置默认值
	if c.LogRootDir == "" {
		// 根据操作系统设置默认日志目录
		homeDir, err := os.UserHomeDir()
		if err == nil {
			c.LogRootDir = filepath.Join(homeDir, "xyz_log")
			logrus.Infof("使用默认日志目录: %s", c.LogRootDir)
		} else {
			c.LogRootDir = "/var/log/xyz" // 备选目录
			logrus.Infof("无法获取用户主目录，使用备选日志目录: %s", c.LogRootDir)
		}
	} else {
		logrus.Infof("使用指定的日志目录: %s", c.LogRootDir)
	}

	// 确保配置目录存在
	if err := os.MkdirAll(c.ConfigDir, 0755); err != nil {
		logrus.Warnf("创建配置目录失败: %v", err)
	}

	// 解析时间范围
	if c.StartTime == nil {
		// 如果没有指定开始时间，使用默认的30分钟
		now := time.Now()
		startTime := now.Add(-30 * time.Minute)
		c.StartTime = &startTime
		if c.EndTime == nil {
			endTime := now
			c.EndTime = &endTime
		}
	}
}
