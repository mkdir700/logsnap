package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestZipDirectory(t *testing.T) {
	sourceDir := "./testdata"
	destZip := "./testresult/test.zip"
	// 创建一个空目录
	os.MkdirAll(sourceDir, 0755)
	os.MkdirAll(filepath.Dir(destZip), 0755)

	// 创建一个测试文件
	filePath := filepath.Join(sourceDir, "test.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	err := ZipDirectory(sourceDir, destZip)
	if err != nil {
		t.Fatalf("ZipDirectory failed: %v", err)
	}

	// 验证压缩文件是否存在
	if _, err := os.Stat(destZip); os.IsNotExist(err) {
		t.Fatalf("压缩文件不存在: %s", destZip)
	}

	// 验证压缩文件是否非空
	info, err := os.Stat(destZip)
	if err != nil {
		t.Fatalf("获取压缩文件信息失败: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("压缩文件为空: %s", destZip)
	}

	// 打开压缩文件并验证内容
	zipFile, err := zip.OpenReader(destZip)
	if err != nil {
		t.Fatalf("打开压缩文件失败: %v", err)
	}
	defer zipFile.Close()
	
	// 验证压缩文件中是否有一个文件
	if len(zipFile.File) != 1 {
		t.Fatalf("压缩文件中应有一个文件, 但实际有 %d 个文件", len(zipFile.File))
	}

	// 验证文件内容是否正确
	content, err := zipFile.File[0].Open()
	if err != nil {
		t.Fatalf("打开压缩文件内容失败: %v", err)
	}
	defer content.Close()
	contentBytes, err := io.ReadAll(content)
	if err != nil {
		t.Fatalf("读取压缩文件内容失败: %v", err)
	}
	if string(contentBytes) != "test" {
		t.Fatalf("压缩文件内容不正确: %s", string(contentBytes))
	}
}

