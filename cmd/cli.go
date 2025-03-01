package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

// Execute 设置并运行CLI应用
func Execute() error {
	app := &cli.App{
		Name:                 "logsnap",
		Usage:                "收集、打包并上传日志文件",
		EnableBashCompletion: true, // 启用 Bash 自动补全支持
		Commands: []*cli.Command{
			{
				Name:      "collect",
				Aliases:   []string{"c"},
				Usage:     "收集指定时间范围内的日志并上传",
				ArgsUsage: "[时间范围]",
				Action:    collectAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "time",
						Aliases: []string{"t"},
						Value:   "30m",
						Usage:   "收集最近多长时间的日志 (例如: 30m, 1h, 2d)",
					},
					&cli.StringFlag{
						Name:    "start-time",
						Aliases: []string{"s"},
						Value:   "",
						Usage:   "日志收集的开始时间 (格式: YYYY-MM-DD HH:MM:SS)",
					},
					&cli.StringFlag{
						Name:    "end-time",
						Aliases: []string{"e"},
						Value:   "",
						Usage:   "日志收集的结束时间 (格式: YYYY-MM-DD HH:MM:SS，默认为当前时间)",
					},
					&cli.PathFlag{
						Name:    "log-dir",
						Aliases: []string{"l"},
						Value:   "~/xyz_log",
						Usage:   "日志目录路径 (默认: ~/xyz_log)",
					},
					&cli.BoolFlag{
						Name:    "upload",
						Aliases: []string{"u"},
						Value:   false,
						Usage:   "是否上传到云端",
					},
					&cli.BoolFlag{
						Name:    "keep-local-snapshot",
						Aliases: []string{"k"},
						Value:   false,
						Usage:   "是否保留本地日志快照",
					},
					&cli.StringFlag{
						Name:    "output-dir",
						Aliases: []string{"o"},
						Value:   "",
						Usage:   "输出目录 (可选，默认当前目录)",
					},
					&cli.StringSliceFlag{
						Name:    "program",
						Aliases: []string{"p"},
						Usage:   "要收集的程序日志，例如：xyz-studio-max, 不指定则全部)",
					},
					&cli.BoolFlag{
						Name:  "today",
						Usage: "收集今天的日志 (从当天 00:00:00 开始)",
					},
					&cli.BoolFlag{
						Name:  "yesterday",
						Usage: "收集昨天的日志 (从昨天 00:00:00 到 23:59:59)",
					},
					&cli.BoolFlag{
						Name:  "this-week",
						Usage: "收集本周的日志 (从本周一 00:00:00 开始)",
					},
					&cli.BoolFlag{
						Name:  "skip-version-check",
						Usage: "跳过版本检查",
						Value: false,
					},
					&cli.StringFlag{
						Name:  "config-dir",
						Usage: "配置目录路径 (默认: ~/.logsnap)",
						Value: "",
					},
					&cli.BoolFlag{
						Name:  "simple",
						Usage: "使用简单模式，不显示终端动画",
						Value: true,
					},
					&cli.BoolFlag{
						Name:    "interactive",
						Aliases: []string{"I"},
						Usage:   "启用交互模式，通过UI配置选项",
						Value:   false,
					},
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "检查并更新程序到最新版本",
				Action:  updateAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "强制更新，不询问确认",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "check-only",
						Aliases: []string{"c"},
						Usage:   "仅检查是否有更新，不执行更新操作",
						Value:   false,
					},
					&cli.StringFlag{
						Name:  "config-dir",
						Usage: "配置目录路径 (默认: ~/.logsnap)",
						Value: "",
					},
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "显示当前版本信息",
				Action:  versionAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "simple",
						Aliases: []string{"s"},
						Usage:   "使用简单模式显示版本信息，不使用TUI界面",
					},
				},
			},
			{
				Name:   "supported-programs",
				Usage:  "显示支持的程序列表",
				Action: supportedProgramsAction,
			},
			{
				Name:   "completion",
				Usage:  "生成自动补全脚本",
				Action: completionAction,
				Subcommands: []*cli.Command{
					{
						Name:   "bash",
						Usage:  "生成 Bash 自动补全脚本",
						Action: bashCompletionAction,
					},
					{
						Name:   "zsh",
						Usage:  "生成 Zsh 自动补全脚本",
						Action: zshCompletionAction,
					},
					{
						Name:   "fish",
						Usage:  "生成 Fish 自动补全脚本",
						Action: fishCompletionAction,
					},
					{
						Name:   "powershell",
						Usage:  "生成 PowerShell 自动补全脚本",
						Action: powershellCompletionAction,
					},
					{
						Name:   "install",
						Usage:  "自动检测 shell 并安装补全脚本",
						Action: installCompletionAction,
					},
				},
			},
		},
	}

	return app.Run(os.Args)
}
