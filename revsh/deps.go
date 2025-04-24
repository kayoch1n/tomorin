package revsh

import (
	_ "embed"

	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed deps
var depsContent string

var depsTmpl = template.Must(template.New("deps").Parse(depsContent))

var ErrNoPublicKey = errors.New("no public key available")

func Dependencies(config *Config, wr io.Writer) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	var matches []string
	matches, err = filepath.Glob(filepath.Join(dirname, ".ssh", "*.pub"))
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return ErrNoPublicKey
	}
	sort.Strings(matches)

	var publicKey []byte
	publicKey, err = os.ReadFile(matches[0])
	if err != nil {
		return err
	}
	var depends string
	if len(config.Depends) > 0 {
		depends = strings.Join(config.Depends, "\n")
	}

	return depsTmpl.Execute(wr, map[string]interface{}{
		"pubkey":  strings.TrimSpace(string(publicKey)),
		"depends": depends,
	})
}
