/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

func TestTheRoutingErrorReportsTheParameterItCarries(t *testing.T) {
	t.Parallel()

	err := exception.NewHttpRoutingInvalidRouteParameterError("slug")

	if err.GetName() != "slug" {
		t.Errorf("GetName must return the parameter name, but returned: %q", err.GetName())
	}

	if err.Error() != "No parameter named 'slug' exists on this route" {
		t.Errorf("the error must name the parameter, but reads: %q", err.Error())
	}
}

func TestTheRoutingErrorSatisfiesTheSubComponentContract(t *testing.T) {
	t.Parallel()

	var throwable contract.HttpRoutingThrowable = exception.NewHttpRoutingInvalidRouteParameterError("slug")

	if !throwable.IsHttpRoutingThrowable() || !throwable.IsHttpThrowable() {
		t.Error("the error must mark itself for the routing sub-component and the component, but did not")
	}

	if throwable.GetTraceCode() == "" {
		t.Error("GetTraceCode must return a code, but is empty")
	}
}
