/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds every error that the validation component raises.
package exception

import (
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

// ValidationRuntimeError is the validation component's base runtime error.
type ValidationRuntimeError struct {
	throwableexception.ValkyrjaRuntimeError
}

// NewValidationRuntimeError builds the validation component's base runtime
// error.
func NewValidationRuntimeError(message string, cause error) ValidationRuntimeError {
	return ValidationRuntimeError{
		ValkyrjaRuntimeError: throwableexception.NewValkyrjaRuntimeError(message, cause),
	}
}

// IsValidationThrowable marks the error as one that the validation component
// raised.
func (e *ValidationRuntimeError) IsValidationThrowable() bool {
	return true
}

// ValidationInvalidArgumentError is the validation component's base
// invalid-argument error.
type ValidationInvalidArgumentError struct {
	throwableexception.ValkyrjaInvalidArgumentError
}

// NewValidationInvalidArgumentError builds the validation component's base
// invalid-argument error.
func NewValidationInvalidArgumentError(message string, cause error) ValidationInvalidArgumentError {
	return ValidationInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(message, cause),
	}
}

// IsValidationThrowable marks the error as one that the validation component
// raised.
func (e *ValidationInvalidArgumentError) IsValidationThrowable() bool {
	return true
}

// ValidationRuleFailureError reports that a subject failed a rule.
//
// The message is the one that the rule carries, so a caller prints it to the
// person who gave the subject.
type ValidationRuleFailureError struct {
	ValidationInvalidArgumentError
}

// NewValidationRuleFailureError builds the error that carries the message of the
// rule.
func NewValidationRuleFailureError(message string) *ValidationRuleFailureError {
	return &ValidationRuleFailureError{
		ValidationInvalidArgumentError: NewValidationInvalidArgumentError(message, nil),
	}
}
