# solislog

`solislog` is a small template-based contextual logger for Go, inspired by the developer experience of Python's Loguru while keeping the API simple and Go-friendly.

The project focuses on readable human-oriented logs, contextual fields, multiple output handlers, and a small API that does not force a JSON-first workflow.

## Features

* Multiple handlers per logger
* Per-handler log level filtering
* Per-handler output templates
* Built-in template fields: `{time}`, `{level}`, `{message}`
* Custom contextual fields through `{extra[key]}`
* `Bind(...)` for creating child loggers with merged extra fields
* `Contextualize(...)` and `FromContext(...)` for passing loggers through `context.Context`
* Simple log methods: `Debug`, `Info`, `Warning`, `Error`, `Fatal`
* Logger methods are safe for concurrent use by multiple goroutines

## Status

`solislog` is currently an early design-stage library. The current model is:

```text
Logger = shared core + extra
shared core = handlers
Handler = writer + level + template
Bind = same core + merged extra
Contextualize = Bind + context.Context
```

The shared core serializes handler access, so a base logger and all loggers created from it with `Bind(...)` can be used safely from multiple goroutines.

The goal is not to compete with `zap`, `zerolog`, or `slog` on performance. The goal is to explore a small logger with pleasant developer experience and readable output.

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
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, ""),
	)

	logger.Info("hello from solislog")
}
```

An empty template uses the default format:

```text
{time} | {level} | {message}\n
```

Example output:

```text
2026-04-29T18:45:00+03:00 | INFO | hello from solislog
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
	solislog.NewHandler(
		os.Stdout,
		solislog.InfoLevel,
		"{time} | {level} | {extra[source]} | {extra[id]} | {message}\n",
	),
)
```

### Handler

A `Handler` defines where records are written, which level it accepts, and how records are rendered.

```go
solislog.NewHandler(
	os.Stdout,
	solislog.InfoLevel,
	"{time} | {level} | {message}\n",
)
```

A handler accepts any `io.Writer`, so file logging, buffers, custom writers, and rotation wrappers can be provided outside of `solislog`.

### Extra fields

Extra fields are stored as:

```go
type Extra map[string]string
```

They can be referenced from templates with `{extra[key]}`.

```go
logger := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
	},
	solislog.NewHandler(
		os.Stdout,
		solislog.InfoLevel,
		"{time} | {level} | source={extra[source]} | id={extra[id]} | {message}\n",
	),
)

logger.Info("base message")
```

## Binding extra fields

Use `Bind(...)` to create a child logger with additional or overridden extra fields.

```go
base := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
	},
	solislog.NewHandler(
		os.Stdout,
		solislog.InfoLevel,
		"{time} | {level} | source={extra[source]} | id={extra[id]} | {message}\n",
	),
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
		solislog.NewHandler(
			os.Stdout,
			solislog.InfoLevel,
			"{time} | {level} | {extra[source]} | {extra[id]} | {message}\n",
		),
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

A single logger can write the same record through multiple handlers. Each handler has its own writer, level, and template.

```go
logger := solislog.NewLogger(
	solislog.Extra{
		"source": "telegram",
		"id":     "-1",
		"path":   "/unknown",
	},
	solislog.NewHandler(
		os.Stdout,
		solislog.InfoLevel,
		"handler 1 -> {time} | {level} | source={extra[source]} | id={extra[id]} | {message}\n",
	),
	solislog.NewHandler(
		os.Stdout,
		solislog.InfoLevel,
		"handler 2 -> {time} | {level} | source={extra[source]} | path={extra[path]} | {message}\n",
	),
)

logger.Info("base message")

requestLogger := logger.Bind(solislog.Extra{
	"id":   "123",
	"path": "/api/users",
})

requestLogger.Info("request message")
```

Outputs are handled by `Handler` values, while contextual data is handled by `Bind(...)` and `Contextualize(...)`.

## Concurrent use

`Logger` methods are safe to call from multiple goroutines.

A base logger and all loggers created from it with `Bind(...)` share the same core, so their writes are serialized through that shared core.

```go
logger := solislog.NewLogger(
	nil,
	solislog.NewHandler(os.Stdout, solislog.InfoLevel, "{level} | {message}\n"),
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
```

Unknown built-in fields, empty placeholders, unclosed placeholders, and unexpected closing braces currently panic during template parsing.

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
├── logger_test.go
└── demo/
    ├── 1.go
    ├── 2.go
    ├── 3.go
    ├── 4.go
    └── example.go
```

## Current limitations

The project is still intentionally small. The following features are not part of the current version:

* no color support
* no hooks
* no async logging
* no caller/file/line fields
* no JSON output mode
* no file rotation built in
* no advanced template formatting or alignment
* no middleware helpers
* no write error handling API

File rotation and other output-specific behavior should be provided through custom `io.Writer` implementations.

## Roadmap

Likely next steps:

* add more tests for `Bind`, context helpers, level filtering, template rendering, and multi-handler behavior
* improve concurrency-related tests
* improve README examples as the API stabilizes

## License

MIT License.
