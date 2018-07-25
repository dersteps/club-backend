package model

import "github.com/globalsign/mgo/bson"

type MemberState int

const (
	Inactive   MemberState = 0
	FullMember MemberState = 666
	Honourable MemberState = 667
	Prospect   MemberState = 2
	Undefined  MemberState = 3
)

type NumberType int

const (
	Phone  NumberType = 0
	Mobile NumberType = 1
	Work   NumberType = 2
	Other  NumberType = 3
)

type PhoneNumber struct {
	Number  string     `bson:"number" json:"number"`
	Primary bool       `bson:"primary" json:"primary"`
	Type    NumberType `bson:"type" json:"type"`
}

type Mail struct {
	Address string `bson:"address" json:"address"`
	Primary bool   `bson:"primary" json:"primary"`
}

type Address struct {
	Street  string `bson:"street" json:"street"`
	Number  string `bson:"number" json:"number"`
	Zip     string `bson:"zip" json:"zip"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	Country string `bson:"country" json:"country"`
}

type Date struct {
	Day   int `bson:"day" json:"day"`
	Month int `bson:"month" json:"month"`
	Year  int `bson:"year" json:"year"`
}

type Member struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	MemberID    string        `bson:"member_id" json:"member_id"`
	Name        string        `bson:"name" json:"name"`
	MiddleNames []string      `bson:"middlenames" json:"middlenames"`
	Surname     string        `bson:"surname" json:"surname"`
	Address     Address       `bson:"address" json:"address"`
	Emails      []Mail        `bson:"emails" json:"emails"`
	Numbers     []PhoneNumber `bson:"phonenumbers" json:"phonennumbers"`
	Birthday    Date          `bson:"birthday" json:"birthday"`
	State       MemberState   `bson:"state" json:"state"`
	Functions   []Function    `bson:"functions" json:"functions"`
}
