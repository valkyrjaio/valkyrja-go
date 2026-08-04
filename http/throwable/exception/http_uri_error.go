/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

import "strconv"

// HttpUriInvalidArgumentError is the URI sub-component's base invalid-argument
// error.
type HttpUriInvalidArgumentError struct {
	HttpInvalidArgumentError
}

// newUriInvalidArgumentError builds a concrete invalid-argument error of the URI.
func newUriInvalidArgumentError(message string) HttpUriInvalidArgumentError {
	return HttpUriInvalidArgumentError{
		HttpInvalidArgumentError: NewHttpInvalidArgumentError(message, nil),
	}
}

// HttpUriInvalidPortError reports a port outside the range that a URI carries.
type HttpUriInvalidPortError struct {
	HttpUriInvalidArgumentError

	port int
}

// NewHttpUriInvalidPortError builds the error for a port.
func NewHttpUriInvalidPortError(port int) *HttpUriInvalidPortError {
	return &HttpUriInvalidPortError{
		HttpUriInvalidArgumentError: newUriInvalidArgumentError(
			"Invalid port `" + strconv.Itoa(port) + "` specified; must be a valid TCP/UDP port",
		),
		port: port,
	}
}

// GetPort returns the port that the error reports.
func (e *HttpUriInvalidPortError) GetPort() int {
	return e.port
}

// HttpUriInvalidPathError reports a path that carries a query string.
type HttpUriInvalidPathError struct {
	HttpUriInvalidArgumentError

	path string
}

// NewHttpUriInvalidPathError builds the error for a path.
func NewHttpUriInvalidPathError(path string) *HttpUriInvalidPathError {
	return &HttpUriInvalidPathError{
		HttpUriInvalidArgumentError: newUriInvalidArgumentError(
			"Invalid path of `" + path + "` provided; must not contain a query string",
		),
		path: path,
	}
}

// GetPath returns the path that the error reports.
func (e *HttpUriInvalidPathError) GetPath() string {
	return e.path
}

// HttpUriInvalidQueryError reports a query string that carries a fragment.
type HttpUriInvalidQueryError struct {
	HttpUriInvalidArgumentError

	query string
}

// NewHttpUriInvalidQueryError builds the error for a query string.
func NewHttpUriInvalidQueryError(query string) *HttpUriInvalidQueryError {
	return &HttpUriInvalidQueryError{
		HttpUriInvalidArgumentError: newUriInvalidArgumentError(
			"Invalid query string of `" + query + "` provided; must not contain a URI fragment",
		),
		query: query,
	}
}

// GetQuery returns the query string that the error reports.
func (e *HttpUriInvalidQueryError) GetQuery() string {
	return e.query
}

// HttpRequestInvalidArgumentError is the request sub-component's base
// invalid-argument error.
type HttpRequestInvalidArgumentError struct {
	HttpInvalidArgumentError
}

// HttpRequestInvalidRequestTargetError reports a request target that carries
// whitespace.
type HttpRequestInvalidRequestTargetError struct {
	HttpRequestInvalidArgumentError

	requestTarget string
}

// NewHttpRequestInvalidRequestTargetError builds the error for a request target.
func NewHttpRequestInvalidRequestTargetError(requestTarget string) *HttpRequestInvalidRequestTargetError {
	return &HttpRequestInvalidRequestTargetError{
		HttpRequestInvalidArgumentError: HttpRequestInvalidArgumentError{
			HttpInvalidArgumentError: NewHttpInvalidArgumentError(
				"Invalid request target provided; cannot contain whitespace",
				nil,
			),
		},
		requestTarget: requestTarget,
	}
}

// GetRequestTarget returns the request target that the error reports.
func (e *HttpRequestInvalidRequestTargetError) GetRequestTarget() string {
	return e.requestTarget
}
