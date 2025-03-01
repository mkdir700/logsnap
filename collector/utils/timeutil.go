package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ParseArchiveTimeStamp(timeStr string) (time.Time, error) {
	// 首先尝试标准格式
	fileTime, err := time.ParseInLocation("2006-01-02_15-04-05_000000", timeStr, time.Local)
	if err == nil {
		return fileTime, nil
	}

	// 如果标准格式失败，尝试提取日期和时间部分
	parts := strings.Split(timeStr, "_")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("时间戳格式错误: %s", timeStr)
	}

	// 重组日期和时间部分
	dateTimeStr := parts[0] + "_" + parts[1]
	fileTime, err = time.ParseInLocation("2006-01-02_15-04-05", dateTimeStr, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	// 尝试解析微秒部分（如果有）
	if len(parts) > 2 {
		microsec, err := strconv.Atoi(parts[2])
		if err == nil && microsec >= 0 && microsec <= 999999 {
			// 添加微秒到时间
			fileTime = fileTime.Add(time.Duration(microsec) * time.Microsecond)
		}
	}

	return fileTime, nil
}

// TimeProvider 是一个接口，定义了获取时间的方法
// 任何实现了这个接口的类型都可以使用时间相关的工具函数
type TimeProvider interface {
	GetStartTime() time.Time
}

// SortByTime 对实现了 TimeProvider 接口的切片按时间排序
// 零时间值（IsZero()为true）的元素被视为最新的，排在最后
func SortByTime[T TimeProvider](items []T) {
	sort.Slice(items, func(i, j int) bool {
		// 当前日志文件始终是最新的（排在最后）
		if items[i].GetStartTime().IsZero() {
			return false
		}
		if items[j].GetStartTime().IsZero() {
			return true
		}
		return items[i].GetStartTime().Before(items[j].GetStartTime())
	})
}

// FilterByTimeRange 筛选时间范围内的项目
// 参数:
//   - items: 要筛选的项目切片
//   - startTime: 开始时间
//   - endTime: 结束时间
//   - getEndTime: 可选函数，用于获取项目的结束时间。如果为nil，则使用下一个项目的开始时间
//
// 返回:
//   - 筛选后的项目切片
func FilterByTimeRange[T TimeProvider](
	items []T,
	startTime, endTime time.Time,
	getEndTime func(items []T, index int) time.Time,
) []T {
	var result []T

	for i, item := range items {
		// 零时间值的项目总是包含在结果中
		if item.GetStartTime().IsZero() {
			result = append(result, item)
			continue
		}

		// 确定项目的结束时间
		var itemEndTime time.Time
		if getEndTime != nil {
			itemEndTime = getEndTime(items, i)
		} else if i+1 < len(items) && !items[i+1].GetStartTime().IsZero() {
			// 默认使用下一个项目的开始时间作为当前项目的结束时间
			itemEndTime = items[i+1].GetStartTime()
		} else {
			// 如果是最后一个项目，使用当前时间
			itemEndTime = time.Now()
		}

		// 检查时间范围是否重叠
		if (item.GetStartTime().Before(endTime) || item.GetStartTime().Equal(endTime)) &&
			(itemEndTime.After(startTime) || itemEndTime.Equal(startTime)) {
			result = append(result, item)
		}
	}

	return result
}
