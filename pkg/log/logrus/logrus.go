package logrus

import (
	"fmt"
	"os"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel string
}

type LogHandler struct {
	Logger *logrus.Logger
	Cfg    *Config
}

type Fields = logrus.Fields

func (log *LogHandler) Init() (err error) {

	l := logrus.New()

	//debug, info, warning, error, fatal, panic
	level, err := logrus.ParseLevel(log.Cfg.LogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l.SetLevel(level)

	return err
}

func (log *LogHandler) SetFluentdFormatter() {
	log.Logger.Formatter = &joonix.FluentdFormatter{}
}

func (log *LogHandler) Debug(msg string, fields Fields) {
	log.Logger.WithFields(fields).Debug(msg)
}

func (log *LogHandler) Info(msg string, fields Fields) {
	log.Logger.WithFields(fields).Info(msg)
}

func (log *LogHandler) Warning(msg string, fields Fields) {
	log.Logger.WithFields(fields).Warning(msg)
}

func (log *LogHandler) Error(msg string, fields Fields) {
	log.Logger.WithFields(fields).Error(msg)
}

func (log *LogHandler) Panic(msg string, fields Fields) {
	log.Logger.WithFields(fields).Panic(msg)
}

func (log *LogHandler) Fatal(msg string, fields Fields) {
	log.Logger.WithFields(fields).Fatal(msg)
}
