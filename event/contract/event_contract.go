/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds every contract of the event component.
//
// The component keeps one `contract` package, for the reason that the container
// component keeps one: Go resolves an import cycle at the package level, and the
// contracts here name each other.
package contract

// EventContract is a thing that the dispatcher dispatches.
//
// The other ports identify an event by its class: PHP reads `$event::class` and
// Java holds a `Class<?>`. Go has no class, so an event states its own
// identifier. That identifier is the key that the collection files a listener
// under, and it takes the binding-key format, `Valkyrja.{Component}.{SubComponent}.{Name}`.
type EventContract interface {
	// GetEventID returns the identifier of the event.
	GetEventID() string
}

// StoppableEventContract is an event that stops the listeners after it.
//
// The PHP port takes this contract from PSR-14. Go has no PSR, so the framework
// declares it.
type StoppableEventContract interface {
	EventContract

	// IsPropagationStopped reports whether the dispatcher stops before it runs
	// the next listener.
	IsPropagationStopped() bool
}

// ArgumentsCapableEventContract is an event that the dispatcher fills with the
// arguments that the caller gave.
type ArgumentsCapableEventContract interface {
	EventContract

	// WithArguments returns a copy of the event that holds the arguments.
	WithArguments(arguments []any) EventContract
}

// DispatchCollectableEventContract is an event that collects what each listener
// returned.
type DispatchCollectableEventContract interface {
	EventContract

	// AddDispatch records what one listener returned.
	AddDispatch(dispatch any)

	// GetDispatches returns what each listener returned, in the order that the
	// dispatcher ran them.
	GetDispatches() []any
}
