package service

// ConfigProvider 定义配置提供者接口
type ConfigProvider interface {
	GetVersion() string
	GetRemoteConfigEnabled() bool
	GetRemoteConfigURL() string
	GetRemoteConfigInterval() int
	GetRemoteConfigLastCheck() string
	SetRemoteConfigLastCheck(lastCheck string)
}

// LocalConfigProvider 实现远程包中的ConfigProvider接口
type LocalConfigProvider struct {
	Version           string
	RemoteEnabled     bool
	RemoteURL         string
	RemoteInterval    int
	RemoteLastCheck   string
	RemoteAutoUpgrade bool
}

// GetVersion 实现ConfigProvider接口
func (l *LocalConfigProvider) GetVersion() string {
	return l.Version
}

// GetRemoteConfigEnabled 实现ConfigProvider接口
func (l *LocalConfigProvider) GetRemoteConfigEnabled() bool {
	return l.RemoteEnabled
}

// GetRemoteConfigURL 实现ConfigProvider接口
func (l *LocalConfigProvider) GetRemoteConfigURL() string {
	return l.RemoteURL
}

// GetRemoteConfigInterval 实现ConfigProvider接口
func (l *LocalConfigProvider) GetRemoteConfigInterval() int {
	return l.RemoteInterval
}

// GetRemoteConfigLastCheck 实现ConfigProvider接口
func (l *LocalConfigProvider) GetRemoteConfigLastCheck() string {
	return l.RemoteLastCheck
}

// SetRemoteConfigLastCheck 实现ConfigProvider接口
func (l *LocalConfigProvider) SetRemoteConfigLastCheck(lastCheck string) {
	l.RemoteLastCheck = lastCheck
}

// NewLocalConfigProvider 创建新的本地配置提供者
func NewLocalConfigProvider() *LocalConfigProvider {
	// 使用固定的版本信息
	currentVersion := "0.0.1" // 这个版本应该从构建时注入

	return &LocalConfigProvider{
		Version:           currentVersion,
		RemoteEnabled:     true,
		RemoteURL:         "https://example.com/logsnap-config.json",
		RemoteInterval:    60, // 60分钟
		RemoteAutoUpgrade: false,
	}
}
