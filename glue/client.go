package glue

import (
	"devtools/comerr"
	"devtools/db/redisdb"
	"encoding/json"
	"errors"
	"log"
)

var (
	readRespTimeout = errors.New("read response from pip timeout")
)

type Client struct {
	rdsWrapper *redisdb.RedisWrapper
}

func NewClient(conf *redisdb.RedisConfig) (*Client, error) {
	pool, err := conf.NewPool()
	if err != nil {
		return nil, err
	}

	return &Client{rdsWrapper: redisdb.NewWrapper(pool)}, nil
}

func (this *Client) Do(req *Request) error {
	if req == nil {
		return comerr.ParamInvalid
	}
	if req.err != nil {
		return req.err
	}

	serverKey, err := req.GetServerKey()
	if err != nil {
		return err
	}

	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// send request
	if _, err = this.rdsWrapper.RPush(serverKey, buf); err != nil {
		return err
	}
	// get response
	if req.IsCallback != 0 {
		go func(req *Request) {
			clientKey, err := req.GetClientKey()
			if err != nil {
				log.Println(err.Error())

				return
			}

			// clear client pip
			this.rdsWrapper.DelKey(clientKey)

			rply, err := this.rdsWrapper.BLPop(clientKey, req.GetTimeout())
			if err != nil {
				log.Println(err.Error())

				return
			}
			if rply == nil {
				log.Println(readRespTimeout.Error())

				return
			}

			buf := rply.([]interface{})[1].([]byte)
			resp := &Response{}
			if err = json.Unmarshal(buf, resp); err != nil {
				log.Println(err.Error())

				return
			}

			go req.callback(resp)
		}(req)
	}

	return nil
}
