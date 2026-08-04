/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package matcher_test

import (
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/matcher"
)

const (
	staticName   = "users.index"
	staticPath   = "/users"
	dynamicName  = "users.show"
	dynamicRegex = `^/users/(?P<id>\d+)$`
)

// newHandler builds a handler that returns an empty response.
func newHandler() contract.HttpHandlerFunc {
	return func(_ containercontract.ContainerContract, _ contract.RouteContract) contract.ResponseContract {
		return response.NewResponse(nil, 0, nil)
	}
}

// newCollection builds a collection that holds a static and a dynamic route.
func newCollection() contract.RouteCollectionContract {
	built := collection.NewRouteCollection()

	built.Add(data.NewRoute(staticPath, staticName, newHandler()))
	built.Add(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", dynamicName, newHandler()),
		dynamicRegex,
		data.NewParameter("id", `(?P<id>\d+)`),
	))

	return built
}

func TestMatchStaticFindsTheStaticRoute(t *testing.T) {
	t.Parallel()

	route := matcher.NewMatcher(newCollection()).MatchStatic(staticPath, constant.RequestMethodGet)

	if route == nil || route.GetName() != staticName {
		t.Error("MatchStatic must find the static route, but did not")
	}
}

func TestMatchStaticIsNilForAnUnknownPath(t *testing.T) {
	t.Parallel()

	if matcher.NewMatcher(newCollection()).MatchStatic("/missing", constant.RequestMethodGet) != nil {
		t.Error("MatchStatic must be nil for an unknown path, but is not")
	}
}

func TestMatchDynamicFindsTheDynamicRoute(t *testing.T) {
	t.Parallel()

	route := matcher.NewMatcher(newCollection()).MatchDynamic("/users/42", constant.RequestMethodGet)

	if route == nil || route.GetName() != dynamicName {
		t.Fatal("MatchDynamic must find the dynamic route, but did not")
	}

	dynamic, isDynamic := route.(contract.DynamicRouteContract)
	if !isDynamic {
		t.Fatalf("MatchDynamic must return a dynamic route, but returned: %T", route)
	}

	if dynamic.GetParameters()[0].GetValue() != "42" {
		t.Errorf("the parameter must carry the value from the path, but carries: %v",
			dynamic.GetParameters()[0].GetValue())
	}
}

func TestMatchDynamicIsNilWhereThePathDoesNotMatch(t *testing.T) {
	t.Parallel()

	built := matcher.NewMatcher(newCollection())

	if built.MatchDynamic("/users/abc", constant.RequestMethodGet) != nil {
		t.Error("MatchDynamic must be nil where the path does not match, but is not")
	}

	if built.MatchDynamic("/users/42/edit", constant.RequestMethodGet) != nil {
		t.Error("the anchors must make the match exact, but a longer path matched")
	}
}

func TestMatchReadsTheStaticRouteFirst(t *testing.T) {
	t.Parallel()

	built := matcher.NewMatcher(newCollection())

	if built.Match(staticPath, constant.RequestMethodGet).GetName() != staticName {
		t.Error("Match must find the static route, but did not")
	}

	if built.Match("/users/42", constant.RequestMethodGet).GetName() != dynamicName {
		t.Error("Match must fall back to the dynamic route, but did not")
	}

	if built.Match("/missing", constant.RequestMethodGet) != nil {
		t.Error("Match must be nil where no route matches, but is not")
	}
}

func TestAParameterTakesItsDefaultWhereThePathCarriesNone(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", dynamicName, newHandler()),
		`^/users(?:/(?P<id>\d+))?$`,
		data.NewParameter("id", `(?P<id>\d+)`).WithDefault("1").WithIsOptional(true),
	))

	route := matcher.NewMatcher(built).MatchDynamic("/users", constant.RequestMethodGet)

	dynamic, isDynamic := route.(contract.DynamicRouteContract)
	if !isDynamic {
		t.Fatalf("MatchDynamic must return a dynamic route, but returned: %T", route)
	}

	if dynamic.GetParameters()[0].GetValue() != "1" {
		t.Errorf("the parameter must take its default, but carries: %v", dynamic.GetParameters()[0].GetValue())
	}
}

func TestAParameterThatTheRegexDoesNotNameTakesItsDefault(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", dynamicName, newHandler()),
		dynamicRegex,
		data.NewParameter("slug", `(?P<slug>\w+)`).WithDefault("none"),
	))

	route := matcher.NewMatcher(built).MatchDynamic("/users/42", constant.RequestMethodGet)

	dynamic, isDynamic := route.(contract.DynamicRouteContract)
	if !isDynamic {
		t.Fatalf("MatchDynamic must return a dynamic route, but returned: %T", route)
	}

	if dynamic.GetParameters()[0].GetValue() != "none" {
		t.Errorf("a parameter that the regular expression does not name must take its default, but carries: %v",
			dynamic.GetParameters()[0].GetValue())
	}
}

func TestMatchDynamicSkipsARegexThatRe2Rejects(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", dynamicName, newHandler()),
		`^/users/(?=\d+)$`,
		data.NewParameter("id", ""),
	))

	if matcher.NewMatcher(built).MatchDynamic("/users/42", constant.RequestMethodGet) != nil {
		t.Error("a pattern that RE2 rejects must match nothing, but matched")
	}
}

func TestMatchDynamicSkipsARegexWithNoRoute(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.SetFromData(data.NewHttpRoutingData(
		nil,
		nil,
		map[constant.RequestMethod]map[string]string{
			constant.RequestMethodGet: {dynamicRegex: "missing"},
		},
	))

	if matcher.NewMatcher(built).MatchDynamic("/users/42", constant.RequestMethodGet) != nil {
		t.Error("a regular expression whose route is absent must match nothing, but matched")
	}
}

func TestTheMatcherSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.MatcherContract = matcher.NewMatcher(newCollection())

	if built.Match(staticPath, constant.RequestMethodGet) == nil {
		t.Error("the contract must find the route, but did not")
	}
}
