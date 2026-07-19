# Proposal: add-theme-ui-colors

## Summary

Extender el sistema de temas para que los 22 estilos lipgloss del chrome de UI (header, tabs, bordes, cursor marks, tareas, git status, diff, barras de progreso) usen colores definidos por el tema activo, en lugar de valores ANSI 0-15 hardcodeados.

## Motivation

La fase 1 del sistema de temas (`--theme`) cubrió glamour, chroma, fondo del viewport y fondos de diff. Pero los estilos de UI (tab inactive white, cursor marks white, secciones amarillas, etc.) siguen hardcodeados en `styles.go`. En el tema `light`, el viewport es blanco (`#ffffff`) y varios elementos con foreground `Color("15")` (blanco) se vuelven invisibles. El tema `dracula` también sufre, aunque menos, porque ciertos colores (amarillo brillante sobre fondo violáceo) no armonizan.

## Scope

- Definir un struct `ThemeColors` con ~11 roles de color semánticos
- Añadir `ThemeStyles` (struct de `lipgloss.Style` pre-construidos) derivado de `ThemeColors`
- Función `BuildStyles(ThemeColors) ThemeStyles`
- Reemplazar los 22 `var` de `styles.go` y los 3 de `tasks.go` por usos de `m.theme.Styles.Xxx`
- Definir paletas para `dark`, `none`, `light`, `dracula`
- Los archivos `view.go`, `index.go`, `tasks.go`, `git.go`, `gitdiff.go` se tocan

## Non-goals

- Cambiar la estructura de glamour/chroma/ViewBg — ya funciona
- Temas definidos en archivos externos (YAML)
- Cambio de tema en runtime
- Añadir/quitar temas (los 4 existentes se mantienen)
