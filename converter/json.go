package converter

import (
	"encoding/json"
)

// JsonDec decodes string in json format
func JsonDec(s string, v interface{}) error {
	// return json.Unmarshal([]byte(s), &v) //so nicht ist ja schon ein Pointer
	return json.Unmarshal([]byte(s), v)
}

// JsonMustDec useit as JsonDec but panics on error
func JsonMustDec(s string, v interface{}) {
	if err := JsonDec(s, v); err != nil {
		panic(err)
	}
}

// JsonEnc encodes value to json string
func JsonEnc(v interface{}) (s string, err error) {
	var b []byte
	b, err = json.Marshal(v)
	s = string(b)
	return
}

// JsonMustEnc useit as JsonEnc but panics on error
func JsonMustEnc(v interface{}) string {
	s, err := JsonEnc(v)
	if err != nil {
		panic(err)
	}
	return s
}
