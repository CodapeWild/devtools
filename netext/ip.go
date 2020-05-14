package netext

import (
	"bytes"
	"net"
	"strconv"
	"strings"
)

type IpRange struct {
	Start, End net.IP
}

func NewIpRange(start, end net.IP) *IpRange {
	if start == nil || end == nil {
		return nil
	}

	return &IpRange{
		Start: start,
		End:   end,
	}
}

var privateRanges = []*IpRange{
	NewIpRange(net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")),
	NewIpRange(net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")),
	NewIpRange(net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")),
	NewIpRange(net.ParseIP("198.18.0.0"), net.ParseIP("198.19.255.255")),
	NewIpRange(net.ParseIP("192.0.0.0"), net.ParseIP("192.0.0.255")),
	NewIpRange(net.ParseIP("100.64.0.0"), net.ParseIP("100.127.255.255")),
}

func IsPrivate(ip net.IP) bool {
	if ip == nil {
		return false
	}

	for _, r := range privateRanges {
		if bytes.Compare(ip, r.Start) >= 0 && bytes.Compare(ip, r.End) <= 0 {
			return true
		}
	}

	return false
}

func Ip4ToInt(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}

	var ipInt uint32
	for _, b := range ip[12:] {
		ipInt = ipInt<<8 + uint32(b)
	}

	return ipInt
}

func SplitHostPort(addr string) (host string, port int, err error) {
	i := strings.Index(addr, "://")
	if i > 0 {
		addr = addr[i+3:]
	}

	var portstr string
	host, portstr, err = net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	port, err = strconv.Atoi(portstr)

	return
}
