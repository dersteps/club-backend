package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/dersteps/club-backend/model"
	"gopkg.in/mgo.v2/bson"

	"github.com/apcera/termtables"

	"github.com/gin-gonic/gin"

	"github.com/dersteps/club-backend/config"
	"github.com/dersteps/club-backend/dao"
)

// Config object for convenient access.
var cfg = config.Config{}
var userDAO = dao.UsersDAO{}
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

func ensureAdmin() {
	/*[admin]
	username="administrator"
	email="stefan.matyba@googlemail.com"
	password="ksljdfosdjfisdjfiosdf" #<- hash!*/
	adminUser, err := userDAO.FindByName(cfg.Admin.Username)
	if err != nil {
		// user not found
		log.Println("No admin user found, creating one")
		admin := model.User{
			Email:    cfg.Admin.Mail,
			Name:     "Administrator",
			Password: string(sha256.New().Sum([]byte(cfg.Admin.Password))),
			Username: cfg.Admin.Username,
			Roles:    []string{RoleAdmin, RoleUserAdmin, RoleUser},
		}
		admin.ID = bson.NewObjectId()
		if err = userDAO.Insert(admin); err != nil {
			panic(err)
		}
		return
	}

	log.Printf("Admin user found: %s\n", adminUser)

}

// init is a special function called by go automagically.
// Initializes basic stuff.
func init() {
	// Read config
	err := cfg.Read("config.toml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	authMiddleware.Key = []byte(cfg.API.Secret)

	table := termtables.CreateTable()
	table.AddHeaders("Item", "Value")
	table.AddRow("DB Server", cfg.Database.Host)
	table.AddRow("DB Name", cfg.Database.Name)
	table.AddRow("DB User", cfg.Database.Username)
	table.AddRow("DB Timeout", fmt.Sprintf("%d seconds", cfg.Database.Timeout))
	table.AddRow("DB Password", "LOL, just kidding")
	table.AddRow("Server", cfg.Server.Host)
	table.AddRow("Port", cfg.Server.Port)

	fmt.Println(table.Render())

	// Init database info
	dbInfo.Database = cfg.Database.Name
	dbInfo.Password = cfg.Database.Password
	dbInfo.Server = cfg.Database.Host
	dbInfo.Timeout = cfg.Database.Timeout
	dbInfo.Username = cfg.Database.Username

	userDAO.Connect(dbInfo)

	ensureAdmin()

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
