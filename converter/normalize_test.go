package converter

import "testing"

func Test_Normalize(t *testing.T) {
	str := []string{"", "0", "A", "Q", "demo.test.at", "/", "|", "\\"}
	for _, s := range str {
		a := Normalize(s)
		b := Denormalize(a)
		if b != s {
			t.Errorf("Invalid %q %q %q", s, a, b)
		}
	}
}
