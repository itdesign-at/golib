// Package encoding implements utility routines for decoding and encoding yaml and json files.
package encoding

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var ErrExtensionUnsupported = errors.New("file extension not supported")

const defaultFilePermission = 0660

// UnmarshalFile decodes a JSON, YAML or GOB file (depending on its extension) into the "out" value.
// json, yaml and gob are supported
func UnmarshalFile(file string, out any) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	switch filepath.Ext(file) {
	case ".yaml":
		err = yaml.NewDecoder(f).Decode(out)
	case ".json":
		err = json.NewDecoder(f).Decode(out)
	case ".gob":
		err = gob.NewDecoder(f).Decode(out)
	default:
		return ErrExtensionUnsupported
	}

	return err
}

// MarshalFile encodes the value "in" into file as JSON, YAML or GOB, depending on its extension.
// json and yaml are supported
// If the file does not exist, MarshalFile creates it with permissions 0660;
// otherwise UnmarshalFile truncates it before writing, without changing permissions.
func MarshalFile(file string, in any) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, defaultFilePermission)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	switch filepath.Ext(file) {
	case ".yaml":
		err = yaml.NewEncoder(f).Encode(in)
	case ".json":
		err = json.NewEncoder(f).Encode(in)
	case ".gob":
		err = gob.NewEncoder(f).Encode(in)
	default:
		return ErrExtensionUnsupported
	}

	return err
}
