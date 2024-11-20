package main

import (
	"fmt"
	"os"
	"src/post_relay/cmd"
	"src/post_relay/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type App struct {
}

func (*App) start() {
	// log := logrus.New()
	// log.SetOutput(&lumberjack.Logger{
	// 	Filename:   "./logs/application.log",
	// 	MaxSize:    10,   // Max size in MB
	// 	MaxBackups: 3,    // Max number of old log files to keep
	// 	MaxAge:     28,   // Max age in days to keep a log file
	// 	Compress:   true, // Compress old log files
	// })

	log := logrus.New()
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/application.log",
		MaxSize:    10,   // Max size in MB
		MaxBackups: 3,    // Max number of old log files to keep
		MaxAge:     28,   // Max age in days to keep a log file
		Compress:   true, // Compress old log files
	})

	fmt.Printf("Version: %s Commit: %s \n", config.Version, config.Commit)
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}

	log.Info("Application started")
}

func main() {

	App := &App{}

	go App.start()

	select {}
}
