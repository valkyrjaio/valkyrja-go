/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

type HttpRoutingInvalidArgumentError struct {
	HttpInvalidArgumentError
}

// newHttpRoutingInvalidArgumentError builds the HTTP routing sub-component's
// base invalid-argument error.
func newHttpRoutingInvalidArgumentError(message string) HttpRoutingInvalidArgumentError {
	return HttpRoutingInvalidArgumentError{
		HttpInvalidArgumentError: NewHttpInvalidArgumentError(message, nil),
	}
}

// IsHttpRoutingThrowable marks the error as one that the HTTP routing
// sub-component raised.
func (e *HttpRoutingInvalidArgumentError) IsHttpRoutingThrowable() bool {
	return true
}

type HttpRoutingInvalidRouteParameterError struct {
	HttpRoutingInvalidArgumentError

	name string
}

// NewHttpRoutingInvalidRouteParameterError builds the error for the parameter
// that the route does not hold.
func NewHttpRoutingInvalidRouteParameterError(name string) *HttpRoutingInvalidRouteParameterError {
	return &HttpRoutingInvalidRouteParameterError{
		HttpRoutingInvalidArgumentError: newHttpRoutingInvalidArgumentError(
			"No parameter named '" + name + "' exists on this route",
		),
		name: name,
	}
}

// GetName returns the parameter name that the error reports.
func (e *HttpRoutingInvalidRouteParameterError) GetName() string {
	return e.name
}
