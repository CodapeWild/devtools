package mysqldb

import (
	"devtools/file"
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestMysqlDb(t *testing.T) {
	conf := &MysqlConfig{}
	err := file.ReadJsonFile("./config.json", conf)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db, err := conf.NewDb()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	log.Println(db.Exec(`create table tb_user_test(name varchar(30), age int);`))
}
