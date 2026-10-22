package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseOptionsReadOnly(t *testing.T) {
	var output bytes.Buffer
	opts, args, err := parseOptions([]string{"--read-only", "/tmp/change"}, &output)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if !opts.readOnly {
		t.Error("expected --read-only to enable read-only mode")
	}
	if len(args) != 1 || args[0] != "/tmp/change" {
		t.Fatalf("expected path argument to be preserved, got %v", args)
	}
}

func TestParseOptionsConfigPath(t *testing.T) {
	var output bytes.Buffer
	opts, _, err := parseOptions([]string{"--config", "/tmp/dossier.toml"}, &output)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if opts.configPath != "/tmp/dossier.toml" {
		t.Fatalf("configPath = %q", opts.configPath)
	}
}

func TestParseOptionsTracksExplicitKeyStyle(t *testing.T) {
	var output bytes.Buffer
	opts, _, err := parseOptions([]string{"--keystyle", "nvim"}, &output)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if !opts.keyStyleSet || opts.keyStyleName != "nvim" {
		t.Fatalf("expected explicit nvim keystyle, got %+v", opts)
	}
}

func TestParseOptionsTracksExplicitTheme(t *testing.T) {
	var output bytes.Buffer
	opts, _, err := parseOptions([]string{"--theme", "light"}, &output)
	if err != nil {
		t.Fatalf("parseOptions: %v", err)
	}
	if !opts.themeSet || opts.themeName != "light" {
		t.Fatalf("expected explicit light theme, got %+v", opts)
	}
}

func TestParseOptionsHelpAliases(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		t.Run(arg, func(t *testing.T) {
			var output bytes.Buffer
			opts, _, err := parseOptions([]string{arg}, &output)
			if err != nil {
				t.Fatalf("parseOptions: %v", err)
			}
			if !opts.showHelp {
				t.Errorf("expected %s to request help", arg)
			}
		})
	}
}

func TestUsageUsesConventionalFlagSpelling(t *testing.T) {
	var output bytes.Buffer
	writeUsage(&output)
	usage := output.String()

	for _, want := range []string{"-h, --help", "--config <path>", "--keystyle <name>", "--read-only", "--theme <name>", "gruvbox-light-soft", "--version"} {
		if !strings.Contains(usage, want) {
			t.Errorf("expected usage to contain %q, got:\n%s", want, usage)
		}
	}
	if strings.Contains(usage, "\n  -help") {
		t.Errorf("usage should not advertise single-dash long help:\n%s", usage)
	}
}
