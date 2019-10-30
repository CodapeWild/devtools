package access

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	db_access = "db_access"
)

const (
	col_access_control = "col_access_control"
	col_group          = "col_group"
	col_in_group       = "col_in_group"
	col_relationship   = "col_relationship"
	col_auth_record    = "col_auth_record"
)

const (
	Access_Public = iota + 1
	Access_Private
	Access_Only
	Access_Group
	Access_AuthCode
	Access_Relationship
)

const (
	AccCtrl_Normal = iota + 1
	AccCtrl_Blocked
)

type MAccessControl struct {
	Id          bson.ObjectId `bson:"_id"`
	CreaterId   string        `bson:"creater_id"`
	ResourceId  string        `bson:"resource_id"`
	AccessTag   int           `bson:"access_tag"`
	AccessParam string        `bson:"access_param"`
	Created     int64         `bson:"created"`
	Updated     int64         `bson:"updated"`
	State       int           `bson:"state"`
}

const (
	Group_Allow = iota + 1
	Group_Deny
)

type MGroup struct {
	Id        bson.ObjectId `bson:"_id"`
	CreaterId string        `bson:"creater_id"`
	GroupType int           `bson:"group_type"`
	Name      string        `bson:"name"`
	Created   int64         `bson:"created"`
}

type MInGroup struct {
	GroupId     string `bson:"group_id"`
	RequestId   string `bson:"request_id"`
	RequestedId string `bson:"requested_id"`
}

const (
	Bind_By_Followed = iota + 1
	Bind_By_WhiteList
	Bind_By_BlackList
)

const (
	Relation_Requested = iota + 1
	Relation_Established
	Relation_Rejected
	Relation_Broke
)

type MRelationship struct {
	Id          bson.ObjectId `bson:"_id"`
	RequestId   string        `bson:"request_id"`
	RequestedId string        `bson:"requested_id"`
	BindBy      int           `bson:"bind_by"`
	Created     int64         `bson:"created"`
	Updated     int64         `bson:"updated"`
	State       int           `bson:"state"`
}

const (
	AuthCode_Applied = iota + 1
	AuthCode_Empowered
	AuthCode_Rejected
	AuthCode_Expired
)

type MAuthRecord struct {
	Id         bson.ObjectId `bson:"_id"`
	RequestId  string        `bson:"request_id"`
	ResourceId string        `bson:"resource_id"`
	AuthCode   string        `bson:"auth_code"`
	Created    int64         `bson:"created"`
	Updated    int64         `bson:"updated"`
	State      int           `bson:"state"`
}
