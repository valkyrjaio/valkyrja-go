/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds every error that the HTTP component raises.
//
// The component always ships a runtime base and an invalid-argument base. Each
// sub-component of HTTP embeds one of them, and each concrete error embeds its
// own sub-component's base.
package exception

import (
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

type HttpRuntimeError struct {
	throwableexception.ValkyrjaRuntimeError
}

// NewHttpRuntimeError builds the HTTP component's base runtime error.
func NewHttpRuntimeError(message string, cause error) HttpRuntimeError {
	return HttpRuntimeError{
		ValkyrjaRuntimeError: throwableexception.NewValkyrjaRuntimeError(message, cause),
	}
}

// IsHttpThrowable marks the error as one that the HTTP component raised.
func (e *HttpRuntimeError) IsHttpThrowable() bool {
	return true
}

type HttpInvalidArgumentError struct {
	throwableexception.ValkyrjaInvalidArgumentError
}

// NewHttpInvalidArgumentError builds the HTTP component's base invalid-argument
// error.
func NewHttpInvalidArgumentError(message string, cause error) HttpInvalidArgumentError {
	return HttpInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(message, cause),
	}
}

// IsHttpThrowable marks the error as one that the HTTP component raised.
func (e *HttpInvalidArgumentError) IsHttpThrowable() bool {
	return true
}
