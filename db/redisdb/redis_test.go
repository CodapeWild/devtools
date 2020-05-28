package redisdb

import (
	"devtools/file"
	"log"
	"os"
	"testing"
)

var (
	conf  = &RedisConfig{}
	rdsdb *RedisWrapper
)

func config() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := file.ReadJsonFile("./config.json", conf)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if pool, err := conf.NewPool(); err != nil {
		log.Fatalln(err.Error())
	} else {
		rdsdb = NewWrapper(pool)
	}
}

type mix struct {
	Name string `redis:"name"`
	Age  int    `redis:"age"`
}

func TestHMSet(t *testing.T) {
	config()

	_, err := rdsdb.HMSet("tnt", &mix{
		Name: "tf",
		Age:  123,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	// _, err = rdsdb.Set("name", "tnt", 0)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	if _, err = rdsdb.HMSet("tnt", map[string]interface{}{"name": "abc", "age": 321}); err != nil {
		log.Fatalln(err.Error())
	}

	data := &mix{}
	if err = rdsdb.HScanStruct("tnt", data); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(data)

	log.Println(rdsdb.Keys("*"))
}
