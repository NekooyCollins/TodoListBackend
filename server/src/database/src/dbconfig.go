package database

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

// Configuration loads db config file
type Configuration struct {
	DBUsername string
	DBPasswd   string
	DBPort     string
	DBHost     string
	DBName     string
}

// GetDBConfig reverse sequence from config file
func GetDBConfig(params ...string) (Configuration, error) {
	configuration := Configuration{}
	env := "dev"
	if len(params) > 0 {
		env = params[0]
	}
	fileName := fmt.Sprintf("./%s_config.json", env)
	err := gonfig.GetConf(fileName, &configuration)
	return configuration, err
}
