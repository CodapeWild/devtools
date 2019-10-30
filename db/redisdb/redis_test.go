package redisdb

import (
	"devtools/file"
	"log"
	"os"
	"testing"
)

var (
	conf       = &RedisConfig{}
	rdsWrapper *RedisWrapper
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := file.ReadJsonFile("./config.json", conf)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if pool, err := conf.NewPool(); err != nil {
		log.Fatalln(err.Error())
	} else {
		rdsWrapper = NewWrapper(pool)
	}
}

type mix struct {
	Name  string
	Age   int
	Title map[string]string
}

func TestHMSet(t *testing.T) {
	log.Println(rdsWrapper.HMSet("tnt", &mix{
		Name: "tf",
		Age:  123,
		Title: map[string]string{
			"dskjfks": "djksf",
			"dsjff":   "idsjfioweu",
		},
	}))
	log.Println(rdsWrapper.HMSet("tnt", map[interface{}]interface{}{"Age": 666}))
	log.Println(rdsWrapper.HMSet("tnt", map[string]string{
		"dskjfks": "djksf",
		"dsjff":   "idsjfioweu",
	}))
}
