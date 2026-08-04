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

const parameterName = "count"

func TestEachBaseErrorMarksItselfAsTheComponentsOwn(t *testing.T) {
	t.Parallel()

	cause := errors.New("the cause")

	runtimeError := exception.NewCliRuntimeError("A failure", cause)
	invalidError := exception.NewCliInvalidArgumentError("A failure", cause)

	if runtimeError.Error() != "A failure" || invalidError.Error() != "A failure" {
		t.Error("each base error must report its message, but did not")
	}

	if !errors.Is(&runtimeError, cause) || !errors.Is(&invalidError, cause) {
		t.Error("each base error must unwrap to its cause, but did not")
	}

	if !runtimeError.IsCliThrowable() || !invalidError.IsCliThrowable() {
		t.Error("each base error must mark itself as the CLI component's own, but did not")
	}
}

func TestEachRoutingErrorReportsWhatWentWrong(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err     error
		message string
	}{
		"an argument that holds a value the command rejects": {
			err:     exception.NewCliRoutingArgumentValuesValidationError(parameterName),
			message: parameterName + " is invalid",
		},
		"an option that holds a value the command rejects": {
			err:     exception.NewCliRoutingOptionValuesValidationError(parameterName),
			message: parameterName + " is invalid",
		},
		"an option that takes no value": {
			err:     exception.NewCliRoutingInvalidOptionWithValueError(parameterName),
			message: parameterName + " should have no value",
		},
		"a command that the application does not hold": {
			err:     exception.NewCliRoutingRouteNotFoundError("cache:clear"),
			message: "Command cache:clear does not exist",
		},
	}

	for name, test := range tests {
		if test.err.Error() != test.message {
			t.Errorf("%s must report %q, but reported %q", name, test.message, test.err.Error())
		}

		var throwable throwablecontract.ValkyrjaThrowable
		if !errors.As(test.err, &throwable) || throwable.GetTraceCode() == "" {
			t.Errorf("%s must carry a trace code, but carried none", name)
		}
	}
}

func TestARoutingErrorMarksItselfAsTheSubComponentsOwn(t *testing.T) {
	t.Parallel()

	err := exception.NewCliRoutingArgumentValuesValidationError(parameterName)

	if !err.IsCliRoutingThrowable() || !err.IsCliThrowable() {
		t.Error("a routing error must mark itself as the routing sub-component's own, but did not")
	}
}
