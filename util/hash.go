package util

import (
	"crypto/md5"
	"fmt"
)

// MD5 returns the md5 hash of the given string
func MD5(value string) string {
	md5 := md5.Sum([]byte(value))
	return fmt.Sprintf("%x", md5)
}
