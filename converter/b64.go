package converter

import "encoding/base64"

// B64Dec decodes string in base64 format
func B64Dec(b64s string) (s string, err error) {
	var b []byte
	b, err = base64.StdEncoding.DecodeString(b64s)
	s = string(b)
	return
}

// B64MustDec useit as B64dec but panics on error
func B64MustDec(b64s string) string {
	s, err := B64Dec(b64s)
	if err != nil {
		panic(err)
	}
	return s
}

// B64Enc encodes string to base64 format
func B64Enc(s string) (b64s string) {
	b64s = base64.StdEncoding.EncodeToString([]byte(s))
	return
}
