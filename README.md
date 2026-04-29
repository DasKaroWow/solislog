# solislog

`solislog` is a small template-based contextual logger for Go, inspired by the developer experience of Python's Loguru while keeping the API simple and Go-friendly.

The project focuses on readable human-oriented logs, contextual fields, multiple output handlers, optional JSON output, and a small API that stays close to normal Go patterns.

## Features

* Multiple handlers per logger
* Per-handler log level filtering
* Per-handler templates
* Handler options for template, time format, time location, and JSON mode
* Built-in template fields: `{time}`, `{level}`, `{message}`, `{extra}`
* Custom contextual fields through `{extra[key]}`
* Optional JSON output mode using the same template placeholders as field selection
* `Bind(...)` for creating child loggers with merged extra fields
* `Contextualize(...)` and `FromContext(...)` for passing loggers through `context.Context`
* Simple log methods: `Debug`, `Info`, `Warning`, `Error`, `Fatal`
* Logger methods are safe for concurrent use by multiple goroutines

## Status

`solislog` is currently an early design-stage library. The current model is:

```text
Logger = shared core + extra
shared core = handlers
Handler = writer + level + options/template/render mode
Bind = same core + merged extra
Contextualize = Bind + context.Context
```

The shared core serializes handler access, so a base logger and all loggers created from it with `Bind(...)` can be used safely from multiple goroutines.

The goal is not to compete with `zap`, `zerolog`, or `slog` on performance. The goal is to explore a small logger with pleasant developer experience, readable output, contextual logging, and simple structured output when needed.

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
2026-04-30T00:35:19+03:00 | INFO | hello from solislog
```

## Core concepts

### Logger

A `Logger` stores default `extra` fields and points to a shared core. The shared core contains handlers and synchronizes access to them.

```go
logger := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
	},
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{time} | {level} | {extra[source]} | {extra[id]} | {message}\n",
	}),
)
```

### Handler

A `Handler` defines where records are written, which level it accepts, and how records are rendered.

```go
handler := solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
	Template:   "{time} | {level} | {message}\n",
	TimeFormat: time.RFC3339,
	Location:   time.Local,
})
```

A handler accepts any `io.Writer`, so file logging, buffers, custom writers, and rotation wrappers can be provided outside of `solislog`.

### Handler options

`NewHandler` requires an output writer and a level. Everything else is configured through `HandlerOptions`.

```go
type HandlerOptions struct {
	Template   string
	TimeFormat string
	Location   *time.Location
	JSON       bool
}
```

Defaults:

```text
Template   = "{time} | {level} | {message}\n"
TimeFormat = time.RFC3339
Location   = time.Local
JSON       = false
```

Use `nil` when the defaults are enough:

```go
solislog.NewHandler(os.Stdout, solislog.InfoLevel, nil)
```

## Extra fields

Extra fields are stored as:

```go
type Extra map[string]string
```

Individual extra fields can be referenced with `{extra[key]}`.

```go
logger := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
	},
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{time} | {level} | source={extra[source]} | id={extra[id]} | {message}\n",
	}),
)

logger.Info("base message")
```

The full extra map can be referenced with `{extra}`.

```go
logger := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "123",
	},
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{level} | {message} | extra={extra}\n",
	}),
)

logger.Info("hello")
```

Example output:

```text
INFO | hello | extra={"id":"123","source":"telegram"}
```

## Binding extra fields

Use `Bind(...)` to create a child logger with additional or overridden extra fields.

```go
base := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
	},
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "{time} | {level} | source={extra[source]} | id={extra[id]} | {message}\n",
	}),
)

base.Info("base message")

requestLogger := base.Bind(solislog.Extra{
	"id": "123",
})

requestLogger.Info("request message")
base.Info("base message again")
```

`Bind(...)` does not copy or replace handlers. The child logger uses the same shared core and only changes the attached extra fields.

## Contextual logging

`Contextualize(...)` creates a bound logger and stores it in `context.Context`. This is useful at request, update, job, or operation boundaries.

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
			"source": "telegram",
			"id":     "-1",
		},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{time} | {level} | {extra[source]} | {extra[id]} | {message}\n",
		}),
	)

	ctx := context.Background()
	ctx = base.Contextualize(ctx, solislog.Extra{
		"id": "123",
	})

	handle(ctx)
}

func handle(ctx context.Context) {
	logger, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}

	logger.Info("entered handle")
	process(ctx)
}

func process(ctx context.Context) {
	logger, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}

	logger.Info("processing request")
}
```

## Multiple handlers

A single logger can write the same record through multiple handlers. Each handler has its own writer, level, options, and render mode.

```go
logger := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
		"path":   "/unknown",
	},
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "handler 1 -> {time} | {level} | source={extra[source]} | id={extra[id]} | {message}\n",
	}),
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
		Template: "handler 2 -> {time} | {level} | source={extra[source]} | path={extra[path]} | {message}\n",
	}),
)

logger.Info("base message")

requestLogger := logger.Bind(solislog.Extra{
	"id":   "123",
	"path": "/api/users",
})

requestLogger.Info("request message")
```

Outputs are handled by `Handler` values, while contextual data is handled by `Bind(...)` and `Contextualize(...)`.

## JSON output

Set `HandlerOptions.JSON` to `true` to render records as JSON.

In JSON mode, `Template` is used as a field list. Only placeholders are used; plain text between placeholders is ignored. This lets the user choose which fields are included and in what order.

```go
loc, err := time.LoadLocation("Europe/Helsinki")
if err != nil {
	panic(err)
}

logger := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
	},
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
		Template:   "{time} {level} {message} {extra[id]} {extra}",
		JSON:       true,
		TimeFormat: time.RFC3339,
		Location:   loc,
	}),
)

logger.Info("base message")

requestLogger := logger.Bind(solislog.Extra{
	"id":   "123",
	"path": "/api/users",
})

requestLogger.Info("request message")
```

Example output:

```json
{"time":"2026-04-30T00:35:19+03:00","level":"INFO","message":"base message","id":"-1","extra":{"id":"-1","source":"telegram"}}
{"time":"2026-04-30T00:35:19+03:00","level":"INFO","message":"request message","id":"123","extra":{"id":"123","path":"/api/users","source":"telegram"}}
```

JSON field behavior:

```text
{time}       -> "time"
{level}      -> "level"
{message}    -> "message"
{extra}      -> full extra object
{extra[id]}  -> single flat field named "id"
```

For example, `{extra}` becomes a nested JSON object, while `{extra[id]}` becomes a top-level field named `id`.

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

`TimeFormat` uses Go's standard time layout system. If no format is provided, `time.RFC3339` is used. If no location is provided, `time.Local` is used.

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

This guarantee applies to loggers that share the same `solislog` core. If the same raw `io.Writer` is manually shared between completely separate logger instances, synchronization of that shared writer is still the caller's responsibility.

## Template syntax

Built-in fields:

```text
{time}
{level}
{message}
{extra}
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
{time} | {level} | source={extra[source]} | {message}
{time} | {level} | source={extra[source]} | id={extra[id]} | {message}
{time} | {level} | {message} | extra={extra}
```

Unknown built-in fields, empty placeholders, empty extra keys, unclosed placeholders, and unexpected closing braces currently panic during template parsing.

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

`Fatal(...)` logs with `FatalLevel` and then exits the process with status code `1`.

## Running the demo

```bash
go run ./demo
```

On Windows, this also works:

```bash
go run .\demo\.
```

## Running tests

```bash
go test ./...
```

Current tests cover extra cloning/merging, handler defaults/options, template parsing, regular template output, full `{extra}` output, JSON output, and a public smoke test for logger writing.

## Project structure

```text
.
├── context.go
├── extra.go
├── handler.go
├── level.go
├── logger.go
├── record.go
├── template.go
├── extra_test.go
├── handler_test.go
├── json_test.go
├── logger_test.go
├── template_test.go
└── demo/
    ├── 1.go
    ├── 2.go
    ├── 3.go
    ├── 4.go
    ├── 5.go
    └── example.go
```

## Current limitations

The project is still intentionally small. The following features are not part of the current version:

* no colorized console output yet
* no caller/file/line fields yet
* no hooks
* no file rotation helper yet
* no middleware helpers yet
* no write error handling API yet

File rotation and other output-specific behavior can already be provided through custom `io.Writer` implementations.

## Roadmap

Near-term goals:

* colorized console output
* caller/file/line fields
* JSON output refinements if needed

Later ideas:

* hooks
* optional file rotation wrapper around `io.Writer`
* middleware helpers, for example Fiber integration
* write error handling strategy

Not planned:

* async logging with queues or workers
* advanced template formatting or alignment

## License

MIT License.
