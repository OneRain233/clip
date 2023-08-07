package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func GetHash(content string) string {
	c := md5.Sum([]byte(content))
	return hex.EncodeToString(c[:])
}

func GetBase64(content []byte) string {
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
}
