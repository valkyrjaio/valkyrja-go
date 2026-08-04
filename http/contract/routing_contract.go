/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// HttpHandlerFunc runs one route. The router passes the container and the route
// that matched.
type HttpHandlerFunc func(container containercontract.ContainerContract, route RouteContract) ResponseContract

// StructContract names one field of a request or a response.
type StructContract interface {
	// GetName returns the name of the field.
	GetName() string

	// GetValue returns the value of the field.
	GetValue() any
}

// RequestStructContract reads the fields of a request.
type RequestStructContract interface {
	StructContract

	// GetDataFromRequest returns the fields that the struct names, read from
	// the request.
	GetDataFromRequest(request ServerRequestContract) map[string]any

	// DetermineIfRequestContainsExtraData reports whether the request carries a
	// field that the struct does not name.
	DetermineIfRequestContainsExtraData(request ServerRequestContract) bool
}

// ResponseStructContract shapes the data of a response.
type ResponseStructContract interface {
	StructContract

	// GetStructuredData returns the data in the shape that the struct names. It
	// keeps every field where includeAll is true.
	GetStructuredData(data map[string]any, includeAll bool) map[string]any
}

// ParameterContract is one dynamic segment of a route path.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type ParameterContract interface {
	// GetName returns the name of the parameter.
	GetName() string

	// WithName returns a copy of the parameter under another name.
	WithName(name string) ParameterContract

	// GetRegex returns the regular expression that the parameter matches.
	//
	// Go's regular expressions are RE2, so the pattern carries no delimiter, it
	// spells a named group `(?P<name>…)`, and it holds no lookahead, no
	// lookbehind, and no backreference.
	GetRegex() string

	// WithRegex returns a copy of the parameter for another regular
	// expression.
	WithRegex(regex string) ParameterContract

	// HasCast reports whether the parameter casts its value to a type.
	HasCast() bool

	// GetCast returns the type that the parameter casts its value to.
	GetCast() any

	// WithCast returns a copy of the parameter for another cast.
	WithCast(cast any) ParameterContract

	// IsOptional reports whether the path matches without the parameter.
	IsOptional() bool

	// WithIsOptional returns a copy of the parameter with another optional
	// flag.
	WithIsOptional(isOptional bool) ParameterContract

	// ShouldCapture reports whether the router passes the parameter to the
	// handler.
	ShouldCapture() bool

	// WithShouldCapture returns a copy of the parameter with another capture
	// flag.
	WithShouldCapture(shouldCapture bool) ParameterContract

	// GetDefault returns the value that the router uses where the path carries
	// none.
	GetDefault() any

	// WithDefault returns a copy of the parameter for another default value.
	WithDefault(defaultValue any) ParameterContract

	// GetValue returns the value that the path carried.
	GetValue() any

	// WithValue returns a copy of the parameter for another value.
	WithValue(value any) ParameterContract
}

// RouteContract is one route of the HTTP router.
//
// Each middleware is named by its binding key. The other ports pass a class
// reference — PHP a `class-string`, Java a `Class<?>`, TypeScript a constructor
// reference — and Go has none of those, so it passes the same string constant
// that it binds the middleware under.
//
// Warning: a route appends each middleware and never dedupes. A middleware that
// is registered twice runs twice. That is the developer's error, and neither the
// route nor `sindri` corrects it.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type RouteContract interface {
	// GetPath returns the path that the route matches.
	GetPath() string

	// WithPath returns a copy of the route for another path.
	WithPath(path string) RouteContract

	// WithAddedPath returns a copy of the route with the path added to the one
	// it holds.
	WithAddedPath(path string) RouteContract

	// GetName returns the name of the route.
	GetName() string

	// WithName returns a copy of the route under another name.
	WithName(name string) RouteContract

	// WithAddedName returns a copy of the route with the name added to the one
	// it holds.
	WithAddedName(name string) RouteContract

	// GetHandler returns what the route runs.
	GetHandler() HttpHandlerFunc

	// WithHandler returns a copy of the route that runs another handler.
	WithHandler(handler HttpHandlerFunc) RouteContract

	// GetRequestMethods returns each request method that the route matches.
	GetRequestMethods() []constant.RequestMethod

	// HasRequestMethod reports whether the route matches the request method.
	HasRequestMethod(requestMethod constant.RequestMethod) bool

	// WithRequestMethods returns a copy of the route for other request
	// methods.
	WithRequestMethods(requestMethods ...constant.RequestMethod) RouteContract

	// WithAddedRequestMethods returns a copy of the route with the request
	// methods added to the ones it holds.
	WithAddedRequestMethods(requestMethods ...constant.RequestMethod) RouteContract

	// GetRouteMatchedMiddleware returns the binding key of each route-matched
	// middleware.
	GetRouteMatchedMiddleware() []string

	// WithRouteMatchedMiddleware returns a copy of the route with other
	// route-matched middleware.
	WithRouteMatchedMiddleware(middleware ...string) RouteContract

	// WithAddedRouteMatchedMiddleware returns a copy of the route with the
	// route-matched middleware appended.
	WithAddedRouteMatchedMiddleware(middleware ...string) RouteContract

	// GetRouteDispatchedMiddleware returns the binding key of each
	// route-dispatched middleware.
	GetRouteDispatchedMiddleware() []string

	// WithRouteDispatchedMiddleware returns a copy of the route with other
	// route-dispatched middleware.
	WithRouteDispatchedMiddleware(middleware ...string) RouteContract

	// WithAddedRouteDispatchedMiddleware returns a copy of the route with the
	// route-dispatched middleware appended.
	WithAddedRouteDispatchedMiddleware(middleware ...string) RouteContract

	// GetThrowableCaughtMiddleware returns the binding key of each
	// throwable-caught middleware.
	GetThrowableCaughtMiddleware() []string

	// WithThrowableCaughtMiddleware returns a copy of the route with other
	// throwable-caught middleware.
	WithThrowableCaughtMiddleware(middleware ...string) RouteContract

	// WithAddedThrowableCaughtMiddleware returns a copy of the route with the
	// throwable-caught middleware appended.
	WithAddedThrowableCaughtMiddleware(middleware ...string) RouteContract

	// GetSendingResponseMiddleware returns the binding key of each
	// sending-response middleware.
	GetSendingResponseMiddleware() []string

	// WithSendingResponseMiddleware returns a copy of the route with other
	// sending-response middleware.
	WithSendingResponseMiddleware(middleware ...string) RouteContract

	// WithAddedSendingResponseMiddleware returns a copy of the route with the
	// sending-response middleware appended.
	WithAddedSendingResponseMiddleware(middleware ...string) RouteContract

	// GetResponseSentMiddleware returns the binding key of each response-sent
	// middleware.
	GetResponseSentMiddleware() []string

	// WithResponseSentMiddleware returns a copy of the route with other
	// response-sent middleware.
	WithResponseSentMiddleware(middleware ...string) RouteContract

	// WithAddedResponseSentMiddleware returns a copy of the route with the
	// response-sent middleware appended.
	WithAddedResponseSentMiddleware(middleware ...string) RouteContract

	// HasRequestStruct reports whether the route names a request struct.
	HasRequestStruct() bool

	// GetRequestStruct returns the request struct of the route.
	GetRequestStruct() RequestStructContract

	// WithRequestStruct returns a copy of the route for another request
	// struct.
	WithRequestStruct(requestStruct RequestStructContract) RouteContract

	// HasResponseStruct reports whether the route names a response struct.
	HasResponseStruct() bool

	// GetResponseStruct returns the response struct of the route.
	GetResponseStruct() ResponseStructContract

	// WithResponseStruct returns a copy of the route for another response
	// struct.
	WithResponseStruct(responseStruct ResponseStructContract) RouteContract
}

// DynamicRouteContract is a route whose path carries a parameter.
type DynamicRouteContract interface {
	RouteContract

	// GetRegex returns the regular expression that the whole path matches.
	GetRegex() string

	// WithRegex returns a copy of the route for another regular expression.
	WithRegex(regex string) DynamicRouteContract

	// GetParameters returns each parameter of the path.
	GetParameters() []ParameterContract

	// WithParameters returns a copy of the route with other parameters.
	WithParameters(parameters ...ParameterContract) DynamicRouteContract

	// WithAddedParameters returns a copy of the route with the parameters
	// appended.
	WithAddedParameters(parameters ...ParameterContract) DynamicRouteContract
}

// HttpRouteProviderContract registers the HTTP routes of one component.
//
// A provider returns a literal slice, and never a computed one, because `sindri`
// reads the slice from the source rather than by running it.
//
// The other ports also read a route from an annotation, and declare
// `getControllerClasses` for the scanner to read. Go has no annotation, so that
// method is absent and a route is always registered here.
type HttpRouteProviderContract interface {
	// GetRoutes returns each route that the component registers.
	GetRoutes() []RouteContract
}
