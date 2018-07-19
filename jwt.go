package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

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

var m = make(map[string]map[string][]string)

func init() {

	mUsers := make(map[string][]string)
	mUsers["GET"] = []string{RoleAdmin, RoleUser}
	mUsers["POST"] = []string{RoleAdmin, RoleUserAdmin}
	mUsers["PUT"] = []string{RoleAdmin, RoleUserAdmin}
	mUsers["DELETE"] = []string{RoleAdmin, RoleUserAdmin}

	m["(/api/v[0-9]{1,}/)(.*)"] = mUsers

}

func passwordMatches(password, dbPassword string) bool {
	hash := sha256.New().Sum([]byte(password))
	return string(hash) == dbPassword
}

func jwtAuthenticate(username string, password string, c *gin.Context) (interface{}, bool) {
	// Attempt to find user in database, match password and there you go
	user, err := userDAO.FindByName(username)
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
			for _, userRole := range roles {
				for _, requiredRole := range m[reg][reqMethod] {
					if requiredRole == userRole {
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
	user, err := userDAO.FindByName(user.(string))
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

	/*log.Printf("User has roles: %s\n", roles)
	log.Printf("Type of roles value: %s\n", reflect.TypeOf(roles))
	log.Printf("Extracted claims: %s\n", claims)

	log.Printf("Auth check for url %s\n", c.Request.URL)

	for reg := range m {
		log.Printf("Regex: %s, roles: %s", reg, m[reg])
		r, _ := regexp.Compile(reg)
		if r.MatchString(c.Request.URL.String()) {
			log.Printf("Regex matches url pattern: %s\n", reg)

			requiredRoles := m[reg]
			log.Printf("Type of required roles: %s\n", reflect.TypeOf(requiredRoles))
			log.Printf("Required roles: %s\n", requiredRoles)

			for _, userRole := range roles {
				log.Printf("Testing user role '%s'\n", userRole)
				for _, requiredRole := range requiredRoles {
					log.Printf("  -> Required role: %s\n", requiredRole)
					if requiredRole == userRole {
						log.Println("User has sufficient authorization!")
						return true
					}
				}
			}

		}
	}

	//log.Printf("Warning: all users are authorized for everything!")
	return false*/
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

var AuthMiddleware = &jwt.GinJWTMiddleware{
	Realm:         "Hello there",
	Key:           []byte("myapisecret"),
	Timeout:       time.Hour,
	MaxRefresh:    time.Hour,
	Authenticator: jwtAuthenticate,
	Authorizator:  jwtAuthorize,
	Unauthorized:  jwtUnauthorized,
	TokenLookup:   "header:Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
	PayloadFunc:   jwtPayload,
}
