package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// The whole config structure as a single struct.
type Config struct {
	Database database
	Server   server
	Admin    admin
	API      api
}

// The [database] element in the config
type database struct {
	Host     string
	Name     string
	Timeout  int
	Username string
	Password string
}

// The [server] element in the config
type server struct {
	Host string
	Port string
}

type admin struct {
	Username string
	Mail     string
	Password string
}

type api struct {
	Secret string
}

// Reads the config file and creates a Config from it.
func (c *Config) Read(path string) (err error) {
	log.Printf("Parsing config file %s\n", path)
	_, err2 := toml.DecodeFile(path, &c)
	return err2
}

/*[admin]
username="administrator"
email="stefan.matyba@googlemail.com"
password="ksljdfosdjfisdjfiosdf" #<- hash!*/
