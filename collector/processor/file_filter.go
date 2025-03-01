package processor

import (
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// FilterAndSortFiles 通用的文件筛选和排序函数
// 参数:
//   - files: 文件路径列表
//   - filter: 文件信息筛选器
//   - startTime: 开始时间
//   - endTime: 结束时间
//   - getEndTime: 可选函数，用于获取项目的结束时间
//
// 返回:
//   - 筛选和排序后的文件信息列表
//   - 错误信息
func FilterAndSortFiles(
	fileInfos []LogFileInfo,
	filter FileInfoFilter,
	startTime, endTime time.Time,
	getEndTime func(items []LogFileInfo, index int) time.Time,
) ([]LogFileInfo, error) {
	// 按时间排序文件
	SortByTime(fileInfos)

	// 筛选相关文件
	relevantFiles := FilterByTimeRange(fileInfos, startTime, endTime, getEndTime)

	return relevantFiles, nil
}

// SortByTime 对实现了 TimeProvider 接口的切片按时间排序
// 零时间值（IsZero()为true）的元素被视为最新的，排在最后
func SortByTime(items []LogFileInfo) {
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

// FilterFiles 是FilterAndSortFiles的简化版本，用于常见场景
// 如果存在多个文件，则至少会返回一个文件，即使该文件不在时间范围内
// 参数:
//   - files: 文件路径列表
//   - filter: 文件信息筛选器
//   - startTime: 开始时间
//   - endTime: 结束时间
//
// 返回:
//   - 筛选和排序后的文件信息列表
//   - 错误信息
func FilterFiles(
	files []string,
	filter FileInfoFilter,
	startTime, endTime time.Time,
	getEndTime func(items []LogFileInfo, index int) time.Time,
) ([]LogFileInfo, error) {
	// 解析文件信息
	fileInfos, err := filter.ParseFileInfos(files)
	if err != nil {
		return nil, err
	}

	if len(fileInfos) != len(files) {
		logrus.Warnf("文件数量不一致，files: %v, fileInfos: %v", files, fileInfos)
	}
	return FilterAndSortFiles(fileInfos, filter, startTime, endTime, getEndTime)
}
