// Package encrypt 加密包
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// Md5Encode Md5加密
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

// ValidatePasswd 密码校验
func ValidatePasswd(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}

// MakePasswd 创建加密密码
func MakePasswd(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// GenerateToken 创建加密token
func GenerateToken() string {
	str := fmt.Sprintf("%d", time.Now().Unix())
	return Md5Encode(str)
}
