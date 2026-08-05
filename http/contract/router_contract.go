/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	"net/http"

	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// HttpRoutingDataContract is the HTTP routing component's state, as a value that
// the framework stores and reloads. `sindri` generates one for the whole
// application, and the route collection loads it at boot.
//
// The contract names an interface rather than a concrete type, for the reason
// that `ContainerDataContract` does: the data holds a route, and a route names
// contracts in this package.
type HttpRoutingDataContract interface {
	// GetRoutes returns each route, keyed by its own name.
	GetRoutes() map[string]RouteContract

	// GetPaths returns the name of the route at each static path, keyed by
	// request method and then by path.
	GetPaths() map[constant.RequestMethod]map[string]string

	// GetRegexes returns the name of the route at each dynamic path, keyed by
	// request method and then by regular expression.
	GetRegexes() map[constant.RequestMethod]map[string]string
}

// RouteCollectionContract holds every route of the application.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type RouteCollectionContract interface {
	// GetData returns the collection's state.
	GetData() HttpRoutingDataContract

	// SetFromData replaces the collection's state.
	SetFromData(data HttpRoutingDataContract)

	// Add files the route under its own name, its path, and each request
	// method that it matches.
	Add(route RouteContract)

	// HasPath reports whether a route matches the static path and the request
	// method.
	HasPath(path string, method constant.RequestMethod) bool

	// GetByPath returns the route at the static path.
	GetByPath(path string, method constant.RequestMethod) RouteContract

	// HasRegex reports whether a route matches the regular expression and the
	// request method.
	HasRegex(regex string, method constant.RequestMethod) bool

	// GetByRegex returns the route at the regular expression.
	GetByRegex(regex string, method constant.RequestMethod) DynamicRouteContract

	// GetPaths returns the name of the route at each static path.
	GetPaths(method constant.RequestMethod) map[string]string

	// GetRegexes returns the name of the route at each regular expression.
	GetRegexes(method constant.RequestMethod) map[string]string

	// HasName reports whether the collection holds a route under the name.
	HasName(name string) bool

	// GetByName returns the route under the name.
	GetByName(name string) RouteContract

	// GetAll returns every route that matches the request method.
	GetAll(method constant.RequestMethod) map[string]RouteContract
}

// This port declares no route collector. The other ports give
// `RouteCollectorContract` a `getRoutes` that takes the controller classes to
// scan for a routing annotation, and Go has no annotation. Without the classes
// the contract is `getRoutes` alone, which is `HttpRouteProviderContract` — and
// Go compares an interface by its method set, so the two would be one type. A
// route reaches the collection through its provider.

// MatcherContract finds the route that a request path matches.
type MatcherContract interface {
	// Match returns the route that the path matches, and nil where none
	// matches.
	Match(path string, requestMethod constant.RequestMethod) RouteContract

	// MatchStatic returns the static route that the path matches, and nil where
	// none matches.
	MatchStatic(path string, requestMethod constant.RequestMethod) RouteContract

	// MatchDynamic returns the dynamic route that the path matches, and nil
	// where none matches.
	MatchDynamic(path string, requestMethod constant.RequestMethod) RouteContract
}

// ProcessorContract prepares a route before the collection holds it. It compiles
// the regular expression of a dynamic route.
type ProcessorContract interface {
	// Route returns the route, prepared.
	Route(route RouteContract) RouteContract
}

// RouterContract runs the route that a request matches.
type RouterContract interface {
	// Dispatch matches the request to a route and runs it.
	Dispatch(request ServerRequestContract) ResponseContract

	// DispatchRoute runs the route for the request.
	DispatchRoute(request ServerRequestContract, route RouteContract) ResponseContract
}

// UrlContract builds the URL of a named route.
type UrlContract interface {
	// GetUrl returns the URL of the route, with the data filled into each
	// parameter of its path.
	GetUrl(name string, data map[string]string) string
}

// RoutingResponseFactoryContract builds a response that sends the client to a
// named route.
type RoutingResponseFactoryContract interface {
	// CreateRouteRedirectResponse builds a response that sends the client to
	// the route.
	CreateRouteRedirectResponse(
		name string,
		data map[string]string,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) RedirectResponseContract
}

// RequestHandlerContract is the server's entry point for one request.
//
// The TypeScript port types the writer as Node's `ServerResponse`. The Go
// equivalent is `http.ResponseWriter`, which is the one place where this port
// names a type outside the framework.
type RequestHandlerContract interface {
	// Handle returns the response for the request.
	Handle(request ServerRequestContract) ResponseContract

	// Send writes the response to the writer.
	Send(response ResponseContract, writer http.ResponseWriter) RequestHandlerContract

	// Terminate runs what is left after the server sent the response.
	Terminate(request ServerRequestContract, response ResponseContract)

	// Run handles the request, sends the response, and terminates.
	Run(request ServerRequestContract, writer http.ResponseWriter)
}

// ExceptionResponseRequestHandlerContract is a request handler that turns a
// failure into a response.
type ExceptionResponseRequestHandlerContract interface {
	RequestHandlerContract

	// CreateResponseFromThrowable returns the response for the failure.
	CreateResponseFromThrowable(throwable error) ResponseContract
}

// ClientContract sends a request to another server and returns its response.
type ClientContract interface {
	// SendRequest sends the request. It reports a failure where the client
	// cannot reach the server.
	SendRequest(request RequestContract) (ResponseContract, error)
}
