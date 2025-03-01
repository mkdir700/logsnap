package cmd

import (
	"testing"
)

func TestParseTimeArg(t *testing.T) {
	tests := []struct {
		name     string
		timeArg  string
		expected int
		wantErr  bool
	}{
		{"分钟", "30m", 30, false},
		{"小时", "2h", 120, false},
		{"天", "1d", 1440, false},
		{"无效格式", "abc", 0, true},
		{"负数", "-10m", 0, true},
		{"小数", "1.5h", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTimeArg(tt.timeArg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTimeArg() 错误 = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("parseTimeArg() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}
