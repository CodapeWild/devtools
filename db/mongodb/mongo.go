package mongodb

import (
	"gopkg.in/mgo.v2"
)

type MgoConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pswd string `json:"pswd"`
	Db   string `json:"db"`
}

// mongodb://user:pswd@host:port/db
func (this *MgoConfig) NewSession() (*mgo.Session, error) {
	connStr := "mongodb://"
	if this.User != "" && this.Pswd != "" {
		connStr += this.User + ":" + this.Pswd + "@"
	}
	if this.Host == "" {
		this.Host = "127.0.0.1"
	}
	if this.Port == "" {
		this.Port = "27017"
	}
	connStr += this.Host + ":" + this.Port
	if this.Db != "" {
		connStr += "/" + this.Db
	}

	return mgo.Dial(connStr)
}
