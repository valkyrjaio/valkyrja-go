/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

// HttpClientRuntimeError is the HTTP client sub-component's base runtime error.
type HttpClientRuntimeError struct {
	HttpRuntimeError
}

// newHttpClientRuntimeError builds the HTTP client sub-component's base runtime
// error.
func newHttpClientRuntimeError(message string, cause error) HttpClientRuntimeError {
	return HttpClientRuntimeError{
		HttpRuntimeError: NewHttpRuntimeError(message, cause),
	}
}

// IsHttpClientThrowable marks the error as one that the HTTP client
// sub-component raised.
func (e *HttpClientRuntimeError) IsHttpClientThrowable() bool {
	return true
}

// HttpClientRequestFailedError reports that the client did not reach the server.
type HttpClientRequestFailedError struct {
	HttpClientRuntimeError

	uri string
}

// NewHttpClientRequestFailedError builds the error for the URI that the client
// tried to reach.
func NewHttpClientRequestFailedError(uri string, cause error) *HttpClientRequestFailedError {
	return &HttpClientRequestFailedError{
		HttpClientRuntimeError: newHttpClientRuntimeError("The request to "+uri+" failed", cause),
		uri:                    uri,
	}
}

// GetUri returns the URI that the client tried to reach.
func (e *HttpClientRequestFailedError) GetUri() string {
	return e.uri
}
