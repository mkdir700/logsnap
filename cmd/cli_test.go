package cmd

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	// 保存原始参数
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// 测试版本命令
	os.Args = []string{"logsnap", "version", "--simple"}
	err := Execute()
	if err != nil {
		t.Errorf("执行版本命令失败: %v", err)
	}

	// 测试无效命令
	os.Args = []string{"logsnap", "invalid-command"}
	err = Execute()
	if err == nil {
		t.Error("执行无效命令应该返回错误")
	}
}
