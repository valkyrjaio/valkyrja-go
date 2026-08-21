/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

type RequestReceivedResultContract interface {
	// GetRequest returns the request that the next middleware receives.
	GetRequest() ServerRequestContract

	// GetResponse returns the response that ends the run, and nil where the run
	// continues.
	GetResponse() ResponseContract

	// IsResponse reports whether the result ends the run with a response.
	IsResponse() bool
}

type RouteMatchedResultContract interface {
	// GetRoute returns the route that the next middleware receives.
	GetRoute() RouteContract

	// GetResponse returns the response that ends the run, and nil where the run
	// continues.
	GetResponse() ResponseContract

	// IsResponse reports whether the result ends the run with a response.
	IsResponse() bool
}

type RequestReceivedMiddlewareContract interface {
	// RequestReceived runs the middleware.
	RequestReceived(
		request ServerRequestContract,
		handler RequestReceivedHandlerContract,
	) RequestReceivedResultContract
}

type RouteMatchedMiddlewareContract interface {
	// RouteMatched runs the middleware.
	RouteMatched(
		request ServerRequestContract,
		route RouteContract,
		handler RouteMatchedHandlerContract,
	) RouteMatchedResultContract
}

type RouteNotMatchedMiddlewareContract interface {
	// RouteNotMatched runs the middleware.
	RouteNotMatched(
		request ServerRequestContract,
		response ResponseContract,
		handler RouteNotMatchedHandlerContract,
	) ResponseContract
}

type RouteDispatchedMiddlewareContract interface {
	// RouteDispatched runs the middleware.
	RouteDispatched(
		request ServerRequestContract,
		response ResponseContract,
		route RouteContract,
		handler RouteDispatchedHandlerContract,
	) ResponseContract
}

type ThrowableCaughtMiddlewareContract interface {
	// ThrowableCaught runs the middleware.
	ThrowableCaught(
		request ServerRequestContract,
		response ResponseContract,
		throwable error,
		handler ThrowableCaughtHandlerContract,
	) ResponseContract
}

type SendingResponseMiddlewareContract interface {
	// SendingResponse runs the middleware.
	SendingResponse(
		request ServerRequestContract,
		response ResponseContract,
		handler SendingResponseHandlerContract,
	) ResponseContract
}

type ResponseSentMiddlewareContract interface {
	// ResponseSent runs the middleware.
	ResponseSent(
		request ServerRequestContract,
		response ResponseContract,
		handler ResponseSentHandlerContract,
	)
}

type HandlerContract interface {
	// Add appends each middleware, by binding key, after the ones the handler
	// holds.
	Add(middleware ...string)
}

type RequestReceivedHandlerContract interface {
	HandlerContract

	// RequestReceived runs each middleware that the handler holds.
	RequestReceived(request ServerRequestContract) RequestReceivedResultContract
}

type RouteMatchedHandlerContract interface {
	HandlerContract

	// RouteMatched runs each middleware that the handler holds.
	RouteMatched(request ServerRequestContract, route RouteContract) RouteMatchedResultContract
}

type RouteNotMatchedHandlerContract interface {
	HandlerContract

	// RouteNotMatched runs each middleware that the handler holds.
	RouteNotMatched(request ServerRequestContract, response ResponseContract) ResponseContract
}

type RouteDispatchedHandlerContract interface {
	HandlerContract

	// RouteDispatched runs each middleware that the handler holds.
	RouteDispatched(
		request ServerRequestContract,
		response ResponseContract,
		route RouteContract,
	) ResponseContract
}

type ThrowableCaughtHandlerContract interface {
	HandlerContract

	// ThrowableCaught runs each middleware that the handler holds.
	ThrowableCaught(
		request ServerRequestContract,
		response ResponseContract,
		throwable error,
	) ResponseContract
}

type SendingResponseHandlerContract interface {
	HandlerContract

	// SendingResponse runs each middleware that the handler holds.
	SendingResponse(request ServerRequestContract, response ResponseContract) ResponseContract
}

type ResponseSentHandlerContract interface {
	HandlerContract

	// ResponseSent runs each middleware that the handler holds.
	ResponseSent(request ServerRequestContract, response ResponseContract)
}
