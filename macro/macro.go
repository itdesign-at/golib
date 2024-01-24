package macro

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/itdesign-at/golib/converter"
	"github.com/itdesign-at/golib/crypt"
)

// MacroHandler is used to replace shell macros like e.g. --Host "{{.Hosts.DomainController}}"
type MacroHandler struct {
	sync.Mutex

	// Macros contains all macros which can be replaced with go built in texttemplate engine
	Macros map[string]interface{}

	// macroDirs stores all dirs which could contain macros. Used in
	// LoadMacrosFromDirectory() and Reload...() to be more dynamic.
	macroDirs []string

	// built in functions for template engine
	funcMap template.FuncMap

	// allows to set different delims
	//  e.g.
	// {(.file)}
	delims []string
}

// New is the main entry point for replacing macros
func New() *MacroHandler {
	var macroConfig MacroHandler
	macroConfig.delims = make([]string, 0)
	macroConfig.Macros = make(map[string]interface{})
	macroConfig.funcMap = template.FuncMap{
		"encrypt": func(plain string) string {
			aes := crypt.NewSymmetricEncryption().SetPlainText(plain)
			return aes.GetCypherBase64()
		},
		"decrypt": func(cypher string) string {
			aes := crypt.NewSymmetricEncryption().SetCypherBase64(cypher)
			plain, err := aes.GetPlainText()
			if err == nil {
				return plain
			}
			return cypher
		},
		"B64Dec": converter.B64Dec,
		"B64Enc": converter.B64Enc,
		// example from api.yaml
		//- description: 'GET german locales'
		//  methods: [GET]
		//  url: /uiconfig/locales/de
		//  action: read_file
		//  auth: public
		//  file:
		//    path: config/ui/locales/de.yaml
		//  response:
		//    header:
		//      content-Type: application/json
		//    body: '{{ .Response.Body | YamlDec | JsonEnc }}'
		"JsonDec": func(s string) (v interface{}) {
			err := converter.JsonDec(s, &v)
			if err == nil {
				return
			}
			return s
		},
		"JsonEnc": converter.JsonEnc,
		"YamlDec": func(s string) (v interface{}) {
			err := converter.YamlDec(s, &v)
			if err == nil {
				return
			}
			return s
		},
		"YamlEnc": converter.YamlEnc,
		"Enc":     converter.Enc,
		"Dec": func(format string, src string) (dst interface{}) {
			err := converter.Dec(format, src, &dst)
			if err == nil {
				return
			}
			return src
		},
	}
	return &macroConfig
}

// AppendMacro adds exactly one macro. It overwrites an existing one with the key
// bit keeps others intact.
func (m *MacroHandler) AppendMacro(key string, value interface{}) *MacroHandler {
	m.Lock()
	m.Macros[key] = value
	m.Unlock()
	return m
}

// AppendMacros adds macros to m.Macros. It overwrites existing ones with the same
// key but keeps others intact.
func (m *MacroHandler) AppendMacros(macros map[string]interface{}) *MacroHandler {
	m.Lock()
	for k, v := range macros {
		m.Macros[k] = v
	}
	m.Unlock()
	return m
}

// SetMacros allows to set macros from en external point of view. It overwrites
// all existing macros.
func (m *MacroHandler) SetMacros(macros map[string]interface{}) *MacroHandler {
	m.Lock()
	m.Macros = macros
	m.Unlock()
	return m
}

// LoadMacrosFromFile decodes one file and adds content to macros
func (m *MacroHandler) LoadMacrosFromFile(fileName string) *MacroHandler {

	b, err := os.ReadFile(fileName)
	if err != nil {
		return m
	}

	// data holds data from one single yaml or json file
	var data map[string]interface{}

	switch filepath.Ext(fileName) {
	case ".yaml":
		err = yaml.Unmarshal(b, &data)
	case ".json":
		err = json.Unmarshal(b, &data)
	default:
		return m
	}

	if err == nil && data != nil {
		// mutex.Lock() is done in m.AppendMacros()
		m.AppendMacros(data)
	}

	return m
}

// Replace does the macro replace logic. The "input"
// variable defines a valid template like "{{ .Hostname }}", the string
// returned is the replaced content.
//
//	Supported functions:
//	{{encrypt "plaintext"}}
//	{{decrypt .CypherText}}
func (m *MacroHandler) Replace(input string) string {

	var err error
	var t = template.New("m")
	if len(m.delims) == 2 {
		t.Delims(m.delims[0], m.delims[1])
	}
	t.Funcs(m.funcMap)

	_, err = t.Parse(input)
	if err != nil {
		return input
	}

	var b bytes.Buffer
	m.Lock()
	err = t.Execute(&b, m.Macros)
	m.Unlock()
	if err != nil {
		return input
	}

	return b.String()
}

// Delims allows to set optional delims
// documentation: https://golang.org/pkg/text/template/#Template.Delims
// example left delim = {( and right delim = )}
//
//	e.g. {(.file)}
func (m *MacroHandler) Delims(left, right string) *MacroHandler {
	m.delims = []string{left, right}
	return m
}

// Replace is a simple one line helper for macro replacing
//
//	usage: replaced := macro.Replace(input, macros)
func Replace(input string, macros map[string]interface{}) string {
	return New().SetMacros(macros).Replace(input)
}
