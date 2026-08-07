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

	throwablecontract "github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/throwable/exception"
)

func TestEachBaseErrorMarksItselfAsTheComponentsOwn(t *testing.T) {
	t.Parallel()

	cause := errors.New("the cause")

	runtimeError := exception.NewValidationRuntimeError("A failure", cause)
	invalidError := exception.NewValidationInvalidArgumentError("A failure", cause)

	if runtimeError.Error() != "A failure" || invalidError.Error() != "A failure" {
		t.Error("each base error must report its message, but did not")
	}

	if !errors.Is(&runtimeError, cause) || !errors.Is(&invalidError, cause) {
		t.Error("each base error must unwrap to its cause, but did not")
	}

	if !runtimeError.IsValidationThrowable() || !invalidError.IsValidationThrowable() {
		t.Error("each base error must mark itself as the validation component's own, but did not")
	}
}

func TestARuleFailureCarriesTheMessageOfTheRule(t *testing.T) {
	t.Parallel()

	err := exception.NewValidationRuleFailureError("This field is required.")

	if err.Error() != "This field is required." {
		t.Errorf("the failure must carry the message of the rule, but carried: %q", err.Error())
	}

	var throwable throwablecontract.ValkyrjaThrowable
	if !errors.As(err, &throwable) || throwable.GetTraceCode() == "" {
		t.Error("the failure must carry a trace code, but carried none")
	}
}
