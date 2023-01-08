package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	err := os.MkdirAll(dir, 0777)

	if err != nil {
		return "", err
	}
	return dir, nil
}

// FixTags create valid html element tags by replacing the closing tag </> with a valid html closing tag </health>
// This makes it easier to style the text in the UI
// Example: <health>Health</> -> <health>Health</health>
func FixTags(description string) string {
	descriptionCopy := description
	pattern := regexp.MustCompile("<[^\\/]+>")
	tags := pattern.FindAllString(descriptionCopy, -1)
	for _, tag := range tags {
		endTag := strings.Replace(tag, "<", "</", 1)
		descriptionCopy = strings.Replace(descriptionCopy, "</>", endTag, 1)
	}
	return descriptionCopy
}
