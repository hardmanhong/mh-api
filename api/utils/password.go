package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

func GenerateString(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt)[:length], nil
}
func GenerateHashedPassword(password, salt string) (string, error) {
	// 将密码和盐拼接在一起
	saltedPassword := password + salt

	// 计算 MD5 哈希值
	hasher := md5.New()
	hasher.Write([]byte(saltedPassword))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// 返回 32 位的十六进制字符串
	return hashedPassword, nil
}
