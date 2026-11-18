/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/throwable/exception"
	throwablecontract "github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
)

func TestAFileWriteFailureNamesTheFile(t *testing.T) {
	t.Parallel()

	cause := errors.New("the disk is full")

	err := exception.NewCliInteractionFileWriteError("/var/log/app.log", cause)

	if err.Error() != "Unable to write the whole message to the file `/var/log/app.log`" {
		t.Errorf("the error must name the file, but reported: %q", err.Error())
	}

	if err.GetFilepath() != "/var/log/app.log" {
		t.Errorf("the error must carry the file, but carried: %q", err.GetFilepath())
	}

	if !errors.Is(err, cause) {
		t.Error("the error must unwrap to its cause, but did not")
	}

	if !err.IsCliInteractionThrowable() || !err.IsCliThrowable() {
		t.Error("the error must mark itself as the CLI interaction sub-component's own, but did not")
	}
}

func TestAStreamWriteFailureNamesHowMuchTheStreamTook(t *testing.T) {
	t.Parallel()

	cause := errors.New("the pipe is closed")

	err := exception.NewCliInteractionStreamWriteError(3, 11, cause)

	if err.Error() != "Unable to write the whole message to the stream: wrote 3 of 11 bytes" {
		t.Errorf("the error must name how much the stream took, but reported: %q", err.Error())
	}

	if err.GetWritten() != 3 || err.GetLength() != 11 {
		t.Errorf("the error must carry the counts, but carried: %d of %d", err.GetWritten(), err.GetLength())
	}

	if !errors.Is(err, cause) {
		t.Error("the error must unwrap to its cause, but did not")
	}
}

func TestAnUnwritableStreamReportsThatItHoldsNone(t *testing.T) {
	t.Parallel()

	err := exception.NewCliInteractionUnwritableStreamError()

	if err.Error() != "The writer holds no stream to write to" {
		t.Errorf("the error must report that it holds no stream, but reported: %q", err.Error())
	}

	var throwable throwablecontract.ValkyrjaThrowable
	if !errors.As(err, &throwable) || throwable.GetTraceCode() == "" {
		t.Error("the error must carry a trace code, but carried none")
	}
}
