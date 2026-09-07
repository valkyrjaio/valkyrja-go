# Event

## Introduction

The Event component dispatches an event to the listeners that registered for it.
A listener provider registers the listeners of one component, a collection files
them, and the dispatcher runs them in the order that the collection recorded.

## Identifying an Event

The other ports identify an event by its class: PHP reads `$event::class`, and
Java holds a `Class<?>`. Go has no class, so an event states its own identifier,
and that identifier is the key the collection files a listener under. It takes
the binding-key format.

```go
// Right — the event names itself, and nothing reads its type to find the name.
func (e *CacheClearedEvent) GetEventID() string {
	return "valkyrja.app.event.CacheClearedEvent"
}
```

## Event Contracts

An event implements the contracts of the behavior that it wants. Each one is
optional beyond the first.

| Contract                           | Gives the event                          |
| :--------------------------------- | :--------------------------------------- |
| `EventContract`                    | Its identifier                           |
| `StoppableEventContract`           | The power to stop the listeners after it |
| `ArgumentsCapableEventContract`    | The arguments that the caller gave       |
| `DispatchCollectableEventContract` | What each listener returned              |

The PHP port takes the stoppable contract from PSR-14. Go has no PSR, so the
framework declares it.

## Listeners

A listener names the event that it listens for, carries a unique name that the
collection files it under, and holds what it runs:

```go
type ListenerHandlerFunc func(
	container containercontract.ContainerContract,
	arguments map[string]any,
) any
```

A handler receives the container and the arguments. It reads whatever else it
needs from the container, rather than from its own signature.

Every `With` method returns a copy and leaves the receiver unchanged. The other
ports return `static`; Go has no such return type, so each one returns the
contract.

## The Collection

The collection holds a `ListenerFactory` for each listener rather than the
listener itself:

```go
type ListenerFactory func() ListenerContract
```

The generated cache then states how to build a listener, rather than holding one.
That is what lets `sindri` write the collection as source.

## Registering a Listener

A listener provider returns a literal slice. `sindri` reads the slice from the
source rather than by running it, so a computed slice generates nothing.

```go
func (p *CacheEventListenerProvider) GetListeners() []contract.ListenerContract {
	return []contract.ListenerContract{
		data.NewListener(cacheClearedEventID, "cache.cleared.warm", warmHandler),
	}
}
```

The collection turns each one into a factory as it files it, so the generated
cache states how to build a listener rather than holding one.

## Dispatching

```go
dispatched := dispatcher.Dispatch(event)
```

The dispatcher resolves each listener from the container and runs it in order. A
stoppable event that stops propagation ends the run, and a collectable event
holds what each listener returned.

Resolving from the container rather than by reflection is what lets the component
work with no cache at all.

## Service Registration

`EventServiceProvider` publishes these binding keys:

| Binding key                                            | Holds               |
| :----------------------------------------------------- | :------------------ |
| `valkyrja.event.data.EventData`                        | The component state |
| `valkyrja.event.collection.ListenerCollectionContract` | Every listener      |
| `valkyrja.event.dispatcher.EventDispatcherContract`    | The dispatcher      |

The state is a value that the framework stores and reloads. `sindri` generates it
for an application, so the framework loads every listener at boot without reading
a provider.
