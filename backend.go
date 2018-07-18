package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/apcera/termtables"

	"github.com/gin-gonic/gin"

	"github.com/dersteps/club-backend/config"
	"github.com/dersteps/club-backend/dao"
)

// Config object for convenient access.
var cfg = config.Config{}
var userDAO = dao.UsersDAO{}

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
}

// setupRoutes will let the router/gin engine know what routes
// there are and what handler functions are mapped
func setupRoutes(router *gin.Engine) {
	// We'll group the routes, so that we are future proof
	v1 := router.Group("/api/v1")
	v1.GET("/users", authenticate, GetAllUsersV1)
	v1.POST("/users", CreateUserV1)
	v1.PUT("/users", NotImplemented)
	v1.DELETE("/users", NotImplemented)
}

// authenticate is currently a dummy middleware that will be utilized
// to authenticate users later on.
func authenticate(c *gin.Context) {
	log.Printf("Authentication for route %s\n", c.Request.URL)
	// c.Abort will kill the middleware chain
	//c.Abort()
	// Send unauthorized status if needed
	//c.String(http.StatusUnauthorized, "Sorry")
	return
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
