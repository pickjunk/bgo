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

	c, err := os.Open("config.yml")
	if err != nil {
		Log.Warn("config.yml not found")
		return config
	}

	data, err := ioutil.ReadAll(c)
	if err != nil {
		Log.Panic("config.yml not found")
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		Log.Panic("config.yml parse error")
	}

	return config
}
