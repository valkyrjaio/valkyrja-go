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

type EventContract interface {
	// GetEventID returns the identifier of the event.
	GetEventID() string
}

type StoppableEventContract interface {
	EventContract

	// IsPropagationStopped reports whether the dispatcher stops before it runs
	// the next listener.
	IsPropagationStopped() bool
}

type ArgumentsCapableEventContract interface {
	EventContract

	// WithArguments returns a copy of the event that holds the arguments.
	WithArguments(arguments []any) EventContract
}

type DispatchCollectableEventContract interface {
	EventContract

	// AddDispatch records what one listener returned.
	AddDispatch(dispatch any)

	// GetDispatches returns what each listener returned, in the order that the
	// dispatcher ran them.
	GetDispatches() []any
}
