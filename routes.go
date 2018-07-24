// The routes.go file encapsulates all routes for the API router
// in order to remove them from the main file (backend.go)
package main

import (
	"log"
	"net/http"

	"github.com/dersteps/club-backend/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func sendInternalError(ctx *gin.Context) {
	ctx.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

// GetAllUsersV1 retrieves all users from the underlying mongodb database
// and renders them as JSON.
// If the operation succeeds, HTTP200 is returned, otherwise, HTTP500 is returned.
func getAllUsersV1(ctx *gin.Context) {
	users, err := db.FindAllUsers()
	if err != nil {
		sendInternalError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetAllFunctionsV1 will retrieve all club functions from the database and
// render them as JSON.
// If the operation succeeds, HTTP200 is returned, otherwise, HTTP500 is returned.
func getAllFunctionsV1(ctx *gin.Context) {
	functions, err := db.FindAllFunctions()
	if err != nil {
		sendInternalError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, functions)
}

// GetAllMembersV1 will retrieve all members fro mthe database and render them as JSON.
func getAllMembersV1(ctx *gin.Context) {
	members, err := db.FindAllMembers()
	if err != nil {
		sendInternalError(ctx)
		return
	}
	ctx.JSON(http.StatusOK, members)
}

// CreateUserV1 creates a new user from a form based POST request.
// It will attempt to create the new user in the database and render their
// data as JSON.
// If the operation succeeds, HTTP200 is returned, otherwise, HTTP500 is returned.
func createUserV1(ctx *gin.Context) {
	// Get user from form
	nick := ctx.PostForm("username")
	mail := ctx.PostForm("mail")
	pass := ctx.PostForm("password_hash")
	name := ctx.PostForm("name")

	user := model.User{Username: nick, Email: mail, Password: pass, Name: name}
	user.ID = bson.NewObjectId()

	if err := db.InsertUser(user); err != nil {
		log.Fatal("Unable to create new user!")
		sendInternalError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// NotImplemented is a catch-all function for not-yet-implemented routes
func notImplemented(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Not implemented yet")
}
