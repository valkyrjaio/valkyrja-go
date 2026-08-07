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
	"regexp"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

// traceCodePattern is the hexadecimal SHA-256 digest that GetTraceCode returns.
var traceCodePattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

func TestValkyrjaRuntimeErrorReturnsTheMessage(t *testing.T) {
	t.Parallel()

	err := exception.NewValkyrjaRuntimeError("the runtime failed", nil)

	if err.Error() != "the runtime failed" {
		t.Errorf("Error must be 'the runtime failed', but is: %s", err.Error())
	}
}

func TestValkyrjaInvalidArgumentErrorReturnsTheMessage(t *testing.T) {
	t.Parallel()

	err := exception.NewValkyrjaInvalidArgumentError("the argument is invalid", nil)

	if err.Error() != "the argument is invalid" {
		t.Errorf("Error must be 'the argument is invalid', but is: %s", err.Error())
	}
}

func TestUnwrapReturnsNilWhereTheErrorHasNoCause(t *testing.T) {
	t.Parallel()

	err := exception.NewValkyrjaRuntimeError("the runtime failed", nil)

	if err.Unwrap() != nil {
		t.Errorf("Unwrap must be nil, but is: %v", err.Unwrap())
	}
}

func TestUnwrapReturnsTheCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("the cause")
	err := exception.NewValkyrjaInvalidArgumentError("the argument is invalid", cause)

	if !errors.Is(&err, cause) {
		t.Errorf("errors.Is must find the cause, but did not: %v", err.Unwrap())
	}
}

func TestGetTraceCodeReturnsAHexadecimalDigest(t *testing.T) {
	t.Parallel()

	err := exception.NewValkyrjaRuntimeError("the runtime failed", nil)

	if !traceCodePattern.MatchString(err.GetTraceCode()) {
		t.Errorf("GetTraceCode must match %s, but is: %s", traceCodePattern, err.GetTraceCode())
	}
}

func TestGetTraceCodeIsStableForOneError(t *testing.T) {
	t.Parallel()

	err := exception.NewValkyrjaRuntimeError("the runtime failed", nil)

	first := err.GetTraceCode()
	second := err.GetTraceCode()

	if first != second {
		t.Error("GetTraceCode must return the same code for one error, but did not")
	}
}

func TestGetTraceCodeDiffersBetweenTwoSites(t *testing.T) {
	t.Parallel()

	first := exception.NewValkyrjaRuntimeError("the runtime failed", nil)
	second := newErrorFromAnotherSite()

	if first.GetTraceCode() == second.GetTraceCode() {
		t.Error("GetTraceCode must differ between two sites, but did not")
	}
}

func TestTheBaseErrorsSatisfyTheRootContract(t *testing.T) {
	t.Parallel()

	runtimeError := exception.NewValkyrjaRuntimeError("the runtime failed", nil)
	invalidArgumentError := exception.NewValkyrjaInvalidArgumentError("the argument is invalid", nil)

	throwables := []contract.ValkyrjaThrowable{&runtimeError, &invalidArgumentError}

	for _, throwable := range throwables {
		if throwable.GetTraceCode() == "" {
			t.Errorf("GetTraceCode must return a code, but is empty for: %s", throwable.Error())
		}
	}
}

// newErrorFromAnotherSite raises the error from a second call site, so the
// recorded frames differ from the frames that the caller records.
func newErrorFromAnotherSite() exception.ValkyrjaRuntimeError {
	return exception.NewValkyrjaRuntimeError("the runtime failed", nil)
}
