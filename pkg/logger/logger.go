package logger

import (
	"github.com/fir1/story-title/pkg/config"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	TextFormater = "text"
	JSONFormater = "json"
)

const (
	Debug   = "debug"
	Info    = "info"
	Warning = "warning"
	Error   = "error"
)

func NewTextLogger(cnf config.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(cnf.Log.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	switch strings.ToLower(cnf.Log.Formater) {
	case TextFormater:
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.999999999",
		})

	case JSONFormater:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.999999999",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			TimestampFormat: "2006-01-02 15:04:05.999999999",
			FullTimestamp:   true,
		})
	}
	return logger, nil
}
