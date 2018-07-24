package model

import "github.com/globalsign/mgo/bson"

/*type Permissions struct {
	ListUsers  bool `bson:"list_users" json:"list_users"`
	CreateUser bool `bson:"create_user" json:"create_user"`
	DeleteUser bool `bson:"delete_user" json:"delete_user"`
	EditUser   bool `bson:"edit_user" json:"edit_user"`
}
*/
type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Username string        `bson:"username" json:"username"`
	Email    string        `bson:"mail" json:"mail"`
	Password string        `bson:"pass" json:"pass"`
	Name     string        `bson:"name" json:"name"`
	Roles    []string      `bson:"roles" json:"roles"`
}
