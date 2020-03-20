package code

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"
)

// return base64 random string with encoding buffer size is l
func RandBase64(bufSize int) string {
	buf := make([]byte, bufSize)
	rand.Seed(time.Now().UnixNano())
	rand.Read(buf)

	return base64.StdEncoding.EncodeToString(buf)
}

// return a string of random number of length l
func RandNum(l int) string {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, l)
	for i := 0; i < l; i++ {
		buf[i] = '0' + byte(rand.Intn(10))
	}

	return string(buf)
}

func Md5Hex(src, salt string) string {
	sum := md5.Sum([]byte(src + salt))

	return hex.EncodeToString(sum[:])
}

func Sha1Hex(src, salt string) string {
	sum := sha1.Sum([]byte(src + salt))

	return hex.EncodeToString(sum[:])
}
