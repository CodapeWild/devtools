package netext

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
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
	req, err := http.NewRequest(http.MethodGet, "https://www.mm666.club", nil)
	if err != nil {
		log.Panicln(err.Error())
	}
	pxuri, err := url.Parse("socks4://127.0.0.1:1080")
	if err != nil {
		log.Panicln(err.Error())
	}
	clnt := http.Client{Transport: &http.Transport{Dial: SSProxyDialFunc(pxuri, 60*time.Second)}}
	// clnt := http.Client{Transport: &http.Transport{Dial: SSProxyDialFunc("socks5://127.0.0.1:1080", 60*time.Second)}}

	resp, err := clnt.Do(req)
	if err != nil {
		log.Panicln(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		log.Panicln(resp.Status)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Println(string(buf))
}
