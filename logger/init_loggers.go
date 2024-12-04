package logger

import (
	"log"
	"os"
)

var (
	Info *log.Logger
	Err  *log.Logger
)

func InitLoggers() {
	file_name := "logger/server.log"
	file, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open server.log file ", ":", err)
	}
	log.Println("All logs are saved in logger/server.log")

	Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Err = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
