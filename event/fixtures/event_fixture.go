/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package fixtures holds the reusable doubles that the event component's tests
// build on. It mirrors the source tree, and each type carries the `Fixture`
// suffix.
package fixtures

import (
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
)

// The identifiers of the events that these fixtures stand for.
const (
	// EventFixtureID identifies EventFixture.
	EventFixtureID = "Valkyrja.Tests.Event.EventFixture"

	// StoppableEventFixtureID identifies StoppableEventFixture.
	StoppableEventFixtureID = "Valkyrja.Tests.Event.StoppableEventFixture"

	// ArgumentsCapableEventFixtureID identifies ArgumentsCapableEventFixture.
	ArgumentsCapableEventFixtureID = "Valkyrja.Tests.Event.ArgumentsCapableEventFixture"

	// DispatchCollectableEventFixtureID identifies
	// DispatchCollectableEventFixture.
	DispatchCollectableEventFixtureID = "Valkyrja.Tests.Event.DispatchCollectableEventFixture"
)

// EventFixture is an event that carries nothing but its identifier.
type EventFixture struct{}

// GetEventID returns the identifier of the event.
func (e *EventFixture) GetEventID() string {
	return EventFixtureID
}

// StoppableEventFixture is an event that stops the listeners after it once a
// listener sets Stopped.
type StoppableEventFixture struct {
	Stopped bool
}

// GetEventID returns the identifier of the event.
func (e *StoppableEventFixture) GetEventID() string {
	return StoppableEventFixtureID
}

// IsPropagationStopped reports whether a listener stopped the propagation.
func (e *StoppableEventFixture) IsPropagationStopped() bool {
	return e.Stopped
}

// ArgumentsCapableEventFixture is an event that the dispatcher fills with the
// arguments that the caller gave.
type ArgumentsCapableEventFixture struct {
	Arguments []any
}

// GetEventID returns the identifier of the event.
func (e *ArgumentsCapableEventFixture) GetEventID() string {
	return ArgumentsCapableEventFixtureID
}

// SetArguments returns a copy of the event that holds the arguments.
func (e *ArgumentsCapableEventFixture) SetArguments(arguments []any) contract.EventContract {
	copied := *e
	copied.Arguments = arguments

	return &copied
}

// DispatchCollectableEventFixture is an event that collects what each listener
// returned.
type DispatchCollectableEventFixture struct {
	Dispatches []any
}

// GetEventID returns the identifier of the event.
func (e *DispatchCollectableEventFixture) GetEventID() string {
	return DispatchCollectableEventFixtureID
}

// AddDispatch records what one listener returned.
func (e *DispatchCollectableEventFixture) AddDispatch(dispatch any) {
	e.Dispatches = append(e.Dispatches, dispatch)
}

// GetDispatches returns what each listener returned, in the order that the
// dispatcher ran them.
func (e *DispatchCollectableEventFixture) GetDispatches() []any {
	return e.Dispatches
}

// NotAnEventFixture is a value that the container resolves and the dispatcher
// rejects, because it is no event.
type NotAnEventFixture struct{}

// ListenerProviderFixture is a listener provider that registers the listeners it
// holds.
type ListenerProviderFixture struct {
	Listeners []contract.ListenerContract
}

// GetListeners returns each listener that the provider registers.
func (p *ListenerProviderFixture) GetListeners() []contract.ListenerContract {
	return p.Listeners
}
