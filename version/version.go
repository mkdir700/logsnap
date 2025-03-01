package version

// 这些变量将在编译时通过 -ldflags 注入值
var (
	// Version 表示应用程序版本
	Version   string
	// BuildTime 表示构建时间
	BuildTime string
	// GitCommit 表示 Git 提交哈希
	GitCommit string
)

// GetVersion 返回应用程序版本
func GetVersion() string {
	return Version
}

// GetBuildTime 返回构建时间
func GetBuildTime() string {
	return BuildTime
}

// GetGitCommit 返回Git提交哈希
func GetGitCommit() string {
	return GitCommit
} 