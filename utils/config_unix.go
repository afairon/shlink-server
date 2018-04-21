// +build !windows

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var Conf = New()

// Name of the config
const configFile = "shlink.yml"

// Config is a yaml structure for
// managing config.
type Config struct {
	Database db     `yaml:"database"`
	Log      log    `yaml:"log"`
	Server   server `yaml:"server"`
}

type db struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	DB   string `yaml:"db"`
}

type log struct {
	Filename   string `yaml:"logName"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
}

type server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Base string `yaml:"base"`
}

func init() {
	Conf.ReadConfig()
}

// exists return whether the given file or directory exists or not
func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return false, err
	}

	return true, nil
}

// New creates a pointer to the configuration struct
func New() *Config {
	return &Config{}
}

// ReadConfig checks if the configuration file exists.
// If it exists than unmarshal it to a structure.
// Else, create one with default settings and exit(0).
func (c *Config) ReadConfig() {
	if ok, _ := exists("/etc/shlink/" + configFile); ok {
		f, err := ioutil.ReadFile("/etc/shlink/" + configFile)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(f, &c); err != nil {
			panic(err)
		}
	}

	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		// Database
		c.Database.Host = "127.0.0.1" // MongoDB address
		c.Database.Port = "27017"     // MongoDB port
		c.Database.DB = "shlink"      // MongoDB db name

		// Logger
		c.Log.Filename = "logs/shlink.log" // Log filename
		c.Log.MaxSize = 10                 // Size in MB
		c.Log.MaxBackups = 2               // Length in days
		c.Log.MaxAge = 7                   // Duration in days

		// Server
		c.Server.Host = "127.0.0.1"                                     // Listening address
		c.Server.Port = "8080"                                          // Listening port
		c.Server.Base = "http://" + c.Server.Host + ":" + c.Server.Port // External address

		f, _ := yaml.Marshal(c)
		ioutil.WriteFile(configFile, f, 0664)

		fmt.Println("Config file not found.")
		fmt.Printf("Creating %s config file.\n", configFile)
		os.Exit(0)
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		panic(err)
	}

	if !strings.HasSuffix(c.Server.Base, "/") {
		c.Server.Base += "/"
	}
}
