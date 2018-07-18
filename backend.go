package main

import (
	"fmt"
	"log"
	"strings"

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
			Password: cfg.Admin.Password,
			Username: cfg.Admin.Username,
			Role:     "admin",
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

	table := termtables.CreateTable()
	table.AddHeaders("Item", "Value")
	table.AddRow("DB Server", cfg.Database.Host)
	table.AddRow("DB Name", cfg.Database.Name)
	table.AddRow("Server", cfg.Server.Host)
	table.AddRow("Port", cfg.Server.Port)

	fmt.Println(table.Render())

	userDAO.Server = cfg.Database.Host
	userDAO.Database = cfg.Database.Name
	userDAO.Connect()

	ensureAdmin()
}

// setupRoutes will let the router/gin engine know what routes
// there are and what handler functions are mapped
func setupRoutes(router *gin.Engine) {
	// We'll group the routes, so that we are future proof
	router.POST("/login", AuthMiddleware.LoginHandler)

	auth := router.Group("/api/auth")
	auth.Use(AuthMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", AuthMiddleware.RefreshHandler)

	v1 := router.Group("/api/v1")
	v1.Use(AuthMiddleware.MiddlewareFunc())
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
