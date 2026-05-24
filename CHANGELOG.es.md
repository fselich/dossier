**[English](CHANGELOG.md)** | **Español**

# Changelog

## v0.10.0

### Añadido
- La barra de pestañas muestra un color distinto (cyan) para las barras de progreso que llegan al 100%. Este cambio merecía saltar directamente a la versión 1.0, lo sé.
- Nueva vista de información del proyecto: pulsa `i` para ver `openspec/config.yaml` renderizado como markdown. Todavía no se puede editar. Olvidé hacerlo.
- Soporte para ratón: haz clic en las pestañas para cambiarlas, la rueda de scroll funciona en los visores. En cualquier caso no uses ratón, es de cobardes.
- `Tab` / `Shift+Tab` para avanzar y retroceder entre pestañas disponibles. Bienvenido al mundo de incompatibilidades de asignación de teclas entre la app y el sistema de ventanas.
- Flag `--version` / `-v` para mostrar la versión actual. Esto lo hizo sólo la IA, sin que se lo pidiese.

### Cambiado
- La barra de progreso al 100% se muestra en cyan en lugar de verde. Cyan es como azul clarito, por si se me olvida.
- Las releases de goreleaser ahora se crean como borrador y con el changelog automático desactivado. Aburrido.
- La barra de ayuda incluye los atajos `Tab` y de ratón.

### Interno
- Separado `internal/ui/model.go` en seis archivos específicos (`model.go`, `update.go`, `viewport.go`, `index.go`, `tasks.go`, `view.go`). Super aburrido.
