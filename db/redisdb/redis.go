package redisdb

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	Pswd            string `json:"pswd"`
	MaxIdle         int    `json:"max_idl"`
	MaxActive       int    `json:"max_active"`
	DialTimeoutSec  int    `json:"dial_timeout_sec"`
	ReadTimeoutSec  int    `json:"read_timeout_sec"`
	WriteTimeoutSec int    `json:"write_timeout_sec"`
}

func (this *RedisConfig) NewPool() (*redis.Pool, error) {
	var addr string
	if this.Host == "" || this.Port == "" {
		addr = ":6379"
	} else {
		addr = this.Host + ":" + this.Port
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr, redis.DialConnectTimeout(time.Duration(this.DialTimeoutSec)*time.Second), redis.DialReadTimeout(time.Duration(this.ReadTimeoutSec)*time.Second), redis.DialWriteTimeout(time.Duration(this.WriteTimeoutSec)*time.Second))
			if err != nil {
				return nil, err
			}

			if this.Pswd != "" {
				if _, err = conn.Do("auth", this.Pswd); err != nil {
					return nil, err
				}
			}

			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := conn.Do("ping")

			return err
		},
		MaxIdle:     this.MaxIdle,
		MaxActive:   this.MaxActive,
		IdleTimeout: time.Millisecond,
		Wait:        true,
	}

	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("ping"); err != nil {
		return nil, err
	} else {
		return pool, nil
	}
}
