# Application

## Introduction

The Application component is the root of a Valkyrja application. It holds the
container, the configuration, and the provider tree, and it walks that tree to
collect every provider that each component registers.

## The Kernel

`kernel.NewValkyrja` builds the application over a container and a configuration:

```go
app := kernel.NewValkyrja(manager.NewContainer(nil), data.NewConfig(
	&httpprovider.HttpComponentProvider{},
	&cliprovider.CliComponentProvider{},
))
```

The application reads the timezone that the configuration names as it starts.

`ChildApplication` wraps a parent and reads the parent's providers. A worker that
serves several requests in one process builds one, so each request gets its own
container while the provider tree stays shared.

## The Provider Tree

A component provider names the components that it needs, and returns each
provider that its own component registers:

| Method                  | Returns                            |
| :---------------------- | :--------------------------------- |
| `GetComponentProviders` | Each component that this one needs |
| `GetContainerProviders` | Each service provider              |
| `GetEventProviders`     | Each listener provider             |
| `GetCliProviders`       | Each CLI route provider            |
| `GetHttpProviders`      | Each HTTP route provider           |

The application walks the tree once, and it collects a provider that two
components name only once. `GetProviders` returns every component provider,
including the ones a named component pulled in.

The other ports declare a CLI variant and an HTTP variant that override
`getComponentProviders`. Go has no method override, so each variant is a provider
of its own.

## Configuration

`data.Config` is the configuration that every application holds. The other ports
give the constructor a default for each parameter; Go has no default parameter,
so `NewConfig` takes the providers and every other field takes its own default. A
caller changes a field afterwards.

| Field           | Default                             |
| :-------------- | :---------------------------------- |
| `Namespace`     | `App`                               |
| `Dir`           | The directory the process runs from |
| `Version`       | The framework version               |
| `Environment`   | `production`                        |
| `DebugMode`     | `false`                             |
| `Timezone`      | `UTC`                               |
| `Key`           | `some_secret_app_key`               |
| `DataPath`      | `App/Provider/Data`                 |
| `DataNamespace` | `App/Provider/Data`                 |

`CliConfig` and `HttpConfig` embed it and add the settings of their protocol.
Each one names the middleware that runs at every stage, by binding key, and each
list starts empty.

```go
config := data.NewHttpConfig(&httpprovider.HttpComponentProvider{})
config.DebugMode = true
config.RequestReceivedMiddleware = []string{"App.Http.Middleware.TrustProxies"}
```

Warning: a response carries what went wrong only in debug mode. A production
application leaves `DebugMode` false, so a failure never tells a client about the
inside of the application.

## Directories

`application/directory` reads a path inside the application. A caller builds one
over the base directory, and reads each path from it:

```go
built := directory.NewDirectory(config.GetDir())

path := built.GetAppDirectory("Http/Controller")
```

| Method             | Returns                                      |
| :----------------- | :------------------------------------------- |
| `GetPath`          | A path under the base directory              |
| `GetBaseDirectory` | A path under the base directory              |
| `GetAppDirectory`  | A path under the application's own directory |
| `GetDataDirectory` | A path under the generated data directory    |

## Application Information

`application/constant` holds the framework's own version, its build date, and the
terminal art that the CLI prints:

| Constant               | Holds                                    |
| :--------------------- | :--------------------------------------- |
| `Version`              | The framework version                    |
| `VersionBuildDateTime` | When the release was built               |
| `Ascii`                | The word mark, drawn for a terminal      |
| `Icon`                 | The framework mark, drawn for a terminal |

## Service Registration

`ApplicationComponentProvider` names the container, which every application
needs, and registers no service of its own.

| Binding key                                       | Holds                  |
| :------------------------------------------------ | :--------------------- |
| `Valkyrja.Application.Kernel.ApplicationContract` | The application        |
| `Valkyrja.Application.Data.ConfigContract`        | The configuration      |
| `Valkyrja.Application.Data.CliConfigContract`     | The CLI configuration  |
| `Valkyrja.Application.Data.HttpConfigContract`    | The HTTP configuration |

A component reads its own settings from these keys. The container component's
publisher registers every service provider that the application names.
