package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"logsnap/config"
	"logsnap/remote"
	"logsnap/ui"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// updateAction 处理update命令
func updateAction(c *cli.Context) error {
	// 获取配置目录
	configDir := c.String("config-dir")

	// 创建本地配置和远程配置管理器
	localConfig := config.NewLocalConfig()
	if configDir != "" {
		os.Setenv("LOGSNAP_CONFIG_DIR", configDir)
	}
	remoteConfig := remote.NewConfigManager(localConfig)

	// 根据模式选择不同的实现
	if c.Bool("simple") || c.Bool("check-only") {
		// 简单模式或仅检查模式使用命令行界面
		return updateInSimpleMode(c, remoteConfig)
	} else {
		// 默认使用 TUI 界面
		return updateWithTUI(c, remoteConfig)
	}
}

// updateInSimpleMode 在简单模式下执行更新操作
func updateInSimpleMode(c *cli.Context, remoteConfig *remote.ConfigManager) error {
	// 检查是否有更新
	logrus.Infof("正在检查更新...")
	hasUpdate, latestVersion, downloadURL, forceUpdate, updateMessage, err := remoteConfig.CheckForUpdates()
	if err != nil {
		return fmt.Errorf("检查更新失败: %v", err)
	}

	// 获取当前版本
	localConfig := remoteConfig.GetLocalConfig()
	currentVersion := localConfig.GetVersion()

	// 显示版本信息
	logrus.Infof("当前版本: %s", currentVersion)
	logrus.Infof("最新版本: %s", latestVersion)

	// 如果没有更新，直接返回
	if !hasUpdate {
		logrus.Infof("您已经使用的是最新版本，无需更新。")
		return nil
	}

	// 显示更新信息
	logrus.Infof("发现新版本: %s", latestVersion)
	if updateMessage != "" {
		logrus.Infof(updateMessage)
	}

	// 如果只是检查更新，不执行更新操作
	if c.Bool("check-only") {
		logrus.Infof("检查完成，有可用更新。使用 'logsnap update' 命令进行更新。")
		return nil
	}

	// 判断是否需要强制更新
	shouldUpdate := forceUpdate || c.Bool("force")

	// 如果不是强制更新，询问用户是否要更新
	if !shouldUpdate {
		var response string
		logrus.Infof("是否要更新到最新版本? (y/n): ")
		fmt.Scanln(&response)

		// 检查用户响应
		response = strings.ToLower(strings.TrimSpace(response))
		shouldUpdate = response == "y" || response == "yes"
	}

	// 如果用户选择不更新，直接返回
	if !shouldUpdate {
		logrus.Infof("更新已取消。")
		return nil
	}

	// 执行更新操作
	logrus.Infof("开始下载更新...")
	updateFilePath, err := remoteConfig.DownloadUpdate(downloadURL)
	if err != nil {
		return fmt.Errorf("下载更新失败: %v", err)
	}

	logrus.Infof("下载完成，正在安装更新...")
	if err := remoteConfig.InstallUpdate(updateFilePath); err != nil {
		return fmt.Errorf("安装更新失败: %v", err)
	}

	logrus.Infof("更新成功！程序已更新到版本 %s", latestVersion)
	logrus.Infof("请重新启动程序以应用更新。")

	return nil
}

// updateWithTUI 使用 TUI 界面执行更新操作
func updateWithTUI(c *cli.Context, remoteConfig *remote.ConfigManager) error {
	// 临时禁用logrus输出到终端，避免干扰TUI界面
	originalOutput := logrus.StandardLogger().Out
	logrus.SetOutput(io.Discard)

	// 运行自动更新版本检查器 TUI
	ui.RunAutoUpdateVersionChecker(remoteConfig)

	// 恢复logrus输出
	logrus.SetOutput(originalOutput)

	return nil
}
