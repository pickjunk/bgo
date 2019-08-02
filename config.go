package bgo

import (
	"io/ioutil"
	"os"

	sentry "github.com/onrik/logrus/sentry"
	logrus "github.com/sirupsen/logrus"
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
		log.Warn(file + " not found")
		return config
	}

	data, err := ioutil.ReadAll(c)
	if err != nil {
		log.Panic(file + " not found")
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Panic(file + " parse error")
	}

	// sentry
	if sentryDSN, ok := config["sentry"].(string); ok {
		sentryHook, err := sentry.NewHook(sentry.Options{
			Dsn: sentryDSN,
		}, logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel)
		if err != nil {
			log.Fatal(err)
		}
		Log.AddHook(sentryHook)
	}

	return config
}
