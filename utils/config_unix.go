// +build darwin linux freebsd !windows

package utils

import (
	"io/ioutil"
	"strconv"

	"github.com/pelletier/go-toml"
)

type config struct {
	Server   server   `toml:"Server"`
	Database database `toml:"MongoDB"`
}

type server struct {
	Host string
	Port int
	Base string
}

type database struct {
	Host string
	Port int
}

// Config Global variable. All settings are there.
var Config = config{}.readConfig()

func (c config) readConfig() config {
	d, _ := ioutil.ReadFile("/etc/short/config.toml")

	d, err := ioutil.ReadFile("config.toml")
	if err != nil {
		// Server
		c.Server.Host = "127.0.0.1"
		c.Server.Port = 8080
		c.Server.Base = "http://" + c.Server.Host + ":" + strconv.Itoa(c.Server.Port) + "/"

		// MongoDB
		c.Database.Host = "127.0.0.1"
		c.Database.Port = 27017

		d, _ := toml.Marshal(c)
		ioutil.WriteFile("config.toml", d, 0666)

		panic(err)
	}

	if err := toml.Unmarshal(d, &c); err != nil {
		panic(err)
	}

	return c
}
