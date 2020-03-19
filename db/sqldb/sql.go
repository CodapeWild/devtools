package sqldb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MYSQL   = "mysql"
	SQLITE3 = "sqlite3"
)

func NewDb(driver string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
