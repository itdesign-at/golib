package macro

import (
	"fmt"
	"testing"

	"github.com/itdesign-at/golib/converter"
)

func testReplace(t *testing.T, expected string, data string, s map[string]interface{}) {

	var m = New().SetMacros(s)
	replaced := m.Replace(data)

	if replaced != expected {
		t.Fatalf("Replace failed:\n\texpected = '%s'\n\t  actual = '%s'", expected, replaced)
	} else {
		//t.Logf("Replace suceed:\n\texpected = '%s'\n\t  actual = '%s'", expected, replaced)
	}

}

func TestMacroReplace(t *testing.T) {
	var macros = map[string]interface{}{
		"Server": "www.univie.ac.at",
		"Client": "myPc",
		"SubTree": map[string]interface{}{
			"Temperature": 16,
			"Flag":        true,
			"AString":     "Hallo",
		},
	}

	const expected = "host=www.univie.ac.atclient=myPcsubTree=16trueHallo"

	var m = New().SetMacros(macros)
	replaced := m.Replace("host={{.Server}}client={{.Client}}subTree={{.SubTree.Temperature}}{{.SubTree.Flag}}{{.SubTree.AString}}")

	if replaced != expected {
		t.Errorf("Expected: %s but got %s", expected, replaced)
	}
}

func TestMacroAppend(t *testing.T) {
	var macros = map[string]interface{}{
		"Name": "theName",
		"Age":  10,
	}

	var m = New().SetMacros(macros)
	m.AppendMacros(map[string]interface{}{
		"Size":   178,
		"Name":   "NewName",
		"tmpDir": "/tmp",
	})

	if m.Macros["Name"] != "NewName" {
		t.Errorf("Expected: NewName but got %s", m.Macros["Name"])
	}

	if len(m.Macros) != 4 {
		t.Errorf("Expected: 4 but got %d macros", len(m.Macros))
	}

}

func TestMacroFunctionEncryptDecrypt(t *testing.T) {
	var m = New().SetMacros(map[string]interface{}{
		"Community": "public",
		"AESSecure": "HTUViWSeWRmTWEOjhENu7/yvi421m+YMUVzD43Fv04UTsQ==",
	})

	cipher := m.Replace(`{{encrypt .Community}}`)
	plain := m.Replace(`x{{decrypt "` + cipher + `"}}`)
	plain1 := m.Replace(`x{{decrypt .AESSecure}}`)

	if plain != "xpublic" {
		t.Error("Expected public as answer")
	}

	if plain1 != "xpublic" {
		t.Error("Expected public as answer")
	}
}

func TestB64Dec(t *testing.T) {
	testReplace(t, "Hello Bar!", "Hello {{.Foo | B64Dec }}!", map[string]interface{}{"Foo": converter.B64Enc("Bar")})
}

func TestB64GenDec(t *testing.T) {
	for _, v := range []string{"b64", "base64"} {
		testReplace(t, "Hello Bar!", "Hello {{.Foo | Dec .Format }}!", map[string]interface{}{"Foo": converter.B64Enc("Bar"), "Format": v})
	}
}

func TestB64Enc(t *testing.T) {
	testReplace(t, "Hello "+converter.B64Enc("Bar")+"!", "Hello {{.Foo | B64Enc }}!", map[string]interface{}{"Foo": "Bar"})
}

func TestB64GenEnc(t *testing.T) {
	for _, v := range []string{"b64", "base64"} {
		testReplace(t, "Hello Bar!", "Hello {{.Foo | Dec .Format }}!", map[string]interface{}{"Foo": converter.B64Enc("Bar"), "Format": v})
	}
}

func TestJsonDec(t *testing.T) {
	var exp map[string]string
	j := `{"foo": "bar", "bar": "foo"}`
	converter.JsonMustDec(j, &exp)
	testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | JsonDec }}!", map[string]interface{}{"Foo": j})
}

func TestJsonGenDec(t *testing.T) {
	var exp map[string]string
	j := `{"foo": "bar", "bar": "foo"}`
	converter.JsonMustDec(j, &exp)
	for _, v := range []string{"json", "application/json", "application/x-json", "text/json", "text/x-json"} {
		testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | Dec .Format }}!", map[string]interface{}{"Foo": j, "Format": v})
	}
}

func TestJsonEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	testReplace(t, "Hello "+converter.JsonMustEnc(data)+"!", "Hello {{.Foo | JsonEnc }}!", map[string]interface{}{"Foo": data})
}

func TestJsonGenEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	for _, v := range []string{"json", "application/json", "application/x-json", "text/json", "text/x-json"} {
		testReplace(t, "Hello "+converter.JsonMustEnc(data)+"!", "Hello {{.Foo | Enc .Format }}!", map[string]interface{}{"Foo": data, "Format": v})
	}
}

func TestYamlDec(t *testing.T) {
	var exp = make(map[string]string)
	y := `
foo: bar
bar: foo
`
	converter.YamlMustDec(y, &exp)
	testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | YamlDec }}!", map[string]interface{}{"Foo": y})
}

func TestYamlGenDec(t *testing.T) {
	var exp = make(map[string]string)
	y := `
foo: bar
bar: foo
`
	converter.YamlMustDec(y, &exp)
	for _, v := range []string{"yaml", "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml"} {
		testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | Dec .Format }}!", map[string]interface{}{"Foo": y, "Format": v})
	}
}

func TestYamlEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	testReplace(t, "Hello "+converter.YamlMustEnc(data)+"!", "Hello {{.Foo | YamlEnc }}!", map[string]interface{}{"Foo": data})
}

func TestYamlGenEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	for _, v := range []string{"yaml", "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml"} {
		testReplace(t, "Hello "+converter.YamlMustEnc(data)+"!", "Hello {{.Foo | Enc .Format }}!", map[string]interface{}{"Foo": data, "Format": v})
	}
}
