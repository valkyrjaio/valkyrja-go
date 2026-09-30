/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

import "strconv"

type CliInteractionRuntimeError struct {
	CliRuntimeError
}

// newCliInteractionRuntimeError builds the CLI interaction sub-component's base
// runtime error.
func newCliInteractionRuntimeError(message string, cause error) CliInteractionRuntimeError {
	return CliInteractionRuntimeError{
		CliRuntimeError: NewCliRuntimeError(message, cause),
	}
}

// IsCliInteractionThrowable marks the error as one that the CLI interaction
// sub-component raised.
func (e *CliInteractionRuntimeError) IsCliInteractionThrowable() bool {
	return true
}

type CliInteractionFileWriteError struct {
	CliInteractionRuntimeError

	filepath string
}

// NewCliInteractionFileWriteError builds the error for the file that took no
// message.
func NewCliInteractionFileWriteError(filepath string, cause error) *CliInteractionFileWriteError {
	return &CliInteractionFileWriteError{
		CliInteractionRuntimeError: newCliInteractionRuntimeError(
			"Unable to write the whole message to the file `"+filepath+"`",
			cause,
		),
		filepath: filepath,
	}
}

// GetFilepath returns the file that took no message.
func (e *CliInteractionFileWriteError) GetFilepath() string {
	return e.filepath
}

type CliInteractionStreamWriteError struct {
	CliInteractionRuntimeError

	written int
	length  int
}

// NewCliInteractionStreamWriteError builds the error for the stream that took
// part of a message, or none of it.
func NewCliInteractionStreamWriteError(written int, length int, cause error) *CliInteractionStreamWriteError {
	return &CliInteractionStreamWriteError{
		CliInteractionRuntimeError: newCliInteractionRuntimeError(
			"Unable to write the whole message to the stream: wrote "+
				strconv.Itoa(written)+" of "+strconv.Itoa(length)+" bytes",
			cause,
		),
		written: written,
		length:  length,
	}
}

// GetWritten returns how many bytes the stream took.
func (e *CliInteractionStreamWriteError) GetWritten() int {
	return e.written
}

// GetLength returns how many bytes the message carries.
func (e *CliInteractionStreamWriteError) GetLength() int {
	return e.length
}

type CliInteractionUnwritableStreamError struct {
	CliInteractionRuntimeError
}

// NewCliInteractionUnwritableStreamError builds the error for a writer that
// holds no stream to write to.
func NewCliInteractionUnwritableStreamError() *CliInteractionUnwritableStreamError {
	return &CliInteractionUnwritableStreamError{
		CliInteractionRuntimeError: newCliInteractionRuntimeError(
			"The writer holds no stream to write to",
			nil,
		),
	}
}
