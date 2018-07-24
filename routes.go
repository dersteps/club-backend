// The routes.go file encapsulates all routes for the API router
// in order to remove them from the main file (backend.go)
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dersteps/club-backend/model"
	"github.com/dersteps/club-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

type jsonObject map[string]interface{}

func sendInternalError(ctx *gin.Context) {
	ctx.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func sendBadRequest(ctx *gin.Context) {
	ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
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

	// Parse JSON data
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err == nil {
		// Successfully parsed data
		user.ID = bson.NewObjectId()
		if err := db.InsertUser(user); err != nil {
			log.Fatal("Unable to create new user!")
			sendInternalError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, user)
	} else {
		sendBadRequest(ctx)
	}
}

func createFunctionV1(ctx *gin.Context) {
	// Parse JSON data
	var function model.Function
	if err := ctx.ShouldBindJSON(&function); err == nil {
		// Successfully parsed POST body JSON data
		function.ID = bson.NewObjectId()
		if err := db.InsertFunction(function); err != nil {
			util.Fatal(fmt.Sprintf("Unable to insert function: '%s'", err.Error()))
			sendInternalError(ctx)
			return
		}
		// Send created function back as a JSON objecz
		ctx.JSON(http.StatusOK, function)
	} else {
		sendBadRequest(ctx)
	}
}

func extractFunctions(ctx *gin.Context) ([]bson.ObjectId, error) {
	var bodyBytes []byte
	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		return nil, errors.New("Unable to access raw request body")
	}

	// Make sure the request body remains parseable
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	dec := json.NewDecoder(bytes.NewReader(bodyBytes))
	data := map[string]interface{}{}
	dec.Decode(&data)

	ids := data["function_ids"]
	strIDs, err := util.MakeStringSlice(ids)
	if err != nil {
		return nil, err
	}

	// Find the functions in the database, return slice
	var ret []bson.ObjectId
	for _, fID := range strIDs {
		function, err := db.FindFunctionByID(fID)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to fetch function from database: %s", fID))
		}

		util.Info(fmt.Sprintf("New Member has function: %s", function.Name))
		ret = append(ret, function.ID)
	}
	fmt.Printf("Member's functions: \n")
	fmt.Println(ret)
	return ret, nil
}

func createMemberV1(ctx *gin.Context) {
	// Save raw body contents
	if ctx.Request.Body == nil {
		util.Fatal("Empty request body")
		sendBadRequest(ctx)
		return
	}

	_, err := extractFunctions(ctx)
	if err != nil {
		util.Fatal("Unable to parse request body as JSON")
		sendBadRequest(ctx)
		return
	}

	// Parse JSON data
	var member model.Member

	if err := ctx.ShouldBindJSON(&member); err == nil {
		// Successfully parsed POST body JSON data
		member.ID = bson.NewObjectId()

		// Find the functions
		functions, err := extractFunctions(ctx)

		if err != nil {
			util.Fatal(fmt.Sprintf("Extracting functions failed: %s", err.Error()))
			sendInternalError(ctx)
			return
		}

		// Set functions

		if err := db.InsertMember(member); err != nil {
			util.Fatal(fmt.Sprintf("Unable to insert member: '%s'", err.Error()))
			sendInternalError(ctx)
			return
		}

		member.Functions = functions
		if err := db.UpdateMember(member); err != nil {
			util.Fatal(fmt.Sprintf("Unable to update member [functions]: %s", err.Error()))
			sendInternalError(ctx)
			return
		}

		// Send created member back as a JSON object

		ctx.JSON(http.StatusOK, member)
	} else {
		util.Fatal("Nope")
		util.Fatal(err.Error())
		sendBadRequest(ctx)
	}
}

// NotImplemented is a catch-all function for not-yet-implemented routes
func notImplemented(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Not implemented yet")
}
