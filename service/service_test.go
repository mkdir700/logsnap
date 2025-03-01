package service

import (
	"logsnap/remote"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 创建模拟的UploadManager
type MockUploadManager struct {
	mock.Mock
}

func (m *MockUploadManager) Upload(request *UploadRequest) (*UploadResult, error) {
	args := m.Called(request)
	return args.Get(0).(*UploadResult), args.Error(1)
}

// 创建模拟的ConfigProvider
type MockConfigProvider struct {
	mock.Mock
}

func (m *MockConfigProvider) GetRemoteConfigEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigProvider) GetRemoteConfigURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigProvider) GetRemoteConfigInterval() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigProvider) GetVersion() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigProvider) GetRemoteConfigLastCheck() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigProvider) SetRemoteConfigLastCheck(lastCheck string) {
	m.Called(lastCheck)
}

func TestNewService(t *testing.T) {
	// 测试创建服务实例
	config := &Config{
		LogRootDir: "/test/logs",
		ConfigDir:  "/test/config",
	}
	
	// 创建一个有效的UploadConfig
	uploadConfig := &remote.UploadConfig{
		Providers: []remote.UploadConfigProvider{
			{
				Provider:   "s3",
				Endpoint:   "https://example.com",
				Bucket:     "test-bucket",
				AccessKey:  "test-key",
				SecretKey:  "test-secret",
				Region:     "us-east-1",
				FolderPath: "/logs",
			},
		},
		DefaultProvider: "s3",
	}

	service := NewService(config, uploadConfig)

	// 验证服务实例
	assert.NotNil(t, service, "服务实例不应为空")
	assert.Equal(t, config, service.Config, "配置应该被正确设置")
	assert.Equal(t, uploadConfig, service.UploadConfig, "上传配置应该被正确设置")
	assert.NotNil(t, service.uploadManager, "上传管理器不应为空")
	assert.NotNil(t, service.configProvider, "配置提供者不应为空")
	assert.NotNil(t, service.progressReporter, "进度报告器不应为空")

	// 测试默认配置
	defaultService := NewService(nil, uploadConfig)
	assert.NotNil(t, defaultService.Config, "默认配置不应为空")
	assert.NotNil(t, defaultService.Config.StartTime, "默认开始时间不应为空")
	assert.NotNil(t, defaultService.Config.EndTime, "默认结束时间不应为空")
}

func TestSetProgressCallback(t *testing.T) {
	service := NewService(nil, nil)
	
	// 测试设置回调
	called := false
	callback := func(stage string, progress int, message string) {
		called = true
	}
	
	service.SetProgressCallback(callback)
	
	// 验证回调被设置
	reporter, ok := service.progressReporter.(*DefaultProgressReporter)
	assert.True(t, ok, "进度报告器应该是DefaultProgressReporter类型")
	assert.NotNil(t, reporter.callback, "回调函数应该被设置")
	
	// 测试回调执行
	reporter.Report("test", 50, "测试消息")
	assert.True(t, called, "回调函数应该被调用")
	
	// 测试设置为nil
	service.SetProgressCallback(nil)
	reporter, ok = service.progressReporter.(*DefaultProgressReporter)
	assert.True(t, ok, "进度报告器应该是DefaultProgressReporter类型")
	assert.Nil(t, reporter.callback, "回调函数应该被设置为nil")
}

func TestCheckForUpdates(t *testing.T) {
	// 创建模拟的ConfigProvider
	mockProvider := new(MockConfigProvider)
	
	// 设置期望
	mockProvider.On("GetVersion").Return("1.0.0")
	mockProvider.On("GetRemoteConfigLastCheck").Return("0")
	
	// 创建服务实例
	service := NewService(nil, nil)
	service.configProvider = mockProvider
	
	// 测试检查更新
	hasUpdate, version, err := service.CheckForUpdates()
	
	// 验证结果
	assert.NoError(t, err, "检查更新不应返回错误")
	assert.False(t, hasUpdate, "不应该检测到更新")
	assert.Equal(t, "1.0.0", version, "版本应该正确")
	
	// 验证模拟对象的调用
	mockProvider.AssertExpectations(t)
}

func TestSaveAndLoadConfig(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "logsnap-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 创建配置
	config := &Config{
		ConfigDir: tempDir,
	}
	
	// 创建服务实例
	service := NewService(config, nil)
	
	// 测试保存配置
	err = service.SaveConfig()
	assert.NoError(t, err, "保存配置不应返回错误")
	
	// 验证配置文件是否存在
	configPath := filepath.Join(tempDir, "config.json")
	_, err = os.Stat(configPath)
	assert.NoError(t, err, "配置文件应该存在")
	
	// 测试加载配置
	err = service.LoadConfig()
	assert.NoError(t, err, "加载配置不应返回错误")
}
