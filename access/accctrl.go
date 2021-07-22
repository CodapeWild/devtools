package access

import (
	"errors"
	"log"
	"time"

	"github.com/CodapeWild/devtools/comerr"
	"github.com/CodapeWild/devtools/db/mongodb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	RelationshipDuplicateBind  = errors.New("relationship duplicate bind")
	RelationshipNotFound       = errors.New("relationship not found")
	AccessControlNotFound      = errors.New("can not find resource access control")
	AccessControlAlreadyExists = errors.New("access control already exists")
	GroupNotFound              = errors.New("group not found")
	GroupAlreadyExists         = errors.New("group already exists")
)

type AccessController struct {
	mgoWrapper *mongodb.MgoWrapper
}

func NewAccessController(sess *mgo.Session) *AccessController {
	return &AccessController{mgoWrapper: mongodb.NewWrapper(sess, db_access)}
}

func (this *AccessController) FindRelationships(requestId string, bind int, state int, requestedIds ...string) ([]*MRelationship, error) {
	if requestId == "" || bind > Bind_By_BlackList || state > Relation_Broke {
		return nil, comerr.ErrParamInvalid
	}

	query := make(bson.M)
	if bind > 0 {
		query["bind_by"] = bind
		switch bind {
		case Bind_By_Followed, Bind_By_BlackList:
			query["request_id"] = requestId
			if len(requestedIds) != 0 {
				query["requested_id"] = bson.M{"$in": requestedIds}
			}
		case Bind_By_WhiteList:
			if l := len(requestedIds); l == 0 {
				query["$or"] = []bson.M{bson.M{"request_id": requestId}, bson.M{"requested_id": requestId}}
			} else if l == 1 {
				query["$or"] = []bson.M{bson.M{"request_id": requestId, "requested_id": requestedIds[0]}, bson.M{"request_id": requestedIds[0], "requested_id": requestId}}
			} else if l > 1 {
				query["$or"] = []bson.M{bson.M{"request_id": requestId, "requested_id": bson.M{"$in": requestedIds}}, bson.M{"request_id": bson.M{"$in": requestedIds}, "requested_id": requestId}}
			}
		}
	}
	if state > 0 {
		query["state"] = state
	}

	var mrelas []*MRelationship
	if err := this.mgoWrapper.FindAll(col_relationship, query, 0, 0, &mrelas); err != nil {
		return nil, err
	} else {
		return mrelas, nil
	}
}

func (this *AccessController) Follow(requestId, requestedId string) error {
	if requestId == "" || requestedId == "" {
		return comerr.ErrParamInvalid
	}

	if this.relationshipExists(requestId, requestedId, Bind_By_Followed, Relation_Established) {
		return RelationshipDuplicateBind
	}

	return this.mgoWrapper.Insert(col_relationship, &MRelationship{
		Id:          bson.NewObjectId(),
		RequestId:   requestId,
		RequestedId: requestedId,
		BindBy:      Bind_By_Followed,
		Created:     time.Now().Unix(),
		State:       Relation_Established,
	})
}

func (this *AccessController) WhiteListBindRequest(requestId, requestedId string) error {
	if requestId == "" || requestedId == "" {
		return comerr.ErrParamInvalid
	}

	query := bson.M{"$or": []bson.M{bson.M{"request_id": requestId, "requested_id": requestedId}, bson.M{"request_id": requestedId, "requested_id": requestId}}, "bind_by": Bind_By_WhiteList}
	if c, err := this.mgoWrapper.Count(col_relationship, query); err != nil {
		return err
	} else if c != 0 {
		return RelationshipDuplicateBind
	}

	return this.mgoWrapper.Insert(col_relationship, &MRelationship{
		Id:          bson.NewObjectId(),
		RequestId:   requestId,
		RequestedId: requestedId,
		BindBy:      Bind_By_WhiteList,
		Created:     time.Now().Unix(),
		State:       Relation_Requested,
	})
}

func (this *AccessController) KickInBlackList(requestId, requestedId string) error {
	if requestId == "" || requestedId == "" {
		return comerr.ErrParamInvalid
	}

	if this.relationshipExists(requestId, requestedId, Bind_By_BlackList, Relation_Established) {
		return RelationshipDuplicateBind
	}

	return this.mgoWrapper.Insert(col_relationship, &MRelationship{
		Id:          bson.NewObjectId(),
		RequestId:   requestId,
		RequestedId: requestedId,
		BindBy:      Bind_By_BlackList,
		Created:     time.Now().Unix(),
		State:       Relation_Established,
	})
}

func (this *AccessController) UpdateRelationshipState(requestId, requestedId string, bind int, oldState, newState int) error {
	if requestId == "" || requestedId == "" || bind < Bind_By_Followed || bind > Bind_By_BlackList {
		return comerr.ErrParamInvalid
	}

	if !this.relationshipExists(requestId, requestedId, bind, oldState) {
		return RelationshipNotFound
	}

	return this.mgoWrapper.UpSetOne(col_relationship, bson.M{"request_id": requestId, "requested_id": requestedId, "bind_by": bind}, bson.M{"updated": time.Now().Unix(), "state": newState})
}

func (this *AccessController) relationshipExists(requestId, requestedId string, bind int, state int) bool {
	if requestId == "" || requestedId == "" || bind > Bind_By_BlackList || bind > Relation_Broke {
		return false
	}

	var query = make(bson.M)
	if bind > 0 {
		query["bind_by"] = bind
		if bind == Bind_By_WhiteList {
			query["$or"] = []bson.M{bson.M{"request_id": requestId, "requested_id": requestedId}, bson.M{"request_id": requestedId, "requested_id": requestId}}
		} else {
			query["request_id"] = requestId
			query["requested_id"] = requestedId
		}
	}
	if state > 0 {
		query["state"] = state
	}

	if c, err := this.mgoWrapper.Count(col_relationship, query); err != nil {
		return false
	} else {
		return c == 1
	}
}

func (this *AccessController) EmpowerRequest(requestId, resourceId string, authCode string) error {
	if requestId == "" || resourceId == "" || authCode == "" {
		return comerr.ErrParamInvalid
	}

	if !this.accessControlExists(resourceId, Access_AuthCode, AccCtrl_Normal) {
		return AccessControlNotFound
	}

	var query = bson.M{"request_id": requestId, "resource_id": resourceId}
	if c, err := this.mgoWrapper.Count(col_auth_record, query); err != nil {
		return err
	} else {
		if c == 0 {
			return this.mgoWrapper.Insert(col_auth_record, &MAuthRecord{
				Id:         bson.NewObjectId(),
				RequestId:  requestId,
				ResourceId: resourceId,
				AuthCode:   authCode,
				Created:    time.Now().Unix(),
				State:      AuthCode_Applied,
			})
		} else {
			return this.mgoWrapper.UpSetOne(col_auth_record, query, bson.M{"auth_code": authCode, "updated": time.Now().Unix()})
		}
	}
}

func (this *AccessController) EmpowerRespond(requestId, resourceId string, state int) error {
	if requestId == "" || resourceId == "" || state < AuthCode_Applied || state > AuthCode_Expired {
		return comerr.ErrParamInvalid
	}

	return this.mgoWrapper.UpSetOne(col_auth_record, bson.M{"request_id": requestId, "resource_id": resourceId}, bson.M{"updated": time.Now().Unix(), "state": state})
}

func (this *AccessController) accessControlExists(resourceId string, accTag int, state int) bool {
	if resourceId == "" || accTag > Access_Relationship || state > AccCtrl_Blocked {
		return false
	}

	query := bson.M{"resource_id": resourceId}
	if accTag > 0 {
		query["access_tag"] = accTag
	}
	if state > 0 {
		query["state"] = state
	}

	if c, err := this.mgoWrapper.Count(col_access_control, query); err != nil {
		return false
	} else {
		return c == 1
	}
}

func (this *AccessController) CreateAccessControl(createrId string, accTag int, accParam string, resourceIds ...string) error {
	if createrId == "" || len(resourceIds) == 0 || accTag < Access_Public || accTag > Access_Relationship || ((accTag == Access_Only || accTag == Access_Group) && accParam == "") {
		return comerr.ErrParamInvalid
	}

	if accTag == Access_Group {
		if c, err := this.mgoWrapper.Count(col_group, bson.M{"_id": bson.ObjectIdHex(accParam)}); err != nil {
			return err
		} else if c == 0 {
			return GroupNotFound
		}
	}

	var maccs []interface{}
	for k := range resourceIds {
		if this.accessControlExists(resourceIds[k], accTag, AccCtrl_Normal) {
			log.Println(AccessControlAlreadyExists.Error())
			continue
		}
		maccs = append(maccs, &MAccessControl{
			Id:          bson.NewObjectId(),
			CreaterId:   createrId,
			ResourceId:  resourceIds[k],
			AccessTag:   accTag,
			AccessParam: accParam,
			Created:     time.Now().Unix(),
			State:       AccCtrl_Normal,
		})
	}

	return this.mgoWrapper.Insert(col_access_control, maccs...)
}

func (this *AccessController) UpdateAccessControl(conCreaterId, conResourceId string, updAccTag int, updAccParam string) error {
	if conCreaterId == "" || conResourceId == "" || updAccTag < Access_Public || updAccTag > Access_Relationship || ((updAccTag == Access_Only || updAccTag == Access_Group) && updAccParam == "") {
		return comerr.ErrParamInvalid
	}

	if updAccTag == Access_Group {
		if c, err := this.mgoWrapper.Count(col_group, bson.M{"_id": bson.ObjectIdHex(updAccParam)}); err != nil {
			return err
		} else if c == 0 {
			return GroupNotFound
		}
	}

	return this.mgoWrapper.UpSetOne(col_access_control, bson.M{"creater_id": conCreaterId, "resource_id": conResourceId}, bson.M{"access_tag": updAccTag, "access_param": updAccParam})
}

func (this *AccessController) DeleteAccessControl(createrId, resourceId string, soft bool) error {
	if createrId == "" || resourceId == "" {
		return comerr.ErrParamInvalid
	}

	macc := &MAccessControl{}
	err := this.mgoWrapper.FindOne(col_access_control, bson.M{"creater_id": createrId, "resource_id": resourceId}, macc)
	if err != nil {
		return err
	}

	if soft {
		return this.mgoWrapper.UpSetOne(col_access_control, bson.M{"_id": macc.Id}, bson.M{"updated": time.Now().Unix(), "state": AccCtrl_Blocked})
	}

	if macc.AccessTag == Access_AuthCode {
		this.mgoWrapper.RemovAll(col_auth_record, bson.M{"resource_id": resourceId})
	}

	return this.mgoWrapper.RemoveOne(col_auth_record, bson.M{"_id": macc.Id})
}

func (this *AccessController) CreateGroup(createrId string, gtype int, name string, requestedIds ...string) (gid string, err error) {
	if createrId == "" || name == "" || gtype < Group_Allow || gtype > Group_Deny {
		return "", comerr.ErrParamInvalid
	}

	var c int
	c, err = this.mgoWrapper.Count(col_group, bson.M{"creater_id": createrId, "group_type": gtype, "name": name})
	if err != nil {
		return
	} else if c != 0 {
		err = GroupAlreadyExists

		return
	}

	gid = bson.NewObjectId().Hex()
	if err = this.mgoWrapper.Insert(col_group, &MGroup{
		Id:        bson.ObjectIdHex(gid),
		CreaterId: createrId,
		GroupType: gtype,
		Name:      name,
		Created:   time.Now().Unix(),
	}); err != nil {
		return
	}
	if len(requestedIds) != 0 {
		if err = this.AddIntoGroup(gid, createrId, requestedIds...); err != nil {
			return "", err
		}
	}

	return
}

func (this *AccessController) AddIntoGroup(groupId string, requestId string, requestedIds ...string) error {
	if groupId == "" || requestId == "" || len(requestedIds) == 0 {
		return comerr.ErrParamInvalid
	}

	gid := bson.ObjectIdHex(groupId)
	c, err := this.mgoWrapper.Count(col_group, bson.M{"_id": gid})
	if err != nil {
		return err
	} else if c == 0 {
		return GroupNotFound
	}

	var igs = make([]interface{}, len(requestedIds))
	for k := range requestedIds {
		igs[k] = &MInGroup{
			GroupId:     groupId,
			RequestId:   requestId,
			RequestedId: requestedIds[k],
		}
	}

	return this.mgoWrapper.Insert(col_in_group, igs...)
}

func (this *AccessController) RemoveFromGroup(groupId string, requestedIds ...string) error {
	if groupId == "" || len(requestedIds) == 0 {
		return comerr.ErrParamInvalid
	}

	gid := bson.ObjectIdHex(groupId)
	c, err := this.mgoWrapper.Count(col_group, bson.M{"_id": gid})
	if err != nil {
		return err
	} else if c == 0 {
		return GroupNotFound
	}

	_, err = this.mgoWrapper.RemovAll(col_in_group, bson.M{"group_id": groupId, "requested_id": bson.M{"$in": requestedIds}})

	return err
}

func (this *AccessController) DismissGroup(gid string) {
	this.mgoWrapper.RemoveOne(col_group, bson.M{"_id": bson.ObjectIdHex(gid)})
	this.mgoWrapper.RemovAll(col_in_group, bson.M{"group_id": gid})
	this.mgoWrapper.UpSetOne(col_access_control, bson.M{"group_id": gid}, bson.M{"access_tag": Access_Private, "group_id": ""})
}

func (this *AccessController) Access(requestId, resourceId string, authCode string) bool {
	if requestId == "" || resourceId == "" {
		return false
	}

	macc := &MAccessControl{}
	err := this.mgoWrapper.FindOne(col_access_control, bson.M{"resource_id": resourceId}, macc)
	if err != nil {
		log.Println(err.Error())

		return false
	}
	if macc.State != AccCtrl_Normal {
		return false
	}

	if requestId == macc.CreaterId {
		return true
	} else if macc.AccessTag == Access_Public {
		return !this.relationshipExists(macc.CreaterId, requestId, Bind_By_BlackList, Relation_Established)
	} else if macc.AccessTag == Access_Private {
		return requestId == macc.CreaterId
	} else if macc.AccessTag == Access_Only {
		return requestId == macc.AccessParam
	} else if macc.AccessTag == Access_Group && macc.AccessParam != "" {
		if !bson.IsObjectIdHex(macc.AccessParam) {
			return false
		}
		mgp := &MGroup{}
		if err = this.mgoWrapper.FindOne(col_group, bson.M{"_id": bson.ObjectIdHex(macc.AccessParam)}, mgp); err != nil {
			log.Println(err.Error())

			return false
		} else {
			if c, err := this.mgoWrapper.Count(col_in_group, bson.M{"group_id": macc.AccessParam, "request_id": macc.CreaterId, "requested_id": requestId}); err != nil {
				log.Println(err.Error())

				return false
			} else {
				return (mgp.GroupType == Group_Allow && c != 0) || (mgp.GroupType == Group_Deny && c == 0)
			}
		}
	} else if macc.AccessTag == Access_AuthCode && authCode != "" {
		if c, err := this.mgoWrapper.Count(col_auth_record, bson.M{"request_id": requestId, "resource_id": resourceId, "auth_code": authCode, "state": AuthCode_Empowered}); err != nil {
			log.Println(err.Error())

			return false
		} else {
			return c == 1
		}
	} else if macc.AccessTag == Access_Relationship {
		return this.relationshipExists(requestId, macc.CreaterId, Bind_By_WhiteList, Relation_Established)
	} else {
		return false
	}
}
