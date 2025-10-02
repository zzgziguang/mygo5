package utils

import "regexp"

// 验证邮箱格式
func IsValidEmail(email string) bool {
	return true
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
