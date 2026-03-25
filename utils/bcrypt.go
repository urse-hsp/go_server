package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 将明文密码转换为哈希字符串
// 直接返回 string，自动处理 []byte 转换
func HashPassword(password string) (string, error) {
	// 这里统一使用 DefaultCost，也可以配置为 12 或 14
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 校验明文密码是否匹配哈希值
// 返回 bool，更符合业务判断直觉
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
