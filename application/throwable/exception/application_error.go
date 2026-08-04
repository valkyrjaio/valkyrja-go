/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds every error that the application component raises.
//
// The component always ships a runtime base and an invalid-argument base, even
// where it raises neither on its own. Each concrete error embeds one of them.
package exception

import (
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

// ApplicationRuntimeError is the application component's base runtime error.
type ApplicationRuntimeError struct {
	throwableexception.ValkyrjaRuntimeError
}

// IsApplicationThrowable marks the error as one that the application raised.
func (e *ApplicationRuntimeError) IsApplicationThrowable() bool {
	return true
}

// ApplicationInvalidArgumentError is the application component's base
// invalid-argument error.
type ApplicationInvalidArgumentError struct {
	throwableexception.ValkyrjaInvalidArgumentError
}

// IsApplicationThrowable marks the error as one that the application raised.
func (e *ApplicationInvalidArgumentError) IsApplicationThrowable() bool {
	return true
}
