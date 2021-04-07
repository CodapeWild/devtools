package session

import (
	"crypto/md5"
	"encoding/base64"
	"math/rand"
	"time"
)

var (
	DefBufSize int   = 256
	DefExpSec  int64 = 7 * 60 * 60
)

type RandToken struct {
	store              SessStore
	bufSize, tokenSize int
}

type RandTokenSetting func(*RandToken)

func SetStore(store SessStore) RandTokenSetting {
	return func(rt *RandToken) {
		rt.store = store
	}
}

func SetBufSize(size int) RandTokenSetting {
	return func(rt *RandToken) {
		rt.bufSize = size
	}
}

func NewRandToken(settings ...RandTokenSetting) *RandToken {
	rt := &RandToken{
		store:   defStore,
		bufSize: DefBufSize,
	}
	for _, v := range settings {
		v(rt)
	}
	rt.tokenSize = base64.StdEncoding.EncodedLen(md5.Size)

	return rt
}

func (this *RandToken) Generate() (token string, err error) {
	buf := make([]byte, this.bufSize)
	rand.Seed(time.Now().UnixNano())
	if _, err = rand.Read(buf); err != nil {
		return "", err
	} else {
		sum := md5.Sum(buf)

		return base64.StdEncoding.EncodeToString(sum[:]), nil
	}
}

func (this *RandToken) Verify(token string) bool {
	if len(token) != this.tokenSize {
		return false
	} else {
		return this.store.Have(token)
	}
}

func (this *RandToken) Value(token string) (interface{}, error) {
	return this.store.Retrieve(token)
}

func (this *RandToken) BeginSession(token string, value interface{}, expsec int64) error {
	return this.store.Store(token, value, expsec)
}

func (this *RandToken) RefreshSession(token string, expsec int64) error {
	if this.store.Have(token) {
		if value, err := this.store.Retrieve(token); err != nil {
			return err
		} else {
			this.store.Remove(token)

			return this.store.Store(token, value, expsec)
		}
	}

	return nil
}

func (this *RandToken) EndSession(token string) error {
	this.store.Remove(token)

	return nil
}
