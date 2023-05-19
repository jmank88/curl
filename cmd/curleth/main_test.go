package main

import (
	"embed"
	"io/fs"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"curleth": Main,
	}))
}

//go:embed testdata/*
var f embed.FS

func TestScript(t *testing.T) {
	if err := fs.WalkDir(f, "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		if match, err := fs.Glob(os.DirFS(path), "*.txtar"); err != nil {
			return err
		} else if len(match) == 0 {
			return nil
		}
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			testscript.Run(t, testscript.Params{Dir: path})
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
