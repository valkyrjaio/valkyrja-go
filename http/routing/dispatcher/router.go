/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package dispatcher runs the route that a request matches.
package dispatcher

import (
	"net/url"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
)

// Router runs the route that a request matches.
type Router struct {
	container       containercontract.ContainerContract
	matcher         contract.MatcherContract
	responseFactory contract.ResponseFactoryContract

	routeMatchedHandler    contract.RouteMatchedHandlerContract
	routeNotMatchedHandler contract.RouteNotMatchedHandlerContract
	routeDispatchedHandler contract.RouteDispatchedHandlerContract
}

// NewRouter builds a router over a container, a matcher, a response factory, and
// the middleware handler of each stage that the router runs.
func NewRouter(
	container containercontract.ContainerContract,
	matcher contract.MatcherContract,
	responseFactory contract.ResponseFactoryContract,
	routeMatchedHandler contract.RouteMatchedHandlerContract,
	routeNotMatchedHandler contract.RouteNotMatchedHandlerContract,
	routeDispatchedHandler contract.RouteDispatchedHandlerContract,
) *Router {
	return &Router{
		container:              container,
		matcher:                matcher,
		responseFactory:        responseFactory,
		routeMatchedHandler:    routeMatchedHandler,
		routeNotMatchedHandler: routeNotMatchedHandler,
		routeDispatchedHandler: routeDispatchedHandler,
	}
}

// Dispatch matches the request to a route and runs it.
//
// A path that matches no route of the request method, but matches one of another
// method, reports 405. A path that matches no route at all reports 404. Either
// way the route-not-matched middleware runs, so an application shapes the
// response.
func (r *Router) Dispatch(request contract.ServerRequestContract) contract.ResponseContract {
	route := r.matchRoute(request)
	if route == nil {
		return r.routeNotMatchedHandler.RouteNotMatched(request, r.getNotMatchedResponse(request))
	}

	return r.DispatchRoute(request, route)
}

// DispatchRoute runs the route for the request.
//
// The route-matched middleware runs first, and it ends the run where it returns
// a response. The router then binds the route in the container, so a handler
// reads the route that matched, and runs the handler.
func (r *Router) DispatchRoute(
	request contract.ServerRequestContract,
	route contract.RouteContract,
) contract.ResponseContract {
	result := r.routeMatchedHandler.RouteMatched(request, route)
	if result.IsResponse() {
		return result.GetResponse()
	}

	matched := result.GetRoute()

	r.container.SetSingleton(routingconstant.RouteContractServiceID, matched)

	response := r.runHandler(matched)

	return r.routeDispatchedHandler.RouteDispatched(request, response, matched)
}

// runHandler runs what the route holds, and returns an empty response where the
// route holds no handler.
//
// The other ports type the handler as callable, so the value cannot be absent. A
// Go function value can be nil, and calling it panics in the middle of a request.
func (r *Router) runHandler(route contract.RouteContract) contract.ResponseContract {
	handler := route.GetHandler()
	if handler == nil {
		return r.responseFactory.CreateResponse("", constant.StatusCodeNoContent, nil)
	}

	return handler(r.container, route)
}

// matchRoute returns the route that the request matches, and nil where none
// matches.
//
// The path is decoded first, so a route matches a path that the client escaped.
// A path that no decoder reads is matched as it arrived.
func (r *Router) matchRoute(request contract.ServerRequestContract) contract.RouteContract {
	return r.matcher.Match(decodePath(request.GetUri().GetPath()), request.GetMethod())
}

// getNotMatchedResponse returns the response for a request that matches no
// route.
//
// A path that matches a route of another request method reports 405 rather than
// 404, because the path is there and the method is what the route does not take.
//
// Warning: the collection files a route under each method that it takes, and a
// route that takes the any method files under every one of them rather than
// under `RequestMethodAny` itself. Asking the matcher for `RequestMethodAny`
// therefore matches nothing at all, and this walks the methods instead.
func (r *Router) getNotMatchedResponse(request contract.ServerRequestContract) contract.ResponseContract {
	path := decodePath(request.GetUri().GetPath())

	if r.matchesAnotherMethod(path, request.GetMethod()) {
		return r.responseFactory.CreateResponse("", constant.StatusCodeMethodNotAllowed, nil)
	}

	return r.responseFactory.CreateResponse("", constant.StatusCodeNotFound, nil)
}

// matchesAnotherMethod reports whether the path matches a route of a request
// method other than the one that the client used.
func (r *Router) matchesAnotherMethod(path string, method constant.RequestMethod) bool {
	for _, other := range constant.GetAllRequestMethods() {
		if other == method {
			continue
		}

		if r.matcher.Match(path, other) != nil {
			return true
		}
	}

	return false
}

// decodePath returns the path with each percent-encoded triplet decoded, and the
// path as it arrived where no decoder reads it.
func decodePath(path string) string {
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return path
	}

	return decoded
}
