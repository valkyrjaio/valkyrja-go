/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds the routes of the HTTP router.
package data

import (
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// Route is one route of the HTTP router.
//
// The other ports declare a `DynamicRoute` that extends `Route` and adds the
// regular expression and the parameters. Go embeds rather than inherits, and a
// `With` method promoted from an embedded struct copies only that struct — so an
// inherited `With` on a dynamic route would return a static one and drop its
// regular expression. One struct carries both, and `IsDynamic` reports which it
// is.
//
// Each middleware is named by its binding key. The other ports pass a class
// reference, and Go has none, so it passes the same string constant that it
// binds the middleware under.
//
// Warning: a route appends each middleware and never dedupes. A middleware that
// is registered twice runs twice. That is the developer's error, and neither the
// route nor `sindri` corrects it.
type Route struct {
	path           string
	name           string
	handler        contract.HttpHandlerFunc
	requestMethods []constant.RequestMethod

	routeMatchedMiddleware    []string
	routeDispatchedMiddleware []string
	throwableCaughtMiddleware []string
	sendingResponseMiddleware []string
	responseSentMiddleware    []string

	requestStruct  contract.RequestStructContract
	responseStruct contract.ResponseStructContract

	regex      string
	parameters []contract.ParameterContract
}

// NewRoute builds a route at a path, under a name, that runs a handler.
//
// A route with no request method matches the GET method, which is what every
// port does.
func NewRoute(
	path string,
	name string,
	handler contract.HttpHandlerFunc,
	requestMethods ...constant.RequestMethod,
) *Route {
	if len(requestMethods) == 0 {
		requestMethods = []constant.RequestMethod{constant.RequestMethodGet}
	}

	return &Route{
		path:           path,
		name:           name,
		handler:        handler,
		requestMethods: requestMethods,
	}
}

// GetPath returns the path that the route matches.
func (r *Route) GetPath() string {
	return r.path
}

// WithPath returns a copy of the route for another path.
func (r *Route) WithPath(path string) contract.RouteContract {
	copied := *r
	copied.path = path

	return &copied
}

// WithAddedPath returns a copy of the route with the path added to the one it
// holds.
func (r *Route) WithAddedPath(path string) contract.RouteContract {
	copied := *r
	copied.path += path

	return &copied
}

// GetName returns the name of the route.
func (r *Route) GetName() string {
	return r.name
}

// WithName returns a copy of the route under another name.
func (r *Route) WithName(name string) contract.RouteContract {
	copied := *r
	copied.name = name

	return &copied
}

// WithAddedName returns a copy of the route with the name added to the one it
// holds.
func (r *Route) WithAddedName(name string) contract.RouteContract {
	copied := *r
	copied.name += name

	return &copied
}

// GetHandler returns what the route runs.
func (r *Route) GetHandler() contract.HttpHandlerFunc {
	return r.handler
}

// WithHandler returns a copy of the route that runs another handler.
func (r *Route) WithHandler(handler contract.HttpHandlerFunc) contract.RouteContract {
	copied := *r
	copied.handler = handler

	return &copied
}

// GetRequestMethods returns each request method that the route matches.
func (r *Route) GetRequestMethods() []constant.RequestMethod {
	return r.requestMethods
}

// HasRequestMethod reports whether the route matches the request method. A route
// that matches the any method matches every request method.
func (r *Route) HasRequestMethod(requestMethod constant.RequestMethod) bool {
	if slices.Contains(r.requestMethods, constant.RequestMethodAny) {
		return true
	}

	return slices.Contains(r.requestMethods, requestMethod)
}

// WithRequestMethods returns a copy of the route for other request methods.
func (r *Route) WithRequestMethods(requestMethods ...constant.RequestMethod) contract.RouteContract {
	copied := *r
	copied.requestMethods = requestMethods

	return &copied
}

// WithAddedRequestMethods returns a copy of the route with the request methods
// added to the ones it holds.
func (r *Route) WithAddedRequestMethods(requestMethods ...constant.RequestMethod) contract.RouteContract {
	copied := *r
	copied.requestMethods = appendStrings(r.requestMethods, requestMethods)

	return &copied
}

// GetRouteMatchedMiddleware returns the binding key of each route-matched
// middleware.
func (r *Route) GetRouteMatchedMiddleware() []string {
	return r.routeMatchedMiddleware
}

// WithRouteMatchedMiddleware returns a copy of the route with other
// route-matched middleware.
func (r *Route) WithRouteMatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeMatchedMiddleware = middleware

	return &copied
}

// WithAddedRouteMatchedMiddleware returns a copy of the route with the
// route-matched middleware appended.
func (r *Route) WithAddedRouteMatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeMatchedMiddleware = appendStrings(r.routeMatchedMiddleware, middleware)

	return &copied
}

// GetRouteDispatchedMiddleware returns the binding key of each route-dispatched
// middleware.
func (r *Route) GetRouteDispatchedMiddleware() []string {
	return r.routeDispatchedMiddleware
}

// WithRouteDispatchedMiddleware returns a copy of the route with other
// route-dispatched middleware.
func (r *Route) WithRouteDispatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeDispatchedMiddleware = middleware

	return &copied
}

// WithAddedRouteDispatchedMiddleware returns a copy of the route with the
// route-dispatched middleware appended.
func (r *Route) WithAddedRouteDispatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeDispatchedMiddleware = appendStrings(r.routeDispatchedMiddleware, middleware)

	return &copied
}

// GetThrowableCaughtMiddleware returns the binding key of each throwable-caught
// middleware.
func (r *Route) GetThrowableCaughtMiddleware() []string {
	return r.throwableCaughtMiddleware
}

// WithThrowableCaughtMiddleware returns a copy of the route with other
// throwable-caught middleware.
func (r *Route) WithThrowableCaughtMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.throwableCaughtMiddleware = middleware

	return &copied
}

// WithAddedThrowableCaughtMiddleware returns a copy of the route with the
// throwable-caught middleware appended.
func (r *Route) WithAddedThrowableCaughtMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.throwableCaughtMiddleware = appendStrings(r.throwableCaughtMiddleware, middleware)

	return &copied
}

// GetSendingResponseMiddleware returns the binding key of each sending-response
// middleware.
func (r *Route) GetSendingResponseMiddleware() []string {
	return r.sendingResponseMiddleware
}

// WithSendingResponseMiddleware returns a copy of the route with other
// sending-response middleware.
func (r *Route) WithSendingResponseMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.sendingResponseMiddleware = middleware

	return &copied
}

// WithAddedSendingResponseMiddleware returns a copy of the route with the
// sending-response middleware appended.
func (r *Route) WithAddedSendingResponseMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.sendingResponseMiddleware = appendStrings(r.sendingResponseMiddleware, middleware)

	return &copied
}

// GetResponseSentMiddleware returns the binding key of each response-sent
// middleware.
func (r *Route) GetResponseSentMiddleware() []string {
	return r.responseSentMiddleware
}

// WithResponseSentMiddleware returns a copy of the route with other
// response-sent middleware.
func (r *Route) WithResponseSentMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.responseSentMiddleware = middleware

	return &copied
}

// WithAddedResponseSentMiddleware returns a copy of the route with the
// response-sent middleware appended.
func (r *Route) WithAddedResponseSentMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.responseSentMiddleware = appendStrings(r.responseSentMiddleware, middleware)

	return &copied
}

// HasRequestStruct reports whether the route names a request struct.
func (r *Route) HasRequestStruct() bool {
	return r.requestStruct != nil
}

// GetRequestStruct returns the request struct of the route.
func (r *Route) GetRequestStruct() contract.RequestStructContract {
	return r.requestStruct
}

// WithRequestStruct returns a copy of the route for another request struct.
func (r *Route) WithRequestStruct(requestStruct contract.RequestStructContract) contract.RouteContract {
	copied := *r
	copied.requestStruct = requestStruct

	return &copied
}

// HasResponseStruct reports whether the route names a response struct.
func (r *Route) HasResponseStruct() bool {
	return r.responseStruct != nil
}

// GetResponseStruct returns the response struct of the route.
func (r *Route) GetResponseStruct() contract.ResponseStructContract {
	return r.responseStruct
}

// WithResponseStruct returns a copy of the route for another response struct.
func (r *Route) WithResponseStruct(responseStruct contract.ResponseStructContract) contract.RouteContract {
	copied := *r
	copied.responseStruct = responseStruct

	return &copied
}

// appendStrings returns the values with the added ones after them, in a slice of
// its own, so a copy never shares a backing array with the route it came from.
func appendStrings[T ~string](values []T, added []T) []T {
	combined := make([]T, 0, len(values)+len(added))
	combined = append(combined, values...)
	combined = append(combined, added...)

	return combined
}
