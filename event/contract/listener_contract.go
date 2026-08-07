/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

type ListenerHandlerFunc func(container containercontract.ContainerContract, arguments map[string]any) any

type ListenerFactory func() ListenerContract

type ListenerContract interface {
	// GetEventID returns the identifier of the event that the listener listens
	// for.
	GetEventID() string

	// WithEventID returns a copy of the listener for another event.
	WithEventID(eventID string) ListenerContract

	// GetName returns the unique name of the listener. The collection files the
	// listener under this name.
	GetName() string

	// WithName returns a copy of the listener under another name.
	WithName(name string) ListenerContract

	// GetHandler returns what the listener runs.
	GetHandler() ListenerHandlerFunc

	// WithHandler returns a copy of the listener that runs another handler.
	WithHandler(handler ListenerHandlerFunc) ListenerContract
}

type EventDataContract interface {
	// GetEvents returns each event identifier with the names of its listeners,
	// in the order that the collection recorded them.
	GetEvents() map[string][]string

	// GetListeners returns a factory for each listener, keyed by the listener's
	// name.
	GetListeners() map[string]ListenerFactory
}

//nolint:interfacebloat // Parity with the PHP reference implementation.
type ListenerCollectionContract interface {
	// GetData returns the collection's state.
	GetData() EventDataContract

	// SetFromData replaces the collection's state.
	SetFromData(data EventDataContract)

	// HasListener reports whether the collection holds the listener.
	HasListener(listener ListenerContract) bool

	// HasListenerByID reports whether the collection holds a listener under the
	// name.
	HasListenerByID(listenerID string) bool

	// AddListener files the listener under its own name, for its own event.
	AddListener(listener ListenerContract)

	// RemoveListener removes the listener from its own event.
	RemoveListener(listener ListenerContract)

	// RemoveListenerByID removes the listener from every event.
	RemoveListenerByID(listenerID string)

	// HasListenersForEvent reports whether the collection holds a listener for
	// the event.
	HasListenersForEvent(event EventContract) bool

	// HasListenersForEventByID reports whether the collection holds a listener
	// for the event identifier.
	HasListenersForEventByID(eventID string) bool

	// GetListenersForEvent returns each listener for the event, in the order
	// that the collection recorded them.
	GetListenersForEvent(event EventContract) []ListenerContract

	// GetListenersForEventByID returns each listener for the event identifier,
	// in the order that the collection recorded them.
	GetListenersForEventByID(eventID string) []ListenerContract

	// SetListenersForEvent files each listener for the event.
	SetListenersForEvent(event EventContract, listeners ...ListenerContract)

	// SetListenersForEventByID files each listener for the event identifier.
	SetListenersForEventByID(eventID string, listeners ...ListenerContract)

	// RemoveListenersForEvent removes every listener for the event.
	RemoveListenersForEvent(event EventContract)

	// RemoveListenersForEventByID removes every listener for the event
	// identifier.
	RemoveListenersForEventByID(eventID string)

	// GetListeners returns every listener that the collection holds.
	GetListeners() []ListenerContract

	// GetEvents returns the identifier of each event that has a listener.
	GetEvents() []string

	// GetEventsWithListeners returns each event identifier with its listeners.
	GetEventsWithListeners() map[string][]ListenerContract
}

type EventDispatcherContract interface {
	// Dispatch runs each listener for the event and returns the event.
	Dispatch(event EventContract) EventContract

	// DispatchIfHasListeners runs each listener for the event where the event
	// has one, and returns the event either way.
	DispatchIfHasListeners(event EventContract) EventContract

	// DispatchByID builds the event from the container and dispatches it. It
	// reports a failure where the container resolves no event for the
	// identifier.
	DispatchByID(eventID string, arguments []any) (EventContract, error)

	// DispatchByIDIfHasListeners builds the event from the container and
	// dispatches it where the event has a listener.
	DispatchByIDIfHasListeners(eventID string, arguments []any) (EventContract, error)

	// DispatchListeners runs each listener that it receives, and stops early
	// where a stoppable event stops the propagation.
	DispatchListeners(event EventContract, listeners ...ListenerContract) EventContract

	// DispatchListener runs one listener and records what it returned on an
	// event that collects it.
	DispatchListener(event EventContract, listener ListenerContract) EventContract
}

type ListenerProviderContract interface {
	// GetListeners returns each listener that the component registers.
	GetListeners() []ListenerContract
}
