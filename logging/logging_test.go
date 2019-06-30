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
