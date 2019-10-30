package session

import (
	"crypto/md5"
	"devtools/comerr"
	"devtools/db/redisdb"
	"encoding/base64"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	rds_keeper = "keeper_"
)

var (
	defBufBits              = 256
	defExpSec         int64 = 12 * 24 * 60 * 60
	defGenFailedTimes       = 3
)

var (
	GenCookieFailed = errors.New("generate cookie failed")
	CookieInvalid   = errors.New("cookie invalid")
	CookieExpired   = errors.New("cookie expired")
)

type Keeper struct {
	bufBits, cookieBits  int
	expsec               int64
	genCookieFailedTimes int
	rdsWrapper           *redisdb.RedisWrapper
}

type KeeperSetting func(*Keeper)

func SetBufBits(bits int) KeeperSetting {
	return func(kr *Keeper) {
		kr.bufBits = bits
	}
}

func SetExpSec(sec int64) KeeperSetting {
	return func(kr *Keeper) {
		kr.expsec = sec
	}
}

func SetGenCookieFailedTimes(times int) KeeperSetting {
	return func(kr *Keeper) {
		kr.genCookieFailedTimes = times
	}
}

func NewKeeper(rdsWrapper *redisdb.RedisWrapper, settings ...KeeperSetting) (*Keeper, error) {
	kr := &Keeper{
		bufBits:              defBufBits,
		cookieBits:           base64.StdEncoding.EncodedLen(md5.Size),
		expsec:               defExpSec,
		genCookieFailedTimes: defGenFailedTimes,
		rdsWrapper:           rdsWrapper,
	}

	for _, s := range settings {
		s(kr)
	}

	return kr, nil
}

func genKeeperKey(cookie string) string {
	return rds_keeper + cookie
}

func (this *Keeper) GenCookie() string {
	var err error
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, this.bufBits)
	for i := 0; i < this.genCookieFailedTimes; i++ {
		if _, err = rand.Read(buf); err == nil {
			s := md5.Sum(buf)

			return base64.StdEncoding.EncodeToString(s[:])
		}
	}
	panic(err)
}

func (this *Keeper) StartSessionWithTimeout(value string, expsec int64) (cookie string, err error) {
	if value == "" {
		return "", comerr.ParamInvalid
	}
	if expsec <= 0 {
		expsec = this.expsec
	}

	cookie = this.GenCookie()
	_, err = this.rdsWrapper.Set(genKeeperKey(cookie), value, expsec)

	return
}

func (this *Keeper) StartSession(value string) (cookie string, err error) {
	return this.StartSessionWithTimeout(value, this.expsec)
}

/*
	sec>0  session keeps alive for sec seconds
	sec<=0 session closed
*/
func (this *Keeper) GetSessionTTL(cookie string) (sec int64) {
	if len(cookie) != this.cookieBits {
		log.Println(CookieInvalid.Error())

		return -1
	}

	var err error
	sec, err = this.rdsWrapper.TTL(genKeeperKey(cookie))
	if err != nil {
		log.Println(err.Error())

		return -1
	}

	return
}

func (this *Keeper) CookieValue(cookie string) (value string, err error) {
	if len(cookie) != this.cookieBits {
		return "", CookieInvalid
	}

	return redis.String(this.rdsWrapper.Get(genKeeperKey(cookie)))
}

func (this *Keeper) VerifyCookie(cookie, value string) bool {
	if len(cookie) != this.cookieBits {
		return false
	}

	if v, err := this.CookieValue(cookie); err != nil {
		return false
	} else {
		return v == value
	}
}

/*
	expsec < 0 session expired immediately
	expsec = 0 session expire reset by this.expsec
	expsec > 0 session expire reset by expsec
*/
func (this *Keeper) ExpireSession(cookie string, expsec int64) error {
	if len(cookie) != this.cookieBits {
		return CookieInvalid
	}

	if expsec < 0 {
		return this.rdsWrapper.DelKey(genKeeperKey(cookie))
	} else if expsec == 0 {
		return this.rdsWrapper.Expire(genKeeperKey(cookie), this.expsec)
	} else {
		return this.rdsWrapper.Expire(genKeeperKey(cookie), expsec)
	}
}
