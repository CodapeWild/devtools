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

// return a random number of length l
func RandNumInt64(l uint, neg bool) int64 {
	if l == 0 {
		return 0
	}

	var num float64 = float64(rand.Intn(9)+1) + rand.Float64()
	var i uint
	for i = 1; i < l; i++ {
		num *= 10
	}
	if neg {
		num = -num
	}

	return int64(num)
}

// return a random string of number of length l
func RandNumString(l uint, neg bool) string {
	if l == 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())

	var (
		buf []byte
		i   uint
	)
	buf = make([]byte, l)
	buf[i] = '1' + byte(rand.Intn(9))
	i++
	for ; i < l; i++ {
		buf[i] = '0' + byte(rand.Intn(10))
	}

	if neg {
		return "-" + string(buf)
	} else {
		return string(buf)
	}
}

func Md5Hex(src, salt string) string {
	sum := md5.Sum([]byte(src + salt))

	return hex.EncodeToString(sum[:])
}

func Sha1Hex(src, salt string) string {
	sum := sha1.Sum([]byte(src + salt))

	return hex.EncodeToString(sum[:])
}
