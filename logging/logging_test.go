package logging

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func Test_LoggingLevel(t *testing.T) {
	seikai := [...]string{
		"",
		`^hoge: [\d/]+ [\d:\.]+ [\w\.:]+ \[INF\] xyz$`,
		`^hoge: [\d/]+ [\d:\.]+ [\w\.:]+ \[WRN\] xyz$`,
		`^hoge: [\d/]+ [\d:\.]+ [\w\.:]+ \[ERR\] xyz$`,
	}

	writer := new(bytes.Buffer)
	l := NewLogger("hoge", INFO)
	l.SetOutput(writer)
	defer l.Close()

	l.Debug("abc")
	result := writer.String()
	result = strings.TrimSpace(result)
	if result != seikai[0] {
		t.Errorf("Wrong output: expected: %s, actual: %s", "<nil>", result)
	}

	writer.Reset()
	l.Info("xyz")
	result = writer.String()
	result = strings.TrimSpace(result)
	re := regexp.MustCompile(seikai[1])
	matched := re.MatchString(result)
	if !matched {
		t.Errorf("Wrong output: expected: %s, actual: %s", seikai[1], result)
	}

	writer.Reset()
	l.Warn("xyz")
	result = writer.String()
	result = strings.TrimSpace(result)
	re = regexp.MustCompile(seikai[2])
	matched = re.MatchString(result)
	if !matched {
		t.Errorf("Wrong output: expected: %s, actual: %s", seikai[2], result)
	}

	writer.Reset()
	l.Error("xyz")
	result = writer.String()
	result = strings.TrimSpace(result)
	re = regexp.MustCompile(seikai[3])
	matched = re.MatchString(result)
	if !matched {
		t.Errorf("Wrong output: expected: %s, actual: %s", seikai[3], result)
	}
}

func Test_OutputConfig(t *testing.T) {
	conf1 := OutputConfig{
		outputType: STDOUT,
	}
	conf2 := FileOutputConfig{
		OutputConfig: OutputConfig{
			outputType: FILE,
		},
		filename: "./test.log",
	}

	config := []interface{}{conf1, conf2}

	l := NewLogger("hoge", INFO)
	err := l.SetOutputByOutputConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	l.Info("hello!")
}

func Test_OutputConfigEmtpty(t *testing.T) {
	config := []interface{}{}

	l := NewLogger("hoge", INFO)
	err := l.SetOutputByOutputConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	l.Close()
}
