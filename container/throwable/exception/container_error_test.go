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

	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/throwable/exception"
)

func TestInvalidReferenceErrorReportsTheBindingKey(t *testing.T) {
	t.Parallel()

	err := exception.NewContainerInvalidReferenceError("io.valkyrja.container.ContainerContract")

	if err.GetID() != "io.valkyrja.container.ContainerContract" {
		t.Errorf("GetID must be the binding key, but is: %s", err.GetID())
	}
}

func TestInvalidReferenceErrorStatesThatTheServiceIsNotFound(t *testing.T) {
	t.Parallel()

	err := exception.NewContainerInvalidReferenceError("MissingID")

	if err.Error() != "Service with `MissingID` not found" {
		t.Errorf("Error must name the binding key, but is: %s", err.Error())
	}
}

func TestInvalidPublishCallbackErrorReportsTheBindingKey(t *testing.T) {
	t.Parallel()

	err := exception.NewContainerInvalidPublishCallbackError("InvalidID")

	if err.GetID() != "InvalidID" {
		t.Errorf("GetID must be the binding key, but is: %s", err.GetID())
	}
}

func TestInvalidPublishCallbackErrorStatesThatThePublisherIsInvalid(t *testing.T) {
	t.Parallel()

	err := exception.NewContainerInvalidPublishCallbackError("InvalidID")

	if err.Error() != "InvalidID should have a valid callable" {
		t.Errorf("Error must name the binding key, but is: %s", err.Error())
	}
}

func TestEachContainerErrorSatisfiesTheComponentContract(t *testing.T) {
	t.Parallel()

	throwables := []contract.ContainerThrowable{
		exception.NewContainerInvalidReferenceError("MissingID"),
		exception.NewContainerInvalidPublishCallbackError("InvalidID"),
	}

	for _, throwable := range throwables {
		if !throwable.IsContainerThrowable() {
			t.Errorf("IsContainerThrowable must be true, but is false for: %s", throwable.Error())
		}

		if throwable.GetTraceCode() == "" {
			t.Errorf("GetTraceCode must return a code, but is empty for: %s", throwable.Error())
		}
	}
}

func TestErrorsAsFindsTheConcreteError(t *testing.T) {
	t.Parallel()

	var err error = exception.NewContainerInvalidReferenceError("MissingID")

	target, found := errors.AsType[*exception.ContainerInvalidReferenceError](err)
	if !found {
		t.Fatal("errors.AsType must find the concrete error, but did not")
	}

	if target.GetID() != "MissingID" {
		t.Errorf("GetID must be the binding key, but is: %s", target.GetID())
	}
}

func TestErrorsAsSeparatesTheTwoConcreteErrors(t *testing.T) {
	t.Parallel()

	var err error = exception.NewContainerInvalidPublishCallbackError("InvalidID")

	if _, found := errors.AsType[*exception.ContainerInvalidReferenceError](err); found {
		t.Error("errors.AsType must not find the reference error in a publish callback error, but did")
	}
}
