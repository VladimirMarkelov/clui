package clui

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)
//InitLogger --
func InitLogger() {
	file, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("----------------------------------")
}
//Logger --
func Logger() *log.Logger {
	if logger == nil {
		InitLogger()
	}

	return logger
}
