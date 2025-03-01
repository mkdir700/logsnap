package service

import (
	"fmt"
	"strconv"
	"strings"
)

// Version 定义版本结构
type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion 解析版本字符串为版本对象
func ParseVersion(versionStr string) (*Version, error) {
	// 移除前缀 'v' 或 'V'
	versionStr = strings.TrimSpace(versionStr)
	if strings.HasPrefix(versionStr, "v") || strings.HasPrefix(versionStr, "V") {
		versionStr = versionStr[1:]
	}

	// 分割版本号
	parts := strings.Split(versionStr, ".")
	if len(parts) < 1 || len(parts) > 3 {
		return nil, fmt.Errorf("无效的版本格式: %s", versionStr)
	}

	// 解析主版本号
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("无效的主版本号: %s", parts[0])
	}

	// 解析次版本号
	minor := 0
	if len(parts) > 1 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("无效的次版本号: %s", parts[1])
		}
	}

	// 解析修订版本号
	patch := 0
	if len(parts) > 2 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("无效的修订版本号: %s", parts[2])
		}
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// String 返回版本的字符串表示
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare 比较两个版本
// 返回值: -1 表示 v < other, 0 表示 v == other, 1 表示 v > other
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	return 0
}

// IsNewer 检查当前版本是否比指定版本更新
func (v *Version) IsNewer(other *Version) bool {
	return v.Compare(other) > 0
}

// IsOlder 检查当前版本是否比指定版本更旧
func (v *Version) IsOlder(other *Version) bool {
	return v.Compare(other) < 0
}

// IsEqual 检查当前版本是否与指定版本相同
func (v *Version) IsEqual(other *Version) bool {
	return v.Compare(other) == 0
}

// GetCurrentVersion 获取当前应用版本
func GetCurrentVersion() *Version {
	// 这里可以从配置文件或环境变量中获取版本信息
	// 暂时返回一个硬编码的版本
	return &Version{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}
}
