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

	"github.com/valkyrjaio/valkyrja-go/v26/log/throwable/exception"
	throwablecontract "github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
)

func TestEachBaseErrorMarksItselfAsTheComponentsOwn(t *testing.T) {
	t.Parallel()

	cause := errors.New("the cause")

	runtimeError := exception.NewLogRuntimeError("A failure", cause)
	invalidError := exception.NewLogInvalidArgumentError("A failure", cause)

	if runtimeError.Error() != "A failure" || invalidError.Error() != "A failure" {
		t.Error("each base error must report its message, but did not")
	}

	if !errors.Is(&runtimeError, cause) || !errors.Is(&invalidError, cause) {
		t.Error("each base error must unwrap to its cause, but did not")
	}

	if !runtimeError.IsLogThrowable() || !invalidError.IsLogThrowable() {
		t.Error("each base error must mark itself as the log component's own, but did not")
	}
}

func TestAnInvalidSeverityNamesWhatTheCallerGave(t *testing.T) {
	t.Parallel()

	err := exception.NewLogInvalidLogLevelError("verbose")

	if err.Error() != "Invalid log level: verbose" {
		t.Errorf("the error must name the severity, but reported: %q", err.Error())
	}

	var throwable throwablecontract.ValkyrjaThrowable
	if !errors.As(err, &throwable) || throwable.GetTraceCode() == "" {
		t.Error("the error must carry a trace code, but carried none")
	}
}
