/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package fixtures holds the reusable doubles that the CLI interaction tests
// build on.
package fixtures

import "errors"

// ErrWriteFailed is what a failing writer reports.
var ErrWriteFailed = errors.New("the destination took no bytes")

// FailingWriterFixture is a stream that reports a failure rather than taking a
// message.
type FailingWriterFixture struct{}

// Write reports the failure.
func (w *FailingWriterFixture) Write(_ []byte) (int, error) {
	return 0, ErrWriteFailed
}

// ShortWriterFixture is a stream that takes part of a message and reports no
// failure, which is the short write that `io.Writer` permits a caller to see.
type ShortWriterFixture struct{}

// Write takes one byte of the message and reports no failure.
func (w *ShortWriterFixture) Write(text []byte) (int, error) {
	if len(text) == 0 {
		return 0, nil
	}

	return 1, nil
}
