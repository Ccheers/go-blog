package common

import (
	"crypto/md5"
	"encoding/hex"
)

// COMMON MD5 FUNC
func Md5(s string) string {
	c := md5.New()
	c.Write([]byte(s))
	cipherStr := c.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
