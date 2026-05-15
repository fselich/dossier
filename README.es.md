**[English](README.md)** | **Español**

# dossier

Interfaz de terminal controlada por teclado para leer y navegar artefactos de proyectos [OpenSpec](https://github.com/openspec) — propuestas, diseños, specs y tareas.

> Desarrollado con OpenSpec. Este repositorio contiene 12 archivos de spec de proyecto y más de 20 cambios archivados que documentan el historial completo de desarrollo de la herramienta.

---

## Funcionalidades

- Navega todos los cambios activos y sus artefactos desde una única interfaz
- Renderiza Markdown con resaltado de sintaxis completo
- Alterna casillas de verificación (`- [ ]` / `- [x]`) en el propio archivo `tasks.md`
- Recarga en vivo ante cambios en disco (sondeo cada 500 ms)
- Abre cualquier artefacto en `$EDITOR`
- Acepta una ruta como argumento para ver un cambio concreto sin necesitar un proyecto completo

---

## Instalación

**Requisitos:** Go 1.25 o posterior, terminal con soporte de color ANSI.

```bash
# Desde el código fuente
git clone https://github.com/fselich/dossier
cd dossier
make build    # genera ./dossier
make install  # instala mediante go install

# Usando go install
go install github.com/fselich/dossier/cmd/dossier@latest
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

### Referencia de teclado

#### Modo normal (viendo un cambio)

| Tecla | Acción |
|---|---|
| `h` / `l` | Cambio anterior / siguiente |
| `1` | Pestaña de propuesta |
| `2` | Pestaña de diseño |
| `3` | Pestaña de specs (pulsando de nuevo se cicla entre varios archivos) |
| `4` | Pestaña de tareas |
| `j` / `↓` | Desplazar hacia abajo (o mover cursor de tareas hacia abajo) |
| `k` / `↑` | Desplazar hacia arriba (o mover cursor de tareas hacia arriba) |
| `Space` | Alternar tarea bajo el cursor (solo en pestaña de tareas) |
| `e` | Abrir artefacto en `$EDITOR` |
| `a` / `Esc` | Entrar en modo índice |
| `q` / `Ctrl+C` | Salir |

#### Modo índice (navegador de cambios y specs)

| Tecla | Acción |
|---|---|
| `j` / `↓` | Mover cursor hacia abajo |
| `k` / `↑` | Mover cursor hacia arriba |
| `Enter` | Abrir el cambio, spec o cambio archivado seleccionado |
| `Space` | Expandir / contraer una spec de proyecto |
| `q` / `Esc` / `Ctrl+C` | Salir |

#### Modo archivo (viendo un cambio archivado)

| Tecla | Acción |
|---|---|
| `1`–`4` | Cambiar pestaña de artefacto |
| `j` / `k` | Desplazar |
| `a` / `Esc` | Volver al índice |
| `q` / `Ctrl+C` | Salir |

#### Modo visor de spec

| Tecla | Acción |
|---|---|
| `j` / `k` | Desplazar |
| `Esc` | Volver al índice |
| `q` / `Ctrl+C` | Salir |

En modo foco de requisitos:

| Tecla | Acción |
|---|---|
| `h` / `l` | Requisito anterior / siguiente |
| `j` / `k` | Desplazar |
| `Esc` | Volver al índice |
| `q` / `Ctrl+C` | Salir |

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
