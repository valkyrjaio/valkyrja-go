/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package collection_test

import (
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
)

const (
	routeName    = "users.index"
	routePath    = "/users"
	dynamicName  = "users.show"
	dynamicRegex = `^/users/(?P<id>\d+)$`
)

// newHandler builds a handler that returns an empty response.
func newHandler() contract.HttpHandlerFunc {
	return func(_ containercontract.ContainerContract, _ contract.RouteContract) contract.ResponseContract {
		return response.NewResponse(nil, 0, nil)
	}
}

// newStaticRoute builds a static route for the request methods.
func newStaticRoute(methods ...constant.RequestMethod) *data.Route {
	return data.NewRoute(routePath, routeName, newHandler(), methods...)
}

// newDynamicRoute builds a dynamic route for the GET method.
func newDynamicRoute() *data.Route {
	return data.NewDynamicRoute(
		data.NewRoute("/users/{id}", dynamicName, newHandler()),
		dynamicRegex,
		data.NewParameter("id", `(?P<id>\d+)`),
	)
}

func TestAddFilesAStaticRouteUnderItsPath(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(newStaticRoute())

	if !built.HasPath(routePath, constant.RequestMethodGet) {
		t.Error("the collection must file the route under its path, but did not")
	}

	if built.GetByPath(routePath, constant.RequestMethodGet).GetName() != routeName {
		t.Error("the collection must return the route at the path, but did not")
	}

	if !built.HasName(routeName) || built.GetByName(routeName) == nil {
		t.Error("the collection must file the route under its name, but did not")
	}
}

func TestAddFilesADynamicRouteUnderItsRegex(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(newDynamicRoute())

	if !built.HasRegex(dynamicRegex, constant.RequestMethodGet) {
		t.Error("the collection must file the route under its regular expression, but did not")
	}

	if built.GetByRegex(dynamicRegex, constant.RequestMethodGet) == nil {
		t.Error("the collection must return the route at the regular expression, but did not")
	}

	if built.HasPath("/users/{id}", constant.RequestMethodGet) {
		t.Error("a dynamic route must not file under its path, but did")
	}
}

func TestARouteFilesUnderEachMethodThatItMatches(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(newStaticRoute(constant.RequestMethodGet, constant.RequestMethodPost))

	if !built.HasPath(routePath, constant.RequestMethodGet) ||
		!built.HasPath(routePath, constant.RequestMethodPost) {
		t.Error("the route must file under each method that it matches, but did not")
	}

	if built.HasPath(routePath, constant.RequestMethodPut) {
		t.Error("the route must not file under a method that it does not match, but did")
	}
}

func TestARouteThatMatchesTheAnyMethodFilesUnderEveryMethod(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(newStaticRoute(constant.RequestMethodAny))

	for _, method := range constant.GetAllRequestMethods() {
		if !built.HasPath(routePath, method) {
			t.Errorf("the route must file under %q, but did not", method)
		}
	}
}

func TestEachLookupReportsAnUnknownRoute(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()

	if built.HasPath("/missing", constant.RequestMethodGet) {
		t.Error("HasPath must be false for an unknown path, but is true")
	}

	if built.GetByPath("/missing", constant.RequestMethodGet) != nil {
		t.Error("GetByPath must be nil for an unknown path, but is not")
	}

	if built.HasRegex("^/missing$", constant.RequestMethodGet) {
		t.Error("HasRegex must be false for an unknown regular expression, but is true")
	}

	if built.GetByRegex("^/missing$", constant.RequestMethodGet) != nil {
		t.Error("GetByRegex must be nil for an unknown regular expression, but is not")
	}

	if built.HasName("missing") || built.GetByName("missing") != nil {
		t.Error("the collection must report an unknown name, but did not")
	}
}

func TestGetByRegexIsNilWhereTheRouteIsStatic(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(newStaticRoute())

	if built.GetByRegex(routePath, constant.RequestMethodGet) != nil {
		t.Error("GetByRegex must be nil where the route is static, but is not")
	}
}

func TestGetPathsAndGetRegexesReadOneMethod(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(newStaticRoute())
	built.Add(newDynamicRoute())

	if len(built.GetPaths(constant.RequestMethodGet)) != 1 {
		t.Error("GetPaths must return the static path, but did not")
	}

	if len(built.GetRegexes(constant.RequestMethodGet)) != 1 {
		t.Error("GetRegexes must return the regular expression, but did not")
	}

	if len(built.GetPaths(constant.RequestMethodPost)) != 0 {
		t.Error("GetPaths must be empty for a method that no route matches, but is not")
	}
}

func TestGetAllReturnsEveryRouteOfOneMethod(t *testing.T) {
	t.Parallel()

	built := collection.NewRouteCollection()
	built.Add(newStaticRoute())
	built.Add(newDynamicRoute())

	all := built.GetAll(constant.RequestMethodGet)

	if len(all) != 2 {
		t.Errorf("GetAll must return each route, but returned: %d", len(all))
	}

	if all[routeName] == nil || all[dynamicName] == nil {
		t.Error("GetAll must return both the static and the dynamic route, but did not")
	}

	if len(built.GetAll(constant.RequestMethodPost)) != 0 {
		t.Error("GetAll must be empty for a method that no route matches, but is not")
	}
}

func TestGetDataAndSetFromDataRoundTrip(t *testing.T) {
	t.Parallel()

	source := collection.NewRouteCollection()
	source.Add(newStaticRoute())
	source.Add(newDynamicRoute())

	target := collection.NewRouteCollection()
	target.SetFromData(source.GetData())

	if !target.HasPath(routePath, constant.RequestMethodGet) {
		t.Error("the loaded collection must hold the static path, but did not")
	}

	if !target.HasRegex(dynamicRegex, constant.RequestMethodGet) {
		t.Error("the loaded collection must hold the regular expression, but did not")
	}

	if !target.HasName(routeName) || !target.HasName(dynamicName) {
		t.Error("the loaded collection must hold each route by name, but did not")
	}
}

func TestTheCollectionSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.RouteCollectionContract = collection.NewRouteCollection()

	built.Add(newStaticRoute())

	if !built.HasName(routeName) {
		t.Error("the contract must file the route, but did not")
	}
}
