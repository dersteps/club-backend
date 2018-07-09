package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Username string        `bson:"username" json:"username"`
	Email    string        `bson:"mail" json:"mail"`
	Password string        `bson:"pass" json:"pass"`
	Name     string        `bson:"name" json:"name"`
}
