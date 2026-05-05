# solislog

[![Go Reference](https://pkg.go.dev/badge/github.com/DasKaroWow/solislog.svg)](https://pkg.go.dev/github.com/DasKaroWow/solislog)
[![Go Report Card](https://goreportcard.com/badge/github.com/DasKaroWow/solislog)](https://goreportcard.com/report/github.com/DasKaroWow/solislog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`solislog` is a small template-based contextual logger for Go.

It focuses on readable console output, contextual fields, simple templates, optional colors, optional JSON output, and an API that stays close to normal Go patterns.

The goal is not to compete with `zap`, `zerolog`, or `slog` on performance. The goal is to keep logging simple, readable, and pleasant to use in small and medium Go projects.

## Features

- Multiple handlers per logger
- Per-handler log level filtering
- Per-handler templates
- Colorized text output with tags like `<red>...</red>` and `<level>...</level>`
- Built-in template fields: `{time}`, `{level}`, `{message}`, `{extra}`
- Caller metadata fields: `{file}`, `{path}`, `{line}`, `{function}`, `{caller}`
- Custom contextual fields through `{extra[key]}`
- Optional JSON output mode
- `BeforeHook` and `AfterHook` for per-handler record processing
- `ErrorHandler` for write errors
- `Bind(...)` for creating child loggers with merged extra fields
- `Contextualize(...)` and `FromContext(...)` for passing loggers through `context.Context`
- Simple log methods: `Debug`, `Info`, `Warning`, `Error`, `Fatal`
- Safe concurrent use by multiple goroutines

## Installation

```bash
go get github.com/DasKaroWow/solislog
```

## Quick start

```go
package main

import (
    "os"

    "github.com/DasKaroWow/solislog"
)

func main() {
    logger := solislog.NewLogger(
        nil,
        solislog.NewHandler(os.Stdout, solislog.InfoLevel, nil),
    )

    logger.Info("hello from solislog")
}
```

Passing `nil` handler options uses the default template:

```text
{time} | {level} | {message}\n
```

Example output:

```text
2026-05-05T15:30:00+03:00 | INFO | hello from solislog
```

## Colored output

Templates support ANSI color tags:

```go
logger := solislog.NewLogger(
    solislog.Extra{
        "service": "api",
        "env":     "dev",
    },
    solislog.NewHandler(os.Stdout, solislog.DebugLevel, &solislog.HandlerOptions{
        Template: "<gray>{time}</gray> | <level>{level}</level> | service={extra[service]} env={extra[env]} | {message}\n",
    }),
)

logger.Debug("debug message")
logger.Info("server started")
logger.Warning("slow request")
logger.Error("request failed")
```

Supported colors:

```text
<black>...</black>
<red>...</red>
<green>...</green>
<yellow>...</yellow>
<blue>...</blue>
<magenta>...</magenta>
<cyan>...</cyan>
<white>...</white>
<gray>...</gray>
```

Special color tag:

```text
<level>...</level>
```

`<level>` chooses a color based on the record level:

```text
DEBUG   gray
INFO    cyan
WARNING yellow
ERROR   red
FATAL   magenta
```

Example:

```text
<level>{level}</level> | {message}
```

## Templates

Built-in fields:

```text
{time}
{level}
{message}
{extra}
{file}
{path}
{line}
{function}
{caller}
```

Extra fields:

```text
{extra[source]}
{extra[id]}
{extra[path]}
```

Template examples:

```text
{time} | {level} | {message}
{caller} | {level} | {message}
{time} | <level>{level}</level> | {message}
<gray>{time}</gray> | <gray>{caller}</gray> | <level>{level}</level> | source={extra[source]} | {message}
{level} | {message} | extra={extra}
```

Escaping is done with `\`:

```text
\<red\>     renders literal <red>
\{level\}   renders literal {level}
```

Invalid templates panic during handler creation. This includes unknown placeholders, unknown colors, empty placeholders, empty extra keys, unclosed placeholders, unclosed color tags, and mismatched color tags.


## Caller metadata

`solislog` can render metadata about the source location where the log call was made.

```text
{file}      base file name, for example main.go
{path}      full source file path
{line}      source line number
{function}  full Go function name
{caller}    compact source location in the form file:line
```

Example:

```go
logger := solislog.NewLogger(
    nil,
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "<gray>{caller}</gray> | <level>{level}</level> | {message}\n",
    }),
)

logger.Info("server started")
```

Example output:

```text
main.go:14 | INFO | server started
```

Caller metadata also works in JSON mode:

```go
logger := solislog.NewLogger(
    nil,
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        JSON:     true,
        Template: "{level} {message} {caller} {file} {line} {function}",
    }),
)
```

## Extra fields

Extra fields are stored as:

```go
type Extra map[string]string
```

Individual extra fields can be rendered with `{extra[key]}`:

```go
logger := solislog.NewLogger(
    solislog.Extra{
        "service": "api",
        "env":     "dev",
    },
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "<level>{level}</level> | service={extra[service]} env={extra[env]} | {message}\n",
    }),
)

logger.Info("server started")
```

The full extra map can be rendered with `{extra}`:

```go
logger := solislog.NewLogger(
    solislog.Extra{
        "service": "api",
        "env":     "dev",
    },
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "{level} | {message} | extra={extra}\n",
    }),
)

logger.Info("hello")
```

Example output:

```text
INFO | hello | extra={"env":"dev","service":"api"}
```

## Binding extra fields

Use `Bind(...)` to create a child logger with additional or overridden extra fields.

```go
base := solislog.NewLogger(
    solislog.Extra{
        "service": "api",
    },
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "<level>{level}</level> | service={extra[service]} request_id={extra[request_id]} | {message}\n",
    }),
)

requestLogger := base.Bind(solislog.Extra{
    "request_id": "req-123",
})

requestLogger.Info("request received")
base.Info("base logger still has no request_id")
```

`Bind(...)` does not copy or replace handlers. The child logger uses the same shared core and only changes the attached extra fields.

If a key already exists, the bound value overrides it for the child logger only.

## Contextual logging

`Contextualize(...)` creates a bound logger and stores it in `context.Context`.

This is useful at request, update, job, or operation boundaries.

```go
package main

import (
    "context"
    "os"

    "github.com/DasKaroWow/solislog"
)

func main() {
    base := solislog.NewLogger(
        solislog.Extra{
            "service": "api",
        },
        solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
            Template: "<level>{level}</level> | service={extra[service]} request_id={extra[request_id]} user_id={extra[user_id]} | {message}\n",
        }),
    )

    requestLogger := base.Bind(solislog.Extra{
        "request_id": "req-123",
    })

    ctx := context.Background()
    ctx = requestLogger.Contextualize(ctx, solislog.Extra{
        "user_id": "42",
    })

    handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
    logger, ok := solislog.FromContext(ctx)
    if !ok {
        return
    }

    logger.Info("request received")
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    logger, ok := solislog.FromContext(ctx)
    if !ok {
        return
    }

    logger.Info("processing request")
}
```

## JSON output

Set `HandlerOptions.JSON` to `true` to render records as JSON.

In JSON mode, the template is used as a field list. Plain text is ignored. Only placeholders become JSON fields.

```go
loc, err := time.LoadLocation("Europe/Helsinki")
if err != nil {
    panic(err)
}

logger := solislog.NewLogger(
    solislog.Extra{
        "service": "api",
        "env":     "dev",
    },
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        JSON:       true,
        TimeFormat: time.RFC3339,
        Location:   loc,
        Template:   "{time} {level} {message} {extra[service]} {extra[env]} {extra}",
    }),
)

logger.Info("json message")
```

Example output:

```json
{"time":"2026-05-05T15:30:00+03:00","level":"INFO","message":"json message","service":"api","env":"dev","extra":{"env":"dev","service":"api"}}
```

JSON field behavior:

```text
{time}       -> "time"
{level}      -> "level"
{message}    -> "message"
{file}       -> "file"
{path}       -> "path"
{line}       -> "line"
{function}   -> "function"
{caller}     -> "caller"
{extra}      -> full extra object
{extra[id]}  -> flat field named "id"
```

For example:

```text
Template: "{level} {extra[id]} {extra}"
```

renders fields like:

```json
{"level":"INFO","id":"123","extra":{"id":"123"}}
```

Color tags are ignored in JSON mode:

```text
<red>{level}</red> <level>{message}</level>
```

is equivalent to:

```text
{level} {message}
```

for JSON output.


## Hooks

Handlers can define hooks for custom per-handler processing.

`BeforeHook` runs before the record is rendered. It receives a mutable `*solislog.Record`, so it can change the message or add extra fields.

```go
logger := solislog.NewLogger(
    solislog.Extra{
        "service": "api",
    },
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "<gray>{caller}</gray> | <level>{level}</level> | service={extra[service]} hook={extra[hook]} | {message}\n",
        BeforeHook: func(record *solislog.Record) {
            record.Message = strings.ToUpper(record.Message)
            record.Extra["hook"] = "before"
        },
    }),
)

logger.Info("hook changed this message")
```

`AfterHook` runs after rendering and receives both the record and the rendered log line.

```go
logger := solislog.NewLogger(
    nil,
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "{level} | {message}\n",
        AfterHook: func(record *solislog.Record, rendered string) {
            // Count metrics, inspect the final line, or mirror it elsewhere.
            _ = rendered
        },
    }),
)
```

When a logger has multiple handlers, each handler gets its own cloned record before running `BeforeHook`. A hook on one handler does not mutate the record used by another handler.

`AfterHook` callbacks run after the logger unlocks its shared core, so hooks can safely log again if needed.

## Write error handling

`ErrorHandler` can be used to observe write errors from a handler's `io.Writer`.

```go
logger := solislog.NewLogger(
    nil,
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "{level} | {message}\n",
        ErrorHandler: func(err error, rendered string) {
            // Handle or report the write error.
            _ = err
            _ = rendered
        },
    }),
)
```

If `ErrorHandler` is nil, write errors are ignored.

## Multiple handlers

A single logger can write the same record through multiple handlers. Each handler has its own writer, level, template, time settings, and output mode.

```go
logger := solislog.NewLogger(
    solislog.Extra{
        "service": "api",
    },
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "<level>{level}</level> | {message}\n",
    }),
    solislog.NewHandler(os.Stdout, solislog.ErrorLevel, &solislog.HandlerOptions{
        Template: "<red>{level}</red> | service={extra[service]} | {message}\n",
    }),
)

logger.Info("server started")
logger.Error("request failed")
```

The first handler receives `INFO` and above. The second handler receives only `ERROR` and above.

## Time format and location

Each handler can configure its own time format and location.

```go
loc, err := time.LoadLocation("Europe/Helsinki")
if err != nil {
    panic(err)
}

logger := solislog.NewLogger(
    nil,
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template:   "{time} | {level} | {message}\n",
        TimeFormat: time.DateTime,
        Location:   loc,
    }),
)

logger.Info("hello")
```

`TimeFormat` uses Go's standard time layout system.

Defaults:

```text
Template   = "{time} | {level} | {message}\n"
TimeFormat = time.RFC3339
Location   = time.Local
JSON       = false
```

## Log levels

Current levels:

```go
solislog.DebugLevel
solislog.InfoLevel
solislog.WarningLevel
solislog.ErrorLevel
solislog.FatalLevel
```

A handler writes records whose level is equal to or higher than the handler's configured level.

```go
logger := solislog.NewLogger(
    nil,
    solislog.NewHandler(os.Stdout, solislog.WarningLevel, &solislog.HandlerOptions{
        Template: "<level>{level}</level> | {message}\n",
    }),
)

logger.Info("ignored")
logger.Warning("written")
logger.Error("written")
```

`Fatal(...)` logs with `FatalLevel` and then exits the process with status code `1`.

## Concurrent use

`Logger` methods are safe to call from multiple goroutines.

A base logger and all loggers created from it with `Bind(...)` share the same core, so their writes are serialized through that shared core.

```go
logger := solislog.NewLogger(
    nil,
    solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
        Template: "{level} | {message}\n",
    }),
)

go logger.Info("from goroutine 1")
go logger.Info("from goroutine 2")
```

This guarantee applies to loggers that share the same `solislog` core.

If the same raw `io.Writer` is manually shared between completely separate logger instances, synchronization of that shared writer is still the caller's responsibility.

## Handler options

`NewHandler` requires an output writer and a level. Everything else is configured through `HandlerOptions`.

```go
type HandlerOptions struct {
    Template     string
    TimeFormat   string
    Location     *time.Location
    JSON         bool
    ErrorHandler ErrorHandlerFunc
    BeforeHook   BeforeHookFunc
    AfterHook    AfterHookFunc
}
```

A handler accepts any `io.Writer`, so file logging, buffers, custom writers, and rotation wrappers can be provided outside of `solislog`.

```go
handler := solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
    Template:   "{time} | <level>{level}</level> | {message}\n",
    TimeFormat: time.RFC3339,
    Location:   time.Local,
    JSON:       false,
})
```

## Running the demo

```bash
go run ./demo
```

On Windows:

```powershell
go run .\demo\.
```

The demo is split by topic:

```text
demo/
├── example.go
├── text.go
├── context.go
├── json.go
└── hooks.go
```

## Running tests

```bash
go test ./...
```

Tests cover:

- extra cloning and merging
- handler defaults and options
- tokenizer behavior
- template parsing
- color tags
- text rendering
- JSON rendering
- JSON ignoring colors
- caller metadata in text and JSON
- hooks
- write error handling
- concurrent logger usage

## Project structure

```text
.
├── context.go
├── extra.go
├── handler.go
├── level.go
├── lexer.go
├── logger.go
├── record.go
├── template.go
├── template_stack.go
├── extra_test.go
├── handler_test.go
├── lexer_test.go
├── logger_color_test.go
├── logger_concurrency_test.go
├── logger_handlingErrors_test.go
├── logger_hooks_test.go
├── logger_json_test.go
├── logger_metadata_test.go
├── logger_test.go
├── template_test.go
└── demo/
    ├── example.go
    ├── text.go
    ├── context.go
    ├── json.go
    └── hooks.go
```

## Current limitations

The project is still intentionally small. The following features are not part of the current version:

- no file rotation helper yet
- no middleware helpers yet
- no async logging with queues or workers
- no advanced template formatting or alignment

File rotation and other output-specific behavior can already be provided through custom `io.Writer` implementations.

## Roadmap

Near-term ideas:
- optional file rotation wrapper around `io.Writer`
- middleware helpers, for example Fiber integration

Not planned for now:

- async logging with queues or workers
- complex structured field types
- advanced template formatting or alignment

## License

MIT License.