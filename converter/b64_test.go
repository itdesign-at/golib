package converter

import "testing"

func TestB64Dec(t *testing.T) {
	for _, plain := range []string{"", "Hello World"} {
		enc := B64Enc(plain)
		dec, err := B64Dec(enc)
		if err != nil {
			t.Error(err)
		}
		if plain != dec {
			t.Errorf("mismatch between plain '%s' and decoded '%s'", plain, dec)
		}
	}
}
