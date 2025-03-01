package cmd

import (
	"fmt"
	"regexp"
	"strconv"
)

// parseTimeArg 解析时间参数，支持多种时间单位
func parseTimeArg(timeArg string) (int, error) {
	if timeArg == "" {
		return 30, nil // 默认30分钟
	}

	// 使用正则表达式解析时间参数
	re := regexp.MustCompile(`^(\d+)([mhdw])$`)
	matches := re.FindStringSubmatch(timeArg)

	if matches == nil {
		return 0, fmt.Errorf("无效的时间格式: %s, 请使用如 30m, 1h, 2d 的格式", timeArg)
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("无效的时间值: %s", matches[1])
	}

	unit := matches[2]

	// 转换为分钟
	switch unit {
	case "m": // 分钟
		return value, nil
	case "h": // 小时
		return value * 60, nil
	case "d": // 天
		return value * 60 * 24, nil
	case "w": // 周
		return value * 60 * 24 * 7, nil
	default:
		return 0, fmt.Errorf("不支持的时间单位: %s", unit)
	}
}
