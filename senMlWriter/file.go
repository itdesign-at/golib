package senMlWriter

import (
	"bytes"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/mainflux/senml"
)

// SenMl2File writes the senml.Pack to a file with the given name template.
// The template can contain any golang time fields that ar replaced during
// runtime.
// It returns the name of the file written to and nil on success.
func (w *Writer) SenMl2File(fileNameTemplate string) (string, error) {

	// Write to tmp file first to avoid partial writes.
	// Other processes should not read the file while it is being written.

	bt := time.Now()
	bn := ""
	// replace now with the first record time
	if len(w.p.Records) > 0 {
		bt = time.UnixMicro(int64(w.p.Records[0].BaseTime) * 1000000)
		bn = w.p.Records[0].BaseName
	}

	fileName := bt.Format(fileNameTemplate) + ".tmp"

	t := template.New("x")
	_, err := t.Parse(fileName)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var tmp bytes.Buffer
	if err := t.Execute(&tmp, map[string]any{"baseName": strings.TrimSuffix(bn, "/")}); err == nil {
		fileName = tmp.String()
	}

	b, err := senml.Encode(w.p, senml.JSON)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(fileName, b, 0640)
	if err != nil {
		return "", err
	}

	finalName := fileName[:len(fileName)-4] // remove .tmp suffix from final file name
	return finalName, os.Rename(fileName, finalName)
}
