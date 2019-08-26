package bgo

import (
	"os"

	stack "github.com/Gurpartap/logrus-stack"
	tf "github.com/pickjunk/bgo/text_formatter"
	logrus "github.com/sirupsen/logrus"
)

// Log bgo Logger
var Log = initLogger()

// inner log for bgo core
var log = Log.WithField("prefix", "bgo")

func initLogger() *logrus.Logger {
	l := logrus.New()

	l.SetFormatter(&tf.TextFormatter{
		MultilineFields: []string{"schema", "stack", "sql"},
		FullTimestamp:   true,
	})

	if os.Getenv("ENV") == "production" {
		l.SetLevel(logrus.InfoLevel)
	} else {
		l.SetLevel(logrus.DebugLevel)
	}

	l.SetOutput(os.Stdout)

	callerLevels := []logrus.Level{}
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
	stackHook := stack.NewHook(callerLevels, stackLevels)
	l.AddHook(stackHook)

	return l
}
