package dao

/*
import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/dersteps/club-backend/model"
)

type FunctionsDAO struct {
}

const (
	COLLECTION = "functions"
)

var db *mgo.Database

// Establish database connection
func (dao *FunctionsDAO) Connect(info DBInfo) {

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

func (dao *FunctionsDAO) FindAll() ([]model.Function, error) {
	var functions []model.Function
	err := db.C(COLLECTION).Find(bson.M{}).All(&functions)
	return functions, err
}

func (dao *FunctionsDAO) FindByID(id string) (model.Function, error) {
	var function model.Function
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&function)
	return function, err
}

func (dao *FunctionsDAO) FindByName(name string) (model.Function, error) {
	var function model.Function
	err := db.C(COLLECTION).Find(bson.M{"username": name}).One(&function)
	return function, err
}

func (dao *FunctionsDAO) Insert(function model.Function) error {
	err := db.C(COLLECTION).Insert(&function)
	return err
}

func (dao *FunctionsDAO) Delete(function model.Function) error {
	err := db.C(COLLECTION).Remove(&function)
	return err
}

func (dao *FunctionsDAO) Update(function model.Function) error {
	err := db.C(COLLECTION).Update(function.ID, &function)
	return err
}
*/
