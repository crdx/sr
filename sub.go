package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type substitution struct {
	pattern     *regexp.Regexp
	fixedString string
	replacement string
	isFixed     bool
}

func (self *substitution) replace(s string) string {
	if self.isFixed {
		return strings.ReplaceAll(s, self.fixedString, self.replacement)
	}
	return self.pattern.ReplaceAllString(s, self.replacement)
}

func process(whole bool, path string, subs []substitution) {
	fileData, ok := read(path)
	if !ok {
		return
	}

	originalFileData := fileData

	if whole {
		for _, sub := range subs {
			fileData = sub.replace(fileData)
		}
	} else {
		fileDataLines := strings.Split(fileData, "\n")
		for _, sub := range subs {
			for i, line := range fileDataLines {
				fileDataLines[i] = sub.replace(line)
			}
		}
		fileData = strings.Join(fileDataLines, "\n")
	}

	if fileData != originalFileData {
		fmt.Println(diff(originalFileData, fileData, path))
	}
}

func read(path string) (string, bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", false
	}

	if fileInfo.IsDir() {
		return "", false
	}

	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", false
	}

	s := string(b)

	if !utf8.ValidString(s) {
		return "", false
	}

	return s, true
}

func diff(a, b string, filename string) string {
	edits := myers.ComputeEdits(span.URIFromPath(filename), a, b)
	return fmt.Sprint(gotextdiff.ToUnified(filename, filename, a, edits))
}
