package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParseArchiveTimeStamp 测试parseArchiveTimeStamp方法
func TestParseArchiveTimeStamp(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "标准格式时间戳",
			timeStr:  "2023-01-02_15-04-05_123456",
			expected: time.Date(2023, 1, 2, 15, 4, 5, 123456000, time.Local),
			wantErr:  false,
		},
		{
			name:     "没有微秒的时间戳",
			timeStr:  "2023-01-02_15-04-05",
			expected: time.Date(2023, 1, 2, 15, 4, 5, 0, time.Local),
			wantErr:  false,
		},
		{
			name:     "无效格式时间戳",
			timeStr:  "invalid_timestamp",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArchiveTimeStamp(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArchiveTimeStamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.expected) {
				t.Errorf("ParseArchiveTimeStamp() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// 测试用的结构体，实现TimeProvider接口
type TestTimeItem struct {
	startTime time.Time
}

func (t TestTimeItem) GetStartTime() time.Time {
	return t.startTime
}

// TestSortByTime 测试SortByTime函数
func TestSortByTime(t *testing.T) {
	// 准备测试数据
	item1 := TestTimeItem{startTime: time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local)}
	item2 := TestTimeItem{startTime: time.Date(2023, 1, 2, 10, 0, 0, 0, time.Local)}
	item3 := TestTimeItem{startTime: time.Date(2023, 1, 3, 10, 0, 0, 0, time.Local)}
	itemZero := TestTimeItem{startTime: time.Time{}} // 零时间值

	// 测试常规排序
	t.Run("常规排序", func(t *testing.T) {
		items := []TestTimeItem{item3, item1, item2}
		SortByTime(items)
		assert.Equal(t, []TestTimeItem{item1, item2, item3}, items)
	})

	// 测试包含零时间值的排序
	t.Run("包含零时间值的排序", func(t *testing.T) {
		items := []TestTimeItem{item3, itemZero, item1, item2}
		SortByTime(items)
		assert.Equal(t, []TestTimeItem{item1, item2, item3, itemZero}, items)
	})

	// 测试全部为零时间值的排序
	t.Run("全部为零时间值", func(t *testing.T) {
		items := []TestTimeItem{itemZero, itemZero, itemZero}
		SortByTime(items)
		assert.Equal(t, []TestTimeItem{itemZero, itemZero, itemZero}, items)
	})

	// 测试空切片
	t.Run("空切片", func(t *testing.T) {
		items := make([]TestTimeItem, 0)
		SortByTime(items)
		assert.Equal(t, make([]TestTimeItem, 0), items)
	})
}

// TestFilterByTimeRange 测试FilterByTimeRange函数
func TestFilterByTimeRange(t *testing.T) {
	// 准备测试数据
	now := time.Now()
	item1 := TestTimeItem{startTime: now.Add(-3 * time.Hour)} // 3小时前
	item2 := TestTimeItem{startTime: now.Add(-2 * time.Hour)} // 2小时前
	item3 := TestTimeItem{startTime: now.Add(-1 * time.Hour)} // 1小时前
	itemZero := TestTimeItem{startTime: time.Time{}}          // 零时间值

	// 测试完全包含的时间范围
	t.Run("完全包含的时间范围", func(t *testing.T) {
		startTime := now.Add(-4 * time.Hour)
		endTime := now
		items := []TestTimeItem{item1, item2, item3}

		// 不指定getEndTime函数
		result := FilterByTimeRange(items, startTime, endTime, nil)
		assert.Equal(t, []TestTimeItem{item1, item2, item3}, result)

		// 指定getEndTime函数
		getEndTime := func(items []TestTimeItem, index int) time.Time {
			return items[index].startTime.Add(30 * time.Minute)
		}
		result = FilterByTimeRange(items, startTime, endTime, getEndTime)
		assert.Equal(t, []TestTimeItem{item1, item2, item3}, result)
	})

	// 测试部分包含的时间范围
	t.Run("部分包含的时间范围", func(t *testing.T) {
		startTime := now.Add(-2 * time.Hour).Add(30 * time.Minute) // 1.5小时前
		endTime := now
		items := []TestTimeItem{item1, item2, item3}

		// 使用自定义的结束时间函数
		getEndTime := func(items []TestTimeItem, index int) time.Time {
			return items[index].startTime.Add(1 * time.Hour)
		}

		result := FilterByTimeRange(items, startTime, endTime, getEndTime)
		assert.Equal(t, []TestTimeItem{item2, item3}, result)
	})

	// 测试包含零时间值
	t.Run("包含零时间值", func(t *testing.T) {
		startTime := now.Add(-4 * time.Hour)
		endTime := now
		items := []TestTimeItem{item1, itemZero, item3}

		result := FilterByTimeRange(items, startTime, endTime, nil)
		assert.Equal(t, []TestTimeItem{item1, itemZero, item3}, result)
	})

	// 测试不在范围内
	t.Run("不在范围内", func(t *testing.T) {
		startTime := now.Add(-30 * time.Minute) // 30分钟前
		endTime := now
		items := []TestTimeItem{
			{startTime: now.Add(-5 * time.Hour)}, // 5小时前
			{startTime: now.Add(-4 * time.Hour)}, // 4小时前
		}

		// 使用自定义的结束时间函数
		getEndTime := func(items []TestTimeItem, index int) time.Time {
			return items[index].startTime.Add(1 * time.Hour)
		}

		result := FilterByTimeRange(items, startTime, endTime, getEndTime)
		if len(result) == 0 {
			// 测试通过
		} else {
			t.Errorf("expected empty result but got %v", result)
		}
	})
}
