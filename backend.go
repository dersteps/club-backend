package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dersteps/club-backend/model"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	//"github.com/dersteps/club-backend/model"
	"github.com/dersteps/club-backend/config"
	"github.com/dersteps/club-backend/dao"
)

// Config object for convenient access.
var cfg = config.Config{}
var userDAO = dao.UsersDAO{}

func AllUsersEndpoint(w http.ResponseWriter, r *http.Request) {
	users, err := userDAO.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
}

func CreateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = bson.NewObjectId()
	if err := userDAO.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, user)

}

func UpdateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Not yet implemented!")
}

func DeleteUserEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Not yet implemented!")
}

func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Not yet implemented!")
}

func APIUsageEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Tell caller about API usage")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func setupRoutes(r *mux.Router) {
	r.HandleFunc("/", APIUsageEndpoint).Methods("GET")
	r.HandleFunc("/users", AllUsersEndpoint).Methods("GET")
	r.HandleFunc("/users", CreateUserEndpoint).Methods("POST")
	r.HandleFunc("/users", UpdateUserEndpoint).Methods("PUT")
	r.HandleFunc("/users", DeleteUserEndpoint).Methods("DELETE")
	r.HandleFunc("/users/{id}", FindUserEndpoint).Methods("GET")
}

// Init is a special function called by go automagically.
// Initializes basic stuff.
func init() {
	// Read config
	err := cfg.Read("config.toml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Printf("Club backend config\n-------------------\n")
	log.Printf("DB Server: %s\n", cfg.Database.Host)
	log.Printf("DB Name: %s\n", cfg.Database.Name)
	log.Printf("Server listen: %s:%s\n", cfg.Server.Host, cfg.Server.Port)

	userDAO.Server = cfg.Database.Host
	userDAO.Database = cfg.Database.Name
	userDAO.Connect()
}

func main() {
	r := mux.NewRouter()
	setupRoutes(r)
	parts := []string{cfg.Server.Host, cfg.Server.Port}
	listenAddr := strings.Join(parts, ":")

	if err := http.ListenAndServe(listenAddr, handlers.LoggingHandler(os.Stdout, r)); err != nil {
		log.Fatal(err)
		panic(err)
	}

}
