package helpers

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/moby/moby/builder/dockerignore"
	"github.com/moby/moby/pkg/archive"
)

func Tarball(dir string) ([]byte, error) {
	sym, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(sym)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(abs, ".dockerignore"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	excludes, err := dockerignore.ReadAll(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	defer os.Chdir(cwd)

	if err := os.Chdir(abs); err != nil {
		return nil, err
	}

	opts := &archive.TarOptions{
		Compression:     archive.Gzip,
		ExcludePatterns: excludes,
		IncludeFiles:    []string{"."},
	}

	r, err := archive.TarWithOptions(abs, opts)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}
