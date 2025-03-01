package service

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"logsnap/remote"
	"logsnap/uploader"
)

// UploadRequest 定义上传请求结构
type UploadRequest struct {
	Files       []*LogFile           // 要上传的文件列表
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
	if request == nil || len(request.Files) == 0 {
		return nil, fmt.Errorf("没有要上传的文件")
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "logsnap_upload_")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建临时压缩文件
	zipFilePath := filepath.Join(tempDir, "logs.zip")
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return nil, fmt.Errorf("创建压缩文件失败: %v", err)
	}
	defer zipFile.Close()

	// 创建zip写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 计算总大小
	var totalSize int64
	for _, file := range request.Files {
		if !file.IsDir {
			totalSize += file.Size
		}
	}

	// 初始化进度
	if request.Reporter != nil {
		request.Reporter.Report("compress", 0, fmt.Sprintf("开始压缩 %d 个文件", len(request.Files)))
	}

	// 添加文件到压缩包
	var processedSize int64
	for i, file := range request.Files {
		if file.IsDir {
			continue
		}

		// 更新进度
		if request.Reporter != nil {
			progress := float64(i) / float64(len(request.Files)) * 50
			request.Reporter.Report("compress", int(progress), fmt.Sprintf("正在压缩 %s", file.Name))
		}

		// 打开源文件
		srcFile, err := os.Open(file.Path)
		if err != nil {
			return nil, fmt.Errorf("打开文件失败 %s: %v", file.Path, err)
		}

		// 获取文件信息
		fileInfo, err := srcFile.Stat()
		if err != nil {
			srcFile.Close()
			return nil, fmt.Errorf("获取文件信息失败 %s: %v", file.Path, err)
		}

		// 创建zip文件头
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			srcFile.Close()
			return nil, fmt.Errorf("创建文件头失败 %s: %v", file.Name, err)
		}

		// 设置压缩方法
		header.Method = zip.Deflate

		// 使用相对路径作为文件名
		header.Name = file.RelativePath()

		// 创建zip文件
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			srcFile.Close()
			return nil, fmt.Errorf("创建压缩文件失败 %s: %v", file.Name, err)
		}

		// 复制文件内容
		_, err = io.Copy(writer, srcFile)
		srcFile.Close()
		if err != nil {
			return nil, fmt.Errorf("写入文件内容失败 %s: %v", file.Name, err)
		}

		processedSize += file.Size
	}

	// 关闭zip写入器
	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭压缩文件失败: %v", err)
	}

	// 更新进度
	if request.Reporter != nil {
		request.Reporter.Report("upload", 50, "压缩完成，开始上传")
	}

	// 创建上传器
	uploaderInstance := uploader.NewUploader(*m.uploadConfig)

	// 执行上传操作
	url, err := uploaderInstance.Upload(zipFilePath)
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
		FileCount:    len(request.Files),
		TotalSize:    totalSize,
	}, nil
}
