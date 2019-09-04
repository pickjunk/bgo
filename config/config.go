package config

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

// Config data
var Config = initConfig()

func initConfig() *gjson.Result {
	// do not depend on bgo/log here, just new a standalone logger for config
	// to prevent circular dependency
	log := zlog.With().Str("component", "bgo.config").Logger()
	log = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

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
