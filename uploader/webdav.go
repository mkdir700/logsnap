package uploader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"logsnap/remote"

	"github.com/sirupsen/logrus"
)

// WebdavUploader 实现WebDAV存储上传
type WebdavUploader struct {
	config remote.UploadConfigProvider
}

func NewWebdavUploader(config remote.UploadConfigProvider) *WebdavUploader {
	return &WebdavUploader{config: config}
}

func (w *WebdavUploader) Upload(localPath, objectKey string) (string, error) {
	// 读取文件内容
	fileContent, err := os.ReadFile(localPath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// 构建WebDAV URL
	webdavURL := fmt.Sprintf("%s/%s", w.config.Endpoint, objectKey)
	webdavURL = filepath.ToSlash(webdavURL) // 确保URL使用正斜杠

	// 创建HTTP请求
	req, err := http.NewRequest("PUT", webdavURL, bytes.NewReader(fileContent))
	if err != nil {
		return "", fmt.Errorf("创建WebDAV请求失败: %w", err)
	}

	// 设置基本认证
	if w.config.Username != "" {
		req.SetBasicAuth(w.config.Username, w.config.Password)
	}

	// 设置内容类型和长度
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("WebDAV上传请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 打印响应
	body, _ := io.ReadAll(resp.Body)
	logrus.Infof("WebDAV上传响应: %s", string(body))

	// 检查响应状态
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("WebDAV上传失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	logrus.Infof("成功上传到WebDAV: %s -> %s\n", localPath, webdavURL)
	return webdavURL, nil
}
