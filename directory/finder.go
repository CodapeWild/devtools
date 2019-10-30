package directory

import (
	"devtools/comerr"
	"devtools/db/mongodb"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Finder interface {
	Set(*MFile) error
	Get(string) (*MFile, error)
}

const (
	db_directory = "db_directory"
)

const (
	col_file = "col_file"
)

const (
	File_Normal = iota + 1
	File_Blocked
)

type MFile struct {
	Id      bson.ObjectId `bson:"_id"`
	Path    string        `bson:"path"`
	Media   MediaType     `bson:"media"`
	Size    int64         `bson:"size"`
	Span    int64         `bson:"span"`
	Created int64         `bson:"created"`
	State   int           `bson:"state"`
}

type MgoFinder struct {
	mgoWrapper *mongodb.MgoWrapper
}

func NewMgoFinder(mgoSess *mgo.Session) *MgoFinder {
	return &MgoFinder{mgoWrapper: mongodb.NewWrapper(db_directory, mgoSess)}
}

func (this *MgoFinder) Set(mf *MFile) error {
	return this.mgoWrapper.Insert(col_file, mf)
}

func (this *MgoFinder) Get(id string) (*MFile, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, comerr.ParamInvalid
	}

	mf := &MFile{}
	if err := this.mgoWrapper.FindOne(col_file, bson.M{"_id": bson.ObjectIdHex(id)}, mf); err != nil {
		return nil, err
	}

	return mf, nil
}
