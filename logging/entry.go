package logging

// Entry is used to customize some properties of Logger, while sharing the other
// properties with other Entry instances.
type Entry struct {
	name     string
	logger   *Logger
	minLevel LogLevel
}

// Debug writes the DEBUG level log.
func (e *Entry) Debug(format string, a ...interface{}) {
	e.printf(DEBUG, format, a...)
}

func (e *Entry) Info(format string, a ...interface{}) {
	e.printf(INFO, format, a...)
}

func (e *Entry) Warn(format string, a ...interface{}) {
	e.printf(WARN, format, a...)
}

func (e *Entry) Error(format string, a ...interface{}) {
	e.printf(ERROR, format, a...)
}

func (e *Entry) printf(level LogLevel, format string, a ...interface{}) {
	if e.shouldSkip(level) {
		return
	}
	e.logger.printf1(level, format, a...)
}

func (e *Entry) shouldSkip(target LogLevel) bool {
	if e.minLevel == NONE {
		return e.logger.minLevel > target
	}
	return e.minLevel > target
}
