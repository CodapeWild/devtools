package mongodb

import (
	"devtools/comerr"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoWrapper struct {
	sess *mgo.Session
	Db   string
}

func NewWrapper(db string, sess *mgo.Session) *MgoWrapper {
	return &MgoWrapper{
		sess: sess,
		Db:   db,
	}
}

func (this *MgoWrapper) ConnDB() *mgo.Database {
	return this.sess.Clone().DB(this.Db)
}

func (this *MgoWrapper) Count(col string, query interface{}) (int, error) {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Find(query).Count()
}

func (this *MgoWrapper) Insert(col string, doc ...interface{}) error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Insert(doc...)
}

func (this *MgoWrapper) FindById(col string, id interface{}, rslt interface{}) error {
	if rslt == nil {
		return comerr.ParamInvalid
	}
	if reflect.TypeOf(rslt).Kind() != reflect.Ptr {
		return comerr.ParamTypeInvalid
	}

	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).FindId(id).One(rslt)
}

func (this *MgoWrapper) FindOne(col string, query interface{}, rslt interface{}) error {
	if rslt == nil {
		return comerr.ParamInvalid
	}
	if reflect.TypeOf(rslt).Kind() != reflect.Ptr {
		return comerr.ParamTypeInvalid
	}

	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Find(query).One(rslt)
}

func (this *MgoWrapper) FindAll(col string, query interface{}, skip, limit int, rslt interface{}) error {
	if t := reflect.TypeOf(rslt); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		return comerr.ParamTypeInvalid
	}

	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Find(query).Skip(skip).Limit(limit).All(rslt)
}

func (this *MgoWrapper) FindTopOne(col string, query interface{}, orderBy string, rslt interface{}) error {
	if rslt == nil {
		return comerr.ParamInvalid
	}
	if reflect.TypeOf(rslt).Kind() != reflect.Ptr {
		return comerr.ParamTypeInvalid
	}

	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Find(query).Sort(orderBy).One(rslt)
}

func (this *MgoWrapper) FindAndSort(col string, query interface{}, orderBy string, skip, limit int, rslt interface{}) error {
	if t := reflect.TypeOf(rslt); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		return comerr.ParamTypeInvalid
	}

	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Find(query).Sort(orderBy).Skip(skip).Limit(limit).All(rslt)
}

func (this *MgoWrapper) UpdateOne(col string, query interface{}, update interface{}) error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Update(query, update)
}

func (this *MgoWrapper) UpSetOne(col string, query interface{}, set interface{}) error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Update(query, bson.M{"$set": set})
}

func (this *MgoWrapper) UpdateAll(col string, query interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).UpdateAll(query, update)
}

func (this *MgoWrapper) UpSert(col string, query interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Upsert(query, update)
}

func (this *MgoWrapper) CreateIndex(col string, index *mgo.Index) error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).EnsureIndex(*index)
}

func (this *MgoWrapper) DropIndex(col string, keys ...string) error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).DropIndex(keys...)
}

func (this *MgoWrapper) RemoveOne(col string, query interface{}) error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).Remove(query)
}

func (this *MgoWrapper) RemovAll(col string, query interface{}) (*mgo.ChangeInfo, error) {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).RemoveAll(query)
}

func (this *MgoWrapper) DropCol(col string) error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.C(col).DropCollection()
}

func (this *MgoWrapper) DropDb() error {
	db := this.ConnDB()
	defer db.Session.Close()

	return db.DropDatabase()
}
