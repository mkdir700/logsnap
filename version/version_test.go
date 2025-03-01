package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	// 保存原始值
	originalVersion := Version
	defer func() { Version = originalVersion }()

	// 设置测试值
	Version = "1.0.0"
	assert.Equal(t, "1.0.0", GetVersion())
}

func TestGetBuildTime(t *testing.T) {
	// 保存原始值
	originalBuildTime := BuildTime
	defer func() { BuildTime = originalBuildTime }()

	// 设置测试值
	BuildTime = "2023-01-01T12:00:00Z"
	assert.Equal(t, "2023-01-01T12:00:00Z", GetBuildTime())
}

func TestGetGitCommit(t *testing.T) {
	// 保存原始值
	originalGitCommit := GitCommit
	defer func() { GitCommit = originalGitCommit }()

	// 设置测试值
	GitCommit = "abc123"
	assert.Equal(t, "abc123", GetGitCommit())
}
