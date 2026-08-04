/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

import "strconv"

// HttpStreamRuntimeError is the stream sub-component's base runtime error.
type HttpStreamRuntimeError struct {
	HttpRuntimeError
}

// HttpStreamInvalidArgumentError is the stream sub-component's base
// invalid-argument error.
type HttpStreamInvalidArgumentError struct {
	HttpInvalidArgumentError
}

// newStreamRuntimeError builds a concrete runtime error of the stream.
func newStreamRuntimeError(message string) HttpStreamRuntimeError {
	return HttpStreamRuntimeError{HttpRuntimeError: NewHttpRuntimeError(message, nil)}
}

// HttpStreamInvalidLengthError reports a read of fewer than zero bytes.
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

// HttpStreamUnreadableStreamError reports a read of a stream that no reader
// reads.
type HttpStreamUnreadableStreamError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamUnreadableStreamError builds the error.
func NewHttpStreamUnreadableStreamError() *HttpStreamUnreadableStreamError {
	return &HttpStreamUnreadableStreamError{
		HttpStreamRuntimeError: newStreamRuntimeError("Stream is not readable"),
	}
}

// HttpStreamUnwritableStreamError reports a write to a stream that no writer
// writes.
type HttpStreamUnwritableStreamError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamUnwritableStreamError builds the error.
func NewHttpStreamUnwritableStreamError() *HttpStreamUnwritableStreamError {
	return &HttpStreamUnwritableStreamError{
		HttpStreamRuntimeError: newStreamRuntimeError("Stream is not writable"),
	}
}

// HttpStreamUnseekableStreamError reports a seek in a stream that no caller
// seeks in.
type HttpStreamUnseekableStreamError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamUnseekableStreamError builds the error.
func NewHttpStreamUnseekableStreamError() *HttpStreamUnseekableStreamError {
	return &HttpStreamUnseekableStreamError{
		HttpStreamRuntimeError: newStreamRuntimeError("Stream is not seekable"),
	}
}

// HttpStreamStreamSeekError reports a seek to a position outside the stream.
type HttpStreamStreamSeekError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamStreamSeekError builds the error.
func NewHttpStreamStreamSeekError() *HttpStreamStreamSeekError {
	return &HttpStreamStreamSeekError{
		HttpStreamRuntimeError: newStreamRuntimeError("Position is outside of the stream"),
	}
}

// HttpStreamStreamTellError reports a read of the pointer of a stream that is
// closed.
type HttpStreamStreamTellError struct {
	HttpStreamRuntimeError
}

// NewHttpStreamStreamTellError builds the error.
func NewHttpStreamStreamTellError() *HttpStreamStreamTellError {
	return &HttpStreamStreamTellError{
		HttpStreamRuntimeError: newStreamRuntimeError("Could not read the position of the stream"),
	}
}
