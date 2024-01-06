package converter

import (
	"testing"
)

func Test_GetRandomString(t *testing.T) {
	for _, l := range []int{-1, 0} {
		s := GetRandomString(l)
		if s != "" {
			t.Errorf("Expect an empty string from length %d", l)
		}
	}
	for _, l := range []int{8, 32, 1024} {
		s1 := GetRandomString(l)
		s2 := GetRandomString(l)
		if len(s1) == len(s2) && s1 != s2 {
			continue
		} else {
			t.Errorf("Random error %d %q %q", l, s1, s2)
		}
	}
}
