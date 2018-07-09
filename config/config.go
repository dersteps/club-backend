package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// The whole config structure as a single struct.
type Config struct {
	Database database
	Server   server
}

// The [database] element in the config
type database struct {
	Host string
	Name string
}

// The [server] element in the config
type server struct {
	Host string
	Port string
}

// Reads the config file and creates a Config from it.
func (c *Config) Read(path string) (err error) {
	log.Printf("Parsing config file %s\n", path)
	_, err2 := toml.DecodeFile(path, &c)
	return err2
}
