/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

import "strconv"

type HttpUriInvalidArgumentError struct {
	HttpInvalidArgumentError
}

// newUriInvalidArgumentError builds a concrete invalid-argument error of the URI.
func newUriInvalidArgumentError(message string) HttpUriInvalidArgumentError {
	return HttpUriInvalidArgumentError{
		HttpInvalidArgumentError: NewHttpInvalidArgumentError(message, nil),
	}
}

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

type HttpRequestInvalidArgumentError struct {
	HttpInvalidArgumentError
}

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
