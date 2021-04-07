package session

import "errors"

var (
	GenTokenFailed = errors.New("generate session token failed")
	InvalidToken   = errors.New("invalid token")
)

type SessToken interface {
	Generate(token string, err error)
	Verify(token string) bool
	Value(token string) (value interface{}, err error)

	BeginSession(token string, value interface{}, expsec int64) (err error)
	RefreshSession(token string, expsec int64) (err error)
	EndSession(token string) (err error)
}
