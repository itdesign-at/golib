package converter

import "gopkg.in/yaml.v3"

// YamlDec decodes string in yaml format
func YamlDec(s string, v interface{}) error {
	// return yaml.Unmarshal([]byte(s), &v) // so nicht ist ja schon ein Pointer
	return yaml.Unmarshal([]byte(s), v)
}

// YamlMustDec useit as YamlDec but panics on error
func YamlMustDec(s string, v interface{}) {
	if err := YamlDec(s, v); err != nil {
		panic(err)
	}
}

// YamlEnc encodes value to yaml string
func YamlEnc(v interface{}) (s string, err error) {
	var b []byte
	b, err = yaml.Marshal(v)
	s = string(b)
	return
}

// YamlMustEnc useit as YamlEnc but panics on error
func YamlMustEnc(v interface{}) string {
	s, err := YamlEnc(v)
	if err != nil {
		panic(err)
	}
	return s
}
