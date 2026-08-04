/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

// RequestReceivedResultContract is what a request-received middleware returns:
// either the request that the next middleware receives, or a response that ends
// the run.
//
// The TypeScript port returns `ServerRequestContract | ResponseContract`, and Go
// has no union. The Java port answers with a `RequestReceivedResult` record. Go
// cannot put that struct in this package, because the taxonomy holds a `contract`
// segment to interfaces, and cannot put it in `data`, because `data` imports this
// package. The result is therefore a contract of its own, with a struct behind it.
type RequestReceivedResultContract interface {
	// GetRequest returns the request that the next middleware receives.
	GetRequest() ServerRequestContract

	// GetResponse returns the response that ends the run, and nil where the run
	// continues.
	GetResponse() ResponseContract

	// IsResponse reports whether the result ends the run with a response.
	IsResponse() bool
}

// RouteMatchedResultContract is what a route-matched middleware returns: either
// the route that the next middleware receives, or a response that ends the run.
type RouteMatchedResultContract interface {
	// GetRoute returns the route that the next middleware receives.
	GetRoute() RouteContract

	// GetResponse returns the response that ends the run, and nil where the run
	// continues.
	GetResponse() ResponseContract

	// IsResponse reports whether the result ends the run with a response.
	IsResponse() bool
}

// RequestReceivedMiddlewareContract runs when the server receives a request,
// before the router matches a route.
type RequestReceivedMiddlewareContract interface {
	// RequestReceived runs the middleware.
	RequestReceived(
		request ServerRequestContract,
		handler RequestReceivedHandlerContract,
	) RequestReceivedResultContract
}

// RouteMatchedMiddlewareContract runs when the router matches a route.
type RouteMatchedMiddlewareContract interface {
	// RouteMatched runs the middleware.
	RouteMatched(
		request ServerRequestContract,
		route RouteContract,
		handler RouteMatchedHandlerContract,
	) RouteMatchedResultContract
}

// RouteNotMatchedMiddlewareContract runs when the router matches no route.
type RouteNotMatchedMiddlewareContract interface {
	// RouteNotMatched runs the middleware.
	RouteNotMatched(
		request ServerRequestContract,
		response ResponseContract,
		handler RouteNotMatchedHandlerContract,
	) ResponseContract
}

// RouteDispatchedMiddlewareContract runs after the route handler returns.
type RouteDispatchedMiddlewareContract interface {
	// RouteDispatched runs the middleware.
	RouteDispatched(
		request ServerRequestContract,
		response ResponseContract,
		route RouteContract,
		handler RouteDispatchedHandlerContract,
	) ResponseContract
}

// ThrowableCaughtMiddlewareContract runs when something in the run reports a
// failure.
type ThrowableCaughtMiddlewareContract interface {
	// ThrowableCaught runs the middleware.
	ThrowableCaught(
		request ServerRequestContract,
		response ResponseContract,
		throwable error,
		handler ThrowableCaughtHandlerContract,
	) ResponseContract
}

// SendingResponseMiddlewareContract runs before the server sends the response.
type SendingResponseMiddlewareContract interface {
	// SendingResponse runs the middleware.
	SendingResponse(
		request ServerRequestContract,
		response ResponseContract,
		handler SendingResponseHandlerContract,
	) ResponseContract
}

// ResponseSentMiddlewareContract runs after the server sent the response.
type ResponseSentMiddlewareContract interface {
	// ResponseSent runs the middleware.
	ResponseSent(
		request ServerRequestContract,
		response ResponseContract,
		handler ResponseSentHandlerContract,
	)
}

// HandlerContract holds the middleware of one stage and runs them in order.
//
// The other ports type the added middleware by the stage's own middleware
// contract, and pass a constructor reference. Go has no constructor reference,
// so a middleware is added by its binding key, which is what the framework uses
// as a class reference everywhere.
//
// Warning: a handler appends each middleware and never dedupes. A middleware
// that is added twice runs twice. That is the developer's error, and the
// framework does not correct it, because the generated cache must match what the
// runtime collects.
type HandlerContract interface {
	// Add appends each middleware, by binding key, after the ones the handler
	// holds.
	Add(middleware ...string)
}

// RequestReceivedHandlerContract runs the request-received middleware.
type RequestReceivedHandlerContract interface {
	HandlerContract

	// RequestReceived runs each middleware that the handler holds.
	RequestReceived(request ServerRequestContract) RequestReceivedResultContract
}

// RouteMatchedHandlerContract runs the route-matched middleware.
type RouteMatchedHandlerContract interface {
	HandlerContract

	// RouteMatched runs each middleware that the handler holds.
	RouteMatched(request ServerRequestContract, route RouteContract) RouteMatchedResultContract
}

// RouteNotMatchedHandlerContract runs the route-not-matched middleware.
type RouteNotMatchedHandlerContract interface {
	HandlerContract

	// RouteNotMatched runs each middleware that the handler holds.
	RouteNotMatched(request ServerRequestContract, response ResponseContract) ResponseContract
}

// RouteDispatchedHandlerContract runs the route-dispatched middleware.
type RouteDispatchedHandlerContract interface {
	HandlerContract

	// RouteDispatched runs each middleware that the handler holds.
	RouteDispatched(
		request ServerRequestContract,
		response ResponseContract,
		route RouteContract,
	) ResponseContract
}

// ThrowableCaughtHandlerContract runs the throwable-caught middleware.
type ThrowableCaughtHandlerContract interface {
	HandlerContract

	// ThrowableCaught runs each middleware that the handler holds.
	ThrowableCaught(
		request ServerRequestContract,
		response ResponseContract,
		throwable error,
	) ResponseContract
}

// SendingResponseHandlerContract runs the sending-response middleware.
type SendingResponseHandlerContract interface {
	HandlerContract

	// SendingResponse runs each middleware that the handler holds.
	SendingResponse(request ServerRequestContract, response ResponseContract) ResponseContract
}

// ResponseSentHandlerContract runs the response-sent middleware.
type ResponseSentHandlerContract interface {
	HandlerContract

	// ResponseSent runs each middleware that the handler holds.
	ResponseSent(request ServerRequestContract, response ResponseContract)
}
