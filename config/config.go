package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

// Config 存储整个应用的配置
type Config struct {
	Logs         []LogConfig  `mapstructure:"logs"`
	RemoteConfig RemoteConfig `mapstructure:"remote_config"`
	Version      string       `mapstructure:"version"`
}

// LogConfig 存储单个日志文件的配置
type LogConfig struct {
	Name       string `mapstructure:"name"`
	Path       string `mapstructure:"path"`
	TimeFormat string `mapstructure:"time_format"`
	TimeRegex  string `mapstructure:"time_regex"`
}

// RemoteConfig 远程配置信息
type RemoteConfig struct {
	Enabled          bool   `mapstructure:"enabled"`            // 是否启用远程配置
	URL              string `mapstructure:"url"`                // 远程配置URL
	Interval         int    `mapstructure:"interval"`           // 检查间隔(分钟)
	LastCheck        string `mapstructure:"last_check"`         // 上次检查时间
	AllowAutoUpgrade bool   `mapstructure:"allow_auto_upgrade"` // 是否允许自动升级
}

// LoadConfig 从文件中加载配置
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configPath)

	// 环境变量配置
	v.SetEnvPrefix("LOGSNAP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		// 如果配置文件不存在，则创建默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return createDefaultConfig(configPath)
		}
		return nil, fmt.Errorf("读取配置文件 %w", err)
	}

	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已更改:", e.Name)
	})

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置 %w", err)
	}

	// 替换配置中的环境变量
	replaceEnvVars(&config)

	return &config, nil
}

// 替换配置中的环境变量 (${VAR_NAME} 格式)
func replaceEnvVars(cfg *Config) {
	// 替换路径中的环境变量
	for i := range cfg.Logs {
		cfg.Logs[i].Path = replaceEnvVar(cfg.Logs[i].Path)
	}

	// 替换远程配置URL中的环境变量
	cfg.RemoteConfig.URL = replaceEnvVar(cfg.RemoteConfig.URL)
}

// 替换单个字符串中的环境变量
func replaceEnvVar(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envVar := value[2 : len(value)-1]
		if envValue := os.Getenv(envVar); envValue != "" {
			return envValue
		}
	}
	return value
}

// 创建默认配置
func createDefaultConfig(configPath string) (*Config, error) {
	// 创建文件夹
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 创建默认配置
	defaultConfig := &Config{
		Logs: []LogConfig{
			{
				Name:       "app",
				Path:       "/var/log/app/app.log",
				TimeFormat: "2006-01-02 15:04:05",
				TimeRegex:  "\\[(.*?)\\]",
			},
		},
		RemoteConfig: RemoteConfig{
			Enabled:          true,
			URL:              "https://example.com/logsnap-config.json",
			Interval:         60, // 1小时检查一次
			AllowAutoUpgrade: false,
		},
		Version: "0.0.1",
	}

	// 写入默认配置
	v := viper.New()
	v.SetConfigFile(configPath)

	for i, log := range defaultConfig.Logs {
		v.Set(fmt.Sprintf("logs.%d.name", i), log.Name)
		v.Set(fmt.Sprintf("logs.%d.path", i), log.Path)
		v.Set(fmt.Sprintf("logs.%d.time_format", i), log.TimeFormat)
		v.Set(fmt.Sprintf("logs.%d.time_regex", i), log.TimeRegex)
	}

	v.Set("remote_config.enabled", defaultConfig.RemoteConfig.Enabled)
	v.Set("remote_config.url", defaultConfig.RemoteConfig.URL)
	v.Set("remote_config.interval", defaultConfig.RemoteConfig.Interval)
	v.Set("remote_config.allow_auto_upgrade", defaultConfig.RemoteConfig.AllowAutoUpgrade)
	v.Set("version", defaultConfig.Version)

	if err := v.WriteConfig(); err != nil {
		return nil, fmt.Errorf("写入默认配置失败: %w", err)
	}

	logrus.Infof("已创建默认配置文件: %s\n", configPath)
	return defaultConfig, nil
}
