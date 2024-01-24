package converter

import (
	"errors"
	"strings"
)

// Dec decodes string in specific format
// supported format strings are:
// base64: "b64", "base64"
// json: "json", "application/json", "application/x-json", "text/json", "text/x-json"
// yaml: "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml"
func Dec(format string, src string, dst interface{}) (err error) {
	switch strings.ToLower(format) {
	case "b64", "base64":
		s := dst.(*interface{})
		*s, err = B64Dec(src)
	case "json", "application/json", "application/x-json", "text/json", "text/x-json":
		err = JsonDec(src, dst)
	case "yaml", "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml":
		err = YamlDec(src, dst)
	default:
		err = errors.New("unsupported format")
	}
	return
}

// MustDec useit as Dec but panics on error
func MustDec(format string, src string, dst interface{}) {
	if err := Dec(format, src, dst); err != nil {
		panic(err)
	}
}

// Enc encodes string to specific format
// supported format strings are:
// base64: "b64", "base64"
// json: "json", "application/json", "application/x-json", "text/json", "text/x-json"
// yaml: "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml"
func Enc(format string, src interface{}) (dst string, err error) {
	switch strings.ToLower(format) {
	case "b64", "base64":
		dst = B64Enc(src.(string))
	case "json", "application/json", "application/x-json", "text/json", "text/x-json":
		dst, err = JsonEnc(src)
	case "yaml", "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml":
		dst, err = YamlEnc(src)
	default:
		err = errors.New("unsupported format")
	}
	return
}

// MustEnc useit as Enc but panics on error
func MustEnc(format string, src interface{}) (dst string) {
	var err error
	if dst, err = Enc(format, src); err != nil {
		panic(err)
	}
	return
}
