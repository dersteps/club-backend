package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/globalsign/mgo/bson"

	"github.com/logrusorgru/aurora"

	jwt "github.com/appleboy/gin-jwt"

	"github.com/gin-gonic/gin"

	"github.com/dersteps/club-backend/config"
	"github.com/dersteps/club-backend/dao"
	"github.com/dersteps/club-backend/model"
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

var (
	app         = kingpin.New("crew", "The crew server application")
	appNoBanner = app.Flag("no-banner", "Disable the banner").Bool()
	appDebug    = app.Flag("debug", "Sets debug mode for REST API framework (gin-gonic)").Bool()

	add = app.Command("add", "Adds an entity to the database")

	addFunction      = add.Command("function", "Used to add a function to the database")
	addFunctionName  = addFunction.Arg("name", "The name of the function").Required().String()
	addFunctionBoard = addFunction.Flag("board", "If given, this is a board function").Bool()

	addUser             = add.Command("user", "Adds a user to the database")
	addUserName         = addUser.Arg("name", "The username of the user to add").Required().String()
	addUserMail         = addUser.Arg("mail", "The user's email address").Required().String()
	addUserPasswordHash = addUser.Arg("pwd", "The SHA-256 hash of the user's password").Required().String()
	addUserRoles        = addUser.Arg("role", "The user's roles").Required().Strings()

	show = app.Command("show", "Display information about something")

	showRoles     = show.Command("roles", "Displays all roles the server knows about")
	showUsers     = show.Command("users", "Displays all users")
	showFunctions = show.Command("functions", "Displays all functions")
	showMembers   = show.Command("members", "Displays all members")

	serve = app.Command("server", "Starts the API server")
)

func ensureDatabaseConnection() {
	util.Info("Connecting to the database...")
	// Init database info
	dbInfo.Database = cfg.Database.Name
	dbInfo.Password = cfg.Database.Password
	dbInfo.Server = cfg.Database.Host
	dbInfo.Timeout = cfg.Database.Timeout
	dbInfo.Username = cfg.Database.Username

	util.Info(fmt.Sprintf("Mongo DB server is %s, database %s, username %s",
		aurora.BgBlue(cfg.Database.Host).Bold(),
		aurora.BgBlue(cfg.Database.Name).Bold(),
		aurora.BgBlue(cfg.Database.Username).Bold()))

	err := db.Connect(dbInfo)
	if err != nil {
		/*
			util.Fatal(fmt.Sprintf("Unable to connect to '%s/%s' as '%s'!",
				aurora.BgRed(dbInfo.Server).Bold(),
				aurora.BgRed(dbInfo.Database).Bold(),
				aurora.BgRed(dbInfo.Username).Bold()))*/
		tmp := fmt.Sprintf("Unable to connect to %s/%s as %s.", dbInfo.Server, dbInfo.Database, dbInfo.Username)
		util.Fatal(fmt.Sprintf("%s", aurora.BgRed(tmp).Bold()))
		util.Fatal(fmt.Sprintf("%s", aurora.BgRed("Please check your config file for the username/password and make sure mongodb is running")))
		os.Exit(1)
	}
}

func ensureConfig() {
	// Read config
	err := cfg.Read("config.toml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

}

func initAPI() {
	authMiddleware.Key = []byte(cfg.API.Secret)
	util.Info(fmt.Sprintf("Will provide my service at %s:%s", aurora.BgGreen(cfg.Server.Host).Bold(), aurora.BgBlue(cfg.Server.Port).Bold()))
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
	v1.GET("/users", getAllUsersV1)
	v1.POST("/users", createUserV1)
	v1.PUT("/users", notImplemented)
	v1.DELETE("/users", notImplemented)

	v1.GET("/functions", getAllFunctionsV1)
	v1.POST("/functions", createFunctionV1)
	v1.PUT("/functions", notImplemented)
	v1.DELETE("/functions", notImplemented)

	v1.GET("/members", getAllMembersV1)
	v1.POST("/members", createMemberV1)
	v1.PUT("/members", notImplemented)
	v1.DELETE("/members", notImplemented)
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

// main is the application's entry point of course.
func main() {

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if !*appNoBanner {
		banner()
	}
	ensureConfig()
	ensureDatabaseConnection()
	ensureAdmin()

	// yields string

	switch command {
	case serve.FullCommand():
		{

			if *appDebug == false {
				gin.SetMode(gin.ReleaseMode)
			}

			util.Info("Starting the REST API...")
			ensureAdmin()
			initAPI()

			// Start API backend
			// Setup a default gin router with logging
			router := gin.Default()

			// Let the router know about the routes
			setupRoutes(router)

			// Listen address is read from the config file
			listenAddr := strings.Join([]string{cfg.Server.Host, cfg.Server.Port}, ":")

			// Run the router
			router.Run(listenAddr)
		}

	case addUser.FullCommand():
		{
			// Initialize stuff
			util.Info("Adding a user...")

			log.Printf("Username: %s\n", *addUserName)
			log.Printf("Usermail: %s\n", *addUserMail)
			log.Printf("Roles: %s\n", *addUserRoles)
		}

	case addFunction.FullCommand():
		{
			log.Println("Adding a function, apparently...")
			log.Printf("Function name: '%s'\n", *addFunctionName)
			log.Printf("Is board member: '%t'\n", *addFunctionBoard)
		}

	case showRoles.FullCommand():
		{
			util.Info("Displaying all roles")
		}
	case showFunctions.FullCommand():
		{
			util.Info("Displaying all functions")
		}
	case showMembers.FullCommand():
		{
			util.Info("Displaying all members")
		}
	case showUsers.FullCommand():
		{
			util.Info("Displaying all users")
		}

	}

}
