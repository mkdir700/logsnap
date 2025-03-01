package main

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"logsnap/cmd"
)

// 这些变量将在编译时通过 -ldflags 注入值
var (
	// Version 表示应用程序版本
	Version   string
	// BuildTime 表示构建时间
	BuildTime string
	// GitCommit 表示 Git 提交哈希
	GitCommit string
)

// ANSI 颜色代码
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldPurple = "\033[1;35m"
	BoldCyan   = "\033[1;36m"
	BoldWhite  = "\033[1;37m"
)

type Formatter struct {
	// 是否启用颜色
	ForceColors bool
	// 是否禁用时间戳
	DisableTimestamp bool
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 获取时间戳
	var timestamp string
	if !f.DisableTimestamp {
		timestamp = entry.Time.Format("2006-01-02 15:04:05")
	}

	// 根据日志级别选择颜色
	var levelColor string
	if f.ForceColors {
		switch entry.Level {
		case logrus.DebugLevel:
			levelColor = Purple
		case logrus.InfoLevel:
			levelColor = Blue
		case logrus.WarnLevel:
			levelColor = Yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = BoldRed
		default:
			levelColor = Reset
		}
	}

	// 构建日志内容
	var newLog string

	// HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		funcName := filepath.Base(entry.Caller.Function)

		if f.ForceColors {
			// 带颜色的格式
			if !f.DisableTimestamp {
				newLog = fmt.Sprintf("%s[%s]%s %s[%s]%s %s[%s:%d %s]%s %s\n",
					Cyan, timestamp, Reset,
					levelColor, entry.Level, Reset,
					BoldYellow, fName, entry.Caller.Line, funcName, Reset,
					entry.Message)
			} else {
				newLog = fmt.Sprintf("%s[%s]%s %s[%s:%d %s]%s %s\n",
					levelColor, entry.Level, Reset,
					BoldYellow, fName, entry.Caller.Line, funcName, Reset,
					entry.Message)
			}
		} else {
			// 不带颜色的格式
			if !f.DisableTimestamp {
				newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
					timestamp, entry.Level, fName, entry.Caller.Line, funcName, entry.Message)
			} else {
				newLog = fmt.Sprintf("[%s] [%s:%d %s] %s\n",
					entry.Level, fName, entry.Caller.Line, funcName, entry.Message)
			}
		}
	} else {
		// 没有调用者信息的格式
		if f.ForceColors {
			if !f.DisableTimestamp {
				newLog = fmt.Sprintf("%s[%s]%s %s[%s]%s %s\n",
					Cyan, timestamp, Reset,
					levelColor, entry.Level, Reset,
					entry.Message)
			} else {
				newLog = fmt.Sprintf("%s[%s]%s %s\n",
					levelColor, entry.Level, Reset,
					entry.Message)
			}
		} else {
			if !f.DisableTimestamp {
				newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
			} else {
				newLog = fmt.Sprintf("[%s] %s\n", entry.Level, entry.Message)
			}
		}
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func main() {
	// 使用自定义格式化器
	logrus.SetFormatter(&Formatter{
		ForceColors:      true,
		DisableTimestamp: false,
	})
	// 启用调用者报告
	logrus.SetReportCaller(true)
	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)

	err := cmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
