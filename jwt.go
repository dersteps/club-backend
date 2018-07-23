package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/dersteps/club-backend/model"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

/*
	string	-> string	-> map[string][]string
	url		-> GET  -> []
			-> POST -> []


*/

var RoleAdmin = "admin"
var RoleUserAdmin = "useradmin"
var RoleUser = "user"
var RoleAny = "*"

var m = make(map[string]map[string][]string)

func init() {

	mUsers := make(map[string][]string)
	mUsers["GET"] = []string{RoleAdmin, RoleUser}
	mUsers["POST"] = []string{RoleAdmin, RoleUserAdmin}
	mUsers["PUT"] = []string{RoleAdmin, RoleUserAdmin}
	mUsers["DELETE"] = []string{RoleAdmin, RoleUserAdmin}

	mAuth := make(map[string][]string)
	mAuth["GET"] = []string{RoleAny}
	mAuth["PUT"] = []string{RoleAny}
	mAuth["POST"] = []string{RoleAny}

	mFunctions := make(map[string][]string)
	// Anyone can read what functions there are in the club, no problem
	mFunctions["GET"] = []string{RoleAny}
	mFunctions["POST"] = []string{RoleAdmin}
	mFunctions["PUT"] = []string{RoleAdmin}
	mFunctions["DELETE"] = []string{RoleAdmin}

	m["(/api/v[0-9]{1,}/)(users)"] = mUsers

	m["/api/auth/*"] = mAuth
	m["/api/login"] = mAuth
	m["(/api/v[0-9]{1,}/)(function)"] = mAuth

}

func passwordMatches(password, dbPassword string) bool {
	hash := sha256.New().Sum([]byte(password))
	return string(hash) == dbPassword
}

func jwtAuthenticate(username string, password string, c *gin.Context) (interface{}, bool) {
	// Attempt to find user in database, match password and there you go
	user, err := db.FindUserByName(username)
	if err != nil {
		log.Printf("Error while searching for user %s: %s\n", username, err.Error())
		return nil, false
	}

	log.Printf("Found user %s!\n", username)
	log.Println(user)

	return user, passwordMatches(password, user.Password)
}

func makeStringSlice(from interface{}) ([]string, error) {
	slice := []string{}
	if reflect.TypeOf(from).Kind() == reflect.Slice {
		tmp := reflect.ValueOf(from)
		for i := 0; i < tmp.Len(); i++ {
			slice = append(slice, fmt.Sprintf("%v", tmp.Index(i)))
		}
		return slice, nil
	} else {
		return nil, errors.New("Unable to convert to string slice")
	}

}

func isAuthorized(reqURL string, reqMethod string, roles []string) bool {
	for reg := range m {
		r, _ := regexp.Compile(reg)
		if r.MatchString(reqURL) {
			for _, requiredRole := range m[reg][reqMethod] {
				if requiredRole == RoleAny {
					return true
				}
				for _, userRole := range roles {
					if userRole == requiredRole {
						return true
					}
				}
			}
		}
	}
	return false
}

func jwtAuthorize(user interface{}, c *gin.Context) bool {
	//log.Printf("USER in authorize method: %s\n", user)
	user, err := db.FindUserByName(user.(string))
	if err != nil {
		log.Printf("Error while searching for user %s: %s\n", user, err.Error())
		return false
	}

	claims := jwt.ExtractClaims(c)
	roles, err := makeStringSlice(claims["roles"])

	if err != nil {
		log.Printf("Unable to determine user's roles from token, bailing!")
		return false
	}

	return isAuthorized(c.Request.URL.String(), c.Request.Method, roles)
}

func jwtUnauthorized(c *gin.Context, code int, message string) {
	log.Println("Unauthorized -> Deny")
	c.JSON(code, gin.H{"code": code, "message": message})
}

func jwtPayload(user interface{}) jwt.MapClaims {
	log.Printf("USER in PAYLOAD: %s\n", user)
	modelUser := user.(model.User)
	return map[string]interface{}{
		"roles": modelUser.Roles,
	}
}
