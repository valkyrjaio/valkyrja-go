/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package exception holds the two base errors that every framework error
// derives from.
//
// The other ports declare these two types abstract. Go has no abstract type, so
// each base is a struct that a component error embeds. A component never
// returns a base error on its own; it returns one of its own concrete errors.
package exception

import (
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"strconv"
	"strings"
)

// traceDepth is the number of call frames that an error records. The frames
// identify the site that raised the error, so a shallow trace is sufficient.
const traceDepth = 32

// traceSkip drops the frames that belong to the constructor itself, so the
// first recorded frame is the site that raised the error.
const traceSkip = 3

// valkyrjaError carries the state that every framework error holds. It is
// unexported, so only a base error in this package embeds it.
type valkyrjaError struct {
	message string
	cause   error
	trace   []uintptr
}

// newValkyrjaError records the message, the cause, and the call site.
func newValkyrjaError(message string, cause error) valkyrjaError {
	programCounters := make([]uintptr, traceDepth)
	recorded := runtime.Callers(traceSkip, programCounters)

	return valkyrjaError{
		message: message,
		cause:   cause,
		trace:   programCounters[:recorded],
	}
}

// Error returns the message.
func (e *valkyrjaError) Error() string {
	return e.message
}

// Unwrap returns the cause, and nil where the error has no cause.
func (e *valkyrjaError) Unwrap() error {
	return e.cause
}

// GetTraceCode returns a stable identifier for the site that raised the error.
//
// The other ports hash the error name and the stack with MD5. This port hashes
// the recorded frames with SHA-256, because `gosec` rejects MD5. The code is an
// opaque identifier, so the digest only has to be stable, and no port can match
// another port's value in any case.
func (e *valkyrjaError) GetTraceCode() string {
	digest := sha256.Sum256([]byte(e.getTraceSignature()))

	return hex.EncodeToString(digest[:])
}

// getTraceSignature renders the recorded frames as `function:line` lines.
func (e *valkyrjaError) getTraceSignature() string {
	frames := runtime.CallersFrames(e.trace)

	signature := &strings.Builder{}

	for {
		frame, more := frames.Next()

		signature.WriteString(frame.Function)
		signature.WriteByte(':')
		signature.WriteString(strconv.Itoa(frame.Line))
		signature.WriteByte('\n')

		if !more {
			break
		}
	}

	return signature.String()
}

// ValkyrjaRuntimeError is the base for every runtime error in the framework. A
// runtime error reports a failure that the caller cannot prevent by passing a
// different argument.
type ValkyrjaRuntimeError struct {
	valkyrjaError
}

// NewValkyrjaRuntimeError builds the base runtime error that a component error
// embeds.
func NewValkyrjaRuntimeError(message string, cause error) ValkyrjaRuntimeError {
	return ValkyrjaRuntimeError{valkyrjaError: newValkyrjaError(message, cause)}
}

// ValkyrjaInvalidArgumentError is the base for every invalid-argument error in
// the framework. An invalid-argument error reports a value that the caller
// passed, and the caller prevents it by passing a different value.
type ValkyrjaInvalidArgumentError struct {
	valkyrjaError
}

// NewValkyrjaInvalidArgumentError builds the base invalid-argument error that a
// component error embeds.
func NewValkyrjaInvalidArgumentError(message string, cause error) ValkyrjaInvalidArgumentError {
	return ValkyrjaInvalidArgumentError{valkyrjaError: newValkyrjaError(message, cause)}
}
