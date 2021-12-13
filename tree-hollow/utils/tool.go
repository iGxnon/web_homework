package utils

import (
	"regexp"
)

// CheckPwdSafe 密码要求: 大于等于6个字符，不能是纯数字，必须有大小写字符和特殊字符，不能有空白字符
// 写的很拉胯
func CheckPwdSafe(password string) bool {

	if len(password) < 6 {
		return false
	}

	// 检测是否是纯数字
	rets := regexp.MustCompile(`\d`)
	alls := rets.FindAllStringSubmatch(password, -1)
	if len(alls) == len(password) {
		return false
	}

	// 检测是否有大写字符
	rets = regexp.MustCompile(`[A-Z]`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) == 0 {
		return false
	}

	// 检测是否有小写字符
	rets = regexp.MustCompile(`[a-z]`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) == 0 {
		return false
	}

	// 检测是否有空白字符
	rets = regexp.MustCompile(`\s`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) != 0 {
		return false
	}

	// 检测是否有特殊字符
	rets = regexp.MustCompile(`\W`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) == 0 {
		return false
	}

	return true
}
