## Context

El sistema de temas actual (fase 1) cubre glamour, chroma, ViewBg, y fondos de diff. Los ~25 estilos lipgloss del chrome de UI usan ANSI 0-15 hardcodeados en `styles.go` y `tasks.go`. Esto rompe el tema light (blanco sobre blanco invisible) y limita cualquier tema futuro.

## Goals / Non-Goals

**Goals:**
- Definir color roles semánticos que todos los temas comparten
- Construir estilos lipgloss desde esos colores una vez al inicio
- Reemplazar todas las referencias a estilos hardcodeados por `m.theme.Styles.Xxx`
- Paletas de color para dark, none (hereda dark), light, dracula

**Non-Goals:**
- Cambiar glamour/chroma/ViewBg — ya funcionan con temas
- Carga de temas desde archivos externos
- Runtime theme switching
- Nuevos temas

## Decisions

### 1. ThemeColors: struct plano con roles semánticos

11 roles de color cubren todos los usos actuales:

| Campo | Uso | Dark | Light |
|-------|-----|------|-------|
| `PrimaryFg` | Texto principal, cursor marks | `"15"` (blanco) | `"0"` (negro) |
| `MutedFg` | Help, separadores, done, dots | `"8"` (gris oscuro) | `"7"` (gris claro) |
| `MidFg` | Tareas pendientes | `"7"` (gris claro) | `"8"` (gris oscuro) |
| `AccentBlue` | Headers | `"12"` (azul brillante) | `"4"` (azul) |
| `AccentYellow` | Secciones, git modified | `"11"` (amarillo) | `"3"` (amarillo oscuro) |
| `AccentCyan` | Progreso, git renamed | `"6"` (cyan) | `"6"` |
| `AccentGreen` | Completado, git added | `"2"` (verde) | `"2"` |
| `AccentRed` | Errores, diff removed | `"9"` (rojo brillante) | `"1"` (rojo) |
| `AccentMagenta` | Git untracked | `"5"` (magenta) | `"5"` |
| `ActiveBg` | Fondo ítem activo | `"4"` (azul) | `"4"` |
| `ActiveFg` | Texto sobre ActiveBg | `"15"` (blanco) | `"15"` |

Struct plano, no mapa: compile-time safety, no magic strings, consistente con el patrón de `Theme` existente.

### 2. ThemeStyles: estilos lipgloss pre-construidos

```go
type ThemeStyles struct {
    Header          lipgloss.Style  // Bold + AccentBlue
    TabActive       lipgloss.Style  // Bold + ActiveFg + ActiveBg + Padding(0,1)
    TabInactive     lipgloss.Style  // PrimaryFg + Padding(0,1)
    TabDisabled     lipgloss.Style  // MutedFg + Padding(0,1)
    IndexActive     lipgloss.Style  // Bold + ActiveFg + ActiveBg
    Section         lipgloss.Style  // Bold + AccentYellow
    TaskCursorMark  lipgloss.Style  // Bold + PrimaryFg
    TaskDone        lipgloss.Style  // MutedFg
    TaskPending     lipgloss.Style  // MidFg
    Help            lipgloss.Style  // MutedFg
    Error           lipgloss.Style  // Bold + AccentRed
    ProgressDone    lipgloss.Style  // AccentCyan
    ProgressComplete lipgloss.Style // AccentGreen
    ProgressEmpty   lipgloss.Style  // MutedFg
    Separator       lipgloss.Style  // MutedFg
    GitModified     lipgloss.Style  // AccentYellow
    GitAdded        lipgloss.Style  // AccentGreen
    GitDeleted      lipgloss.Style  // MutedFg
    GitRenamed      lipgloss.Style  // AccentCyan
    GitUntracked    lipgloss.Style  // AccentMagenta
    GitDot          lipgloss.Style  // MutedFg
    DiffRemoved     lipgloss.Style  // AccentRed
    TaskCodeCyan    lipgloss.Style  // AccentCyan
    TaskCodeDone    lipgloss.Style  // Underline + MutedFg
}
```

`BuildStyles(c ThemeColors) ThemeStyles` construye todos de una vez. Llamada una vez por tema, resultado guardado en `Theme.Styles`.

El `boldStyle` (solo Bold sin color) en tasks.go no necesita tematización.

### 3. Paletas por tema

```
         PrimaryFg  MutedFg  MidFg  AccBlue  AccYell  AccCyan  AccGrn  AccRed  AccMgnt  ActBg  ActFg
dark     "15"       "8"      "7"    "12"     "11"     "6"      "2"     "9"     "5"      "4"    "15"
none     "15"       "8"      "7"    "12"     "11"     "6"      "2"     "9"     "5"      "4"    "15"
light    "0"        "7"      "8"    "4"      "3"      "6"      "2"     "1"     "5"      "4"    "15"
dracula  "15"       "8"      "7"    "12"     "3"      "6"      "2"     "9"     "5"      "4"    "15"
```

- `dark` = preserva comportamiento actual
- `none` = hereda dark (misma paleta, ViewBg nil)
- `light` = colores adaptados para fondo blanco: negro para texto, azul/amarillo/rojo oscurecidos
- `dracula` = casi igual a dark, solo `AccentYellow` baja de `"11"` a `"3"` (el amarillo brillante choca con el fondo `#282a36`)

### 4. Reemplazo de estilos: directo, archivo por archivo

Cada `xxxStyle.Render(...)` se reemplaza por `m.theme.Styles.Xxx.Render(...)`. No se necesita plumbing adicional: todas las funciones de render ya tienen acceso a `Model`.

Archivos afectados:
- `view.go`: `headerStyle`, `tabActiveStyle`, `tabInactiveStyle`, `tabDisabledStyle`, `sectionStyle`, `helpStyle`, `errStyle`, `separatorStyle`, `gitStatusAdded`, `gitStatusRenamed`
- `index.go`: `helpStyle`, `sectionStyle`, `progressDoneStyle`, `progressCompleteStyle`, `progressEmptyStyle`, `indexActiveStyle`, `taskPendingStyle`
- `tasks.go`: `taskCursorMarkStyle`, `taskDoneStyle`, `taskPendingStyle`, `sectionStyle`, `helpStyle`, `progressDoneStyle`, `progressCompleteStyle`, `progressEmptyStyle` + `doneCodeStyle`, `cyanStyle`
- `git.go`: `taskCursorMarkStyle`, `helpStyle`, `gitStatusDeleted`, `gitStatusUntracked`, `gitStatusAdded`, `gitStatusRenamed`, `gitStatusModified`, `gitStatusDot`, `sectionStyle`
- `gitdiff.go`: `helpStyle`, `gitStatusRenamed`, `gitStatusAdded`, `diffRemoved`

### 5. Inline styles en tasks.go

`doneCodeStyle`, `cyanStyle` y `boldStyle` son estilos inline para markdown en tareas. `boldStyle` (solo Bold) se queda como variable de paquete. Los otros dos se mueven a `ThemeStyles.TaskCodeDone` y `ThemeStyles.TaskCodeCyan`.

### 6. none theme hereda colores de dark

`NoneTheme` comparte `DarkColors`. La única diferencia es `ViewBg = nil`.

## Risks / Trade-offs

- **[Low] Cambio mecánico**: el reemplazo es `s/headerStyle/m.theme.Styles.Header/g` (con variantes). Riesgo de olvidar alguna ocurrencia. Mitigación: grep exhaustivo antes y después.
- **[Low] Tests**: `view_test.go` tiene referencias a estilos. Se actualizan para usar `m.theme.Styles.Xxx` o se construye un `Theme` con `BuildStyles` en los tests.
- **[None] Regresión visual en dark**: los colores de `DarkColors` son idénticos a los valores hardcodeados actuales. Sin cambios visuales para el tema por defecto.
