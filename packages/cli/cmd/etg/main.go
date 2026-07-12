package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mattrbdr/envtogithub/packages/cli/internal/envfile"
	"github.com/mattrbdr/envtogithub/packages/cli/internal/ui"
)

type runner interface {
	Run(name string, args ...string) ([]byte, error)
}

type commandRunner struct{}

func (commandRunner) Run(name string, args ...string) ([]byte, error) {
	command := exec.Command(name, args...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, stderr.String())
	}
	return output, nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr, commandRunner{}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer, commands runner) error {
	if len(args) > 0 {
		switch args[0] {
		case "help", "--help", "-h":
			printHelp(stdout)
			return nil
		case "info":
			return runInfo(args[1:], stdout, stderr, commands)
		case "list":
			return runList(args[1:], stdout, stderr, commands)
		case "sync":
			return runSync(args[1:], stdout, stderr, commands)
		}
	}
	return runSync(args, stdout, stderr, commands)
}

func runSync(args []string, stdout, stderr io.Writer, commands runner) error {
	flags := flag.NewFlagSet("etg", flag.ContinueOnError)
	flags.SetOutput(stderr)
	dryRun := flags.Bool("dry-run", false, "print the secrets that would be sent")
	directory := flags.String("dir", ".", "directory containing env.*.to.github files")
	repository := flags.String("repo", "", "GitHub repository as owner/name (defaults to the current repository)")
	colorMode := flags.String("color", "auto", "color output: auto, always, or never")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		return errors.New("provide at most one directory")
	}
	if flags.NArg() == 1 {
		if *directory != "." {
			return errors.New("use either --dir or a directory argument, not both")
		}
		*directory = flags.Arg(0)
	}
	color, err := shouldUseColor(*colorMode, stdout)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(*directory, "env.*.to.github"))
	if err != nil {
		return fmt.Errorf("find env files: %w", err)
	}
	sort.Strings(files)
	if len(files) == 0 {
		return fmt.Errorf("no env.*.to.github files found in %s", *directory)
	}
	output := ui.New(stdout, color)

	if !*dryRun {
		*repository, err = resolveRepository(*repository, commands)
		if err != nil {
			return err
		}
	}

	output.Heading("Syncing secrets")
	secretCount := 0
	for _, file := range files {
		environment, err := envfile.EnvironmentFromFilename(file)
		if err != nil {
			return err
		}
		handle, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("open %s: %w", file, err)
		}
		entries, parseErr := envfile.Parse(handle)
		closeErr := handle.Close()
		if parseErr != nil {
			return fmt.Errorf("parse %s: %w", file, parseErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close %s: %w", file, closeErr)
		}

		output.Environment(environment, filepath.Base(file))
		for _, entry := range entries {
			secretCount++
			if *dryRun {
				output.DryRun(entry.Key, len(entry.Value))
				continue
			}
			if _, err := commands.Run("gh", "secret", "set", entry.Key, "--body", entry.Value, "--env", environment, "--repo", *repository); err != nil {
				return fmt.Errorf("set %s for %s: %w", entry.Key, environment, err)
			}
			output.Success(entry.Key)
		}
	}
	output.Summary(len(files), secretCount, *dryRun)
	return nil
}

func runInfo(args []string, stdout, stderr io.Writer, commands runner) error {
	repository, output, err := remoteOutput("info", args, stdout, stderr, commands)
	if err != nil {
		return err
	}

	raw, err := commands.Run("gh", "api", "repos/"+repository+"/environments", "--paginate", "--jq", ".environments[].name")
	if err != nil {
		return fmt.Errorf("list GitHub environments: %w", err)
	}
	environments := outputLines(raw)
	output.Label("Repository", repository)
	fmt.Fprintln(stdout)
	output.Heading("Environments")
	for _, environment := range environments {
		output.Item(environment)
	}
	fmt.Fprintln(stdout)
	output.Count(len(environments), "environment")
	return nil
}

func runList(args []string, stdout, stderr io.Writer, commands runner) error {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		args = append(args[1:], args[0])
	}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.SetOutput(stderr)
	repository := flags.String("repo", "", "GitHub repository as owner/name (defaults to the current repository)")
	colorMode := flags.String("color", "auto", "color output: auto, always, or never")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: etg list <environment> [--repo owner/repo]")
	}
	color, err := shouldUseColor(*colorMode, stdout)
	if err != nil {
		return err
	}
	resolvedRepository, err := resolveRepository(*repository, commands)
	if err != nil {
		return err
	}

	environment := flags.Arg(0)
	raw, err := commands.Run("gh", "secret", "list", "--env", environment, "--repo", resolvedRepository, "--json", "name", "--jq", ".[].name")
	if err != nil {
		return fmt.Errorf("list secrets for %s: %w", environment, err)
	}
	keys := outputLines(raw)
	output := ui.New(stdout, color)
	output.Label("Secrets", environment)
	for _, key := range keys {
		output.Item(key)
	}
	fmt.Fprintln(stdout)
	output.Count(len(keys), "secret")
	return nil
}

func remoteOutput(command string, args []string, stdout, stderr io.Writer, commands runner) (string, ui.Output, error) {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	repository := flags.String("repo", "", "GitHub repository as owner/name (defaults to the current repository)")
	colorMode := flags.String("color", "auto", "color output: auto, always, or never")
	if err := flags.Parse(args); err != nil {
		return "", ui.Output{}, err
	}
	if flags.NArg() != 0 {
		return "", ui.Output{}, fmt.Errorf("usage: etg %s [--repo owner/repo]", command)
	}
	color, err := shouldUseColor(*colorMode, stdout)
	if err != nil {
		return "", ui.Output{}, err
	}
	resolvedRepository, err := resolveRepository(*repository, commands)
	if err != nil {
		return "", ui.Output{}, err
	}
	return resolvedRepository, ui.New(stdout, color), nil
}

func resolveRepository(repository string, commands runner) (string, error) {
	if repository != "" {
		return repository, nil
	}
	output, err := commands.Run("gh", "repo", "view", "--json", "nameWithOwner", "-q", ".nameWithOwner")
	if err != nil {
		return "", fmt.Errorf("could not detect the GitHub repository (is gh authenticated?): %w", err)
	}
	repository = string(bytes.TrimSpace(output))
	if repository == "" {
		return "", errors.New("GitHub repository detection returned no repository")
	}
	return repository, nil
}

func outputLines(raw []byte) []string {
	var lines []string
	for _, line := range strings.Split(string(raw), "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

func printHelp(writer io.Writer) {
	fmt.Fprint(writer, `etg — sync local env files with GitHub environment secrets

Usage:
  etg [--dry-run] [--dir directory] [--repo owner/repo]
  etg sync [--dry-run] [--dir directory] [--repo owner/repo]
  etg info [--repo owner/repo]
  etg list <environment> [--repo owner/repo]

Commands:
  info              List environments available in the GitHub repository.
  list              List secret names in one GitHub environment.
  sync              Sync env.*.to.github files (the default command).
  help              Show this help message.

Options:
  --dry-run         Preview secrets without sending values.
  --color MODE      Color output: auto, always, or never.
  --repo REPO       GitHub repository (defaults to the current repository).

Example:
  etg --dry-run
`)
}

func shouldUseColor(mode string, writer io.Writer) (bool, error) {
	switch strings.ToLower(mode) {
	case "always":
		return true, nil
	case "never":
		return false, nil
	case "auto":
		file, ok := writer.(*os.File)
		if !ok || os.Getenv("NO_COLOR") != "" {
			return false, nil
		}
		info, err := file.Stat()
		if err != nil {
			return false, nil
		}
		return info.Mode()&os.ModeCharDevice != 0, nil
	default:
		return false, fmt.Errorf("invalid color mode %q (use auto, always, or never)", mode)
	}
}
