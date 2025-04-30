package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{}) // output JSON
	log.SetOutput(os.Stdout)                  // output to console

	// Log to file
	// file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening log file: %v", err)
	// }
	// Add file output to log
	// log.SetOutput(file)

	log.SetLevel(logrus.InfoLevel) // default INFO
	return log
}
