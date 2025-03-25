package service

import (
	"fmt"
	"time"

	"logsnap/remote"
	"logsnap/uploader"
)

// UploadRequest 定义上传请求结构
type UploadRequest struct {
	File        *LogFile             // 要上传的文件
	Config      *remote.UploadConfig // 配置信息
	Reporter    ProgressReporter     // 进度报告器
	Description string               // 上传描述
	Tags        []string             // 标签
}

// UploadResult 定义上传结果结构
type UploadResult struct {
	Success      bool      // 是否成功
	Message      string    // 消息
	URL          string    // 上传后的URL
	UploadedTime time.Time // 上传时间
	FileCount    int       // 文件数量
	TotalSize    int64     // 总大小
}

// UploadManager 定义上传管理器接口
type UploadManager interface {
	Upload(request *UploadRequest) (*UploadResult, error)
}

// DefaultUploadManager 默认上传管理器实现
type DefaultUploadManager struct {
	uploadConfig *remote.UploadConfig
}

// NewUploadManager 创建新的上传管理器
func NewUploadManager(uploadConfig *remote.UploadConfig) UploadManager {
	return &DefaultUploadManager{
		uploadConfig: uploadConfig,
	}
}

// Upload 执行上传操作
func (m *DefaultUploadManager) Upload(request *UploadRequest) (*UploadResult, error) {
	// 计算总大小
	var totalSize int64
	if request.File != nil {
		if !request.File.IsDir {
			totalSize += request.File.Size
		}
	}

	// 更新进度
	if request.Reporter != nil {
		request.Reporter.Report("upload", 50, "压缩完成，开始上传")
	}

	// 创建上传器
	uploaderInstance := uploader.NewUploader(*m.uploadConfig)

	// 执行上传操作
	url, err := uploaderInstance.Upload(request.File.Path)
	if err != nil {
		if request.Reporter != nil {
			request.Reporter.Report("upload", 100, fmt.Sprintf("上传失败: %v", err))
		}
		return nil, fmt.Errorf("上传文件失败: %v", err)
	}

	// 更新进度
	if request.Reporter != nil {
		request.Reporter.Report("upload", 100, "上传完成")
	}

	// 返回上传结果
	return &UploadResult{
		Success:      true,
		Message:      "上传成功",
		URL:          url,
		UploadedTime: time.Now(),
		FileCount:    1,
		TotalSize:    totalSize,
	}, nil
}
