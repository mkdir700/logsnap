package remote

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"logsnap/config"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// DownloadConfig 存储从远程获取的完整配置
type DownloadConfig struct {
	Version        string                 `json:"version"`
	LatestVersions map[string]string      `json:"latest_versions"`
	DownloadURLs   map[string]DownloadURL `json:"download_urls"`
	ForceUpdate    bool                   `json:"force_update"`   // 是否强制更新
	UpdateMessage  string                 `json:"update_message"` // 更新提示消息
}

type UploadConfig struct {
	Providers       []UploadConfigProvider `json:"providers"`
	DefaultProvider string                 `json:"default_provider"`
}

type UploadConfigProvider struct {
	Provider   string `json:"provider"`    // 提供商: s3, local, webdav
	Endpoint   string `json:"endpoint"`    // 服务端点
	Bucket     string `json:"bucket"`      // 存储桶名称
	AccessKey  string `json:"access_key"`  // 访问密钥
	SecretKey  string `json:"secret_key"`  // 访问密钥
	Region     string `json:"region"`      // 区域 (S3适用)
	FolderPath string `json:"folder_path"` // 上传目录路径
	Username   string `json:"username"`    // 用户名 (WebDAV适用)
	Password   string `json:"password"`    // 密码 (WebDAV适用)
}

func (u *UploadConfig) GetDefaultProvider() *UploadConfigProvider {
	for _, provider := range u.Providers {
		if provider.Provider == u.DefaultProvider {
			return &provider
		}
	}
	return nil
}

func (u *UploadConfig) GetProvider(provider string) *UploadConfigProvider {
	for _, p := range u.Providers {
		if p.Provider == provider {
			return &p
		}
	}
	return nil
}

// DownloadURL 下载URL信息
type DownloadURL struct {
	Windows string `json:"windows"`
	Linux   string `json:"linux"`
	Darwin  string `json:"darwin"`
}

// ConfigManager 负责管理远程配置
type ConfigManager struct {
	localConfig  *config.LocalConfig
	httpClient   *http.Client
	remoteConfig *DownloadConfig
	uploadConfig *UploadConfig
	configDir    string
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager(cfg *config.LocalConfig) *ConfigManager {
	// 获取配置目录
	configDir := os.Getenv("LOGSNAP_CONFIG_DIR")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configDir = filepath.Join(homeDir, ".logsnap")
		} else {
			configDir = "/etc/logsnap" // 备选目录
		}
	}

	return &ConfigManager{
		localConfig: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		configDir: configDir,
	}
}

// GetLocalConfig 返回本地配置对象
func (cm *ConfigManager) GetLocalConfig() *config.LocalConfig {
	return cm.localConfig
}

func (cm *ConfigManager) GetDownloadConfig() (*DownloadConfig, error) {
	downloadConfigURL := cm.localConfig.GetDownloadConfigURL()
	resp, err := cm.httpClient.Get(downloadConfigURL)
	if err != nil {
		return nil, fmt.Errorf("获取远程配置失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取远程配置失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取远程配置失败: %w", err)
	}

	// 解析JSON
	var downloadConfig DownloadConfig
	if err := json.Unmarshal(body, &downloadConfig); err != nil {
		return nil, fmt.Errorf("解析远程配置失败: %w", err)
	}

	return &downloadConfig, nil
}

// FetchUploadConfig 获取远程配置
func (cm *ConfigManager) FetchUploadConfig() (*UploadConfig, error) {
	uploadConfigURL := cm.localConfig.GetUploadConfigURL()
	resp, err := cm.httpClient.Get(uploadConfigURL)
	if err != nil {
		return nil, fmt.Errorf("获取上传配置失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取上传配置失败: %w", err)
	}

	// 解析JSON
	var uploadConfig UploadConfig
	if err := json.Unmarshal(body, &uploadConfig); err != nil {
		return nil, fmt.Errorf("解析上传配置失败: %w", err)
	}

	// 更新本地配置的最后检查时间
	cm.updateLastCheckTime()

	return &uploadConfig, nil
}

// GetUploadConfig 获取上传配置（优先使用远程配置）
func (cm *ConfigManager) GetUploadConfig() (*UploadConfig, error) {
	uploadConfig, err := cm.FetchUploadConfig()
	if err != nil {
		return nil, fmt.Errorf("获取上传配置失败: %w", err)
	}
	return uploadConfig, nil
}

// CheckForUpdates 检查是否有更新
func (cm *ConfigManager) CheckForUpdates() (bool, string, string, bool, string, error) {
	downloadConfig, err := cm.GetDownloadConfig()
	if err != nil {
		return false, "", "", false, "", fmt.Errorf("检查更新失败: %w", err)
	}

	// 获取当前版本
	currentVersion := cm.localConfig.GetVersion()

	// 获取最新版本
	latestVersion, ok := downloadConfig.LatestVersions["stable"]
	if !ok {
		return false, "", "", false, "", fmt.Errorf("无法获取最新版本信息")
	}

	// 比较版本号（此处使用简单字符串比较，实际可能需要更复杂的版本比较）
	hasUpdate := latestVersion != currentVersion

	// 获取下载URL
	var downloadURL string
	if downloadURLs, ok := downloadConfig.DownloadURLs["stable"]; ok {
		// 根据当前系统选择下载URL
		switch runtime.GOOS {
		case "windows":
			downloadURL = downloadURLs.Windows
		case "linux":
			downloadURL = downloadURLs.Linux
		case "darwin":
			downloadURL = downloadURLs.Darwin
		}
	}

	// 获取是否强制更新的标志和消息
	forceUpdate := downloadConfig.ForceUpdate
	updateMessage := downloadConfig.UpdateMessage

	return hasUpdate, latestVersion, downloadURL, forceUpdate, updateMessage, nil
}

// 更新最后检查时间
func (cm *ConfigManager) updateLastCheckTime() {
	// 更新本地配置的最后检查时间
	cm.localConfig.SetRemoteConfigLastCheck(time.Now().Format(time.RFC3339))
}

// DownloadUpdate 下载更新
func (cm *ConfigManager) DownloadUpdate(downloadURL string) (string, error) {
	logrus.Infof("开始下载更新: %s", downloadURL)

	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "logsnap_update")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 解析URL获取文件名
	fileName := filepath.Base(downloadURL)
	filePath := filepath.Join(tempDir, fileName)

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建下载文件失败: %w", err)
	}
	defer out.Close()

	// 发起HTTP请求
	resp, err := cm.httpClient.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("下载更新失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载更新失败，状态码: %d", resp.StatusCode)
	}

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("写入下载文件失败: %w", err)
	}

	logrus.Infof("更新下载完成: %s", filePath)
	return filePath, nil
}

// InstallUpdate 安装更新
func (cm *ConfigManager) InstallUpdate(updateFilePath string) error {
	logrus.Infof("开始安装更新: %s", updateFilePath)

	// 检查文件是否存在
	if _, err := os.Stat(updateFilePath); err != nil {
		return fmt.Errorf("更新文件不存在: %w", err)
	}

	// 创建解压目录
	extractDir := filepath.Join(os.TempDir(), "logsnap_extract")

	// 清理之前的解压目录
	os.RemoveAll(extractDir)
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("创建解压目录失败: %w", err)
	}

	// 根据文件类型进行不同的处理
	if strings.HasSuffix(updateFilePath, ".zip") {
		// 解压ZIP文件
		if err := unzipFile(updateFilePath, extractDir); err != nil {
			return fmt.Errorf("解压更新文件失败: %w", err)
		}
	} else {
		return fmt.Errorf("不支持的更新文件格式: %s", updateFilePath)
	}

	// 获取解压后的可执行文件
	executablePath, err := findExecutable(extractDir)
	if err != nil {
		return fmt.Errorf("查找可执行文件失败: %w", err)
	}

	// 获取当前程序的路径
	currentExecPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取当前程序路径失败: %w", err)
	}

	// 备份当前程序
	backupPath := currentExecPath + ".bak"
	if err := os.Rename(currentExecPath, backupPath); err != nil {
		return fmt.Errorf("备份当前程序失败: %w", err)
	}

	// 复制新程序到当前程序位置
	if err := copyFile(executablePath, currentExecPath); err != nil {
		// 恢复备份
		os.Rename(backupPath, currentExecPath)
		return fmt.Errorf("复制新程序失败: %w", err)
	}

	// 设置执行权限
	if err := os.Chmod(currentExecPath, 0755); err != nil {
		logrus.Warnf("设置执行权限失败: %v", err)
	}

	logrus.Infof("更新安装完成")
	return nil
}

// unzipFile 解压ZIP文件
func unzipFile(zipFile, destDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// 遍历ZIP文件中的所有文件
	for _, f := range r.File {
		// 构造目标文件路径
		fpath := filepath.Join(destDir, f.Name)

		// 检查文件路径是否在目标目录下（避免路径穿越攻击）
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("非法的文件路径: %s", fpath)
		}

		// 如果是目录，创建目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// 确保父目录存在
		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)

		// 创建文件
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		// 打开压缩文件
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		// 复制内容
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// findExecutable 在目录中查找可执行文件
func findExecutable(dir string) (string, error) {
	// 程序名称（基于当前操作系统）
	var execName string
	switch runtime.GOOS {
	case "windows":
		execName = "logsnap.exe"
	default:
		execName = "logsnap"
	}

	// 遍历目录查找可执行文件
	var foundPath string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == execName {
			foundPath = path
			return filepath.SkipDir // 找到后停止查找
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("未找到可执行文件: %s", execName)
	}

	return foundPath, nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 打开源文件
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 复制内容
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// 刷新缓冲区到磁盘
	return out.Sync()
}
