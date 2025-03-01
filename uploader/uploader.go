package uploader

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"logsnap/remote"

	"github.com/sirupsen/logrus"
)

const (
	ProviderS3        = "s3"
	ProviderLocal     = "local"
	ProviderWebdav    = "webdav"
	ProviderCloudreve = "cloudreve"
)

type CloudUploaderInterface interface {
	Upload(localPath, objectKey string) (string, error)
}

// Uploader 负责将日志上传到云存储
type Uploader struct {
	config *remote.UploadConfig
}

// NewUploader 创建新的上传器
func NewUploader(cfg remote.UploadConfig) *Uploader {
	return &Uploader{
		config: &cfg,
	}
}

// Upload 上传指定的日志包到云存储
func (u *Uploader) Upload(filePath string) (string, error) {
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("上传文件不存在或无法访问: %w", err)
	}

	// 根据配置的提供商选择对应的上传实现
	var uploader CloudUploaderInterface
	provider := u.config.GetProvider(u.config.DefaultProvider)
	if provider == nil {
		return "", errors.New("不支持的云存储提供商: " + u.config.DefaultProvider)
	}
	switch provider.Provider {
	case ProviderWebdav:
		uploader = NewWebdavUploader(*provider)
		logrus.Infof("使用WebDAV上传器")
	case ProviderS3:
		uploader = NewS3Uploader(*provider)
		logrus.Infof("使用S3上传器")
	case ProviderLocal:
		uploader = NewLocalUploader(*provider)
		logrus.Infof("使用本地上传器")
	case ProviderCloudreve:
		uploader = NewCloudreveUploader(*provider)
		logrus.Infof("使用Cloudreve上传器")
	default:
		return "", errors.New("不支持的云存储提供商: " + provider.Provider)
	}

	// 生成云端对象键
	fileName := filepath.Base(filePath)

	// 计算文件的md5
	md5 := md5.Sum([]byte(filePath))
	md5Str := hex.EncodeToString(md5[:])
	objectKey := filepath.Join(provider.FolderPath, time.Now().Format("2006/01/02"), md5Str+"_"+fileName)

	// 执行上传
	return uploader.Upload(filePath, objectKey)
}
