package main

import (
	"crypto/sha256"
	"log"
	"regexp"

	"github.com/dersteps/club-backend/model"
	"github.com/dersteps/club-backend/util"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

var RoleAdmin = "admin"
var RoleUserAdmin = "useradmin"
var RoleUser = "user"
var RoleAny = "*"
var RoleNone = "x"

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
	mFunctions["GET"] = []string{RoleUser}
	mFunctions["POST"] = []string{RoleAdmin}
	mFunctions["PUT"] = []string{RoleAdmin}
	mFunctions["DELETE"] = []string{RoleAdmin}

	mMembers := make(map[string][]string)
	mMembers["GET"] = []string{RoleUser}
	mMembers["POST"] = []string{RoleAdmin}
	mMembers["PUT"] = []string{RoleAdmin}
	mMembers["DELETE"] = []string{RoleAdmin}

	m["/api/v[0-9]{1,}/users"] = mUsers

	m["/api/auth/*"] = mAuth
	m["/api/login"] = mAuth
	m["/api/v[0-9]{1,}/functions"] = mFunctions
	m["/api/v[0-9]{1,}/members"] = mMembers
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

	return user, passwordMatches(password, user.Password)
}

func isAuthorized(username string, reqURL string, reqMethod string, roles []string) bool {
	//util.Info(fmt.Sprintf("Testing auth for '%s' on [%s] '%s'...", username, reqMethod, reqURL))
	for reg := range m {
		r, _ := regexp.Compile(reg)
		if r.MatchString(reqURL) {
			for _, requiredRole := range m[reg][reqMethod] {
				if requiredRole == RoleAny {
					//util.Info(fmt.Sprintf("Anyone can access [%s] %s, letting '%s' in", reqMethod, reqURL, username))
					return true
				}
				for _, userRole := range roles {
					if userRole == requiredRole {
						//util.Info(fmt.Sprintf("'%s' has role '%s', which enables access to [%s] %s", username, userRole, reqMethod, reqURL))
						return true
					}
				}
			}
		}
	}
	//util.Warn(fmt.Sprintf("User '%s' is not authorized to access [%s] %s", username, reqMethod, reqURL))
	return false
}

func jwtAuthorize(user interface{}, c *gin.Context) bool {
	//log.Printf("USER in authorize method: %s\n", user)
	userObject, err := db.FindUserByName(user.(string))

	if err != nil {
		log.Printf("Error while searching for user %s: %s\n", userObject, err.Error())
		return false
	}

	claims := jwt.ExtractClaims(c)
	roles, err := util.MakeStringSlice(claims["roles"])

	if err != nil {
		log.Printf("Unable to determine user's roles from token, bailing!")
		return false
	}

	return isAuthorized(userObject.Username, c.Request.URL.String(), c.Request.Method, roles)
}

func jwtUnauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"code": code, "message": message})
}

func jwtPayload(user interface{}) jwt.MapClaims {
	log.Printf("USER in PAYLOAD: %s\n", user)
	modelUser := user.(model.User)
	return map[string]interface{}{
		"roles": modelUser.Roles,
	}
}
