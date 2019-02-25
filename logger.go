package bgo

import (
	"os"

	stack "github.com/Gurpartap/logrus-stack"
	tf "github.com/pickjunk/text_formatter"
	log "github.com/sirupsen/logrus"
)

// Logger struct
type Logger struct {
	*log.Logger
}

// Log instance
var Log = initLogger()

func initLogger() *Logger {
	if os.Getenv("ENV") == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&tf.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}

	log.SetOutput(os.Stdout)

	callerLevels := []log.Level{}
	stackLevels := []log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel}
	stackHook := stack.NewHook(callerLevels, stackLevels)
	log.AddHook(stackHook)

	return &Logger{
		log.StandardLogger(),
	}
}
