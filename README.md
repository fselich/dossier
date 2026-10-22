**English** | **[Español](README.es.md)**

# dossier

A keyboard-driven terminal UI for reading and navigating [OpenSpec](https://github.com/openspec) project artifacts — proposals, designs, specs, and tasks.

> Built with OpenSpec. This repository contains 12 project-level spec files and 20+ archived changes that document the complete development history of the tool.

<p align="center">
  <img src="docs/dossier.gif" alt="dossier demo" />
</p>

---

## Features

- Starts in a schema-aware index for active changes, canonical specs, requirements, and history
- Creates, validates, archives, reactivates, and safely undoes OpenSpec lifecycle actions from the index
- Renders Markdown with full syntax highlighting
- Toggles task checkboxes (`- [ ]` / `- [x]`) in the selected schema-defined `tasks` artifact output
- Live-reloads on disk changes (500 ms polling)
- Opens any artifact in `$EDITOR`
- Offers `--read-only` mode for safe review without change lifecycle, task, editor, or Git index mutations
- Supports XDG TOML configuration for custom UI, Glamour, and Chroma colors
- Supports context-specific custom keybindings with dynamic help labels
- Accepts a path argument to view a single change directory without a full project

---

## Installation

**Requirements:** terminal with ANSI color support. Go 1.25 or later. The `openspec` executable is required for create, validate, and archive actions; navigation, inspection, reactivation, and eligible session undo remain available without it.

```bash
go install github.com/aemonge/dossier/cmd/dossier@latest

# From source
git clone https://github.com/aemonge/dossier
cd dossier
make build    # produces ./dossier
make install  # installs via go install
```

---

## Usage

Run from the root of an OpenSpec project:

```bash
dossier
```

View a single change directory by path:

```bash
dossier /path/to/openspec/changes/my-change
```

Review without allowing dossier to toggle tasks, launch an editor, or stage/unstage Git files:

```bash
dossier --read-only
dossier --read-only /path/to/openspec/changes/my-change
```

Use `dossier -h` or `dossier --help` for all command-line options. Read-only mode remains visible in the help bar; navigation, live reload, Git status, and diff viewing stay available.

### Configuration

Dossier loads optional TOML configuration from the XDG path:

```text
${XDG_CONFIG_HOME:-~/.config}/dossier/config.toml
```

Use another file with `dossier --config <path>`. A missing default file is ignored; a missing explicit file or invalid setting produces a startup error. `--theme` and `--keystyle` select built-in bases while preserving sparse configuration overrides. Built-in themes are `none`, `dark`, `light`, `dracula`, and `gruvbox-light-soft`; the default key style is `nvim`:

```bash
dossier --read-only --theme gruvbox-light-soft --keystyle nvim
```

The same bases can be selected in TOML and extended with sparse overrides. Dossier starts in the index by default; set `start_view = "change"` for the legacy first-active-change view. Explicit path launches always open their target directly.

```toml
[ui]
start_view = "index" # index or change

[theme]
base = "gruvbox-light-soft"

[keys]
style = "nvim"

# Optional: replace only this action; all other nvim bindings remain.
[keys.viewer]
open = ["e", "o"]
```

Copy the complete light-terminal example, which documents every Dossier color and key action plus Glamour/Chroma token overrides:

```bash
mkdir -p "${XDG_CONFIG_HOME:-$HOME/.config}/dossier"
cp examples/config.example.toml "${XDG_CONFIG_HOME:-$HOME/.config}/dossier/config.toml"
```

A Pi-aligned Gruvbox Light Soft example is also available at `examples/config.gruvbox-light-soft.toml`.

Keybindings are context-specific, so a key may be reused in different modes but cannot be assigned to two actions in the same mode. Help labels automatically reflect the configured bindings.

### OpenSpec hierarchy and dynamic artifacts

The index separates **Active Work**, **Canonical Specs**, and **History**. Active and archived changes show their selected schema (for example `spec-driven`, `bugfix`, or `spike`) and expand into that schema's artifact IDs and concrete output files. Artifact dependencies remain sibling metadata (`requires …`) rather than being presented as a false ownership tree.

Change-local `specs/<capability>/spec.md` files are **delta specs** owned by a change. Project-level `openspec/specs/<capability>/spec.md` files are **canonical specs** and expand into requirements independently of any one change.

When the OpenSpec CLI is available, Dossier enriches active work with authoritative schema order, status, dependencies, and resolved outputs. If the CLI or a schema is unavailable, safe filesystem discovery keeps readable Markdown artifacts navigable and labels the degraded state. Hidden files and symlink escapes are not discovered.

### Keyboard reference

#### Normal mode (viewing a change)

| Key | Action |
|---|---|
| `j` / `down` | Scroll down (or move task cursor down) |
| `k` / `up` | Scroll up (or move task cursor up) |
| `PgDown` / `Ctrl+D` | Scroll one page down (except the Tasks tab) |
| `PgUp` / `Ctrl+U` | Scroll one page up (except the Tasks tab) |
| `h` / `l` | Previous / next artifact tab |
| `Tab` / `Shift+Tab` | Next / previous change |
| `1`–`4` | Select the schema artifact at that position |
| `5` | Open the synthetic active-only code destination when available |
| `Space` | Toggle task under cursor (tasks tab only) |
| `e` | Open artifact in `$EDITOR` |
| `?` | Open configuration information |
| `q` / `Esc` | Enter index mode |
| `Q` | Quit |

#### Index mode (change and spec navigator)

| Key | Action |
|---|---|
| `j` / `down` | Move cursor down |
| `k` / `up` | Move cursor up |
| `Enter` | Toggle expandable rows; inspect output and requirement leaves |
| `Space` | Toggle expandable rows without inspecting |
| `i` | Inspect the selected change, artifact output, spec, or requirement |
| `n` | Choose a workflow schema and create a change |
| `a` | Confirm archive for active work or make archived work active |
| `e` | Edit a safe active artifact output or canonical spec |
| `v` | Validate an active change or canonical spec through OpenSpec |
| `u` | Confirm undo of the latest lifecycle action in this session |
| `/` | Filter index items |
| `s` | Toggle index sorting |
| `?` | Open configuration information |
| `q` / `Q` / `Esc` | Quit from the root index |

While editing an index filter or action prompt, normal text updates the input, `Backspace` deletes, `Enter` accepts, and `Esc` cancels. Archive updates canonical project specs through the official OpenSpec CLI. Reactivation moves the exact archived directory back to Active Work without reversing already-published specs.

Only the latest successful archive or reactivation can be undone, and only during the current Dossier session. Before undo writes anything, Dossier fingerprints every affected post-operation path; external changes or destination collisions refuse the entire undo. Undoing archive restores the active change and exact pre-archive spec contents, while undoing reactivation returns the change to its original date-prefixed archive path.

Read-only mode omits and blocks `n`, `a`, `e`, and `u`; inspect and validate remain available.

#### Git tab and diff viewer

| Key | Action |
|---|---|
| `j` / `k` | Select a changed file, or scroll an open diff vertically |
| `d` / `Enter` / `e` | Open or close the selected diff |
| `[` / `]` | Previous / next changed file while viewing a diff |
| `H` / `L` | Scroll an open diff horizontally left / right |
| `PgDown` / `Ctrl+D` | Scroll an open diff one page down |
| `PgUp` / `Ctrl+U` | Scroll an open diff one page up |
| `s` | Stage / unstage the selected file (disabled in read-only mode) |
| `Tab` / `Shift+Tab` | Next / previous change |
| `q` / `Esc` | Close the diff, or return to the index |
| `Q` | Quit |

#### Archive mode (viewing an archived change)

| Key | Action |
|---|---|
| `j` / `k` | Scroll |
| `PgDown` / `Ctrl+D` | Scroll one page down (except the Tasks tab) |
| `PgUp` / `Ctrl+U` | Scroll one page up (except the Tasks tab) |
| `h` / `l` | Previous / next artifact tab |
| `1`–`4` | Select an artifact tab directly |
| `q` / `Esc` | Return to index |
| `Q` | Quit |

#### Spec viewer mode

| Key | Action |
|---|---|
| `j` / `k` | Scroll |
| `PgDown` / `Ctrl+D` | Scroll one page down |
| `PgUp` / `Ctrl+U` | Scroll one page up |
| `q` / `Esc` | Return to index |
| `Q` | Quit |

In requirement focus mode:

| Key | Action |
|---|---|
| `h` / `l` | Previous / next requirement |
| `j` / `k` | Scroll |
| `PgDown` / `Ctrl+D` | Scroll one page down |
| `PgUp` / `Ctrl+U` | Scroll one page up |
| `q` / `Esc` | Return to index |
| `Q` | Quit |

#### Configuration viewer

| Key | Action |
|---|---|
| `j` / `k` | Scroll |
| `PgDown` / `Ctrl+D` | Scroll one page down |
| `PgUp` / `Ctrl+U` | Scroll one page up |
| `q` / `?` / `Esc` | Return to the previous view |

---

## Project structure

dossier expects an `openspec/` directory at the project root. The layout below is an illustrative `spec-driven` workflow; custom schemas may define arbitrary Markdown artifact IDs, nested paths, and multi-file outputs.

```
openspec/
├── changes/
│   ├── <change-name>/
│   │   ├── .openspec.yaml   # Required: identifies the directory as a change
│   │   ├── proposal.md
│   │   ├── design.md
│   │   ├── tasks.md         # GFM checkbox syntax: - [ ] / - [x]
│   │   └── specs/
│   │       └── <spec-name>/
│   │           └── spec.md
│   └── archive/
│       └── YYYY-MM-DD-<name>/
└── specs/
    └── <spec-name>/
        └── spec.md          # Requirements parsed from: ### Requirement: <name>
```
