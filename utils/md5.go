package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 小写
func Md5Encode(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))

	return hex.EncodeToString(hash.Sum(nil))

}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))

}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)

}

// 解密
func ValidPassword(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password

}

func GetUUID() string {
	str := fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), rand.Int31n(100))
	return str
}
