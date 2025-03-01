package cmd

import (
	"fmt"
	"io"
	"runtime"

	"logsnap/config"
	"logsnap/remote"
	"logsnap/ui"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// versionAction 处理version命令
func versionAction(c *cli.Context) error {
	// 创建本地配置
	localConfig := config.NewLocalConfig()

	// 创建远程配置管理器
	remoteConfig := remote.NewConfigManager(localConfig)

	// 检查是否使用简单模式
	if c.Bool("simple") {
		return showVersionInSimpleMode(remoteConfig)
	}

	// 使用 TUI 界面显示版本信息
	return showVersionWithTUI(remoteConfig)
}

// showVersionInSimpleMode 在简单模式下显示版本信息
func showVersionInSimpleMode(remoteConfig *remote.ConfigManager) error {
	// 获取当前版本
	currentVersion := remoteConfig.GetLocalConfig().GetVersion()

	// 显示版本信息
	fmt.Printf("LogSnap 版本: %s\n", currentVersion)
	fmt.Printf("系统: %s\n", runtime.GOOS)
	fmt.Printf("架构: %s\n", runtime.GOARCH)

	// 尝试获取远程配置，检查是否有更新
	hasUpdate, latestVersion, _, _, _, err := remoteConfig.CheckForUpdates()
	if err != nil {
		fmt.Printf("检查更新失败: %v\n", err)
	} else if hasUpdate {
		fmt.Printf("发现新版本: %s (可以使用 'logsnap update' 命令更新)\n", latestVersion)
	} else {
		fmt.Println("当前已是最新版本")
	}

	return nil
}

// showVersionWithTUI 使用 TUI 界面显示版本信息
func showVersionWithTUI(remoteConfig *remote.ConfigManager) error {
	// 临时禁用logrus输出到终端，避免干扰TUI界面
	originalOutput := logrus.StandardLogger().Out
	logrus.SetOutput(io.Discard)

	// 运行版本显示 TUI
	ui.RunShowVersion(remoteConfig)

	// 恢复logrus输出
	logrus.SetOutput(originalOutput)

	return nil
}
