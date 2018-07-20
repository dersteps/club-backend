package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"

	jwt "github.com/appleboy/gin-jwt"

	"github.com/gin-gonic/gin"

	"github.com/dersteps/club-backend/config"
	"github.com/dersteps/club-backend/dao"
	"github.com/dersteps/club-backend/util"
)

// Config object for convenient access.
var cfg = config.Config{}
var db = dao.DAO{}
var dbInfo = dao.DBInfo{}

var authMiddleware = &jwt.GinJWTMiddleware{
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

// init is a special function called by go automagically.
// Initializes basic stuff.
func init() {

	banner()

	// Read config
	err := cfg.Read("config.toml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	authMiddleware.Key = []byte(cfg.API.Secret)

	util.Info(fmt.Sprintf("Will provide my service at %s:%s", aurora.BgGreen(cfg.Server.Host).Bold(), aurora.BgBlue(cfg.Server.Port).Bold()))
	util.Info(fmt.Sprintf("Mongo DB server is %s, database %s, username %s",
		aurora.BgBlue(cfg.Database.Host).Bold(), aurora.BgBlue(cfg.Database.Name).Bold(), aurora.BgBlue(cfg.Database.Username).Bold()))

	// Init database info
	dbInfo.Database = cfg.Database.Name
	dbInfo.Password = cfg.Database.Password
	dbInfo.Server = cfg.Database.Host
	dbInfo.Timeout = cfg.Database.Timeout
	dbInfo.Username = cfg.Database.Username

	db.Connect(dbInfo)

	InitDatabase()
}

// setupRoutes will let the router/gin engine know what routes
// there are and what handler functions are mapped
func setupRoutes(router *gin.Engine) {
	// We'll group the routes, so that we are future proof
	router.POST("/api/login", authMiddleware.LoginHandler)

	auth := router.Group("/api/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware.MiddlewareFunc())
	v1.GET("/users", GetAllUsersV1)
	v1.POST("/users", CreateUserV1)
	v1.PUT("/users", NotImplemented)
	v1.DELETE("/users", NotImplemented)

	v1.GET("/functions", NotImplemented)
	v1.POST("/functions", NotImplemented)
	v1.PUT("/functions", NotImplemented)
	v1.DELETE("/functions", NotImplemented)
}

// main is the application's entry point of course.
func main() {

	// Setup a default gin router with logging
	router := gin.Default()

	// Let the router know about the routes
	setupRoutes(router)

	// Listen address is read from the config file
	parts := []string{cfg.Server.Host, cfg.Server.Port}
	listenAddr := strings.Join(parts, ":")

	// Run the router
	router.Run(listenAddr)
}
