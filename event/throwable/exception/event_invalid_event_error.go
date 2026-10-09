/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

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
