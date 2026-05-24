**[English](CHANGELOG.md)** | **Español**

# Changelog

## v0.11.0

### Arreglado
- El ratón dejaba de funcionar al volver del editor externo (`e`). Resulta que Bubble Tea v1 no guardaba el estado del ratón al suspender la terminal. Ahora ya funciona, pero da igual porque nadie debería de usar ratón.
- La app fallaba al iniciar si no existían los directorios `archive/`, `specs/` o `changes/`. Ahora devuelve listas vacías como debe ser, sin montar un drama.
- El fondo de la app era negro en vez del color por defecto de tu terminal. `NoColor` significa "sin color", no "negro". Quién lo diría.
- El `go.mod` tenía todas las dependencias marcadas como indirectas. Todas. Incluida Bubble Tea, que es literalmente de lo que va la app.

### Cambiado
- Migración completa a Bubble Tea v2, Bubbles v2, Lip Gloss v2 y Glamour v2. Nuevos imports, nueva API declarativa para `View()`, mensajes de teclado y ratón partidos en tipos separados. Unas 1300 líneas tocadas. No preguntes por quien.
- `renderWithBackground()` y `bgSGRRestore()` eliminados. Bubble Tea v2 maneja el fondo por su cuenta. Una función menos que mantener.

### Añadido
- Tests unitarios. Sí, por fin. ~30 tests entre `loader_test.go`, `tasks_test.go` y `view_test.go`. Cobertura del 74% en `openspec`. Los de UI son más difíciles, no me juzgues.
- CI en GitHub Actions: `go vet`, `go test -race` y coverage en cada push y PR a `main`. Ahora los fallos se ven antes de mergear, no después.

### Interno
- El paquete `openspec` acepta un path raíz explícito en todas sus funciones (`LoadFrom`, `LoadConfigFrom`, etc.) en vez de llamar a `os.Getwd()` internamente. Más testeable, menos acoplado al estado global.
- Todas las funciones del loader devuelven `error` en vez de tragarse los fallos silenciosamente. Los errores de YAML malformado ya no se ignoran como si nada.

## v0.10.0

### Añadido
- La barra de pestañas muestra un color distinto (cyan) para las barras de progreso que llegan al 100%. Este cambio merecía saltar directamente a la versión 1.0, lo sé.
- Nueva vista de información del proyecto: pulsa `i` para ver `openspec/config.yaml` renderizado como markdown. Todavía no se puede editar. Olvidé hacerlo.
- Soporte para ratón: haz clic en las pestañas para cambiarlas, la rueda de scroll funciona en los visores. En cualquier caso no uses ratón, es de cobardes.
- `Tab` / `Shift+Tab` para avanzar y retroceder entre pestañas disponibles. Bienvenido al mundo de incompatibilidades de asignación de teclas entre la app y el sistema de ventanas.
- Flag `--version` / `-v` para mostrar la versión actual. Esto lo hizo sólo la IA, sin que se lo pidiese.

### Cambiado
- La barra de progreso al 100% se muestra en cyan en lugar de verde. Cyan es como azul clarito, por si se me olvida.
- Las releases de goreleaser ahora son totalmente automáticas (sin borradores). Aburrido.
- La barra de ayuda incluye los atajos `Tab` y de ratón.

### Interno
- Separado `internal/ui/model.go` en seis archivos específicos (`model.go`, `update.go`, `viewport.go`, `index.go`, `tasks.go`, `view.go`). Super aburrido.
