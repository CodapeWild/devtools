package sqldb

import (
	"log"
	"testing"
)

func TestNewDb(t *testing.T) {
	log.Println(NewDb(SQLITE3, "./testdb.db"))
}
