package util

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const ParsedDataDir = "parsedData"

// Check if the error exists (err != nil). If so the error and relating data is printed and the program panics.
// Params:
// err: error - The error itself
// data: ...any - all data objects to print
func Check(err error, data ...any) {
	if err == nil {
		return
	}
	fmt.Printf("[Critical] Error{%s} Data [%+v\n]", err, data)
	panic(err)

}

func WriteInfo[T any](file string, infos []T) error {
	dir, err := CreateDir(".", ParsedDataDir)

	f, err := os.Create(filepath.Join(dir, file))
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	err = enc.Encode(infos)
	if err != nil {
		return err
	}

	return f.Close()
}

func CreateDir(paths ...string) (string, error) {
	dir := filepath.Join(paths...)
	err := os.MkdirAll(dir, fs.ModeDir)

	if err != nil {
		return "", err
	}
	return dir, nil
}

func Ternary[T any](check bool, a, b T) T {
	if check {
		return a
	} else {
		return b
	}
}