package main

import (
	"fmt"
	"os"
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

func process(whole bool, path string, substitutions *[]substitution) {
	fileData, ok := read(path)
	if !ok {
		return
	}

	originalFileData := fileData

	if whole {
		for _, substitution := range *substitutions {
			fileData = substitution.replace(fileData)
		}
	} else {
		fileDataLines := strings.Split(fileData, "\n")
		for _, substitution := range *substitutions {
			for i, line := range fileDataLines {
				fileDataLines[i] = substitution.replace(line)
			}
		}
		fileData = strings.Join(fileDataLines, "\n")
	}

	if fileData != originalFileData {
		diff := diff(originalFileData, fileData, path)
		fmt.Println(diff)
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

	b, err := os.ReadFile(path)
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
