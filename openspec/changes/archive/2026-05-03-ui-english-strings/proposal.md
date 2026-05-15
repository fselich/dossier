## Why

All UI strings are currently in Spanish, making the tool inconsistent for an English-language codebase and harder to use for non-Spanish speakers. Standardising to English aligns the interface with the language of the specs, code, and tooling around it.

## What Changes

- Index section headers: "Activos" → "Active Changes", "Archivados" → "Archived Changes", "Specs" → "Specifications"
- Index empty-state messages: "No hay changes activos" → "No active changes", "No hay changes archivados" → "No archived changes", "No hay specs disponibles" → "No specifications available"
- Index helpbar: `j/k: navegar  Enter: abrir  Esc: salir` → `j/k: navigate  Enter: open  Esc: quit`
- Archive viewer header badge: `[archivo]` → `[archive]`
- Archive viewer helpbar: `a/Esc: índice` → `a/Esc: index`
- Spec viewer helpbar: `Esc: índice` → `Esc: index`
- Normal mode helpbars: `Esc: índice` → `Esc: index`, `j/k: navegar` → `j/k: navigate`
- Welcome screen: "No hay changes activos. Crea uno con /opsx:propose" → "No active changes. Create one with /opsx:propose"

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `change-index`: Section header labels and helpbar string are part of requirements
- `tui-viewer`: Welcome message and helpbar strings are part of requirements
- `archive-viewer`: Header badge `[archivo]` and helpbar string are part of requirements
- `index-specs-section`: Empty-state message "No hay specs disponibles" is part of requirements
- `spec-detail-viewer`: Helpbar string is part of requirements

## Impact

- `internal/ui/model.go`: all string literals in `renderIndexContent`, `renderHeader`, `renderHelpBar`, and `emptyView`
- No logic changes, no new dependencies
