/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

// EventInvalidEventError reports a binding key that the container resolves to
// something which is not an event.
//
// No other port raises this. PHP and Java build the event from its own class, so
// the built value is an event by construction. Go builds it through the
// container, which resolves a binding key to any value at all.
type EventInvalidEventError struct {
	EventInvalidArgumentError

	id string
}

// NewEventInvalidEventError builds the error for a binding key.
func NewEventInvalidEventError(id string) *EventInvalidEventError {
	return &EventInvalidEventError{
		EventInvalidArgumentError: newEventInvalidArgumentError(
			"Service with `"+id+"` is not an event",
			nil,
		),
		id: id,
	}
}

// GetID returns the binding key that resolved to something which is not an
// event.
func (e *EventInvalidEventError) GetID() string {
	return e.id
}
