package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTime 解析时间字符串为时间对象
// 支持多种格式：
// - 时间戳（秒或毫秒）
// - ISO 8601格式 (2006-01-02T15:04:05Z)
// - 常见日期时间格式 (2006-01-02 15:04:05)
// - 常见日期格式 (2006-01-02 | 2006/01/02 | 20060102)
// - 常见时间格式 (20060102150405)
func ParseTime(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("空时间字符串")
	}

	// 首先尝试解析特殊格式 YYYY-MM-DD_HH-MM-SS_MICROS
	if strings.Count(timeStr, "_") == 2 {
		parts := strings.Split(timeStr, "_")
		if len(parts) == 3 {
			datePart := parts[0]     // YYYY-MM-DD
			timePart := parts[1]     // HH-MM-SS
			microsPart := parts[2]   // MICROS

			// 转换日期部分格式
			datePart = strings.ReplaceAll(datePart, "-", "")
			// 转换时间部分格式
			timePart = strings.ReplaceAll(timePart, "-", "")
			
			// 组合日期和时间
			basicStr := datePart + timePart
			t, err := time.ParseInLocation("20060102150405", basicStr, time.Local)
			if err == nil {
				// 解析微秒部分
				if micros, err := strconv.ParseInt(microsPart, 10, 64); err == nil {
					return t.Add(time.Duration(micros) * time.Microsecond), nil
				}
			}
		}
	}

	// 继续尝试其他常见格式
	formats := []string{
		time.RFC3339,          // 2006-01-02T15:04:05Z07:00
		"20060102",            // 常见日期格式 YYYYMMDD
		"20060102150405",      // 常见时间格式 YYYYMMDDHHMMSS
		"2006-01-02_15-04-05", // 常见时间格式 YYYY-MM-DD_HH-MM-SS
		"2006-01-02T15:04:05", // ISO 8601 无时区
		"2006-01-02 15:04:05", // 常见格式
		"2006/01/02 15:04:05", // 斜杠分隔
		"2006-01-02",          // 仅日期
		"2006/01/02",          // 仅日期，斜杠分隔
		"15:04:05",            // 仅时间
	}

	loc := time.Local // 使用本地时区
	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, loc); err == nil {
			return t, nil
		}
	}

	// 如果所有日期格式都失败了，再尝试解析为时间戳
	if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
		if timestamp > 1000000000000 {
			return time.UnixMilli(timestamp).In(loc), nil
		}
		return time.Unix(timestamp, 0).In(loc), nil
	}

	return time.Time{}, fmt.Errorf("无法解析时间: %s", timeStr)
}

// FormatTime 格式化时间为指定格式
func FormatTime(t time.Time, format string) string {
	if t.IsZero() {
		return ""
	}

	switch format {
	case "timestamp":
		return strconv.FormatInt(t.Unix(), 10)
	case "timestamp_ms":
		return strconv.FormatInt(t.UnixNano()/1000000, 10)
	case "iso8601":
		return t.Format(time.RFC3339)
	case "date":
		return t.Format("2006-01-02")
	case "time":
		return t.Format("15:04:05")
	default:
		return t.Format("2006-01-02 15:04:05")
	}
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimeString 获取当前时间字符串
func GetCurrentTimeString(format string) string {
	return FormatTime(GetCurrentTime(), format)
}
