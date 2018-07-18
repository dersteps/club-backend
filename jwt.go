package main

import (
	"log"
	"time"

	"github.com/dersteps/club-backend/model"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func jwtAuthenticate(username string, password string, c *gin.Context) (interface{}, bool) {
	// Attempt to find user in database, match password and there you go
	user, err := userDAO.FindByName(username)
	if err != nil {
		log.Printf("Error while searching for user %s: %s\n", username, err.Error())
		return nil, false
	}

	log.Printf("Found user %s!\n", username)
	log.Println(user)

	log.Println("Warning, user passwords are currently not hashed")
	if user.Password == password {
		return user, true
	} else {
		return user, false
	}
}

func jwtAuthorize(user interface{}, c *gin.Context) bool {
	//log.Printf("USER in authorize method: %s\n", user)
	user, err := userDAO.FindByName(user.(string))
	if err != nil {
		log.Printf("Error while searching for user %s: %s\n", user, err.Error())
		return false
	}

	claims := jwt.ExtractClaims(c)
	log.Printf("Extracted claims: %s\n", claims)

	log.Printf("Auth check for url %s\n", c.Request.URL)

	log.Printf("Warning: all users are authorized for everything!")
	return true
}

func jwtUnauthorized(c *gin.Context, code int, message string) {
	log.Println("Unauthorized -> Deny")
	c.JSON(code, gin.H{"code": code, "message": message})
}

func jwtPayload(user interface{}) jwt.MapClaims {
	log.Printf("USER in PAYLOAD: %s\n", user)
	modelUser := user.(model.User)
	return map[string]interface{}{
		"role": modelUser.Role,
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
