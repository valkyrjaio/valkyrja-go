/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
)

const (
	parameterName  = "id"
	parameterRegex = `(?P<id>\d+)`
	dynamicRegex   = `^/users/(?P<id>\d+)$`
)

func TestNewParameterHoldsWhatItReceives(t *testing.T) {
	t.Parallel()

	parameter := data.NewParameter(parameterName, parameterRegex)

	if parameter.GetName() != parameterName || parameter.GetRegex() != parameterRegex {
		t.Error("the parameter must hold its name and its regular expression, but did not")
	}

	if !parameter.ShouldCapture() {
		t.Error("a parameter must default to captured, but did not")
	}

	if parameter.IsOptional() || parameter.HasCast() {
		t.Error("a parameter must default to required and uncast, but did not")
	}

	if parameter.GetCast() != nil || parameter.GetDefault() != nil || parameter.GetValue() != nil {
		t.Error("a parameter must default to no cast, no default, and no value, but did not")
	}
}

func TestEachParameterWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	parameter := data.NewParameter(parameterName, parameterRegex)

	if parameter.WithName("slug").GetName() != "slug" {
		t.Error("WithName must hold the new name, but did not")
	}

	if parameter.WithRegex(`(?P<slug>\w+)`).GetRegex() != `(?P<slug>\w+)` {
		t.Error("WithRegex must hold the new regular expression, but did not")
	}

	if parameter.GetName() != parameterName || parameter.GetRegex() != parameterRegex {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachRemainingParameterWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	parameter := data.NewParameter(parameterName, parameterRegex)

	if !parameter.WithCast("int").HasCast() || parameter.WithCast("int").GetCast() != "int" {
		t.Error("WithCast must hold the new cast, but did not")
	}

	if !parameter.WithIsOptional(true).IsOptional() {
		t.Error("WithIsOptional must hold the new flag, but did not")
	}

	if parameter.WithShouldCapture(false).ShouldCapture() {
		t.Error("WithShouldCapture must hold the new flag, but did not")
	}

	if parameter.WithDefault("1").GetDefault() != "1" {
		t.Error("WithDefault must hold the new default, but did not")
	}

	if parameter.WithValue("2").GetValue() != "2" {
		t.Error("WithValue must hold the new value, but did not")
	}

	if parameter.IsOptional() || parameter.HasCast() || parameter.GetValue() != nil {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestNewDynamicRouteHoldsItsRouteAndItsRegex(t *testing.T) {
	t.Parallel()

	route := data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, newHandler()),
		dynamicRegex,
		data.NewParameter(parameterName, parameterRegex),
	)

	if route.GetRegex() != dynamicRegex {
		t.Errorf("the route must hold its regular expression, but holds: %q", route.GetRegex())
	}

	if route.GetPath() != "/users/{id}" || route.GetName() != routeName {
		t.Error("the route must hold what its route holds, but did not")
	}

	if len(route.GetParameters()) != 1 {
		t.Errorf("the route must hold its parameters, but holds: %d", len(route.GetParameters()))
	}

	if !route.HasRequestMethod(constant.RequestMethodGet) {
		t.Error("the route must match what its route matches, but did not")
	}
}

func TestEachDynamicRouteWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	route := data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, newHandler()),
		dynamicRegex,
		data.NewParameter(parameterName, parameterRegex),
	)

	if route.WithRegex("^/other$").GetRegex() != "^/other$" {
		t.Error("WithRegex must hold the new regular expression, but did not")
	}

	replaced := route.WithParameters(data.NewParameter("slug", `(?P<slug>\w+)`))
	added := route.WithAddedParameters(data.NewParameter("slug", `(?P<slug>\w+)`))

	if len(replaced.GetParameters()) != 1 || replaced.GetParameters()[0].GetName() != "slug" {
		t.Error("WithParameters must replace the parameters, but did not")
	}

	if len(added.GetParameters()) != 2 {
		t.Errorf("WithAddedParameters must append the parameter, but holds: %d", len(added.GetParameters()))
	}

	if route.GetRegex() != dynamicRegex || len(route.GetParameters()) != 1 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestTheDynamicRouteSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var route contract.DynamicRouteContract = data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, newHandler()),
		dynamicRegex,
	)

	if route.GetRegex() != dynamicRegex {
		t.Error("the contract must read the regular expression, but did not")
	}
}

func TestARouteReportsWhetherItIsDynamic(t *testing.T) {
	t.Parallel()

	static := data.NewRoute(routePath, routeName, newHandler())

	if static.IsDynamic() {
		t.Error("a route with no regular expression and no parameter must not be dynamic, but is")
	}

	withRegex := data.NewDynamicRoute(static, dynamicRegex)
	withParameter := data.NewDynamicRoute(static, "", data.NewParameter(parameterName, parameterRegex))

	if !withRegex.IsDynamic() || !withParameter.IsDynamic() {
		t.Error("a route that names a regular expression or a parameter must be dynamic, but is not")
	}
}

func TestAnInheritedWithMethodKeepsTheDynamicState(t *testing.T) {
	t.Parallel()

	route := data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, newHandler()),
		dynamicRegex,
		data.NewParameter(parameterName, parameterRegex),
	)

	renamed, isDynamic := route.WithName("other").(contract.DynamicRouteContract)
	if !isDynamic {
		t.Fatal("an inherited With method must return a route that is still dynamic, but did not")
	}

	if renamed.GetRegex() != dynamicRegex {
		t.Errorf("the copy must keep the regular expression, but holds: %q", renamed.GetRegex())
	}

	if len(renamed.GetParameters()) != 1 {
		t.Errorf("the copy must keep the parameters, but holds: %d", len(renamed.GetParameters()))
	}
}
