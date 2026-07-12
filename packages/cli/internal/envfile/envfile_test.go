package envfile

import (
	"strings"
	"testing"
)

func TestParseSkipsCommentsAndEmptyLines(t *testing.T) {
	entries, err := Parse(strings.NewReader("\n# a comment\nAPI_KEY=abc\n\n"))
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	if len(entries) != 1 || entries[0] != (Entry{Key: "API_KEY", Value: "abc"}) {
		t.Fatalf("Parse returned %#v, want one API_KEY entry", entries)
	}
}

func TestParsePreservesEqualsAndUnquotesValues(t *testing.T) {
	entries, err := Parse(strings.NewReader("TOKEN=part=one=two\nQUOTED=\"hello world\"\nSINGLE='yes'\n"))
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	want := []Entry{
		{Key: "TOKEN", Value: "part=one=two"},
		{Key: "QUOTED", Value: "hello world"},
		{Key: "SINGLE", Value: "yes"},
	}
	if len(entries) != len(want) {
		t.Fatalf("Parse returned %d entries, want %d", len(entries), len(want))
	}
	for i := range want {
		if entries[i] != want[i] {
			t.Errorf("entry %d = %#v, want %#v", i, entries[i], want[i])
		}
	}
}

func TestParseRejectsMalformedNonCommentLine(t *testing.T) {
	_, err := Parse(strings.NewReader("VALID=value\nnot an env assignment\n"))
	if err == nil {
		t.Fatal("Parse did not reject an invalid assignment")
	}
}

func TestEnvironmentFromFilename(t *testing.T) {
	got, err := EnvironmentFromFilename("/tmp/env.production.to.github")
	if err != nil {
		t.Fatalf("EnvironmentFromFilename returned an error: %v", err)
	}
	if got != "production" {
		t.Fatalf("EnvironmentFromFilename = %q, want production", got)
	}
}

func TestEnvironmentFromFilenameRejectsUnexpectedName(t *testing.T) {
	if _, err := EnvironmentFromFilename(".env.production"); err == nil {
		t.Fatal("EnvironmentFromFilename accepted an unexpected filename")
	}
}
