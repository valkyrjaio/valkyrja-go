/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

// CliRoutingInvalidArgumentError is the CLI routing sub-component's base
// invalid-argument error.
type CliRoutingInvalidArgumentError struct {
	CliInvalidArgumentError
}

// newCliRoutingInvalidArgumentError builds the CLI routing sub-component's base
// invalid-argument error.
func newCliRoutingInvalidArgumentError(message string) CliRoutingInvalidArgumentError {
	return CliRoutingInvalidArgumentError{
		CliInvalidArgumentError: NewCliInvalidArgumentError(message, nil),
	}
}

// IsCliRoutingThrowable marks the error as one that the CLI routing
// sub-component raised.
func (e *CliRoutingInvalidArgumentError) IsCliRoutingThrowable() bool {
	return true
}

// CliRoutingArgumentValuesValidationError reports that an argument holds a value
// that the command does not accept.
type CliRoutingArgumentValuesValidationError struct {
	CliRoutingInvalidArgumentError
}

// NewCliRoutingArgumentValuesValidationError builds the error for the argument
// under the name.
func NewCliRoutingArgumentValuesValidationError(name string) *CliRoutingArgumentValuesValidationError {
	return &CliRoutingArgumentValuesValidationError{
		CliRoutingInvalidArgumentError: newCliRoutingInvalidArgumentError(name + " is invalid"),
	}
}

// CliRoutingOptionValuesValidationError reports that an option holds a value
// that the command does not accept.
type CliRoutingOptionValuesValidationError struct {
	CliRoutingInvalidArgumentError
}

// NewCliRoutingOptionValuesValidationError builds the error for the option under
// the name.
func NewCliRoutingOptionValuesValidationError(name string) *CliRoutingOptionValuesValidationError {
	return &CliRoutingOptionValuesValidationError{
		CliRoutingInvalidArgumentError: newCliRoutingInvalidArgumentError(name + " is invalid"),
	}
}

// CliRoutingInvalidOptionWithValueError reports that the caller gave a value to
// an option that takes none.
type CliRoutingInvalidOptionWithValueError struct {
	CliRoutingInvalidArgumentError
}

// NewCliRoutingInvalidOptionWithValueError builds the error for the option under
// the name.
func NewCliRoutingInvalidOptionWithValueError(name string) *CliRoutingInvalidOptionWithValueError {
	return &CliRoutingInvalidOptionWithValueError{
		CliRoutingInvalidArgumentError: newCliRoutingInvalidArgumentError(name + " should have no value"),
	}
}

// CliRoutingRouteNotFoundError reports that the collection holds no command
// under the name.
type CliRoutingRouteNotFoundError struct {
	CliRoutingInvalidArgumentError
}

// NewCliRoutingRouteNotFoundError builds the error for the command under the
// name.
func NewCliRoutingRouteNotFoundError(name string) *CliRoutingRouteNotFoundError {
	return &CliRoutingRouteNotFoundError{
		CliRoutingInvalidArgumentError: newCliRoutingInvalidArgumentError("Command " + name + " does not exist"),
	}
}
