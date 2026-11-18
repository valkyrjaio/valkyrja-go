/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds every error that the log component raises.
package exception

import (
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

type LogRuntimeError struct {
	throwableexception.ValkyrjaRuntimeError
}

// NewLogRuntimeError builds the log component's base runtime error.
func NewLogRuntimeError(message string, cause error) LogRuntimeError {
	return LogRuntimeError{
		ValkyrjaRuntimeError: throwableexception.NewValkyrjaRuntimeError(message, cause),
	}
}

// IsLogThrowable marks the error as one that the log component raised.
func (e *LogRuntimeError) IsLogThrowable() bool {
	return true
}

type LogInvalidArgumentError struct {
	throwableexception.ValkyrjaInvalidArgumentError
}

// NewLogInvalidArgumentError builds the log component's base invalid-argument
// error.
func NewLogInvalidArgumentError(message string, cause error) LogInvalidArgumentError {
	return LogInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(message, cause),
	}
}

// IsLogThrowable marks the error as one that the log component raised.
func (e *LogInvalidArgumentError) IsLogThrowable() bool {
	return true
}

type LogInvalidLogLevelError struct {
	LogInvalidArgumentError
}

// NewLogInvalidLogLevelError builds the error for the severity.
func NewLogInvalidLogLevelError(level string) *LogInvalidLogLevelError {
	return &LogInvalidLogLevelError{
		LogInvalidArgumentError: NewLogInvalidArgumentError("Invalid log level: "+level, nil),
	}
}
