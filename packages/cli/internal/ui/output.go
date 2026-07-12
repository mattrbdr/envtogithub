// Package ui formats etg's terminal output.
package ui

import (
	"fmt"
	"io"
)

const (
	reset  = "\x1b[0m"
	bold   = "\x1b[1m"
	dim    = "\x1b[2m"
	cyan   = "\x1b[36m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
)

type Output struct {
	writer io.Writer
	color  bool
}

func New(writer io.Writer, color bool) Output {
	return Output{writer: writer, color: color}
}

func (output Output) Heading(message string) {
	fmt.Fprintln(output.writer, output.style(bold, message))
}

func (output Output) Label(label, value string) {
	fmt.Fprintf(output.writer, "%s  %s\n", output.style(cyan, label), value)
}

func (output Output) Item(value string) {
	fmt.Fprintf(output.writer, "  %s\n", value)
}

func (output Output) Count(count int, noun string) {
	if count != 1 {
		noun += "s"
	}
	fmt.Fprintf(output.writer, "%d %s found.\n", count, noun)
}

func (output Output) Environment(name, filename string) {
	fmt.Fprintf(output.writer, "%s  %s\n", output.style(cyan, name), output.style(dim, filename))
}

func (output Output) DryRun(key string, length int) {
	fmt.Fprintf(output.writer, "  %s  %s  %s\n", output.status(yellow, "dry-run"), key, output.style(dim, fmt.Sprintf("(%d chars)", length)))
}

func (output Output) Success(key string) {
	fmt.Fprintf(output.writer, "  %s  %s\n", output.status(green, "set"), key)
}

func (output Output) Summary(_ int, secrets int, dryRun bool) {
	word := "secrets"
	if secrets == 1 {
		word = "secret"
	}
	if dryRun {
		fmt.Fprintf(output.writer, "\n%s %d %s ready to sync.\n", output.style(yellow, "Dry run complete."), secrets, word)
		return
	}
	fmt.Fprintf(output.writer, "\n%s %d %s synced.\n", output.style(green, "Done."), secrets, word)
}

func (output Output) style(code, value string) string {
	if !output.color {
		return value
	}
	return code + value + reset
}

func (output Output) status(code, value string) string {
	return output.style(code, fmt.Sprintf("%-7s", value))
}
