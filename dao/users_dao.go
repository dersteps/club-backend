package dao

/*
import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/dersteps/club-backend/model"
)

type UsersDAO struct {
}

const (
	COLLECTION = "users"
)

var db *mgo.Database

// Establish database connection
func (dao *UsersDAO) Connect(info DBInfo) {

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

func (dao *UsersDAO) FindByName(name string) (model.User, error) {
	var user model.User
	err := db.C(COLLECTION).Find(bson.M{"username": name}).One(&user)
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
*/
