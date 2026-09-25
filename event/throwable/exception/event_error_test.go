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

	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/throwable/exception"
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

func TestInvalidEventErrorReportsTheBindingKey(t *testing.T) {
	t.Parallel()

	err := exception.NewEventInvalidEventError("NotAnEventID")

	if err.GetID() != "NotAnEventID" {
		t.Errorf("GetID must be the binding key, but is: %s", err.GetID())
	}
}

func TestInvalidEventErrorStatesThatTheServiceIsNoEvent(t *testing.T) {
	t.Parallel()

	err := exception.NewEventInvalidEventError("NotAnEventID")

	if err.Error() != "Service with `NotAnEventID` is not an event" {
		t.Errorf("Error must name the binding key, but is: %s", err.Error())
	}
}

func TestTheConcreteErrorSatisfiesTheComponentContract(t *testing.T) {
	t.Parallel()

	var throwable contract.EventThrowable = exception.NewEventInvalidEventError("NotAnEventID")

	if !throwable.IsEventThrowable() {
		t.Error("IsEventThrowable must be true, but is false")
	}

	if throwable.GetTraceCode() == "" {
		t.Error("GetTraceCode must return a code, but is empty")
	}
}

func TestEachBaseErrorSatisfiesTheComponentContract(t *testing.T) {
	t.Parallel()

	runtimeError := &exception.EventRuntimeError{
		ValkyrjaRuntimeError: throwableexception.NewValkyrjaRuntimeError("the runtime failed", nil),
	}
	invalidArgumentError := &exception.EventInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(
			"the argument is invalid",
			nil,
		),
	}

	throwables := []contract.EventThrowable{runtimeError, invalidArgumentError}

	for _, throwable := range throwables {
		if !throwable.IsEventThrowable() {
			t.Errorf("IsEventThrowable must be true, but is false for: %s", throwable.Error())
		}
	}
}

func TestErrorsAsFindsTheConcreteError(t *testing.T) {
	t.Parallel()

	var err error = exception.NewEventInvalidEventError("NotAnEventID")

	target, found := errors.AsType[*exception.EventInvalidEventError](err)
	if !found {
		t.Fatal("errors.AsType must find the concrete error, but did not")
	}

	if target.GetID() != "NotAnEventID" {
		t.Errorf("GetID must be the binding key, but is: %s", target.GetID())
	}
}
