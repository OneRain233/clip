package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetHash(content string) string {
	c := md5.Sum([]byte(content))
	return hex.EncodeToString(c[:])
}
