# Log

## Introduction

The Log component writes a message at each severity that the framework reports.
The severities are the eight of [RFC 5424][rfc url], which every port keeps, so a
message reads the same in each one.

The default logger writes to a stream. A null logger is included, which an
application that reports no message uses, and which a test that must not write to
a file uses as well.

## The Logger Contract

`log/contract.LoggerContract` declares one method per severity, and one that
reports a failure:

```go
type LoggerContract interface {
	Throwable(throwable error, message string)

	Debug(message string, context map[string]any)
	Info(message string, context map[string]any)
	Notice(message string, context map[string]any)
	Warning(message string, context map[string]any)
	Error(message string, context map[string]any)
	Critical(message string, context map[string]any)
	Alert(message string, context map[string]any)
	Emergency(message string, context map[string]any)
}
```

`Throwable` writes a failure at the error severity, with what went wrong in the
context. Use it where a caller reads an error and records it, rather than
building the context itself.

## Severities

`log/constant.LogLevel` is a string type, and each severity is a constant of it.

| Constant            | Value       |
| :------------------ | :---------- |
| `LogLevelDebug`     | `debug`     |
| `LogLevelInfo`      | `info`      |
| `LogLevelNotice`    | `notice`    |
| `LogLevelWarning`   | `warning`   |
| `LogLevelError`     | `error`     |
| `LogLevelCritical`  | `critical`  |
| `LogLevelAlert`     | `alert`     |
| `LogLevelEmergency` | `emergency` |

A caller that reads the contract writes through the method of the severity:

```go
built.Warning("The cache is stale.", nil)
```

Each logger also carries a `Log` method, which writes at the severity that a
caller names and reports a failure where the severity is not one of the eight.
It is not on the contract, so a caller reaches it through the concrete logger
rather than through the container:

```go
err := logger.NewStreamLogger(os.Stderr).Log(constant.LogLevelWarning, "The cache is stale.", nil)
if err != nil {
	return err
}
```

## Loggers

| Type           | Writes                                                |
| :------------- | :---------------------------------------------------- |
| `StreamLogger` | Each message to a stream                              |
| `NullLogger`   | Nothing, and reports an unknown severity the same way |

A message reads `SEVERITY: text`, with the context after it as JSON where the
caller gives one:

```text
INFO: The cache is warm. {"entries":128}
```

Warning: a logger has nowhere to report its own failure. A stream that reports
one therefore changes nothing that a caller can act on, and the logger drops it.
Data that no JSON encoder renders writes an empty object rather than nothing.

## Configuration

The component reads `LogConfigContract` for the setting that applies to the whole
component, and `LogStreamConfigContract` for the settings of the stream logger.
An application config type implements the contracts of the adapters that it uses.

| Contract                  | Method              | Default                            |
| :------------------------ | :------------------ | :--------------------------------- |
| `LogConfigContract`       | `GetDefaultLogger`  | `Valkyrja.Log.Logger.StreamLogger` |
| `LogStreamConfigContract` | `GetStreamFilePath` | `""`, the standard error           |

An empty path writes to the standard error of the process. So does a path that no
process can open, for the reason above.

Every property of an adapter contract carries the adapter name, so one
application config type implements several adapter contracts at once:

```go
// Right — the prefix keeps two adapters from colliding on one property name.
type LogStreamConfigContract interface {
	GetStreamFilePath() string
}
```

## Service Registration

`LogServiceProvider` publishes these binding keys:

| Binding key                           | Holds                       |
| :------------------------------------ | :-------------------------- |
| `Valkyrja.Log.Data.LogConfigContract` | The component configuration |
| `Valkyrja.Log.Logger.LoggerContract`  | The active logger           |

The provider binds the application's own config where it implements the contract,
and the framework default where it does not.

[rfc url]: https://datatracker.ietf.org/doc/html/rfc5424
