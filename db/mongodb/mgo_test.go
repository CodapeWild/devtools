package mongodb

import (
	"devtools/file"
	"log"
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type muser struct {
	Id    bson.ObjectId `bson:"_id"`
	Name  string        `bson:"name"`
	Age   int           `bson:"age"`
	Class string        `bson:"class"`
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestMgoDb(t *testing.T) {
	conf := &MgoConfig{}
	err := file.ReadJsonFile("./config.json", conf)
	if err != nil {
		log.Fatalln(err.Error())
	}
	sess, err := conf.NewSession()
	if err != nil {
		log.Fatalln(err.Error())
	}

	dbFoo := NewWrapper(sess, "foo")

	if err = dbFoo.Insert("user", bson.M{"name": "tnt", "age": 123}); err != nil {
		log.Fatalln(err.Error())
	}
	if err = dbFoo.Insert("user", &muser{bson.NewObjectId(), "tf", 225, "1/2"}); err != nil {
		log.Fatalln(err.Error())
	}
	if err = dbFoo.UpdateOne("user", bson.M{"name": "tnt"}, bson.M{"$set": bson.M{"age": 333, "sex": 1}}); err != nil {
		log.Fatalln(err.Error())
	}
	u := &muser{}
	if err = dbFoo.FindOne("user", bson.M{"name": "tnt"}, u); err != nil {
		log.Println(err.Error())
	}
	log.Println(*u)
	if err = dbFoo.FindOne("user", bson.M{"name": "tnt"}, u); err != nil {
		log.Fatalln("@@@", err.Error())
	}
	log.Println(u.Id.String())
	log.Println(u.Id.Hex())
	log.Println(dbFoo.Count("user", bson.M{"ff": "xxoo"}))
}
