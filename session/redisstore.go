package session

import "github.com/CodapeWild/devtools/db/redisdb"

const rds_sess_prefix = "rds_sess_"

type RedisStore struct{ *redisdb.RedisWrapper }

func NewRedisStore(rds *redisdb.RedisWrapper) *RedisStore {
	return &RedisStore{rds}
}

func (this *RedisStore) prefix(token string) string {
	return rds_sess_prefix + token
}

func (this *RedisStore) Store(token string, value interface{}, expsec int64) (err error) {
	_, err = this.Set(this.prefix(token), value, expsec)

	return
}

func (this *RedisStore) Retrieve(token string) (value interface{}, err error) {
	return this.Get(this.prefix(token))
}

func (this *RedisStore) Have(token string) bool {
	return this.IsKeyExists(this.prefix(token))
}

func (this *RedisStore) Remove(token string) {
	this.DelKey(this.prefix(token))
}
