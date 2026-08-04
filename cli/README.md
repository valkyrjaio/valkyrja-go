# Cli

## Introduction

The Cli component runs a command that a caller typed. It has three
sub-components: interaction reads the input and writes the output, routing
matches a command and fills its parameters, and the server is the entry point for
one run.

## Interaction

### Input

`NewInputFromArgs` reads what the caller typed:

```go
built := input.NewInputFromArgs(os.Args)
```

Go's `os.Args` leads with the path that the caller ran, and the command name
follows it. That is the shape that PHP and TypeScript read as well; the Java port
has no caller slot, because the JVM does not pass one.

| The caller typed | The input reads                    |
| :--------------- | :--------------------------------- |
| `--verbose`      | A long option, with no value       |
| `--name=value`   | A long option that carries a value |
| `-q`             | A short option                     |
| `first`          | A positional argument              |

An argument list that names no command gives an input that names none, and the
server then runs the default command.

### Output

An output holds each message, and a writer puts it where the caller reads it.
`WriteMessages` moves every waiting message to the written list.

| Factory method       | Writes                            |
| :------------------- | :-------------------------------- |
| `CreateOutput`       | Through the default writers       |
| `CreatePlainOutput`  | With no format applied            |
| `CreateStreamOutput` | To a stream that the caller names |
| `CreateFileOutput`   | To a file that the caller names   |
| `CreateEmptyOutput`  | Nothing                           |

Warning: a silent output writes nothing at all, and a quiet output writes nothing
while the exit code reports success. A command that must always reach the caller
reports a failure instead of writing.

### Messages and Formats

A message carries text and an optional formatter. A formatter wraps the text in
terminal escape codes.

| Constructor         | Prints                                          |
| :------------------ | :---------------------------------------------- |
| `NewMessage`        | The text, with no format                        |
| `NewErrorMessage`   | White on red                                    |
| `NewSuccessMessage` | White on green                                  |
| `NewWarningMessage` | Black on yellow                                 |
| `NewNewLine`        | A line break alone                              |
| `NewMessages`       | Several messages as one                         |
| `NewBanner`         | The message on a block of its own color         |
| `NewHeader`         | The application, the framework, and the command |

## Routing

### A Command

```go
route := data.NewRoute("cache:clear", "Clear the cache", handler).
	WithHelpText(helpText).
	WithOptions(data.NewOptionParameter("force", "Clear without asking").WithShortNames("f"))
```

A handler receives the container and the route that matched, and returns an
output:

```go
type CliHandlerFunc func(
	container containercontract.ContainerContract,
	route RouteContract,
) OutputContract
```

The route holds a function that builds its help text rather than the message, so
a command builds that text only where the caller asks for it.

### Parameters

An argument parameter takes a positional value; an option parameter takes one
that the caller names. Each one states whether the command needs it, and how many
values it takes.

| Mode                       | The command                       |
| :------------------------- | :-------------------------------- |
| `ArgumentModeRequired`     | Needs the argument                |
| `ArgumentModeOptional`     | Runs without it                   |
| `ArgumentValueModeDefault` | Takes one value at most           |
| `ArgumentValueModeArray`   | Takes as many as the caller gives |

Warning: an argument parameter in the array value mode takes every argument that
is left, so it must be the last one that the command declares. A parameter that
follows it receives nothing.

An option parameter adds short names, a display name for its value, a default
value, and the values that it accepts. An option in the `OptionValueModeNone`
mode takes no value, and one that carries a value reports a failure.

The other ports carry a `Cast` from the Type component, which this port does not
have. A parameter therefore holds the function that converts its values, and that
function reports a failure the way Go does:

```go
type CastFunc func(value string) (any, error)
```

### The Router

The router matches the command by name, fills each parameter from what the caller
typed, and runs the handler. An option reaches a parameter under its long name
and under each short name that the parameter declares.

An input that names no command of the application reaches the route-not-matched
middleware with an output that reports the failure. A parameter that holds a
value the command does not accept ends the run the same way.

## Middleware

Each stage runs its middleware in order, and a middleware receives the handler,
so the middleware decides whether the run continues.

| Stage             | Runs                                        |
| :---------------- | :------------------------------------------ |
| `InputReceived`   | Before the router matches a command         |
| `RouteMatched`    | When the router matches one                 |
| `RouteNotMatched` | When the router matches none                |
| `RouteDispatched` | After the command handler returns           |
| `ThrowableCaught` | When something in the run reports a failure |
| `ProcessExiting`  | Before the process exits                    |

The input-received and route-matched stages return a result contract, because the
TypeScript port returns a union there and Go has no union. The result reports
which side it carries.

Warning: a handler appends each middleware and never dedupes. A middleware that
is added twice runs twice. That is the developer's error, and the framework does
not correct it, because the generated cache must match what the runtime collects.

## The Server

`InputHandler` is the entry point for one input:

```go
handler.Run(input.NewInputFromArgs(os.Args))
```

`Run` handles the input, writes what the command reported, runs the
process-exiting middleware, and exits with the code that the output holds.

A command that panics reaches the caller as an output that reports the failure,
rather than ending the process. The other ports catch a throw here; Go has no
throw, so this recovers a panic instead.

### Built-in Commands

| Command     | Prints                                             |
| :---------- | :------------------------------------------------- |
| `list`      | Every command, in the order of its name            |
| `list:bash` | Every name, separated by a space, for completion   |
| `help`      | The help text of the command that the caller names |
| `version`   | The version of the application                     |

Every command accepts `--quiet`, `--silent`, `--no-interaction`, `--help`, and
`--version`, and each one carries a short name.

## Configuration

The component reads `CliInteractionConfigContract` for the settings that every
output carries, and `application/contract.CliConfigContract` for the middleware
that each stage runs.

| Method          | Default | The output                 |
| :-------------- | :------ | :------------------------- |
| `IsInteractive` | `true`  | Asks the caller a question |
| `IsQuiet`       | `false` | Writes less                |
| `IsSilent`      | `false` | Writes nothing             |

## Service Registration

| Binding key                                                     | Holds                    |
| :-------------------------------------------------------------- | :----------------------- |
| `Valkyrja.Cli.Interaction.Data.CliInteractionConfigContract`    | The interaction config   |
| `Valkyrja.Cli.Interaction.Output.Factory.OutputFactoryContract` | The output factory       |
| `Valkyrja.Cli.Interaction.Input.InputContract`                  | The input of the run     |
| `Valkyrja.Cli.Interaction.Output.OutputContract`                | The output of the run    |
| `Valkyrja.Cli.Routing.Collection.RouteCollectionContract`       | Every command            |
| `Valkyrja.Cli.Routing.Dispatcher.RouterContract`                | The router               |
| `Valkyrja.Cli.Routing.Data.RouteContract`                       | The command that matched |
| `Valkyrja.Cli.Server.Handler.InputHandlerContract`              | The entry point          |

`CliServerCliRoutesProvider` registers the four built-in commands. The other
ports read a command from an annotation on its class; Go has no annotation, so
each command declares its own route and the provider returns them as a literal
slice that `sindri` reads from the source.
