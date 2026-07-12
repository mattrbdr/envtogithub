package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	calls [][]string
}

func (runner *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	runner.calls = append(runner.calls, append([]string{name}, args...))
	return nil, nil
}

type scriptedRunner struct {
	calls   [][]string
	outputs map[string]string
}

func (runner *scriptedRunner) Run(name string, args ...string) ([]byte, error) {
	runner.calls = append(runner.calls, append([]string{name}, args...))
	return []byte(runner.outputs[strings.Join(append([]string{name}, args...), " ")]), nil
}

func TestRunDryRunShowsAReadableSummary(t *testing.T) {
	directory := t.TempDir()
	file := filepath.Join(directory, "env.production.to.github")
	if err := os.WriteFile(file, []byte("API_KEY=abc\nDATABASE_URL=postgres://example\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	runner := &fakeRunner{}
	err := run([]string{"--dry-run", "--color=never", directory}, &stdout, &stderr, runner)
	if err != nil {
		t.Fatalf("run returned an error: %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("dry run called gh: %#v", runner.calls)
	}

	want := "Syncing secrets\n" +
		"production  env.production.to.github\n" +
		"  dry-run  API_KEY  (3 chars)\n" +
		"  dry-run  DATABASE_URL  (18 chars)\n" +
		"\nDry run complete. 2 secrets ready to sync.\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunRejectsUnknownColorMode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"--color=blue"}, &stdout, &stderr, &fakeRunner{})
	if err == nil || !strings.Contains(err.Error(), "color") {
		t.Fatalf("run error = %v, want a color mode error", err)
	}
}

func TestRunHelpDescribesRemoteCommands(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"help"}, &stdout, &stderr, &fakeRunner{})
	if err != nil {
		t.Fatalf("run returned an error: %v", err)
	}
	for _, command := range []string{"etg info", "etg list <environment>", "etg --dry-run"} {
		if !strings.Contains(stdout.String(), command) {
			t.Errorf("help output does not include %q", command)
		}
	}
}

func TestRunInfoListsRemoteEnvironments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	runner := &scriptedRunner{outputs: map[string]string{
		"gh api repos/acme/etg/environments --paginate --jq .environments[].name": "production\nstaging\n",
	}}

	err := run([]string{"info", "--repo", "acme/etg", "--color=never"}, &stdout, &stderr, runner)
	if err != nil {
		t.Fatalf("run returned an error: %v", err)
	}
	want := "Repository  acme/etg\n\nEnvironments\n  production\n  staging\n\n2 environments found.\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunListShowsRemoteSecretNamesOnly(t *testing.T) {
	var stdout, stderr bytes.Buffer
	runner := &scriptedRunner{outputs: map[string]string{
		"gh secret list --env production --repo acme/etg --json name --jq .[].name": "API_KEY\nDATABASE_URL\n",
	}}

	err := run([]string{"list", "production", "--repo", "acme/etg", "--color=never"}, &stdout, &stderr, runner)
	if err != nil {
		t.Fatalf("run returned an error: %v", err)
	}
	want := "Secrets  production\n  API_KEY\n  DATABASE_URL\n\n2 secrets found.\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
