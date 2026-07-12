// Package envfile reads the env.<environment>.to.github format used by etg.
package envfile

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

var assignment = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)
var filename = regexp.MustCompile(`^env\.([^.]+)\.to\.github$`)

type Entry struct {
	Key   string
	Value string
}

// Parse returns the secret assignments in an env file. Comments and empty
// lines are ignored; any other malformed line is reported with its line number.
func Parse(reader io.Reader) ([]Entry, error) {
	scanner := bufio.NewScanner(reader)
	var entries []Entry
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := assignment.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("line %d: expected KEY=VALUE", lineNumber)
		}
		entries = append(entries, Entry{Key: matches[1], Value: unquote(matches[2])})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}
	return entries, nil
}

// EnvironmentFromFilename extracts the GitHub environment name.
func EnvironmentFromFilename(path string) (string, error) {
	name := filepath.Base(path)
	matches := filename.FindStringSubmatch(name)
	if matches == nil || matches[1] == "" {
		return "", fmt.Errorf("%q must follow env.<environment>.to.github", name)
	}
	return matches[1], nil
}

func unquote(value string) string {
	if len(value) >= 2 {
		first, last := value[0], value[len(value)-1]
		if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
