package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}
func GenerateHashedPassword(password, salt string) (string, error) {
	// 将密码和盐拼接在一起
	saltedPassword := password + salt

	// 使用加密算法对拼接后的密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 将加密后的密码转换成字符串并返回
	return string(hashedPassword), nil
}
