package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger
var once sync.Once

func GetLogger() *logrus.Logger {
	once.Do(func() {
		log = logrus.New()
		log.SetOutput(&lumberjack.Logger{
			Filename:   "./logs/application.log",
			MaxSize:    10,   // Max size in MB
			MaxBackups: 3,    // Max number of old log files to keep
			MaxAge:     28,   // Max age in days to keep a log file
			Compress:   true, // Compress old log files
		})
	})
	return log
}
