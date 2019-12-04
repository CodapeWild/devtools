package crawler

import (
	"devtools/comerr"
	"devtools/structure"
	"net"
	"time"
)

type Scheme byte

const (
	Http    Scheme = 1
	Https   Scheme = 2
	Socks4  Scheme = 4
	Socks4A Scheme = 8
	Socks5  Scheme = 16
)

type Anonymity byte

const (
	None   Anonymity = 0
	Medium Anonymity = 1
	Hight  Anonymity = 2
)

type Proxy struct {
	Schemes   Scheme
	Host      string
	Port      int
	User      string
	Pswd      string
	Anon      Anonymity
	Country   string
	City      string
	LastCheck time.Time
}

func (this *Proxy) IsValid() bool {
	return false
}

func (this *Proxy) Dial() net.Conn {
}

const (
	DefProxyPoolMax = 15
)

type ProxyPool struct {
	src     []string
	proxies *structure.LinkList
	max     int
}

func NewProxyPool(max int, proxy bool, src ...string) (*ProxyPool, error) {
	c := len(src)
	if c == 0 {
		return nil, comerr.ParamInvalid
	}

	pp := &ProxyPool{
		src:     make([]string, c),
		proxies: structure.NewLinkList(),
		max:     max,
	}
	if max <= 0 {
		pp.max = DefProxyPoolMax
	}

	return pp, nil
}
