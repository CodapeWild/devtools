package code

import (
	"devtools/db/redisdb"
)

const (
	rds_prefix    = "auth_coder_"
	rds_set_Value = "OK"
)

var (
	defBits         = 6
	defExpSec int64 = 5 * 60
)

type AuthCoder struct {
	bits       int
	expSec     int64
	rdsWrapper *redisdb.RedisWrapper
}

type AuthCoderSetting func(*AuthCoder)

func SetBits(bits int) AuthCoderSetting {
	return func(acr *AuthCoder) {
		acr.bits = bits
	}
}

func SetExpSec(sec int64) AuthCoderSetting {
	return func(acr *AuthCoder) {
		acr.expSec = sec
	}
}

func NewAuthCoder(rdsWrapper *redisdb.RedisWrapper, settings ...AuthCoderSetting) (*AuthCoder, error) {
	acr := &AuthCoder{
		bits:       defBits,
		expSec:     defExpSec,
		rdsWrapper: rdsWrapper,
	}

	for _, s := range settings {
		s(acr)
	}

	return acr, nil
}

func genAuthKey(flavor, auth string) string {
	return rds_prefix + auth + flavor
}

func (this *AuthCoder) SetAuthCodeWithTimeout(flavor string, expsec int64) (auth string, err error) {
	if expsec <= 0 {
		expsec = this.expSec
	}

	auth = RandNum(this.bits)
	if _, err = this.rdsWrapper.Set(genAuthKey(flavor, auth), rds_set_Value, expsec); err != nil {
		return "", err
	}

	return
}

func (this *AuthCoder) SetAuthCode(flavor string) (auth string, err error) {
	return this.SetAuthCodeWithTimeout(flavor, this.expSec)
}

func (this *AuthCoder) VerifyAuthCode(flavor, auth string) (ok bool) {
	if len(auth) != defBits {
		return false
	}

	key := genAuthKey(flavor, auth)
	defer func() {
		if ok {
			this.rdsWrapper.DelKey(key)
		}
	}()

	rply, err := this.rdsWrapper.Get(key)

	return err == nil && rply != nil && string(rply.([]uint8)) == rds_set_Value
}
