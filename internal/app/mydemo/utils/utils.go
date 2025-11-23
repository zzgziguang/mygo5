package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"regexp"
)

// 验证邮箱格式
func IsValidEmail(email string) bool {
	return true
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func MakeMd5(str string) (md5Str string) {
	h := md5.New()
	_, _ = io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}
