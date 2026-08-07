# Container

## Introduction

The Container component resolves a service from a binding key. It holds a
factory, a singleton, or an alias under each key, and a service provider
registers what a component publishes.

Every language works without cache. A provider exposes its publishers as
interface methods, so the framework walks the provider tree and registers
everything at run time. Cache is a cold-start optimization, never a correctness
requirement.

## Binding Keys

Go has no `::class`, so a binding key is a string constant. Each component holds
its own constants file, and the format is
`Valkyrja.{Component}.{SubComponent}.{Name}`:

```go
const RouterContractServiceID = "valkyrja.http.routing.dispatcher.RouterContract"
```

## Binding a Service

| Method          | Holds                                            |
| :-------------- | :----------------------------------------------- |
| `Bind`          | A factory that the container calls on each `Get` |
| `BindSingleton` | A factory that the container calls once          |
| `SetSingleton`  | An instance that the container returns as it is  |
| `SetAlias`      | A key that points at another key                 |

```go
container.Bind(constant.RouterContractServiceID, func(
	container contract.ContainerContract,
	arguments []any,
) any {
	return router.NewRouter(container)
})
```

## Resolving a Service

`Get` takes the binding key, the arguments that a factory reads, and what to do
with a key that the container does not hold:

```go
resolved, err := container.Get(id, nil, constant.NewInstanceOrThrowException)
if err != nil {
	return err
}
```

| Mode                          | Where the container holds no such key |
| :---------------------------- | :------------------------------------ |
| `NewInstanceOrThrowException` | Reports a failure                     |
| `NewInstanceOrNull`           | Returns nil                           |

Go reports a failure with a returned error rather than a throw, so every resolve
carries an error alongside the value.

## Service Providers

A service provider returns one publisher per binding key that it defers. The
container runs a publisher the first time something asks for its key, so an
application never builds a service that it does not resolve.

```go
func (p *HttpMessageServiceProvider) Publishers() map[string]contract.PublishFunc {
	return map[string]contract.PublishFunc{
		constant.ResponseFactoryContractServiceID: PublishResponseFactory,
	}
}

func PublishResponseFactory(container contract.ContainerContract) {
	container.SetSingleton(constant.ResponseFactoryContractServiceID, factory.NewResponseFactory())
}
```

Warning: a provider returns a literal map, never a computed one. `sindri` reads
the map from the source rather than by running it, so a computed map generates
nothing.

A publisher is a package-level function rather than a method, which is one of the
two forms that `sindri` reads.

## The Container's State

`GetData` returns the container's state as a value, and `SetFromData` replaces
it. `sindri` generates that value for an application, so the framework loads
every binding at boot without reading a provider.

The other ports type `GetData` with the concrete `ContainerData`. In Go that is a
cycle: `ContainerData` holds a `ServiceFactory`, and a `ServiceFactory` takes a
`ContainerContract`. The contract therefore names an interface, and the `data`
package implements it.

## Service Registration

`ContainerServiceProvider` publishes one binding key:

| Binding key                             | Holds                     |
| :-------------------------------------- | :------------------------ |
| `valkyrja.container.data.ContainerData` | The container's own state |

Its publisher registers every service provider that the application names, then
records the state.
