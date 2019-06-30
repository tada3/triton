package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

//LogLevel represents the importance of log message.
type LogLevel int

//Values of LogLevel
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// OutputType is type of output.
type OutputType int

// Value of OutputType
const (
	FILE OutputType = iota
	STDOUT
)

// Logger is a wrapper of the standard log.Logger.
type Logger struct {
	logger   *log.Logger
	minLevel LogLevel
	files    []*os.File
}

// OutputConfig is a configuration parameters for output.
type OutputConfig struct {
	outputType OutputType
}

// FileOutputConfig is a specialized type of OutputConfig for file.
type FileOutputConfig struct {
	OutputConfig
	filename string
}

// NewLogger create a new instance of TLogger.
func NewLogger(p string, l LogLevel) *Logger {
	tlogger := Logger{
		logger:   log.New(os.Stdout, p+": ", log.Ldate|log.Lmicroseconds|log.Lshortfile),
		minLevel: l,
		files:    make([]*os.File, 0),
	}
	return &tlogger
}

// SetOutput sets the output by os.Writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// SetOutputByOutputConfig sets the output by the slice of OuputConfig.
func (l *Logger) SetOutputByOutputConfig(configs []interface{}) error {
	writers := make([]io.Writer, 0)
	for _, obj := range configs {
		switch config := obj.(type) {
		case OutputConfig:
			writers = append(writers, os.Stdout)
		case FileOutputConfig:
			file, err := os.OpenFile(config.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
			l.files = append(l.files, file)
			writers = append(writers, file)
		default:
			fmt.Fprintf(os.Stderr, "Invalid config type: %T, %[1]v\n", config)
		}
	}
	writer := io.MultiWriter(writers...)
	l.SetOutput(writer)
	return nil
}

// Close closes the logger.
func (l *Logger) Close() {
	for _, file := range l.files {
		fmt.Printf("Closing file[%s]..\n", file.Name())
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot close the file: %s, %v", file.Name(), err)
		}
	}
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
