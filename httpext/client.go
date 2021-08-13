package httpext

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/CodapeWild/devtools/msgque"
)

var (
	getProxyUrlOnce sync.Once
	getProxyURL     func(*http.Request) (*url.URL, error)
	defTransport    *http.Transport = &http.Transport{
		Proxy: getProxyURL,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
	maxConcurrency uint = 10
	newTicketOnce  sync.Once
	tickets        msgque.TicketQueue
)

func UseProxy(addr string) {
	getProxyUrlOnce.Do(func() {
		if pxurl, err := url.ParseRequestURI(addr); err != nil {
			log.Printf("parse proxy address failed err:%s\n", err.Error())
		} else {
			getProxyURL = func(r *http.Request) (*url.URL, error) {
				return pxurl, nil
			}
		}
	})
}

func SetMaxConcurrency(max uint) {
	maxConcurrency = max
}

func SendRequest(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	newTicketOnce.Do(func() {
		tickets = msgque.NewSimpleTicketQueue(int(maxConcurrency))
	})
}
