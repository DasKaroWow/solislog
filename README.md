# solislog

A small **template-based contextual logger for Go** inspired by the developer experience of Loguru.

`solislog` is an experiment focused on a pleasant human-readable API rather than a high-performance JSON-first design.
The current version already supports:

- log levels
- template-based output
- built-in fields like `time`, `level`, and `message`
- `extra[...]` fields
- contextual logging through `context.Context`
- per-request / per-operation contextualized loggers built from a base logger

This project is currently closer to **v0.1** than a production-ready logging library, but the core idea is already working.

## Why this project exists

In Go, logging often feels like a set of low-level building blocks. Libraries like `zerolog`, `zap`, and `slog` are powerful, but they do not always provide the kind of out-of-the-box developer experience that feels as convenient as Python's `loguru`.

The goal of `solislog` is to explore a simpler and more ergonomic model:

- a readable template string for output
- contextual data through `extra`
- one base logger configuration
- one contextualized logger per request / operation
- human-readable logs without forcing a JSON-first workflow

## Current features

- `Debug`, `Info`, `Warning`, and `Error` methods
- level filtering
- template parsing once at logger creation time
- built-in template fields:
  - `{time}`
  - `{level}`
  - `{message}`
- extra field access:
  - `{extra[id]}`
  - `{extra[source]}`
  - etc.
- contextual logger creation via `Contextualize(...)`
- logger propagation through `context.Context`
- support for default extra fields on the base logger

## What the current design looks like

The current model is intentionally small:

- `Add(...)` creates a **base logger**
- `Contextualize(...)` creates a **new logger instance** with merged `extra` fields and stores it in `context.Context`
- `FromContext(...)` retrieves that contextualized logger
- log calls create a record and render it through a compiled template

This means the base logger acts mostly as a **root configuration object**, while the logger stored in context is the one used during actual request or operation handling.

## Example: simple logger

```go
package main

import (
    "os"

    "github.com/DasKaroWow/solislog"
)

func main() {
    logger := solislog.Add(
        os.Stdout,
        solislog.InfoLevel,
        "{time} | {level} | {message}\n",
        nil,
    )

    _ = logger.Info("hello from solislog")
}
```

Example output:

```text
2026-04-23T01:27:07+03:00 | INFO | hello from solislog
```

## Example: logger with default extra fields

```go
package main

import (
    "os"

    "github.com/DasKaroWow/solislog"
)

func main() {
    logger := solislog.Add(
        os.Stdout,
        solislog.InfoLevel,
        "{time} | {level} | {extra[name]} | {message}\n",
        map[string]string{
            "name": "ivan",
        },
    )

    _ = logger.Info("base logger message")
}
```

## Example: contextual logging through `context.Context`

This is the more important use case for the current project.

The base logger is created once, then a **contextualized logger** is created at the boundary of a request, update, or operation.  
That contextualized logger is stored in `context.Context`, and deeper functions can retrieve it and continue logging with the same `extra` values.

```go
package main

import (
    "context"
    "os"

    "github.com/DasKaroWow/solislog"
)

func main() {
    base := solislog.Add(
        os.Stdout,
        solislog.InfoLevel,
        "{time} | {level} | {extra[name]} | {extra[id]} | {message}\n",
        map[string]string{
            "name": "ivan",
        },
    )

    ctx := context.Background()
    ctx = base.Contextualize(ctx, map[string]string{
        "id": "0",
    })

    handle(ctx)
}

func handle(ctx context.Context) {
    log, ok := solislog.FromContext(ctx)
    if !ok {
        return
    }

    _ = log.Info("entered handle")
    process(ctx)
}

func process(ctx context.Context) {
    log, ok := solislog.FromContext(ctx)
    if !ok {
        return
    }

    _ = log.Info("processing request")
}
```

This shows the important idea: **the extra fields move down through the call chain via context**.

Expected output shape:

```text
2026-04-23T01:27:07+03:00 | INFO | ivan | 0 | entered handle
2026-04-23T01:27:07+03:00 | INFO | ivan | 0 | processing request
```

## Supported template fields

Current built-in fields:

- `{time}`
- `{level}`
- `{message}`

Current extra syntax:

- `{extra[key]}`

Examples:

```text
{time} | {level} | {message}
{time} | {extra[source]} | {message}
{time} | {level} | {extra[name]} | {extra[id]} | {message}
```

## Current project structure

```text
.
├── level.go
├── logger.go
├── record.go
├── template.go
└── demo/
    └── main.go
```

## Limitations of the current version

This is still an early version. Things that are intentionally missing or still rough:

- no color support yet
- no alignment / formatting spec support yet
- no file/function/line built-ins yet
- no JSON output mode
- no hooks
- no async logging
- no middleware helpers yet
- no advanced error / exception formatting
- no polished public API around contextual helpers yet

## Roadmap ideas

Likely next steps:

- cleaner README-quality examples and tests
- polish around contextual logging helpers
- better template validation
- pretty console formatting
- optional color support
- richer built-in fields
- Fiber / HTTP integration experiments

## Project status

This repository is currently a **pet project / design exploration** around logging ergonomics in Go.

The point is not to beat existing logging libraries on performance.  
The point is to build something that feels:

- simple
- readable
- pleasant to use
- contextual by default
- practical in normal Go application code

## License

MIT License.