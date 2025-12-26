package util

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

// IsVersionGreaterThan 比较两个版本字符串，判断 newVersion 是否 > oldVersion
// 支持带 "v" 前缀（如 "v1.2.3"），符合 SemVer 规范
// 如果任一版本格式非法，返回 false 和 error
func IsVersionGreaterThan(newVersion, oldVersion string) error {
	if newVersion == "" || oldVersion == "" {
		return fmt.Errorf("version cannot be empty")
	}

	vOld, err := version.NewVersion(oldVersion)
	if err != nil {
		return fmt.Errorf("invalid old version '%s': %w", oldVersion, err)
	}

	vNew, err := version.NewVersion(newVersion)
	if err != nil {
		return fmt.Errorf("invalid new version '%s': %w", newVersion, err)
	}
	if vNew.GreaterThan(vOld) {
		return nil
	} else {
		return fmt.Errorf("the version number is not self-incrementing")
	}
}
