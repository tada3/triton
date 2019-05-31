package util

import (
	"testing"
)
func Test_ToString(t *testing.T) {
	m := map[int]string{
		1:  "ついたち",
		10: "とおか",
		29: "29日",
	}
	for k, v := range m {
		s := GetDayStr(k)
		if s != v {
			t.Errorf("Invalid result. expect: %s, actual: %s", v, s)
		}
	}
}
