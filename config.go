package bgo

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config data
var Config = initConfig()

func initConfig() map[string]interface{} {
	var config map[string]interface{}

	var file string
	if os.Getenv("ENV") == "testing" {
		file = "config_test.yml"
	} else {
		file = "config.yml"
	}

	c, err := os.Open(file)
	if err != nil {
		Log.Warn(file + " not found")
		return config
	}

	data, err := ioutil.ReadAll(c)
	if err != nil {
		Log.Panic(file + " not found")
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		Log.Panic(file + " parse error")
	}

	return config
}
