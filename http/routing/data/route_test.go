/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"slices"
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
)

const (
	routeName    = "users.index"
	routePath    = "/users"
	firstIDName  = "Valkyrja.Tests.Http.FirstMiddleware"
	secondIDName = "Valkyrja.Tests.Http.SecondMiddleware"
)

// newHandler builds a handler that returns an empty response.
func newHandler() contract.HttpHandlerFunc {
	return func(_ containercontract.ContainerContract, _ contract.RouteContract) contract.ResponseContract {
		return response.NewResponse(nil, 0, nil)
	}
}

// newRoute builds a route at the standard path, under the standard name.
func newRoute() *data.Route {
	return data.NewRoute(routePath, routeName, newHandler())
}

func TestNewRouteHoldsWhatItReceives(t *testing.T) {
	t.Parallel()

	route := newRoute()

	if route.GetPath() != routePath || route.GetName() != routeName {
		t.Error("the route must hold its path and its name, but did not")
	}

	if route.GetHandler() == nil {
		t.Error("the route must hold its handler, but did not")
	}

	if !slices.Equal(route.GetRequestMethods(), []constant.RequestMethod{constant.RequestMethodGet}) {
		t.Errorf("the route must default to the GET method, but matches: %v", route.GetRequestMethods())
	}
}

func TestHasRequestMethodReadsWhatTheRouteMatches(t *testing.T) {
	t.Parallel()

	route := newRoute()

	if !route.HasRequestMethod(constant.RequestMethodGet) {
		t.Error("the route must match the GET method, but did not")
	}

	if route.HasRequestMethod(constant.RequestMethodPost) {
		t.Error("the route must not match the POST method, but did")
	}
}

func TestARouteThatMatchesTheAnyMethodMatchesEveryMethod(t *testing.T) {
	t.Parallel()

	route := data.NewRoute(routePath, routeName, newHandler(), constant.RequestMethodAny)

	for _, method := range constant.GetAllRequestMethods() {
		if !route.HasRequestMethod(method) {
			t.Errorf("the route must match %q, but did not", method)
		}
	}
}

func TestEachPathAndNameMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	route := newRoute()

	if route.WithPath("/other").GetPath() != "/other" {
		t.Error("WithPath must hold the new path, but did not")
	}

	if route.WithAddedPath("/{id}").GetPath() != "/users/{id}" {
		t.Error("WithAddedPath must append the path, but did not")
	}

	if route.WithName("other").GetName() != "other" {
		t.Error("WithName must hold the new name, but did not")
	}

	if route.WithAddedName(".show").GetName() != "users.index.show" {
		t.Error("WithAddedName must append the name, but did not")
	}

	if route.GetPath() != routePath || route.GetName() != routeName {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestWithHandlerAndWithRequestMethodsReturnACopy(t *testing.T) {
	t.Parallel()

	route := newRoute()

	replaced := route.WithRequestMethods(constant.RequestMethodPost)
	added := route.WithAddedRequestMethods(constant.RequestMethodPut)

	if !slices.Equal(replaced.GetRequestMethods(), []constant.RequestMethod{constant.RequestMethodPost}) {
		t.Error("WithRequestMethods must replace the methods, but did not")
	}

	if len(added.GetRequestMethods()) != 2 {
		t.Errorf("WithAddedRequestMethods must append the method, but matches: %v", added.GetRequestMethods())
	}

	if route.WithHandler(nil).GetHandler() != nil {
		t.Error("WithHandler must hold the new handler, but did not")
	}

	if len(route.GetRequestMethods()) != 1 || route.GetHandler() == nil {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachMiddlewareStageHoldsItsOwnList(t *testing.T) {
	t.Parallel()

	route := newRoute().
		WithRouteMatchedMiddleware(firstIDName).
		WithRouteDispatchedMiddleware(firstIDName).
		WithThrowableCaughtMiddleware(firstIDName).
		WithSendingResponseMiddleware(firstIDName).
		WithResponseSentMiddleware(firstIDName)

	stages := map[string][]string{
		"route matched":    route.GetRouteMatchedMiddleware(),
		"route dispatched": route.GetRouteDispatchedMiddleware(),
		"throwable caught": route.GetThrowableCaughtMiddleware(),
		"sending response": route.GetSendingResponseMiddleware(),
		"response sent":    route.GetResponseSentMiddleware(),
	}

	for name, middleware := range stages {
		if !slices.Equal(middleware, []string{firstIDName}) {
			t.Errorf("the %s stage must hold its middleware, but holds: %v", name, middleware)
		}
	}
}

func TestEachAddedMiddlewareStageAppends(t *testing.T) {
	t.Parallel()

	route := newRoute().
		WithRouteMatchedMiddleware(firstIDName).
		WithAddedRouteMatchedMiddleware(secondIDName)

	if !slices.Equal(route.GetRouteMatchedMiddleware(), []string{firstIDName, secondIDName}) {
		t.Errorf("the stage must append the middleware, but holds: %v", route.GetRouteMatchedMiddleware())
	}

	stages := map[string]func() []string{
		"route dispatched": newRoute().WithAddedRouteDispatchedMiddleware(firstIDName, secondIDName).
			GetRouteDispatchedMiddleware,
		"throwable caught": newRoute().WithAddedThrowableCaughtMiddleware(firstIDName, secondIDName).
			GetThrowableCaughtMiddleware,
		"sending response": newRoute().WithAddedSendingResponseMiddleware(firstIDName, secondIDName).
			GetSendingResponseMiddleware,
		"response sent": newRoute().WithAddedResponseSentMiddleware(firstIDName, secondIDName).
			GetResponseSentMiddleware,
	}

	for name, get := range stages {
		if len(get()) != 2 {
			t.Errorf("the %s stage must append each middleware, but holds: %v", name, get())
		}
	}
}

func TestAMiddlewareThatIsAddedTwiceIsHeldTwice(t *testing.T) {
	t.Parallel()

	route := newRoute().
		WithRouteMatchedMiddleware(firstIDName).
		WithAddedRouteMatchedMiddleware(firstIDName)

	if !slices.Equal(route.GetRouteMatchedMiddleware(), []string{firstIDName, firstIDName}) {
		t.Errorf("a middleware that is added twice must be held twice, but holds: %v",
			route.GetRouteMatchedMiddleware())
	}
}

func TestTheRouteHoldsItsStructs(t *testing.T) {
	t.Parallel()

	route := newRoute()

	if route.HasRequestStruct() || route.HasResponseStruct() {
		t.Error("a route must name no struct by default, but named one")
	}

	if route.GetRequestStruct() != nil || route.GetResponseStruct() != nil {
		t.Error("a route must return no struct by default, but returned one")
	}

	withRequest := route.WithRequestStruct(&requestStructFixture{})
	withResponse := route.WithResponseStruct(&responseStructFixture{})

	if !withRequest.HasRequestStruct() || withRequest.GetRequestStruct() == nil {
		t.Error("WithRequestStruct must hold the struct, but did not")
	}

	if !withResponse.HasResponseStruct() || withResponse.GetResponseStruct() == nil {
		t.Error("WithResponseStruct must hold the struct, but did not")
	}
}

// requestStructFixture is a request struct that reads nothing.
type requestStructFixture struct{}

// GetName returns the name of the struct.
func (s *requestStructFixture) GetName() string { return "request" }

// GetValue returns the value of the struct.
func (s *requestStructFixture) GetValue() any { return nil }

// GetDataFromRequest returns no field.
func (s *requestStructFixture) GetDataFromRequest(_ contract.ServerRequestContract) map[string]any {
	return map[string]any{}
}

// DetermineIfRequestContainsExtraData reports no extra field.
func (s *requestStructFixture) DetermineIfRequestContainsExtraData(_ contract.ServerRequestContract) bool {
	return false
}

// responseStructFixture is a response struct that shapes nothing.
type responseStructFixture struct{}

// GetName returns the name of the struct.
func (s *responseStructFixture) GetName() string { return "response" }

// GetValue returns the value of the struct.
func (s *responseStructFixture) GetValue() any { return nil }

// GetStructuredData returns the data as it received it.
func (s *responseStructFixture) GetStructuredData(data map[string]any, _ bool) map[string]any {
	return data
}
