package session

import (
	"devtools/db/redisdb"
	"log"
	"testing"
)

func TestSess(t *testing.T) {
	rdsConf := &redisdb.RedisConfig{
		Host: "127.0.0.1",
		Port: "6379",
	}
	rdsPool, err := rdsConf.NewPool()
	if err != nil {
		log.Panic(err.Error())
	}
	kr, err := NewKeeper(redisdb.NewWrapper(rdsPool))
	if err != nil {
		log.Fatalln(err.Error())
	}

	var (
		cookie = kr.GenCookie()
		value  = "4567890987654567890"
		v      = ""
	)
	log.Println("cookie:", cookie, len(cookie))
	cookie, err = kr.StartSession(value)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("start session:", cookie)
	log.Println("session ttl:", kr.GetSessionTTL(cookie))
	v, err = kr.CookieValue(cookie)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("value from session:", v)
	log.Println("verify cookie:", kr.VerifyCookie(cookie, v))
	log.Println("end session:", kr.ExpireSession(cookie, -1))
	log.Println("session ttl:", kr.GetSessionTTL(cookie))
}
