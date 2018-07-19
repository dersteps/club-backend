// The routes.go file encapsulates all routes for the API router
// in order to remove them from the main file (backend.go)
package main

import (
	"log"
	"net/http"

	"github.com/dersteps/club-backend/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// GetAllUsersV1 retrieves all users from the underlying mongodb database
// and renders them as JSON.
// If the operation succeeds, HTTP200 is returned, otherwise, HTTP500 is returned.
func GetAllUsersV1(ctx *gin.Context) {
	users, err := userDAO.FindAll()
	if err != nil {
		ctx.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// CreateUserV1 creates a new user from a form based POST request.
// It will attempt to create the new user in the database and render their
// data as JSON.
// If the operation succeeds, HTTP200 is returned, otherwise, HTTP500 is returned.
func CreateUserV1(ctx *gin.Context) {
	// Get user from form
	nick := ctx.PostForm("username")
	mail := ctx.PostForm("mail")
	pass := ctx.PostForm("password_hash")
	name := ctx.PostForm("name")

	user := model.User{Username: nick, Email: mail, Password: pass, Name: name}
	user.ID = bson.NewObjectId()

	if err := userDAO.Insert(user); err != nil {
		log.Fatal("Unable to create new user!")
		ctx.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// NotImplemented is a catch-all function for not-yet-implemented routes
func NotImplemented(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Not implemented yet")
}
