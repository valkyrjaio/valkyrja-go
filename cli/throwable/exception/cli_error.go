/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds every error that the CLI component raises.
//
// The component always ships a runtime base and an invalid-argument base. Each
// sub-component of the CLI embeds one of them, and each concrete error embeds
// its own sub-component's base.
package exception

import (
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

// CliRuntimeError is the CLI component's base runtime error.
type CliRuntimeError struct {
	throwableexception.ValkyrjaRuntimeError
}

// NewCliRuntimeError builds the CLI component's base runtime error.
func NewCliRuntimeError(message string, cause error) CliRuntimeError {
	return CliRuntimeError{
		ValkyrjaRuntimeError: throwableexception.NewValkyrjaRuntimeError(message, cause),
	}
}

// IsCliThrowable marks the error as one that the CLI component raised.
func (e *CliRuntimeError) IsCliThrowable() bool {
	return true
}

// CliInvalidArgumentError is the CLI component's base invalid-argument error.
type CliInvalidArgumentError struct {
	throwableexception.ValkyrjaInvalidArgumentError
}

// NewCliInvalidArgumentError builds the CLI component's base invalid-argument
// error.
func NewCliInvalidArgumentError(message string, cause error) CliInvalidArgumentError {
	return CliInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(message, cause),
	}
}

// IsCliThrowable marks the error as one that the CLI component raised.
func (e *CliInvalidArgumentError) IsCliThrowable() bool {
	return true
}
