package logger

import (
	"log"
	"os"
)

var (
	ILog *log.Logger
	ELog *log.Logger
)

func init() {
	ILog = log.New(os.Stdout, "[info]", log.LstdFlags|log.LUTC)
	ELog = log.New(os.Stderr, "[error]", log.LstdFlags|log.LUTC)
}
