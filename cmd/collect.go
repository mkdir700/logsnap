package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"logsnap/config"
	"logsnap/remote"
	"logsnap/service"
	"logsnap/utils"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// collectAction 处理collect命令
func collectAction(c *cli.Context) error {
	// 处理位置参数（如果有）
	timeArg := ""
	if c.Args().Len() > 0 {
		timeArg = c.Args().First()
	} else {
		timeArg = c.String("time")
	}

	// 检查是否提供了时间范围
	startTime := c.String("start-time")
	endTime := c.String("end-time")
	var startTimeVal *time.Time
	var endTimeVal *time.Time

	// 检查互斥选项
	timeOptionCount := 0
	if c.Bool("today") {
		timeOptionCount++
	}
	if c.Bool("yesterday") {
		timeOptionCount++
	}
	if c.Bool("this-week") {
		timeOptionCount++
	}
	if startTime != "" {
		timeOptionCount++
	}
	// 只有当 timeArg 不是默认值且不是从 --time 选项获取的默认值时才计数
	if timeArg != "30m" && timeArg != "00:00:00" && c.Args().Len() > 0 {
		timeOptionCount++
	}

	if timeOptionCount > 1 {
		return fmt.Errorf("时间选项 (--time 位置参数, --start-time, --today, --yesterday, --this-week) 不能同时使用")
	}

	now := utils.GetCurrentTime()

	// 处理便捷时间选项
	if c.Bool("today") {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		startTimeVal = &today
		endTimeVal = &now
	} else if c.Bool("yesterday") {
		yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
		yesterdayEnd := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, time.Local)
		startTimeVal = &yesterday
		endTimeVal = &yesterdayEnd
	} else if c.Bool("this-week") {
		// 计算本周一
		daysFromMonday := (int(now.Weekday()) + 6) % 7 // 将周日视为7，周一为0
		monday := time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday, 0, 0, 0, 0, time.Local)
		startTimeVal = &monday
		endTimeVal = &now
	} else if timeArg != "" && startTime == "" {
		// 处理时间参数 (如 30m, 1h, 2d)
		minutes, err := parseTimeArg(timeArg)
		if err != nil {
			return err
		}

		// 计算开始时间
		start := now.Add(-time.Duration(minutes) * time.Minute)
		startTimeVal = &start
		endTimeVal = &now
	}

	if startTime != "" {
		t, err := utils.ParseTime(startTime)
		if err != nil {
			return err
		}
		startTimeVal = &t
	}

	if endTime == "" {
		endTimeVal = &now
	} else {
		t, err := utils.ParseTime(endTime)
		if err != nil {
			return err
		}
		endTimeVal = &t
	}

	logrus.Infof("startTime: %v", startTimeVal)
	logrus.Infof("endTime: %v", endTimeVal)

	// 设置输出目录
	outputDir := c.String("output-dir")
	if outputDir == "" {
		outputDir = "."
	}

	// 创建配置对象
	serviceConfig := service.Config{
		StartTime:        startTimeVal,
		EndTime:          endTimeVal,
		ShouldUpload:     c.Bool("upload"),
		KeepLocalSnap:    c.Bool("keep-local-snapshot"),
		OutputDir:        outputDir,
		SkipVersionCheck: c.Bool("skip-version-check"),
		ConfigDir:        c.String("config-dir"),
		LogRootDir:       c.String("log-dir"),
		Programs:         c.StringSlice("program"),
	}

	// 如果指定了程序，记录日志
	if len(c.StringSlice("program")) > 0 {
		logrus.Infof("指定收集的程序日志: %v", c.StringSlice("program"))
	} else {
		logrus.Infof("未指定要收集的程序日志，将收集所有支持的日志")
	}

	// 获取远程配置
	localConfig := config.NewLocalConfig()
	remoteConfig := remote.NewConfigManager(localConfig)

	// 简单模式，使用直接调用方式
	return runInSimpleMode(&serviceConfig, remoteConfig)
}

// runInSimpleMode 在简单模式下执行收集操作
func runInSimpleMode(config *service.Config, remoteConfig *remote.ConfigManager) error {
	// 检查版本更新（如果未跳过版本检查）
	if !config.SkipVersionCheck {
		hasUpdate, latestVersion, downloadURL, forceUpdate, updateMessage, err := remoteConfig.CheckForUpdates()
		if err != nil {
			logrus.Warnf("检查版本更新失败: %v", err)
		} else if hasUpdate {
			// 获取当前版本信息
			localCfg := remoteConfig.GetLocalConfig()
			currentVersion := localCfg.GetVersion()
			logrus.Infof("发现新版本: %s (当前版本: %s)", latestVersion, currentVersion)
			if updateMessage != "" {
				logrus.Infof(updateMessage)
			}

			if forceUpdate {
				logrus.Warnf("需要强制更新到最新版本")
				// 下载并安装更新
				if downloadURL != "" {
					logrus.Infof("正在下载更新...")
					updateFilePath, err := remoteConfig.DownloadUpdate(downloadURL)
					if err != nil {
						logrus.Errorf("下载更新失败: %v", err)
					} else {
						logrus.Infof("下载完成，正在安装更新...")
						if err := remoteConfig.InstallUpdate(updateFilePath); err != nil {
							logrus.Errorf("安装更新失败: %v", err)
						} else {
							return fmt.Errorf("程序已更新到最新版本 %s，请重新启动程序", latestVersion)
						}
					}
				}
			} else {
				logrus.Infof("建议更新到最新版本，请访问官方网站下载")
			}
		}
	}

	uploadConfig, err := remoteConfig.GetUploadConfig()
	if err != nil {
		logrus.Warnf("获取上传配置失败: %v", err)
		return err
	}

	if uploadConfig == nil {
		// 抛出错误
		return fmt.Errorf("获取上传配置失败: %v", err)
	}

	snapPath, uploadURL, err := service.CollectAndUploadLogs(config, uploadConfig)
	if err != nil {
		// 特殊处理需要重启的错误
		if strings.Contains(err.Error(), "程序已更新到最新版本") {
			fmt.Println("========= 自动更新完成 =========")
			fmt.Println(err.Error())
			fmt.Println("请重新执行命令继续操作。")
			return nil
		}
		return err
	}

	logrus.Infof("日志收集完成，已保存至: %s", snapPath)
	if config.ShouldUpload && uploadURL != "" {
		if !config.KeepLocalSnap {
			os.RemoveAll(snapPath)
		}
		logrus.Infof("日志已上传至: %s", uploadURL)
	}

	return nil
}
