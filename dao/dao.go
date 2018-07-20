package dao

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/dersteps/club-backend/model"
)

type DBInfo struct {
	Server   string
	Timeout  int
	Database string
	Username string
	Password string
}

type DAO struct {
}

const (
	USERS     = "users"
	FUNCTIONS = "functions"
	MEMBERS   = "members"
)

var db *mgo.Database

// Establish database connection
func (dao *DAO) Connect(info DBInfo) {

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{info.Server},
		Timeout:  time.Duration(info.Timeout) * time.Second,
		Database: info.Database,
		Username: info.Username,
		Password: info.Password,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}

	db = session.DB(dialInfo.Database)
}

func (dao *DAO) FindAllUsers() ([]model.User, error) {
	var users []model.User
	err := db.C(USERS).Find(bson.M{}).All(&users)
	return users, err
}

func (dao *DAO) FindUserByID(id string) (model.User, error) {
	var user model.User
	err := db.C(USERS).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func (dao *DAO) FindUserByName(name string) (model.User, error) {
	var user model.User
	err := db.C(USERS).Find(bson.M{"username": name}).One(&user)
	return user, err
}

func (dao *DAO) InsertUser(user model.User) error {
	err := db.C(USERS).Insert(&user)
	return err
}

func (dao *DAO) DeleteUser(user model.User) error {
	err := db.C(USERS).Remove(&user)
	return err
}

func (dao *DAO) UpdateUser(user model.User) error {
	err := db.C(USERS).Update(user.ID, &user)
	return err
}

// Functions
func (dao *DAO) FindAllFunctions() ([]model.Function, error) {
	var functions []model.Function
	err := db.C(FUNCTIONS).Find(bson.M{}).All(&functions)
	return functions, err
}

func (dao *DAO) FindFunctionByID(id string) (model.Function, error) {
	var function model.Function
	err := db.C(FUNCTIONS).FindId(bson.ObjectIdHex(id)).One(&function)
	return function, err
}

func (dao *DAO) FindFunctionByName(name string) (model.Function, error) {
	var function model.Function
	err := db.C(FUNCTIONS).Find(bson.M{"name": name}).One(&function)
	return function, err
}

func (dao *DAO) InsertFunction(function model.Function) error {
	err := db.C(FUNCTIONS).Insert(&function)
	return err
}

func (dao *DAO) DeleteFunction(function model.Function) error {
	err := db.C(FUNCTIONS).Remove(&function)
	return err
}

func (dao *DAO) UpdateFunction(function model.Function) error {
	err := db.C(FUNCTIONS).Update(function.ID, &function)
	return err
}

// Members
func (dao *DAO) FindAllMembers() ([]model.Member, error) {
	var members []model.Member
	err := db.C(MEMBERS).Find(bson.M{}).All(&members)
	return members, err
}

func (dao *DAO) FindMemberByID(id string) (model.Member, error) {
	var member model.Member
	err := db.C(MEMBERS).FindId(bson.ObjectIdHex(id)).One(&member)
	return member, err
}

func (dao *DAO) FindMemberByName(name string) (model.Member, error) {
	var member model.Member
	err := db.C(MEMBERS).Find(bson.M{"name": name}).One(&member)
	return member, err
}

func (dao *DAO) InsertMember(member model.Member) error {
	err := db.C(MEMBERS).Insert(&member)
	return err
}

func (dao *DAO) DeleteMember(member model.Member) error {
	err := db.C(MEMBERS).Remove(&member)
	return err
}

func (dao *DAO) UpdateMember(member model.Member) error {
	err := db.C(MEMBERS).Update(member.ID, &member)
	return err
}
