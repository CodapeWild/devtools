package redisdb

import (
	"devtools/comerr"
	"reflect"

	"github.com/garyburd/redigo/redis"
)

type RedisWrapper struct {
	pool *redis.Pool
}

func NewWrapper(pool *redis.Pool) *RedisWrapper {
	return &RedisWrapper{pool: pool}
}

func (this *RedisWrapper) Session() redis.Conn {
	return this.pool.Get()
}

func (this *RedisWrapper) Expire(key string, expsec int64) error {
	conn := this.Session()
	defer conn.Close()

	_, err := conn.Do("expire", key, expsec)

	return err
}

/*
	sec==-2 expired
	sec==-1 never
	sec>=0 seconds before expired
*/
func (this *RedisWrapper) TTL(key string) (sec int64, err error) {
	conn := this.Session()
	defer conn.Close()

	return redis.Int64(conn.Do("ttl", key))
}

func (this *RedisWrapper) IsKeyExists(key string) bool {
	if key == "" {
		return false
	}

	conn := this.Session()
	defer conn.Close()

	rply, err := conn.Do("get", key)

	return rply != nil && err == nil
}

func (this *RedisWrapper) Set(key string, value interface{}, expsec int64) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	if expsec > 0 {
		return conn.Do("set", key, value, "ex", expsec)
	} else {
		return conn.Do("set", key, value)
	}
}

func (this *RedisWrapper) Get(key string) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	return conn.Do("get", key)
}

func (this *RedisWrapper) MSet(compound interface{}) (interface{}, error) {
	if compound == nil {
		return nil, comerr.ParamInvalid
	}
	var (
		t    = reflect.TypeOf(compound)
		k    = t.Kind()
		args = redis.Args{}
	)
	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	if k != reflect.Struct && k != reflect.Map {
		return nil, comerr.ParamTypeInvalid
	} else {
		args = args.AddFlat(compound)
	}

	conn := this.Session()
	defer conn.Close()

	return conn.Do("mset", args...)
}

func (this *RedisWrapper) MGet(keys ...string) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	return conn.Do("mget", redis.Args{}.AddFlat(keys)...)
}

func (this *RedisWrapper) DelKey(key string) error {
	conn := this.Session()
	defer conn.Close()

	_, err := conn.Do("del", key)

	return err
}

func (this *RedisWrapper) IsFieldExists(key string, field interface{}) bool {
	conn := this.Session()
	defer conn.Close()

	rply, err := conn.Do("hget", key, field)

	return rply != nil && err == nil
}

func (this *RedisWrapper) HSet(key string, field, value interface{}) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	rply, err := conn.Do("hset", key, field, value)

	return rply, err
}

func (this *RedisWrapper) HGet(key string, field interface{}) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	if field != nil {
		return conn.Do("hget", key, field)
	} else {
		return conn.Do("hgetall", key)
	}
}

func (this *RedisWrapper) HMSet(key string, compound interface{}) (interface{}, error) {
	if compound == nil {
		return nil, comerr.ParamInvalid
	}
	var (
		t    = reflect.TypeOf(compound)
		k    = t.Kind()
		args = redis.Args{}
	)
	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	if k != reflect.Map && k != reflect.Struct {
		return nil, comerr.ParamTypeInvalid
	} else {
		args = args.Add(key).AddFlat(compound)
	}

	conn := this.Session()
	defer conn.Close()

	return conn.Do("hmset", args...)
}

func (this *RedisWrapper) HMGet(key string, fields ...interface{}) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	if fields[0] != nil {
		return conn.Do("hmget", redis.Args{}.Add(key).AddFlat(fields)...)
	} else {
		return conn.Do("hgetall", key)
	}
}

func (this *RedisWrapper) HMScanStruct(key string, dst interface{}) error {
	if dst == nil {
		return comerr.ParamInvalid
	}
	if t := reflect.TypeOf(dst); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return comerr.ParamTypeInvalid
	}

	conn := this.Session()
	defer conn.Close()

	rply, err := conn.Do("hgetall", key)
	if err != nil {
		return err
	}
	src := rply.([]interface{})
	if len(src) == 0 {
		return comerr.NoEntryFound
	}

	return redis.ScanStruct(src, dst)
}

func (this *RedisWrapper) HDel(key string, fields ...interface{}) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	if fields[0] != nil {
		return conn.Do("hdel", redis.Args{}.Add(key).AddFlat(fields))
	} else {
		return conn.Do("del", key)
	}
}

func (this *RedisWrapper) LPush(key string, values ...interface{}) (interface{}, error) {
	return this.push("lpush", key, values...)
}

func (this *RedisWrapper) RPush(key string, values ...interface{}) (interface{}, error) {
	return this.push("rpush", key, values...)
}

func (this *RedisWrapper) LPop(key string) (interface{}, error) {
	return this.pop("lpop", key)
}

func (this *RedisWrapper) RPop(key string) (interface{}, error) {
	return this.pop("rpop", key)
}

func (this *RedisWrapper) RPopLPush(src, dst string) (interface{}, error) {
	conn := this.Session()
	defer conn.Close()

	return conn.Do("rpoplpush", src, dst)
}

func (this *RedisWrapper) BLPop(key string, timeout int) (interface{}, error) {
	if timeout < 0 {
		return nil, comerr.ParamInvalid
	}

	return this.bpop("blpop", key, timeout)
}

func (this *RedisWrapper) BRPop(key string, timeout int) (interface{}, error) {
	if timeout < 0 {
		return nil, comerr.ParamInvalid
	}

	return this.bpop("brpop", key, timeout)
}

func (this *RedisWrapper) BRPopLPush(src, dst string, timeout int) (interface{}, error) {
	if timeout < 0 {
		return nil, comerr.ParamInvalid
	}

	conn := this.Session()
	defer conn.Close()

	return conn.Do("brpoplpush", src, dst, timeout)
}

func (this *RedisWrapper) LRemove(key string, count int, value interface{}) error {
	conn := this.Session()
	defer conn.Close()

	_, err := conn.Do("lrem", key, count, value)

	return err
}

func (this *RedisWrapper) push(cmd string, key string, values ...interface{}) (interface{}, error) {
	if cmd != "lpush" && cmd != "rpush" {
		return nil, comerr.ParamInvalid
	}

	conn := this.Session()
	defer conn.Close()

	return conn.Do(cmd, redis.Args{}.Add(key).AddFlat(values)...)
}

func (this *RedisWrapper) pop(cmd string, key string) (interface{}, error) {
	if cmd != "lpop" && cmd != "rpop" {
		return nil, comerr.ParamInvalid
	}

	conn := this.Session()
	defer conn.Close()

	return conn.Do(cmd, key)
}

func (this *RedisWrapper) bpop(cmd string, key string, timeout int) (interface{}, error) {
	if cmd != "blpop" && cmd != "brpop" {
		return nil, comerr.ParamInvalid
	}

	conn := this.Session()
	defer conn.Close()

	if timeout < 0 {
		timeout = 0
	}

	return conn.Do(cmd, key, timeout)
}
