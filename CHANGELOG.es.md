**[English](CHANGELOG.md)** | **Español**

# Changelog

## v0.16.0

### Añadido

- Nueva pestaña `changes` (tecla `5`) en el visor de cambios que muestra los ficheros modificados del working tree de git. Ficheros modificados, añadidos, untracked, renombrados y borrados aparecen con indicadores de estado (`·M` sin staged, `M·` staged, `MM` en ambos).
- Pulsa `d` (o `Enter`/`e`) sobre un fichero en la pestaña de cambios para ver su diff con syntax highlighting completo.

## v0.15.0

### Añadido

- Las secciones del índice ahora se pueden plegar. Pulsa `Space` en cualquier cabecera de sección para ocultar sus hijos; pulsa otra vez para expandir; pulsa otra vez para plegar; y otra para expandir... Lo acabarás entendiendo.

### Arreglado

- La caché de renderizado ahora se invalida al volver del editor externo (`e`). Antes editabas un archivo y, al volver a la vista de ese artefacto, estaba lleno de basura. Si no lo habías notado es que tus artefactos ya eran basura antes de editarlos.

## v0.14.1

### Arreglado

- Los spans de código en tareas completadas ya no muestran la primera letra de un color diferente. Lipgloss renderiza el texto subrayado carácter por carácter, reseteando el foreground entre ellos. La solución combina el subrayado con el color de primer plano para que cada carácter herede ambos.

## v0.14.0

### Añadido

- Pulsa `/` en la vista de índice para filtrar cambios, specs y archivados por nombre en tiempo real. Escribe para acotar, `Enter` para fijar el filtro, `Esc` para limpiarlo. Un buscador de toda la vida, vaya. 

## v0.13.0

### Interno

- Partido el monolito `handleKeyPress()` en funciones `update*()` por modo, cada una en su propio archivo: `viewer.go`, `index.go`, `spec.go`, `config.go`. `update.go` ahora es un simple despachador. 
- Introducida una interfaz `fileSystem` y un struct `Loader` en `openspec`, así el paquete ya no depende directamente de `os`. Todas las funciones públicas se conservan mediante wrappers compatibles hacia atrás. 
- Añadido `.golangci.yml` con errcheck, staticcheck, govet, unused, gofmt, goimports, y un `Makefile` con objetivos `test`, `lint` y `fmt`. 
- Eliminadas todas las llamadas silenciosas a `log.Printf`. Los errores de carga de archivos y specs ahora se muestran en la barra de ayuda durante 3 segundos mediante `m.errMsg`, exactamente igual que los errores de toggle. 

### Cambiado

- El slice `parts` de la barra de pestañas ahora se preasigna a exactamente 4 entradas, y el slice `items` de tareas se preasigna al número de líneas. Ahora todo va 3 nanosegundos más rápido. Ha merecido la pena el gasto de tokens.
- La función `taskCounts` ya no usa naked returns (que confundían a cualquiera que se desplazase más allá de la línea 491 de index.go).
- Las constantes de layout (`chromeTop`, `chromeHeader`, etc.) reemplazan a los números mágicos en `contentHeight()`. Ahora sabes por qué restaba 6.
- La lógica de recarga que estaba copiada y pegada en dos sitios ahora es un único método `mergeReloadedChange()`. DRY*2.

## v0.12.0

### Arreglado

- Si arrancamos dossier cuando no hay cambios pendientes, te muestra el índice de specs y cambios archivados. 
- Las actualizaciones de contenido de tareas dentro de cambios existentes ahora provocan una actualización en vivo de la lista de índices, en lugar de ignorarlas silenciosamente.
- El placeholder de carga (`"Loading..."` / `"Cargando..."`) fue eliminado. El markdown crudo se muestra de inmediato mientras la versión con estilo se renderiza en segundo plano. Adios al modo involuntario para epilépticos.

### Cambiado

- La lista de cambios en la vista de índice ahora se ordena por fecha de creación (descendente). Antes se ordenaban con el método slsdlp. 

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
