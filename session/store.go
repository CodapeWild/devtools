package session

type SessStore interface {
	Store(token string, value interface{}, expsec int64) (err error)
	Retrieve(token string) (value interface{}, err error)
	Have(token string) bool
	Remove(token string)
}

var defStore = NewMapStore()
