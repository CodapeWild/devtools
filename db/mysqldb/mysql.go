package mysqldb

import (
	"database/sql"
	"devtools/comerr"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver_ = "mysql"
)

type MysqlConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pswd string `json:"pswd"`
	Db   string `json:"db"`
}

// user:pswd@tcp(host:port)/db
func (this *MysqlConfig) NewDb() (*sql.DB, error) {
	if this.Db == "" {
		return nil, comerr.ParamInvalid
	}

	if this.Host == "" {
		this.Host = "localhost"
	}
	if this.Port == "" {
		this.Port = "3306"
	}
	if this.User == "" {
		this.User = "root"
	}

	db, err := sql.Open(driver_, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", this.User, this.Pswd, this.Host, this.Port, this.Db))
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
