package logging

import (
	"io"
	"log"
	"os"
)

//Level is log level
type LogLevel int64

//Values of Level
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger is a wrapper of the standard log.Logger.
type Logger struct {
	logger   *log.Logger
	minLevel LogLevel
}

// NewLogger create a new instance of TLogger.
func NewLogger(p string, l LogLevel) *Logger {
	tlogger := Logger{
		logger:   log.New(os.Stdout, p+": ", log.Ldate|log.Lmicroseconds|log.Lshortfile),
		minLevel: l,
	}
	return &tlogger
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// Debug writes the DEBUG level log.
func (l *Logger) Debug(format string, a ...interface{}) {
	if l.minLevel > DEBUG {
		return
	}
	l.logger.Printf("[DBG] "+format, a...)
}

// Info writes the INFO level log.
func (l *Logger) Info(format string, a ...interface{}) {
	if l.minLevel > INFO {
		return
	}
	l.logger.Printf("[INF] "+format, a...)
}

// Warn writes hte WARN level log.
func (l *Logger) Warn(format string, a ...interface{}) {
	if l.minLevel > WARN {
		return
	}
	l.logger.Printf("[WRN] "+format, a...)
}

// Error writes the ERROR level log.
func (l *Logger) Error(format string, a ...interface{}) {
	l.logger.Printf("[ERR] "+format, a...)
}
