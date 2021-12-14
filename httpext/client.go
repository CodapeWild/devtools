package httpext

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	defTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)

type ClientOption func(clnt *http.Client)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(clnt *http.Client) {
		clnt.Timeout = timeout
	}
}

func WithCookies(u *url.URL, cookies []*http.Cookie) ClientOption {
	return func(clnt *http.Client) {
		clnt.Jar.SetCookies(u, cookies)
	}
}

func WithTransport(transport *http.Transport) ClientOption {
	return func(clnt *http.Client) {
		clnt.Transport = transport
	}
}

func WithInsecureSkipVerify(skip bool) ClientOption {
	return func(clnt *http.Client) {
		trans, ok := clnt.Transport.(*http.Transport)
		if ok && trans.TLSClientConfig.InsecureSkipVerify != skip {
			trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: skip}
			clnt.Transport = trans
		}
	}
}

func SendRequest(req *http.Request, opts ...ClientOption) (*http.Response, error) {
	clnt := &http.Client{Transport: defTransport}
	for _, opt := range opts {
		opt(clnt)
	}

	return clnt.Do(req)
}
