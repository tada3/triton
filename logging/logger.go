package logging

import (
	"log"
	"os"
)

// TLogger is a wrapper of the standard log.Logger.
type TLogger struct {
	logger *log.Logger
}

// NewTLogger create a new instance of TLogger.
func NewTLogger(p string) *TLogger {
	tlogger := TLogger{
		logger: log.New(os.Stdout, p, log.Lmicroseconds|log.Lshortfile),
	}
	return &tlogger
}
