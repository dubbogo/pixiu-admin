package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 md5加密
func Md5(encodeString string) string {
	h := md5.New()
	h.Write([]byte(encodeString))
	return hex.EncodeToString(h.Sum(nil))
}
