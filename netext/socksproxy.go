package netext

import (
	"devtools/comerr"
	"net"
	"net/url"
)

const (
	SOCKS4  = "socks4"
	SOCKS4A = "socks4a"
	SOCKS5  = "socks5"
)

// proxy: proto://ip:port like
// socks4://127.0.0.1:1080
// socks5://user:pswd@127.0.0.1:1080
func SSProxyDialFunc(proxy string) func(network, address string) (net.Conn, error) {
	pxuri, err := url.Parse(proxy)
	if err != nil {
		return dialError(err)
	}

	if pxuri.Scheme != SOCKS4 && pxuri.Scheme != SOCKS5 && pxuri.Scheme != SOCKS4A {
		return dialError(comerr.UnrecognizedProtocol)
	}

	conn, err := net.Dial("tcp", pxuri.Host)
	if err != nil {
		return dialError(err)
	}

	switch pxuri.Scheme {
	case SOCKS4, SOCKS4A:
		return func(_, address string) (net.Conn, error) {
			return dialSocks4(conn, address)
		}
	case SOCKS5:
		return func(_, address string) (net.Conn, error) {
			return dialSocks5(conn, address)
		}
	default:
		return dialError(comerr.UnrecognizedProtocol)
	}
}

func dialSocks4(conn net.Conn, target string) (net.Conn, error) {
	return nil, nil
}

func dialSocks5(conn net.Conn, target string) (net.Conn, error) {
	return nil, nil
}

func dialError(err error) func(network, address string) (net.Conn, error) {
	return func(string, string) (net.Conn, error) {
		return nil, err
	}
}
