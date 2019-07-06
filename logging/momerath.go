package logging

var (
	// The Momerath Logger
	tml *Logger
)

func init() {
	tml = NewLogger("M", INFO)
}

func SetLevel(level LogLevel) {
	tml.SetLevel(level)
}

func SetOutputByOutputConfig(configs []interface{}) error {
	return tml.SetOutputByOutputConfig(configs)
}

func NewEntry(name string) *Entry {
	return tml.NewEntry(name)
}
