package converter

import (
	"bytes"
)

const (
	upperhex = "0123456789ABCDEF"
)

// unhex is copied from /usr/local/go/src/net/url/url.go
func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

// Normalize a string by removing all special characters and replacing it with
// Q + <hex representation>. Ported fom PHP.
func Normalize(in string) string {
	var b bytes.Buffer
	for i := 0; i < len(in); i++ {
		c := in[i]
		// valid chars are passed through
		if 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || 'A' <= c && c < 'Q' || 'R' <= c && c <= 'Z' {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('Q')
		b.WriteByte(upperhex[c>>4])
		b.WriteByte(upperhex[c&15])
	}
	return b.String()
}

// Denormalize does the opposite of Normalize and converts
// Q<hex> strings back to readable
func Denormalize(in string) string {
	var b bytes.Buffer
	for i := 0; i < len(in); i++ {
		c := in[i]
		if c == 'Q' && i+2 < len(in) {
			u := unhex(in[i+1])<<4 | unhex(in[i+2])
			b.WriteByte(u)
			i += 2
			continue
		}
		b.WriteByte(c)
		continue
	}
	return b.String()
}
