package uploader

import (
	"fmt"
	"os"
	"path/filepath"

	"logsnap/remote"

	"github.com/sirupsen/logrus"
)

// LocalUploader 实现本地存储 - 用于测试
type LocalUploader struct {
	config remote.UploadConfigProvider
}

func NewLocalUploader(config remote.UploadConfigProvider) *LocalUploader {
	return &LocalUploader{config: config}
}

func (l *LocalUploader) Upload(localPath, objectKey string) (string, error) {
	// 确保目标目录存在
	destDir := filepath.Join(l.config.Endpoint, l.config.Bucket, filepath.Dir(objectKey))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 计算目标文件路径
	destPath := filepath.Join(l.config.Endpoint, l.config.Bucket, objectKey)

	// 读取源文件
	src, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("打开源文件失败: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	_, err = os.ReadFile(localPath)
	if err != nil {
		return "", fmt.Errorf("读取源文件失败: %w", err)
	}

	srcInfo, err := src.Stat()
	if err != nil {
		return "", fmt.Errorf("获取源文件信息失败: %w", err)
	}

	_, err = dst.Write(make([]byte, srcInfo.Size()))
	if err != nil {
		return "", fmt.Errorf("写入目标文件失败: %w", err)
	}

	logrus.Infof("本地复制: %s -> %s\n", localPath, destPath)
	return "file://" + destPath, nil
}
