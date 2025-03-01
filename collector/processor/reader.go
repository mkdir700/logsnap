package processor

import (
	"archive/zip"
	"io"
	"time"
)

// LogFileInfo 表示日志文件的基本信息
type LogFileInfo struct {
	Path      string
	FileName  string
	StartTime time.Time
	FileType  string // 文件类型，由调用者定义，例如 "zip", "log" 等
	Extra     any    // 额外的特定解析器信息
}

// GetStartTime 实现 TimeProvider 接口
func (l LogFileInfo) GetStartTime() time.Time {
	return l.StartTime
}



// CompositeReadCloser 是一个复合的ReadCloser，用于同时关闭多个资源
type CompositeReadCloser struct {
	reader    io.ReadCloser
	zipReader *zip.ReadCloser
}

func NewCompositeReadCloser(reader io.ReadCloser, zipReader *zip.ReadCloser) *CompositeReadCloser {
	return &CompositeReadCloser{
		reader:    reader,
		zipReader: zipReader,
	}
}

func (c *CompositeReadCloser) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *CompositeReadCloser) Close() error {
	err1 := c.reader.Close()
	err2 := c.zipReader.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
