package access

import (
	"log"
	"os"
	"testing"

	"github.com/CodapeWild/devtools/db/mongodb"
	"gopkg.in/mgo.v2/bson"
)

var (
	accCtrlr     *AccessController
	requestId    = bson.NewObjectId().Hex()
	requestedIds = []string{bson.NewObjectId().Hex(), bson.NewObjectId().Hex(), bson.NewObjectId().Hex()}
	gid          string
	resourceIds  = []string{bson.NewObjectId().Hex(), bson.NewObjectId().Hex(), bson.NewObjectId().Hex()}
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mgoSess, err := (&mongodb.MgoConfig{}).NewSession()
	if err != nil {
		log.Fatalln(err.Error())
	}
	accCtrlr = NewAccessController(mgoSess)
}

func TestOnly(t *testing.T) {
	if err := accCtrlr.CreateAccessControl(requestedIds[2], Access_Only, requestId, resourceIds[1]); err != nil {
		log.Panicln(err.Error())
	}
}

func TestRelationship(t *testing.T) {
	err := accCtrlr.Follow(requestId, requestedIds[0])
	if err != nil {
		log.Panicln(err.Error())
	}
	if err = accCtrlr.Follow(requestId, requestedIds[0]); err != nil {
		log.Println(err.Error())
	}

	if err = accCtrlr.WhiteListBindRequest(requestId, requestedIds[0]); err != nil {
		log.Panicln(err.Error())
	}
	if err = accCtrlr.UpdateRelationshipState(requestId, requestedIds[0], Bind_By_WhiteList, Relation_Requested, Relation_Established); err != nil {
		log.Panicln()
	}
	if err = accCtrlr.WhiteListBindRequest(requestedIds[0], requestId); err != nil {
		log.Println(err.Error())
	}

	if err = accCtrlr.KickInBlackList(requestId, requestedIds[0]); err != nil {
		log.Panicln(err.Error())
	}
	if err = accCtrlr.KickInBlackList(requestId, requestedIds[0]); err != nil {
		log.Println(err.Error())
	}
}

func TestGroup(t *testing.T) {
	var err error
	gid, err = accCtrlr.CreateGroup(requestId, Group_Allow, "friends", requestedIds...)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("group id:", gid)

	if err = accCtrlr.RemoveFromGroup(gid, requestedIds[0], requestedIds[2]); err != nil {
		log.Panicln(err.Error())
	}
}

func TestAccessControl(t *testing.T) {
	err := accCtrlr.CreateAccessControl(requestId, Access_Group, gid, resourceIds...)
	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println(accCtrlr.Access(requestedIds[2], resourceIds[1], ""))
	log.Println(accCtrlr.Access(requestedIds[1], resourceIds[1], ""))
	log.Println(accCtrlr.Access(requestId, resourceIds[2], ""))
	log.Println(accCtrlr.Access(requestedIds[0], resourceIds[2], ""))
	log.Println(accCtrlr.Access(requestedIds[1], resourceIds[2], ""))
	log.Println(accCtrlr.Access(requestedIds[2], resourceIds[2], ""))
}
