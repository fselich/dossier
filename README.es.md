**[English](README.md)** | **Español**

# dossier

Interfaz de terminal controlada por teclado para leer y navegar artefactos de proyectos [OpenSpec](https://github.com/openspec) — propuestas, diseños, specs y tareas.

> Desarrollado con OpenSpec. Este repositorio contiene 12 archivos de spec de proyecto y más de 20 cambios archivados que documentan el historial completo de desarrollo de la herramienta.

<p align="center">
  <img src="docs/dossier.gif" alt="demo de dossier" />
</p>

---

## Funcionalidades

- Navega todos los cambios activos y sus artefactos desde una única interfaz
- Renderiza Markdown con resaltado de sintaxis completo
- Alterna casillas de verificación (`- [ ]` / `- [x]`) en el propio archivo `tasks.md`
- Recarga en vivo ante cambios en disco (sondeo cada 500 ms)
- Abre cualquier artefacto en `$EDITOR`
- Ofrece el modo `--read-only` para revisar sin modificar tareas, archivos mediante el editor ni el índice de Git
- Admite configuración TOML XDG para personalizar los colores de la interfaz, Glamour y Chroma
- Admite atajos personalizados por contexto con ayuda dinámica
- Acepta una ruta como argumento para ver un cambio concreto sin necesitar un proyecto completo

---

## Instalación

**Requisitos:** terminal con soporte de color ANSI. Go 1.25 o posterior.

```bash
go install github.com/aemonge/dossier/cmd/dossier@latest

# Desde el código fuente
git clone https://github.com/aemonge/dossier
cd dossier
make build    # genera ./dossier
make install  # instala mediante go install
```

---

## Uso

Ejecutar desde la raíz de un proyecto OpenSpec:

```bash
dossier
```

Ver un directorio de cambio concreto por ruta:

```bash
dossier /ruta/a/openspec/changes/mi-cambio
```

Revisar sin permitir que dossier alterne tareas, abra un editor o prepare/quite archivos del índice de Git:

```bash
dossier --read-only
dossier --read-only /ruta/a/openspec/changes/mi-cambio
```

Usa `dossier -h` o `dossier --help` para ver todas las opciones. El modo de solo lectura permanece visible en la barra de ayuda; la navegación, recarga en vivo, estado de Git y vista de diferencias siguen disponibles.

### Configuración

Dossier carga configuración TOML opcional desde la ruta XDG:

```text
${XDG_CONFIG_HOME:-~/.config}/dossier/config.toml
```

Usa otro archivo con `dossier --config <ruta>`. La ausencia del archivo predeterminado se ignora; un archivo explícito ausente o una opción inválida produce un error al iniciar. `--theme` y `--keystyle` seleccionan bases incluidas conservando las personalizaciones parciales del archivo. Los temas incluidos son `none`, `dark`, `light`, `dracula` y `gruvbox-light-soft`; el estilo de teclas predeterminado es `nvim`:

```bash
dossier --read-only --theme gruvbox-light-soft --keystyle nvim
```

Las mismas bases pueden seleccionarse en TOML y ampliarse con personalizaciones parciales:

```toml
[theme]
base = "gruvbox-light-soft"

[keys]
style = "nvim"

# Opcional: sustituye solo esta acción; el resto de atajos nvim permanece.
[keys.viewer]
open = ["e", "o"]
```

Puedes copiar el ejemplo completo para terminales claras, que documenta todos los colores y acciones de teclado de Dossier, además de los tokens de Glamour/Chroma:

```bash
mkdir -p "${XDG_CONFIG_HOME:-$HOME/.config}/dossier"
cp examples/config.example.toml "${XDG_CONFIG_HOME:-$HOME/.config}/dossier/config.toml"
```

También hay un ejemplo Gruvbox Light Soft alineado con Pi en `examples/config.gruvbox-light-soft.toml`.

Los atajos dependen del contexto: una tecla puede reutilizarse en modos distintos, pero no asignarse a dos acciones dentro del mismo modo. La barra de ayuda refleja automáticamente los atajos configurados.

### Referencia de teclado

#### Modo normal (viendo un cambio)

| Tecla | Acción |
|---|---|
| `j` / `↓` | Desplazar hacia abajo (o mover cursor de tareas hacia abajo) |
| `k` / `↑` | Desplazar hacia arriba (o mover cursor de tareas hacia arriba) |
| `PgDown` / `Ctrl+D` | Desplazar una página hacia abajo (excepto en la pestaña Tareas) |
| `PgUp` / `Ctrl+U` | Desplazar una página hacia arriba (excepto en la pestaña Tareas) |
| `h` / `l` | Pestaña de artefacto anterior / siguiente |
| `Tab` / `Shift+Tab` | Cambio siguiente / anterior |
| `1`–`5` | Seleccionar directamente una pestaña de artefacto |
| `Space` | Alternar tarea bajo el cursor (solo en pestaña de tareas) |
| `e` | Abrir artefacto en `$EDITOR` |
| `?` | Abrir información de configuración |
| `q` / `Esc` | Entrar en modo índice |
| `Q` | Salir |

#### Modo índice (navegador de cambios y specs)

| Tecla | Acción |
|---|---|
| `j` / `↓` | Mover cursor hacia abajo |
| `k` / `↑` | Mover cursor hacia arriba |
| `l` / `Enter` | Abrir el cambio, spec o cambio archivado seleccionado |
| `Space` | Expandir / contraer una spec de proyecto |
| `/` | Filtrar elementos del índice |
| `s` | Alternar el orden del índice |
| `?` | Abrir información de configuración |
| `q` / `Q` / `Esc` | Salir desde el índice raíz |

Al editar un filtro del índice, el texto normal actualiza el filtro, `Backspace` borra, `Enter` acepta y `Esc` cancela.

#### Pestaña Git y visor de diferencias

| Tecla | Acción |
|---|---|
| `j` / `k` | Seleccionar un archivo cambiado o desplazar verticalmente una diferencia abierta |
| `d` / `Enter` / `e` | Abrir o cerrar la diferencia seleccionada |
| `[` / `]` | Archivo cambiado anterior / siguiente al ver una diferencia |
| `H` / `L` | Desplazar horizontalmente una diferencia abierta |
| `PgDown` / `Ctrl+D` | Desplazar una página hacia abajo en la diferencia |
| `PgUp` / `Ctrl+U` | Desplazar una página hacia arriba en la diferencia |
| `s` | Preparar / quitar del índice el archivo seleccionado (desactivado en modo de solo lectura) |
| `Tab` / `Shift+Tab` | Cambio siguiente / anterior |
| `q` / `Esc` | Cerrar la diferencia o volver al índice |
| `Q` | Salir |

#### Modo archivo (viendo un cambio archivado)

| Tecla | Acción |
|---|---|
| `j` / `k` | Desplazar |
| `PgDown` / `Ctrl+D` | Desplazar una página hacia abajo (excepto en la pestaña Tareas) |
| `PgUp` / `Ctrl+U` | Desplazar una página hacia arriba (excepto en la pestaña Tareas) |
| `h` / `l` | Pestaña de artefacto anterior / siguiente |
| `1`–`4` | Seleccionar directamente una pestaña de artefacto |
| `q` / `Esc` | Volver al índice |
| `Q` | Salir |

#### Modo visor de spec

| Tecla | Acción |
|---|---|
| `j` / `k` | Desplazar |
| `PgDown` / `Ctrl+D` | Desplazar una página hacia abajo |
| `PgUp` / `Ctrl+U` | Desplazar una página hacia arriba |
| `q` / `Esc` | Volver al índice |
| `Q` | Salir |

En modo foco de requisitos:

| Tecla | Acción |
|---|---|
| `h` / `l` | Requisito anterior / siguiente |
| `j` / `k` | Desplazar |
| `PgDown` / `Ctrl+D` | Desplazar una página hacia abajo |
| `PgUp` / `Ctrl+U` | Desplazar una página hacia arriba |
| `q` / `Esc` | Volver al índice |
| `Q` | Salir |

#### Visor de configuración

| Tecla | Acción |
|---|---|
| `j` / `k` | Desplazar |
| `PgDown` / `Ctrl+D` | Desplazar una página hacia abajo |
| `PgUp` / `Ctrl+U` | Desplazar una página hacia arriba |
| `q` / `?` / `Esc` | Volver a la vista anterior |

---

## Estructura de proyecto

dossier espera un directorio `openspec/` en la raíz del proyecto:

```
openspec/
├── changes/
│   ├── <nombre-cambio>/
│   │   ├── .openspec.yaml   # Requerido: identifica el directorio como un cambio
│   │   ├── proposal.md
│   │   ├── design.md
│   │   ├── tasks.md         # Sintaxis de casillas GFM: - [ ] / - [x]
│   │   └── specs/
│   │       └── <nombre-spec>/
│   │           └── spec.md
│   └── archive/
│       └── YYYY-MM-DD-<nombre>/
└── specs/
    └── <nombre-spec>/
        └── spec.md          # Requisitos detectados por: ### Requirement: <nombre>
```
