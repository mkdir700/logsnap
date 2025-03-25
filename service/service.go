package service

import (
	"fmt"
	"logsnap/collector"
	"logsnap/collector/factory"
	"os"
	"path/filepath"

	"logsnap/remote"

	"github.com/sirupsen/logrus"
)

// Service 定义主服务结构
type Service struct {
	Config           *Config
	UploadConfig     *remote.UploadConfig
	uploadManager    UploadManager
	configProvider   ConfigProvider
	progressReporter ProgressReporter
}

// NewService 创建新的服务实例
func NewService(config *Config, uploadConfig *remote.UploadConfig) *Service {
	if config == nil {
		// 创建默认配置
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}

		config = &Config{
			LogRootDir: filepath.Join(homeDir, "xyz_log"),
			ConfigDir:  filepath.Join(homeDir, ".logsnap"),
		}
		config.EnsureDefaultValues()
	}

	if uploadConfig == nil {
	}

	// 创建进度报告器
	progressReporter := &DefaultProgressReporter{}

	service := &Service{
		Config:           config,
		progressReporter: progressReporter,
	}

	// 初始化上传管理器
	service.uploadManager = NewUploadManager(uploadConfig)

	// 初始化配置提供者
	service.configProvider = NewLocalConfigProvider()

	return service
}

// UploadLogSnapFile 上传日志文件
func (s *Service) UploadLogSnapFile(file *LogFile, description string, tags []string) (*UploadResult, error) {
	if file == nil {
		return nil, fmt.Errorf("没有要上传的文件")
	}

	// 创建上传请求
	request := &UploadRequest{
		File:        file,
		Config:      s.UploadConfig,
		Reporter:    s.progressReporter,
		Description: description,
		Tags:        tags,
	}

	// 执行上传
	return s.uploadManager.Upload(request)
}

// SetProgressCallback 设置进度回调函数
func (s *Service) SetProgressCallback(callback ProgressCallback) {
	if callback != nil {
		s.progressReporter = &DefaultProgressReporter{
			callback: callback,
		}
	} else {
		s.progressReporter = &DefaultProgressReporter{}
	}
}

// CheckForUpdates 检查更新
func (s *Service) CheckForUpdates() (bool, string, error) {
	// 获取当前版本
	currentVersion := GetCurrentVersion()

	// 获取远程版本
	remoteVersionStr := s.configProvider.GetVersion()

	// 解析远程版本
	remoteVer, err := ParseVersion(remoteVersionStr)
	if err != nil {
		return false, "", err
	}

	// 比较版本
	if remoteVer.IsNewer(currentVersion) {
		return true, remoteVersionStr, nil
	}

	return false, remoteVersionStr, nil
}

// GetRemoteConfigEnabled 获取远程配置是否启用
func (s *Service) GetRemoteConfigEnabled() bool {
	return s.configProvider.GetRemoteConfigEnabled()
}

// GetRemoteConfigURL 获取远程配置URL
func (s *Service) GetRemoteConfigURL() string {
	return s.configProvider.GetRemoteConfigURL()
}

// SetRemoteConfigURL 设置远程配置URL
func (s *Service) SetRemoteConfigURL(url string) {
	// 这里应该实现设置URL的逻辑
	// 由于ConfigProvider接口没有提供设置URL的方法，我们需要进行类型断言
	if localProvider, ok := s.configProvider.(*LocalConfigProvider); ok {
		localProvider.RemoteURL = url
	}
}

// GetRemoteConfigInterval 获取远程配置检查间隔
func (s *Service) GetRemoteConfigInterval() int {
	return s.configProvider.GetRemoteConfigInterval()
}

// SetRemoteConfigInterval 设置远程配置检查间隔
func (s *Service) SetRemoteConfigInterval(interval int) {
	// 这里应该实现设置间隔的逻辑
	// 由于ConfigProvider接口没有提供设置间隔的方法，我们需要进行类型断言
	if localProvider, ok := s.configProvider.(*LocalConfigProvider); ok {
		localProvider.RemoteInterval = interval
	}
}

// SaveConfig 保存配置
func (s *Service) SaveConfig() error {
	// 这里应该将配置保存到文件
	configPath := filepath.Join(s.Config.ConfigDir, "config.json")

	// 确保目录存在
	err := os.MkdirAll(s.Config.ConfigDir, 0755)
	if err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 保存配置的逻辑应该在这里实现
	// 暂时只记录日志
	logrus.Infof("配置已保存到: %s", configPath)

	return nil
}

// LoadConfig 加载配置
func (s *Service) LoadConfig() error {
	// 这里应该从文件加载配置
	configPath := filepath.Join(s.Config.ConfigDir, "config.json")

	// 检查文件是否存在
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		// 如果文件不存在，则创建默认配置
		return s.SaveConfig()
	} else if err != nil {
		return fmt.Errorf("检查配置文件失败: %v", err)
	}

	// 加载配置的逻辑应该在这里实现
	// 暂时只记录日志
	logrus.Infof("配置已从 %s 加载", configPath)

	return nil
}

// CollectAndUploadLogs 收集并上传日志文件
// 这是为了保持向后兼容性而提供的函数
func CollectAndUploadLogs(config *Config, uploadConfig *remote.UploadConfig) (string, string, error) {
	// 创建服务实例
	service := NewService(config, uploadConfig)

	// 从collector包获取日志处理器
	var processors []collector.LogProcessor

	// 根据配置添加日志处理器
	// 当LogTypes长度为0时，加载所有支持的处理器
	// 否则根据LogTypes中的类型加载相应的处理器
	loadAll := len(config.Programs) == 0

	if loadAll {
		supportedProcessors := factory.GetSupportedProcessorTypes()
		for _, processorType := range supportedProcessors {
			processor, err := factory.CreateProcessor(processorType, config.LogRootDir, "")
			if err != nil {
				logrus.Errorf("创建 %s 处理器失败: %v", processorType, err)
			}
			processors = append(processors, processor)
		}
	} else {
		for _, processorType := range config.Programs {
			processor, err := factory.CreateProcessor(collector.ProcessorType(processorType), config.LogRootDir, "")
			if err != nil {
				return "", "", fmt.Errorf("创建 %s 处理器失败: %v", processorType, err)
			}
			processors = append(processors, processor)
		}
	}

	// 如果没有加载任何处理器，记录警告
	if len(processors) == 0 {
		logrus.Warnf("没有添加任何日志处理器，LogTypes: %v", config.Programs)
	}

	// 创建收集器
	collect := collector.NewCollector(processors, config.OutputDir)

	// 使用collector收集日志
	snapPath, err := collect.Collect(*config.StartTime, *config.EndTime)
	if err != nil {
		return "", "", fmt.Errorf("收集日志失败: %v", err)
	}

	// 如果不需要上传，直接返回结果
	if !config.ShouldUpload {
		return snapPath, "", nil
	}

	// 准备上传文件
	file, err := os.Stat(snapPath)
	if err != nil {
		return snapPath, "", fmt.Errorf("获取日志文件信息失败: %v", err)
	}

	// 创建LogPath和LogFile对象
	logPath := &LogPath{
		Name: "快照",
		Path: filepath.Dir(snapPath),
	}

	logFile := &LogFile{
		Name:     file.Name(),
		Path:     snapPath,
		Size:     file.Size(),
		ModTime:  file.ModTime(),
		IsDir:    file.IsDir(),
		LogPath:  logPath,
		RootPath: file.Name(),
		Selected: true,
	}

	// 上传文件
	result, err := service.UploadLogSnapFile(logFile, "通过CLI上传的日志", nil)
	if err != nil {
		return snapPath, "", fmt.Errorf("上传日志失败: %v", err)
	}

	// 如果需要，删除上传后的文件
	if !config.KeepLocalSnap {
		os.Remove(snapPath)
		logrus.Infof("已删除上传后的文件: %s", snapPath)
	}

	// 返回结果
	return snapPath, result.URL, nil
}
