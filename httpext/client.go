package httpext

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	defTransport *http.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          300,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)

type ClntOption func(clnt *http.Client)

func WithTimeout(timeout time.Duration) ClntOption {
	return func(clnt *http.Client) {
		clnt.Timeout = timeout
	}
}

func WithCookies(u *url.URL, cookies []*http.Cookie) ClntOption {
	return func(clnt *http.Client) {
		clnt.Jar.SetCookies(u, cookies)
	}
}

func WithTransport(transport *http.Transport) ClntOption {
	return func(clnt *http.Client) {
		clnt.Transport = transport
	}
}

func NewTransportWithProxy(proxy *url.URL) *http.Transport {
	transport := defTransport.Clone()
	transport.Proxy = http.ProxyURL(proxy)

	return transport
}

func SendRequest(req *http.Request, opts ...ClntOption) (*http.Response, error) {
	clnt := &http.Client{Transport: defTransport}
	for _, opt := range opts {
		opt(clnt)
	}

	return clnt.Do(req)
}
