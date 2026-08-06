/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package dispatcher holds the event dispatcher.
package dispatcher

import (
	containerconstant "github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/throwable/exception"
)

// EventDispatcher runs the listeners for an event.
type EventDispatcher struct {
	collection contract.ListenerCollectionContract
	container  containercontract.ContainerContract
}

// NewEventDispatcher builds a dispatcher over a collection and a container.
func NewEventDispatcher(
	collection contract.ListenerCollectionContract,
	container containercontract.ContainerContract,
) *EventDispatcher {
	return &EventDispatcher{
		collection: collection,
		container:  container,
	}
}

// Dispatch runs each listener for the event and returns the event.
func (d *EventDispatcher) Dispatch(event contract.EventContract) contract.EventContract {
	return d.DispatchListeners(event, d.collection.GetListenersForEvent(event)...)
}

// DispatchIfHasListeners runs each listener for the event where the event has
// one, and returns the event either way.
func (d *EventDispatcher) DispatchIfHasListeners(event contract.EventContract) contract.EventContract {
	if d.collection.HasListenersForEvent(event) {
		return d.Dispatch(event)
	}

	return event
}

// DispatchByID builds the event from the container and dispatches it.
func (d *EventDispatcher) DispatchByID(eventID string, arguments []any) (contract.EventContract, error) {
	event, err := d.getEventFromID(eventID, arguments)
	if err != nil {
		return nil, err
	}

	return d.Dispatch(event), nil
}

// DispatchByIDIfHasListeners builds the event from the container and dispatches
// it where the event has a listener.
func (d *EventDispatcher) DispatchByIDIfHasListeners(
	eventID string,
	arguments []any,
) (contract.EventContract, error) {
	event, err := d.getEventFromID(eventID, arguments)
	if err != nil {
		return nil, err
	}

	if d.collection.HasListenersForEventByID(eventID) {
		return d.Dispatch(event), nil
	}

	return event, nil
}

// DispatchListeners runs each listener that it receives, and stops early where a
// stoppable event stops the propagation.
func (d *EventDispatcher) DispatchListeners(
	event contract.EventContract,
	listeners ...contract.ListenerContract,
) contract.EventContract {
	for _, listener := range listeners {
		event = d.DispatchListener(event, listener)

		stoppable, isStoppable := event.(contract.StoppableEventContract)
		if isStoppable && stoppable.IsPropagationStopped() {
			return event
		}
	}

	return event
}

// DispatchListener runs one listener and records what it returned on an event
// that collects it.
//
// A listener with no handler runs nothing. The other ports type the handler as
// callable, so the value cannot be absent; a Go function value can be nil, and
// calling it panics in the middle of a dispatch.
func (d *EventDispatcher) DispatchListener(
	event contract.EventContract,
	listener contract.ListenerContract,
) contract.EventContract {
	handler := listener.GetHandler()
	if handler == nil {
		return event
	}

	dispatch := handler(d.container, map[string]any{constant.EventArgumentKey: event})

	collectable, isCollectable := event.(contract.DispatchCollectableEventContract)
	if isCollectable {
		collectable.AddDispatch(dispatch)
	}

	return event
}

// getEventFromID builds the event that the binding key names.
//
// PHP and Java build the event from its own class name. Go cannot construct a
// type from a string, so this port resolves the binding key through the
// container, which is the framework's own answer to "build the thing that this
// identifier names". An application binds each event that it dispatches by
// identifier.
func (d *EventDispatcher) getEventFromID(eventID string, arguments []any) (contract.EventContract, error) {
	resolved, err := d.container.Get(eventID, arguments, containerconstant.NewInstanceOrThrowException)
	if err != nil {
		return nil, err
	}

	event, isEvent := resolved.(contract.EventContract)
	if !isEvent {
		return nil, exception.NewEventInvalidEventError(eventID)
	}

	argumentsCapable, isArgumentsCapable := event.(contract.ArgumentsCapableEventContract)
	if isArgumentsCapable {
		return argumentsCapable.WithArguments(arguments), nil
	}

	return event, nil
}
