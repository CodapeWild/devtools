package code

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/CodapeWild/devtools/db/redisdb"
	"gopkg.in/mgo.v2/bson"
)

func TestAuthCode(t *testing.T) {
	rdsConf := &redisdb.RedisConfig{
		Host: "127.0.0.1",
		Port: "6379",
	}
	rdsPool, err := rdsConf.NewPool()
	if err != nil {
		log.Panic(err.Error())
	}
	acr, err := NewAuthCoder(redisdb.NewWrapper(rdsPool), SetBits(15))
	if err != nil {
		log.Panic(err.Error())
	}
	flavor := "192.122.3.6"
	auth, err := acr.SetAuthCode(flavor)
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println("auth code:", auth, len(auth))

	fmt.Println(acr.VerifyAuthCode(flavor, auth))
}

func TestRandBase64(t *testing.T) {
	for i := 0; i < 100; i++ {
		log.Println(RandBase64(20))
	}
}

func TestSumHex(t *testing.T) {
	log.Println("encoded password:", Md5Hex("1233456", "X#d12s*dsf&^df"))
	log.Println("encoded id:", Sha1Hex(bson.NewObjectId().String(), "X#d12s*dsf&^df"))
}

func TestRandNumInt64(t *testing.T) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			log.Println(RandNumInt64(uint(i), rand.Intn(10) >= 5))
		}
	}
}

func TestRandNumString(t *testing.T) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			log.Println(RandNumString(uint(i), rand.Intn(10) >= 5))
		}
	}
}
