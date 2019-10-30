package redisdb

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	Pswd         string        `json:"pswd"`
	MaxIdle      int           `json:"max_idl"`
	MaxActive    int           `json:"max_active"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
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
			conn, err := redis.Dial("tcp", addr, redis.DialConnectTimeout(3*time.Second), redis.DialReadTimeout(this.ReadTimeout*time.Second), redis.DialWriteTimeout(this.WriteTimeout*time.Second))
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

	if _, err := pool.Get().Do("ping"); err != nil {
		return nil, err
	} else {
		return pool, nil
	}
}
