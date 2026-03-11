package config

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.SetLevel(logrus.InfoLevel)

	logDir := "logs"
	logFile := filepath.Join(logDir, "app.log")

	// create logs folder if not exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("failed to create log directory")
	}

	file, err := os.OpenFile(
		logFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic("failed to open log file")
	}

	logrus.SetOutput(file)
}
