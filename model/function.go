package model

import "github.com/globalsign/mgo/bson"

type Function struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Board bool          `bson:"is_board_function" json:"is_board_function"`
}
