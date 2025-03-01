package utils

import (
	"testing"
	"time"
)

func TestGetCurrentTime(t *testing.T) {
	now := GetCurrentTime()
	if now.Location() != time.Local {
		t.Errorf("GetCurrentTime() 返回的时间不是本地时区")
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "标准日期时间",
			input:    "2024-01-01 12:00:00",
			expected: time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local),
		},
		{
			name:     "仅日期带横杠",
			input:    "2024-01-01",
			expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			name:     "紧凑日期格式",
			input:    "20240202",
			expected: time.Date(2024, 2, 2, 0, 0, 0, 0, time.Local),
		},
		{
			name:     "带T的ISO格式",
			input:    "2024-02-02T15:04:05",
			expected: time.Date(2024, 2, 2, 15, 4, 5, 0, time.Local),
		},
		{
			name:     "Unix时间戳(秒)",
			input:    "1706803200",
			expected: time.Date(2024, 2, 2, 0, 0, 0, 0, time.Local),
		},
		{
			name:     "常见时间格式 YYYY-MM-DD_HH-MM-SS",
			input:    "2024-02-02_15-04-05",
			expected: time.Date(2024, 2, 2, 15, 4, 5, 0, time.Local),
		},
		{
			name:     "常见时间格式 YYYY-MM-DD_HH-MM-SS_MS",
			input:    "2025-03-01_17-13-33_527057",
			expected: time.Date(2025, 3, 1, 17, 13, 33, 527057000, time.Local),
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !result.Equal(tt.expected) {
				t.Errorf("ParseTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}


