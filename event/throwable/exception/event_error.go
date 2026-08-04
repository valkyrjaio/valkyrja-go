/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds every error that the event component raises.
//
// The component always ships a runtime base and an invalid-argument base, even
// where it raises neither on its own. Each concrete error embeds one of them.
package exception

import (
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

// EventRuntimeError is the event component's base runtime error.
type EventRuntimeError struct {
	throwableexception.ValkyrjaRuntimeError
}

// IsEventThrowable marks the error as one that the event component raised.
func (e *EventRuntimeError) IsEventThrowable() bool {
	return true
}

// EventInvalidArgumentError is the event component's base invalid-argument
// error.
type EventInvalidArgumentError struct {
	throwableexception.ValkyrjaInvalidArgumentError
}

// IsEventThrowable marks the error as one that the event component raised.
func (e *EventInvalidArgumentError) IsEventThrowable() bool {
	return true
}

// newEventInvalidArgumentError builds the event component's base
// invalid-argument error.
func newEventInvalidArgumentError(message string, cause error) EventInvalidArgumentError {
	return EventInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(message, cause),
	}
}
