package netext

import (
	"devtools/comerr"
	"errors"
	"net"
	"net/url"
	"time"
)

var (
	SocksRespError        = errors.New("socks server respond improperly")
	SocksConnRejected     = errors.New("socks connection request rejected or failed")
	SocksConnFailed       = errors.New("socks server cannot connect to identd")
	SocksDiffUserIds      = errors.New("client program and identd report different user-ids")
	SocksUnknownError     = errors.New("socks connection report unknown error")
	Socks5SupportFailed   = errors.New("server does not support socks5")
	Socks5CannotAnonymous = errors.New("socks5 can not be anonymous")
	Socks5CannotComplete  = errors.New("can not complete socks5 connection")
)

const (
	SOCKS4  = "socks4"
	SOCKS4A = "socks4a"
	SOCKS5  = "socks5"
)

// proxy: proto://ip:port like
// socks4://127.0.0.1:1080
// socks5://user:pswd@127.0.0.1:1080
func SSProxyDialFunc(pxuri *url.URL, timeout time.Duration) func(network, addr string) (net.Conn, error) {
	if pxuri.Scheme != SOCKS4 && pxuri.Scheme != SOCKS5 && pxuri.Scheme != SOCKS4A {
		return dialError(comerr.UnrecognizedProtocol)
	}

	conn, err := net.Dial("tcp", pxuri.Host)
	if err != nil {
		return dialError(err)
	}

	return func(_, addr string) (net.Conn, error) {
		return (&socksConn{
			conn:    conn,
			scheme:  pxuri.Scheme,
			timeout: timeout,
		}).dial(addr)
	}
}

func dialError(err error) func(network, addr string) (net.Conn, error) {
	return func(string, string) (net.Conn, error) {
		return nil, err
	}
}

type socksConn struct {
	conn    net.Conn
	scheme  string
	timeout time.Duration
}

func (this *socksConn) dial(target string) (net.Conn, error) {
	switch this.scheme {
	case SOCKS4, SOCKS4A:
		return this.dialSocks4(target)
	case SOCKS5:
		return this.dialSocks5(target)
	}

	return nil, comerr.UnrecognizedProtocol
}

func (this *socksConn) query(req []byte) (resp []byte, err error) {
	if this.timeout > 0 {
		if err = this.conn.SetDeadline(time.Now().Add(this.timeout)); err != nil {
			return nil, err
		}
		defer this.conn.SetDeadline(time.Time{})
	}

	if _, err = this.conn.Write(req); err != nil {
		return nil, err
	}

	resp = make([]byte, 1024)
	c, err := this.conn.Read(resp)

	return resp[:c], err
}

func (this *socksConn) dialSocks4(target string) (net.Conn, error) {
	host, port, err := SplitHostPort(target)
	if err != nil {
		return nil, err
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, comerr.HostLookupFailed
	}
	ip := ips[0].To4()
	if len(ip) != net.IPv4len {
		return nil, comerr.ParamInvalid
	}

	req := []byte{
		4, // socks version 4
		1, // commend CONNECT
		byte(port >> 8),
		byte(port),
		ip[0], ip[1], ip[2], ip[3],
		0, // anonymous
	}
	if this.scheme == SOCKS4A {
		req = append(req, []byte(host+"\x00")...)
	}

	resp, err := this.query(req)
	if err != nil {
		return nil, err
	}
	if len(resp) != 8 {
		return nil, SocksRespError
	}

	switch resp[1] {
	case 90:
		return this.conn, nil
	case 91:
		return nil, SocksConnRejected
	case 92:
		return nil, SocksConnFailed
	case 93:
		return nil, SocksDiffUserIds
	default:
		return nil, SocksUnknownError
	}
}

func (this *socksConn) dialSocks5(target string) (net.Conn, error) {
	req := []byte{
		5, // socks version 5
		1, // command CONNECT
		0, // anonyouse
	}
	resp, err := this.query(req)
	if err != nil {
		return nil, err
	}
	if len(resp) != 2 {
		return nil, SocksRespError
	}
	if resp[0] != 5 {
		return nil, Socks5SupportFailed
	}
	if resp[1] != 0 {
		return nil, Socks5CannotAnonymous
	}

	host, port, err := SplitHostPort(target)
	if err != nil {
		return nil, err
	}

	req = append(req, []byte{
		3,               // address type 3 means domain name
		byte(len(host)), // domain lenght
	}...)
	req = append(req, []byte(host)...)
	req = append(req, []byte{
		byte(port >> 8),
		byte(port),
	}...)
	resp, err = this.query(req)
	if err != nil {
		return nil, err
	}
	if len(resp) != 10 {
		return nil, SocksRespError
	}
	if resp[1] != 0 {
		return nil, Socks5CannotComplete
	}

	return this.conn, nil
}
