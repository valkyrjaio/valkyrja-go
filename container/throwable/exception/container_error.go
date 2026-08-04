/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds every error that the container component raises.
//
// The component always ships a runtime base and an invalid-argument base, even
// where it raises neither on its own. Each concrete error embeds one of them.
package exception

import (
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

// ContainerRuntimeError is the container's base runtime error.
type ContainerRuntimeError struct {
	throwableexception.ValkyrjaRuntimeError
}

// IsContainerThrowable marks the error as one that the container raised.
func (e *ContainerRuntimeError) IsContainerThrowable() bool {
	return true
}

// newContainerRuntimeError builds the container's base runtime error.
func newContainerRuntimeError(message string, cause error) ContainerRuntimeError {
	return ContainerRuntimeError{
		ValkyrjaRuntimeError: throwableexception.NewValkyrjaRuntimeError(message, cause),
	}
}

// ContainerInvalidArgumentError is the container's base invalid-argument error.
type ContainerInvalidArgumentError struct {
	throwableexception.ValkyrjaInvalidArgumentError
}

// IsContainerThrowable marks the error as one that the container raised.
func (e *ContainerInvalidArgumentError) IsContainerThrowable() bool {
	return true
}

// newContainerInvalidArgumentError builds the container's base
// invalid-argument error.
func newContainerInvalidArgumentError(message string, cause error) ContainerInvalidArgumentError {
	return ContainerInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(message, cause),
	}
}
