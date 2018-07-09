package dao

import (
	"log"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"

	"github.com/dersteps/club-backend/model"
)

type UsersDAO struct {
	Server   string
	Database string
}

const (
	COLLECTION = "users"
)

var db *mgo.Database

// Establish database connection
func (dao *UsersDAO) Connect() {
	session, err := mgo.Dial(dao.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(dao.Database)
}

func (dao *UsersDAO) FindAll() ([]model.User, error) {
	var users []model.User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func (dao *UsersDAO) FindByID(id string) (model.User, error) {
	var user model.User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func (dao *UsersDAO) Insert(user model.User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

func (dao *UsersDAO) Delete(user model.User) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

func (dao *UsersDAO) Update(user model.User) error {
	err := db.C(COLLECTION).Update(user.ID, &user)
	return err
}
