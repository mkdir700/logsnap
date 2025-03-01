package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocalConfig 模拟LocalConfig
type MockLocalConfig struct {
	mock.Mock
}

func (m *MockLocalConfig) GetVersion() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLocalConfig) GetRemoteConfigEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLocalConfig) GetUploadConfigURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLocalConfig) GetDownloadConfigURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLocalConfig) GetRemoteConfigInterval() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockLocalConfig) GetRemoteConfigLastCheck() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLocalConfig) SetRemoteConfigLastCheck(lastCheck string) {
	m.Called(lastCheck)
}

// 测试UploadConfig.GetDefaultProvider方法
func TestUploadConfigGetDefaultProvider(t *testing.T) {
	// 创建上传配置
	uploadConfig := &UploadConfig{
		Providers: []UploadConfigProvider{
			{
				Provider:   "s3",
				Endpoint:   "https://s3.example.com",
				Bucket:     "test-bucket",
				AccessKey:  "test-key",
				SecretKey:  "test-secret",
				Region:     "us-east-1",
				FolderPath: "/logs",
			},
		},
		DefaultProvider: "s3",
	}
	
	// 获取默认提供商
	provider := uploadConfig.GetDefaultProvider()
	
	// 验证结果
	assert.NotNil(t, provider, "默认提供商不应为空")
	assert.Equal(t, "s3", provider.Provider, "默认提供商类型应该正确")
	
	// 测试不存在的默认提供商
	uploadConfig.DefaultProvider = "not-exist"
	provider = uploadConfig.GetDefaultProvider()
	assert.Nil(t, provider, "不存在的默认提供商应该返回nil")
}

// 测试UploadConfig.GetProvider方法
func TestUploadConfigGetProvider(t *testing.T) {
	// 创建上传配置
	uploadConfig := &UploadConfig{
		Providers: []UploadConfigProvider{
			{
				Provider:   "s3",
				Endpoint:   "https://s3.example.com",
				Bucket:     "test-bucket",
				AccessKey:  "test-key",
				SecretKey:  "test-secret",
				Region:     "us-east-1",
				FolderPath: "/logs",
			},
			{
				Provider:   "webdav",
				Endpoint:   "https://webdav.example.com",
				Username:   "test-user",
				Password:   "test-password",
				FolderPath: "/logs",
			},
		},
		DefaultProvider: "s3",
	}
	
	// 获取指定提供商
	provider := uploadConfig.GetProvider("webdav")
	
	// 验证结果
	assert.NotNil(t, provider, "指定提供商不应为空")
	assert.Equal(t, "webdav", provider.Provider, "指定提供商类型应该正确")
	
	// 测试不存在的提供商
	provider = uploadConfig.GetProvider("not-exist")
	assert.Nil(t, provider, "不存在的提供商应该返回nil")
}

// 测试ConfigManager的基本功能
func TestConfigManagerBasic(t *testing.T) {
	// 跳过此测试，需要更多的模拟工作
	t.Skip("需要更多模拟工作来测试ConfigManager")
}
