package config

import (
	"bytes"
	"io"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/go-ini/ini"
)

func ReadFile(dest interface{}, path string) error {
	out, err := ParseFile(path)
	if err != nil {
		return err
	}
	file, err := ini.Load(out)
	if err != nil {
		return err
	}

	return file.MapTo(dest)
}

func ParseFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	out, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	tmp := template.New("")
	// Use sprig function map to enable loading configuration from the environment variable
	// and other capabilities.
	tmp = tmp.Funcs(sprig.FuncMap())
	tmp, err = tmp.Parse(string(out))
	if err != nil {
		return nil, err
	}

	bb := bytes.NewBuffer(nil)
	if err := tmp.Execute(bb, nil); err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}
