package myLogger

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "Logger: ", log.LstdFlags|log.Lshortfile)
}
