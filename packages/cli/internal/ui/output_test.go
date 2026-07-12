package ui

import (
	"bytes"
	"testing"
)

func TestOutputUsesColorForTerminalStatus(t *testing.T) {
	var buffer bytes.Buffer
	output := New(&buffer, true)

	output.Heading("Syncing secrets")
	output.Environment("production", "env.production.to.github")
	output.DryRun("API_KEY", 9)
	output.Summary(1, 1, true)

	want := "\x1b[1mSyncing secrets\x1b[0m\n" +
		"\x1b[36mproduction\x1b[0m  \x1b[2menv.production.to.github\x1b[0m\n" +
		"  \x1b[33mdry-run\x1b[0m  API_KEY  \x1b[2m(9 chars)\x1b[0m\n" +
		"\n\x1b[33mDry run complete.\x1b[0m 1 secret ready to sync.\n"
	if got := buffer.String(); got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestOutputOmitsEscapeCodesWhenColorIsDisabled(t *testing.T) {
	var buffer bytes.Buffer
	output := New(&buffer, false)

	output.Environment("staging", "env.staging.to.github")
	output.Success("DATABASE_URL")
	output.Summary(1, 1, false)

	want := "staging  env.staging.to.github\n" +
		"  set      DATABASE_URL\n" +
		"\nDone. 1 secret synced.\n"
	if got := buffer.String(); got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
