package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	queueSize = 10
)

//LogLevel represents the importance of log message.
type LogLevel int

//Values of LogLevel
const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	NONE
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
	name           string
	delegate       *log.Logger
	minLevel       LogLevel
	files          []*os.File
	nonFileWriters []io.Writer
	queue          chan string
	done           chan bool
}

// Entry is used to customize some properties of Logger, while sharing the other
// properties with other Entry instances.
type Entry struct {
	name     string
	logger   *Logger
	minLevel LogLevel
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

// NewLogger create a new instance of Logger.
func NewLogger(name string, level LogLevel) *Logger {
	logger := Logger{
		name:     name,
		delegate: log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds),
		minLevel: level,
		files:    make([]*os.File, 0),
		queue:    make(chan string, queueSize),
		done:     make(chan bool),
	}
	go func() {
		for {
			msg, ok := <-logger.queue
			if !ok {
				fmt.Printf("LOG[%s] queue has been closed.\n", logger.name)
				break
			}
			logger.delegate.Output(0, msg)
			fmt.Printf("LOG[%s] queue: %d\n", logger.name, len(logger.queue))
		}
		logger.done <- true
	}()
	return &logger
}

// NewEntry creates a new instance of Entry
func (l *Logger) NewEntry(name string) *Entry {
	return &Entry{
		name:     name,
		logger:   l,
		minLevel: NONE,
	}
}

// SetOutput sets the output by os.Writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.delegate.SetOutput(w)
}

// SetOutputByOutputConfig sets the output by the slice of OuputConfig.
func (l *Logger) SetOutputByOutputConfig(configs []interface{}) error {
	writers := make([]io.Writer, 0)
	nonFileWriters := make([]io.Writer, 0)
	fileWriters := make([]*os.File, 0)
	for _, obj := range configs {
		switch config := obj.(type) {
		case OutputConfig:
			writers = append(writers, os.Stdout)
			nonFileWriters = append(nonFileWriters, os.Stdout)
		case FileOutputConfig:
			file, err := os.OpenFile(config.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
			writers = append(writers, file)
			fileWriters = append(fileWriters, file)
		default:
			fmt.Fprintf(os.Stderr, "LOG[%s] Invalid config type: %T, %[2]v\n", l.name, config)
		}
	}
	writer := io.MultiWriter(writers...)
	l.SetOutput(writer)
	l.nonFileWriters = nonFileWriters
	l.files = fileWriters
	return nil
}

// Close closes the logger.
func (l *Logger) Close() {
	fmt.Printf("LOG[%s] Closing logger..\n", l.name)
	for _, file := range l.files {
		fmt.Printf("LOG[%s] Closing file[%s]..\n", l.name, file.Name())
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot close the file[%s]: %v", file.Name(), err)
		}
	}
	fmt.Printf("LOG[%s] Closing queue..\n", l.name)
	close(l.queue)
	<-l.done
	close(l.done)
	fmt.Printf("LOG[%s] Closed.\n", l.name)
}

// Rotate performs daily file rotation.
func (l *Logger) Rotate() {
	if len(l.files) == 0 {
		return
	}

	suffix := time.Now().AddDate(0, 0, -1).Format("2016-01-02")

	// Update Logger.files
	newfiles := make([]*os.File, 0)
	for _, f := range l.files {
		newfile, err := rotateFile(f, suffix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to rotate the file[%s]: %v\n", f.Name(), err)
			continue
		}
		newfiles = append(newfiles, newfile)
	}
	l.files = newfiles

	// Update output
	writers := make([]io.Writer, 0)
	for _, file := range l.files {
		writers = append(writers, file)
	}
	for _, nonfile := range l.nonFileWriters {
		writers = append(writers, nonfile)
	}
	multiWriter := io.MultiWriter(writers...)
	l.SetOutput(multiWriter)
}

func rotateFile(f *os.File, suffix string) (*os.File, error) {
	// 1. Close current file
	err := f.Close()
	if err != nil {
		return nil, err
	}

	// 2. Rename current file
	filename := f.Name()
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if fileinfo.Size() > 0 {
		err = os.Rename(filename, filename+"."+suffix)
		if err != nil {
			return nil, err
		}
	}

	// 3. Open new file
	newfile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return newfile, nil
}

// Debug writes the DEBUG level log.
func (l *Logger) Debug(format string, a ...interface{}) {
	if l.minLevel > DEBUG {
		return
	}
	l.delegate.Printf("[DBG] "+format, a...)
}

// Info writes the INFO level log.
func (l *Logger) Info(format string, a ...interface{}) {
	if l.minLevel > INFO {
		return
	}
	l.delegate.Printf("[INF] "+format, a...)
}

// Warn writes hte WARN level log.
func (l *Logger) Warn(format string, a ...interface{}) {
	if l.minLevel > WARN {
		return
	}
	l.delegate.Printf("[WRN] "+format, a...)
}

// Error writes the ERROR level log.
func (l *Logger) Error(format string, a ...interface{}) {
	l.delegate.Printf("[ERR] "+format, a...)
}

// Debug writes the DEBUG level log.
func (e *Entry) Debug(format string, a ...interface{}) {
	if e.shouldSkip(DEBUG) {
		return
	}

	msg := fmt.Sprintf("[DBG] "+format, a...)
	e.logger.queue <- msg
}

func (e *Entry) shouldSkip(target LogLevel) bool {
	if e.minLevel == NONE {
		return e.logger.minLevel > target
	}
	return e.minLevel > target
}
