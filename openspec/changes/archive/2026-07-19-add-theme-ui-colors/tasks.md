## Tasks

- [x] Definir `ThemeColors` struct con 11 campos de color semánticos en `internal/ui/themes.go`
- [x] Definir `ThemeStyles` struct con ~24 estilos lipgloss pre-construidos
- [x] Implementar `BuildStyles(c ThemeColors) ThemeStyles`
- [x] Definir `DarkColors`, `LightColors`, `DraculaColors` en `internal/ui/themes.go`
- [x] Actualizar `DarkTheme`, `NoneTheme`, `LightTheme`, `DraculaTheme` con `Colors` y `Styles`
- [x] Eliminar los `var` de `styles.go` (o vaciar el archivo)
- [x] Reemplazar usos en `view.go`: `headerStyle` → `m.theme.Styles.Header`, etc.
- [x] Reemplazar usos en `index.go`: `helpStyle`, `sectionStyle`, `progressDoneStyle`, `progressCompleteStyle`, `progressEmptyStyle`, `indexActiveStyle`, `taskPendingStyle`
- [x] Reemplazar usos en `tasks.go`: todos los estilos + `doneCodeStyle`/`cyanStyle` → `TaskCodeDone`/`TaskCodeCyan`
- [x] Reemplazar usos en `git.go`: estilos de git status y sección
- [x] Reemplazar usos en `gitdiff.go`: `helpStyle`, `gitStatusRenamed`, `gitStatusAdded`, `diffRemoved`
- [x] Actualizar `themes_test.go` para verificar `BuildStyles` y paletas
- [x] Actualizar `view_test.go` para usar `Theme.Styles` en lugar de vars globales
- [x] `make test` pasa
- [x] `make lint` pasa
- [x] Verificación manual: `--theme light`, `--theme dracula`, `--theme dark`, `--theme none`
