package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora"

	"github.com/dersteps/club-backend/model"
	"github.com/dersteps/club-backend/util"
	"github.com/globalsign/mgo/bson"
)

var functionPresident = "President"
var functionVicePresident = "Vice President"
var functionTreasurer = "Treasurer"
var functionScribe = "Scribe"
var functionsSAA = "Bitch at Arms"

var triggerDBInitVariable = "CREW_INIT"

// Assumes the database connection is established!
func InitDatabase() {

	envValue := os.Getenv(triggerDBInitVariable)
	if envValue == "" {
		return
	}

	b, err := strconv.ParseBool(envValue)

	if err != nil || b != true {
		util.Warn(fmt.Sprintf("%s is set to '%s'. Set it to 'true' in order to init the database with default values", triggerDBInitVariable, envValue))
		return
	}

	ensureAdmin()
	ensureDefaultFunctions()
}

func ensureAdmin() {

	adminUser, err := db.FindUserByName(cfg.Admin.Username)
	if err != nil {
		// user not found
		util.Warn("No admin user found, creating one")
		admin := model.User{
			Email:    cfg.Admin.Mail,
			Name:     "Administrator",
			Password: string(sha256.New().Sum([]byte(cfg.Admin.Password))),
			Username: cfg.Admin.Username,
			Roles:    []string{RoleAdmin, RoleUserAdmin, RoleUser},
		}
		admin.ID = bson.NewObjectId()
		if err = db.InsertUser(admin); err != nil {
			panic(err)
		}
		util.Success(fmt.Sprintf("Created admin user %s", aurora.Green(admin.Username)))
		return
	}

	util.Info(fmt.Sprintf("Admin user found: %s", aurora.Green(adminUser.Username)))

}

func ensureFunction(function string) {
	_, err := db.FindFunctionByName(function)
	if err != nil {
		util.Warn(fmt.Sprintf("Default function not found: %s", aurora.Red(function)))
		// Insert that function!
		newFunction := model.Function{
			Name:  function,
			Board: true,
		}
		newFunction.ID = bson.NewObjectId()
		if err = db.InsertFunction(newFunction); err != nil {
			panic(err)
		}
		util.Success(fmt.Sprintf("Created default function %s", aurora.Green(newFunction.Name)))
		return
	}
	util.Info(fmt.Sprintf("Default function present in database: %s", aurora.Green(function)))
}

func ensureDefaultFunctions() {
	functionsMap := []string{
		functionPresident,
		functionVicePresident,
		functionScribe,
		functionTreasurer,
		functionsSAA,
	}

	for _, function := range functionsMap {
		ensureFunction(function)
	}

}
