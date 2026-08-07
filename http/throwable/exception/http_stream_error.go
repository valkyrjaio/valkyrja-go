/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

import "strconv"

type HttpStreamRuntimeError struct {
	HttpRuntimeError
}

type HttpStreamInvalidArgumentError struct {
	HttpInvalidArgumentError
}

// newStreamRuntimeError builds a concrete runtime error of the stream.
func newStreamRuntimeError(message string) HttpStreamRuntimeError {
	return HttpStreamRuntimeError{HttpRuntimeError: NewHttpRuntimeError(message, nil)}
}

type HttpStreamInvalidLengthError struct {
	HttpStreamInvalidArgumentError

	length int
}

// NewHttpStreamInvalidLengthError builds the error for a read length.
func NewHttpStreamInvalidLengthError(length int) *HttpStreamInvalidLengthError {
	return &HttpStreamInvalidLengthError{
		HttpStreamInvalidArgumentError: HttpStreamInvalidArgumentError{
			HttpInvalidArgumentError: NewHttpInvalidArgumentError(
				"Invalid length of "+strconv.Itoa(length)+" provided. Length must be greater than 0",
				nil,
			),
		},
		length: length,
	}
}

// GetLength returns the read length that the error reports.
func (e *HttpStreamInvalidLengthError) GetLength() int {
	return e.length
}

type HttpStreamUnreadableStreamError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamUnreadableStreamError builds the error.
func NewHttpStreamUnreadableStreamError() *HttpStreamUnreadableStreamError {
	return &HttpStreamUnreadableStreamError{
		HttpStreamRuntimeError: newStreamRuntimeError("Stream is not readable"),
	}
}

type HttpStreamUnwritableStreamError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamUnwritableStreamError builds the error.
func NewHttpStreamUnwritableStreamError() *HttpStreamUnwritableStreamError {
	return &HttpStreamUnwritableStreamError{
		HttpStreamRuntimeError: newStreamRuntimeError("Stream is not writable"),
	}
}

type HttpStreamUnseekableStreamError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamUnseekableStreamError builds the error.
func NewHttpStreamUnseekableStreamError() *HttpStreamUnseekableStreamError {
	return &HttpStreamUnseekableStreamError{
		HttpStreamRuntimeError: newStreamRuntimeError("Stream is not seekable"),
	}
}

type HttpStreamStreamSeekError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamStreamSeekError builds the error.
func NewHttpStreamStreamSeekError() *HttpStreamStreamSeekError {
	return &HttpStreamStreamSeekError{
		HttpStreamRuntimeError: newStreamRuntimeError("Position is outside of the stream"),
	}
}

type HttpStreamStreamTellError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamStreamTellError builds the error.
func NewHttpStreamStreamTellError() *HttpStreamStreamTellError {
	return &HttpStreamStreamTellError{
		HttpStreamRuntimeError: newStreamRuntimeError("Could not read the position of the stream"),
	}
}
