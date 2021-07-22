package mongodb

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/CodapeWild/devtools/tlsext"
	"gopkg.in/mgo.v2"
)

var (
	defHost = "127.0.0.1"
	defPort = "27017"
)

type MgoConfig struct {
	Host      string                  `json:"host" toml:"host"`
	Port      string                  `json:"port" toml:"port"`
	User      string                  `json:"user" toml:"user"`
	Pswd      string                  `json:"pswd" toml:"pswd"`
	Db        string                  `json:"db" toml:"db"`
	EnableTls bool                    `json:"enable_tls" toml:"enable_tls"`
	TlsConf   *tlsext.TlsClientConfig `json:"tls_conf" toml:"tls_conf"`
}

// mongodb://user:pswd@host:port/db
func (this *MgoConfig) NewSession() (*mgo.Session, error) {
	connStr := "mongodb://"
	if this.User != "" && this.Pswd != "" {
		connStr += this.User + ":" + this.Pswd + "@"
	}
	if this.Host == "" {
		this.Host = defHost
	}
	if this.Port == "" {
		this.Port = defPort
	}
	connStr += this.Host + ":" + this.Port
	if this.Db != "" {
		connStr += "/" + this.Db
	}

	dialInfo, err := mgo.ParseURL(connStr)
	if err != nil {
		return nil, err
	}

	if this.EnableTls && this.TlsConf != nil {
		if tlsConf, err := this.TlsConf.TlsConfig(); err != nil {
			return nil, err
		} else {
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				return tls.Dial("tcp", addr.String(), tlsConf)
			}
		}
	}

	dialInfo.Direct = true
	dialInfo.Timeout = 3 * time.Second

	return mgo.DialWithInfo(dialInfo)
}
