package md5

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func GetMD5Upper(s string) string {
	return strings.ToUpper(Md5(s))
}
