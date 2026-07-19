## ADDED Requirements

### Requirement: CLI flag parsing uses Go `flag` package
The system SHALL parse CLI arguments using Go's standard `flag` package, accepting a `--theme` string flag (default: `"dark"`), `--version` boolean flag, `--help` boolean flag, and an optional positional argument for the change path. The manual `os.Args` inspection in `cmd/dossier/main.go` SHALL be replaced with `flag`-based parsing.

#### Scenario: `--theme` flag parsed alongside other flags
- **WHEN** the user runs `dossier --theme light --version`
- **THEN** the version is printed and the program exits (theme flag is accepted but not applied since `--version` exits immediately)

#### Scenario: Positional argument after flags
- **WHEN** the user runs `dossier --theme dracula openspec/changes/my-feature`
- **THEN** the TUI opens in single-change mode for `my-feature` with the `dracula` theme

#### Scenario: `--help` shows all flags
- **WHEN** the user runs `dossier --help`
- **THEN** the output includes descriptions for `--theme`, `--version`, `--help`, and the positional `[path]` argument
