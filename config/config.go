package config

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pickjunk/zerolog"
	"github.com/tidwall/gjson"
)

var config = initConfig()

func initConfig() *gjson.Result {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	// do not depend on bgo/log here, just new a standalone logger for config
	// to prevent circular dependency
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	log = log.With().Str("component", "bgo.config").Logger()

	var config gjson.Result

	file := "config.yml"
	if os.Getenv("ENV") != "production" {
		if _, err := os.Stat("config_test.yml"); err == nil {
			file = "config_test.yml"
		}
	}

	c, err := os.Open(file)
	if err != nil {
		log.Warn().Str("file", file).Msg("config not found")
		return &config
	}

	data, err := ioutil.ReadAll(c)
	if err != nil {
		log.Panic().Err(err).Str("file", file).Msg("config read")
	}

	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		log.Panic().Err(err).Str("file", file).Msg("yaml to json")
	}

	config = gjson.ParseBytes(json)

	log.Info().Str("file", file).Msg("config loaded")

	return &config
}

// Get config by path
func Get(path string) gjson.Result {
	return config.Get(path)
}
