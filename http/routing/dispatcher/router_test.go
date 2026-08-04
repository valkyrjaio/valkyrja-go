/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package dispatcher_test

import (
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
	middlewaredata "github.com/valkyrjaio/valkyrja-go/v26/http/middleware/data"
	"github.com/valkyrjaio/valkyrja-go/v26/http/middleware/handler"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/collection"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/dispatcher"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/matcher"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/processor"
)

const (
	usersPath = "/users"
	usersName = "users.index"
	bodyText  = "the users"
)

// responseFactoryFixture builds the responses that the router asks for.
type responseFactoryFixture struct {
	contract.ResponseFactoryContract
}

// CreateResponse builds a response that carries the content.
func (f *responseFactoryFixture) CreateResponse(
	content string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.ResponseContract {
	return response.NewResponseFromContent(content, statusCode, headers)
}

// newRequest builds a GET request for the path.
func newRequest(t *testing.T, path string) contract.ServerRequestContract {
	t.Helper()

	requestUri, err := uri.NewUri(constant.SchemeHttp, "", "", "example.com", 0, path, "", "")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	return request.NewServerRequest(requestUri, constant.RequestMethodGet, nil, nil)
}

// newRouter builds a router over a collection that holds the routes.
func newRouter(
	container containercontract.ContainerContract,
	routes ...contract.RouteContract,
) *dispatcher.Router {
	built := collection.NewRouteCollection()

	for _, route := range routes {
		built.Add(processor.NewProcessor().Route(route))
	}

	return dispatcher.NewRouter(
		container,
		matcher.NewMatcher(built),
		&responseFactoryFixture{},
		handler.NewRouteMatchedHandler(container),
		handler.NewRouteNotMatchedHandler(container),
		handler.NewRouteDispatchedHandler(container),
	)
}

// newRoute builds a route that returns the body text.
func newRoute(path string, methods ...constant.RequestMethod) *data.Route {
	return data.NewRoute(path, usersName, func(
		_ containercontract.ContainerContract,
		_ contract.RouteContract,
	) contract.ResponseContract {
		return response.NewResponseFromContent(bodyText, constant.StatusCodeOk, nil)
	}, methods...)
}

func TestDispatchRunsTheRouteThatTheRequestMatches(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	built := newRouter(container, newRoute(usersPath)).Dispatch(newRequest(t, usersPath))

	if built.GetStatusCode() != constant.StatusCodeOk {
		t.Errorf("the response must carry the status of the handler, but is: %d", built.GetStatusCode())
	}

	if built.GetBody().String() != bodyText {
		t.Errorf("the response must carry the body of the handler, but is: %q", built.GetBody().String())
	}
}

func TestDispatchBindsTheRouteThatMatched(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	newRouter(container, newRoute(usersPath)).Dispatch(newRequest(t, usersPath))

	bound, err := container.GetSingleton(routingconstant.RouteContractServiceID)
	if err != nil {
		t.Fatalf("the router must bind the route that matched, but reported: %v", err)
	}

	route, isRoute := bound.(contract.RouteContract)
	if !isRoute {
		t.Fatalf("the binding must be a route, but is: %T", bound)
	}

	if route.GetName() != usersName {
		t.Errorf("the router must bind the route that matched, but bound: %q", route.GetName())
	}
}

func TestDispatchReportsAPathThatMatchesNoRoute(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	built := newRouter(container, newRoute(usersPath)).Dispatch(newRequest(t, "/missing"))

	if built.GetStatusCode() != constant.StatusCodeNotFound {
		t.Errorf("a path that matches no route must report 404, but reported: %d", built.GetStatusCode())
	}
}

func TestDispatchReportsAMethodThatTheRouteDoesNotMatch(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	built := newRouter(container, newRoute(usersPath, constant.RequestMethodAny)).
		Dispatch(request.NewServerRequest(nil, constant.RequestMethodPost, nil, nil))

	if built.GetStatusCode() != constant.StatusCodeNotFound {
		t.Errorf("a request for another path must report 404, but reported: %d", built.GetStatusCode())
	}
}

func TestDispatchReportsAMethodNotAllowed(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// The route matches the any method, so the path matches for every method.
	// The request asks for a method that the collection files, so the router
	// reads the any method and reports 405 rather than 404.
	built := collection.NewRouteCollection()
	built.Add(newRoute(usersPath, constant.RequestMethodPost))

	router := dispatcher.NewRouter(
		container,
		&anyMatchingMatcherFixture{},
		&responseFactoryFixture{},
		handler.NewRouteMatchedHandler(container),
		handler.NewRouteNotMatchedHandler(container),
		handler.NewRouteDispatchedHandler(container),
	)

	if router.Dispatch(newRequest(t, usersPath)).GetStatusCode() != constant.StatusCodeMethodNotAllowed {
		t.Error("a path that another method matches must report 405, but did not")
	}
}

func TestDispatchDecodesThePath(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	built := newRouter(container, newRoute("/a b")).Dispatch(newRequest(t, "/a%20b"))

	if built.GetBody().String() != bodyText {
		t.Errorf("the router must decode the path, but the response is: %q", built.GetBody().String())
	}
}

func TestDispatchMatchesAPathThatNoDecoderReads(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	built := newRouter(container, newRoute("/a%zzb")).Dispatch(newRequest(t, "/a%zzb"))

	if built.GetStatusCode() != constant.StatusCodeOk {
		t.Errorf("a path that no decoder reads must match as it arrived, but reported: %d",
			built.GetStatusCode())
	}
}

func TestDispatchRouteReturnsAnEmptyResponseForARouteWithNoHandler(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	built := newRouter(container).
		DispatchRoute(newRequest(t, usersPath), data.NewRoute(usersPath, usersName, nil))

	if built.GetStatusCode() != constant.StatusCodeNoContent {
		t.Errorf("a route with no handler must report 204, but reported: %d", built.GetStatusCode())
	}
}

func TestDispatchRouteEndsWhereTheRouteMatchedMiddlewareReturnsAResponse(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	ending := response.NewResponseFromContent("ended", constant.StatusCodeAccepted, nil)

	container.Bind("Valkyrja.Tests.Http.EndingMiddleware",
		func(_ containercontract.ContainerContract, _ []any) any {
			return &endingRouteMatchedFixture{response: ending}
		})

	router := dispatcher.NewRouter(
		container,
		matcher.NewMatcher(collection.NewRouteCollection()),
		&responseFactoryFixture{},
		handler.NewRouteMatchedHandler(container, "Valkyrja.Tests.Http.EndingMiddleware"),
		handler.NewRouteNotMatchedHandler(container),
		handler.NewRouteDispatchedHandler(container),
	)

	built := router.DispatchRoute(newRequest(t, usersPath), newRoute(usersPath))

	if built.GetBody().String() != "ended" {
		t.Errorf("the middleware must end the run, but the response is: %q", built.GetBody().String())
	}
}

func TestTheRouterSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	var router contract.RouterContract = newRouter(container, newRoute(usersPath))

	if router.Dispatch(newRequest(t, usersPath)).GetStatusCode() != constant.StatusCodeOk {
		t.Error("the contract must run the route, but did not")
	}
}

// anyMatchingMatcherFixture matches nothing for a concrete request method, and
// matches for the any method. That is the state that makes the router report
// that the method is not allowed.
type anyMatchingMatcherFixture struct{}

// Match returns a route for the any method only.
func (m *anyMatchingMatcherFixture) Match(
	path string,
	requestMethod constant.RequestMethod,
) contract.RouteContract {
	if requestMethod == constant.RequestMethodAny {
		return newRoute(path)
	}

	return nil
}

// MatchStatic returns what Match returns.
func (m *anyMatchingMatcherFixture) MatchStatic(
	path string,
	requestMethod constant.RequestMethod,
) contract.RouteContract {
	return m.Match(path, requestMethod)
}

// MatchDynamic returns what Match returns.
func (m *anyMatchingMatcherFixture) MatchDynamic(
	path string,
	requestMethod constant.RequestMethod,
) contract.RouteContract {
	return m.Match(path, requestMethod)
}

// endingRouteMatchedFixture is a route-matched middleware that ends the run.
type endingRouteMatchedFixture struct {
	response contract.ResponseContract
}

// RouteMatched ends the run with the response that the fixture holds.
func (m *endingRouteMatchedFixture) RouteMatched(
	_ contract.ServerRequestContract,
	_ contract.RouteContract,
	_ contract.RouteMatchedHandlerContract,
) contract.RouteMatchedResultContract {
	return middlewaredata.NewRouteMatchedResponse(m.response)
}

// rawPathUriFixture is a URI that returns a path which no decoder reads.
//
// A real URI encodes its path, so a decoder always reads one back. This fixture
// is what makes the router's fallback reachable: it hands over the path as it
// arrived, invalid escape and all.
type rawPathUriFixture struct {
	contract.UriContract
}

// GetPath returns a path that carries an invalid escape.
func (u *rawPathUriFixture) GetPath() string { return "/a%zzb" }

// GetHost returns no host, so the request adds no host header.
func (u *rawPathUriFixture) GetHost() string { return "" }

func TestDispatchMatchesARawPathThatNoDecoderReads(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	built := newRouter(container, newRoute("/a%zzb")).Dispatch(
		request.NewServerRequest(&rawPathUriFixture{}, constant.RequestMethodGet, nil, nil),
	)

	if built.GetStatusCode() != constant.StatusCodeOk {
		t.Errorf("a path that no decoder reads must match as it arrived, but reported: %d",
			built.GetStatusCode())
	}
}
