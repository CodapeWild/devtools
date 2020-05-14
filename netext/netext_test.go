package netext

import (
	"log"
	"net"
	"testing"
)

func TestIp(t *testing.T) {
	log.Println(IsPrivate(net.ParseIP("192.168.1.100")))
	log.Println(IsPrivate(net.ParseIP("172.16.1.100")))
	log.Println(IsPrivate(net.ParseIP("10.18.3.6")))
	log.Println(IsPrivate(net.ParseIP("198.18.1.100")))
	log.Println(IsPrivate(net.ParseIP("192.0.0.100")))
	log.Println(IsPrivate(net.ParseIP("100.64.1.100")))
}

func TestSocks(t *testing.T) {

}
