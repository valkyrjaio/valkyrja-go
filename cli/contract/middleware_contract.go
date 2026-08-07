/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

type InputReceivedResultContract interface {
	// GetInput returns the input that the next middleware receives.
	GetInput() InputContract

	// GetOutput returns the output that ends the run, and nil where the run
	// continues.
	GetOutput() OutputContract

	// IsOutput reports whether the result ends the run with an output.
	IsOutput() bool
}

type RouteMatchedResultContract interface {
	// GetRoute returns the route that the next middleware receives.
	GetRoute() RouteContract

	// GetOutput returns the output that ends the run, and nil where the run
	// continues.
	GetOutput() OutputContract

	// IsOutput reports whether the result ends the run with an output.
	IsOutput() bool
}

type InputReceivedMiddlewareContract interface {
	// InputReceived runs the middleware.
	InputReceived(input InputContract, handler InputReceivedHandlerContract) InputReceivedResultContract
}

type RouteMatchedMiddlewareContract interface {
	// RouteMatched runs the middleware.
	RouteMatched(
		input InputContract,
		route RouteContract,
		handler RouteMatchedHandlerContract,
	) RouteMatchedResultContract
}

type RouteNotMatchedMiddlewareContract interface {
	// RouteNotMatched runs the middleware.
	RouteNotMatched(
		input InputContract,
		output OutputContract,
		handler RouteNotMatchedHandlerContract,
	) OutputContract
}

type RouteDispatchedMiddlewareContract interface {
	// RouteDispatched runs the middleware.
	RouteDispatched(
		input InputContract,
		output OutputContract,
		route RouteContract,
		handler RouteDispatchedHandlerContract,
	) OutputContract
}

type ThrowableCaughtMiddlewareContract interface {
	// ThrowableCaught runs the middleware.
	ThrowableCaught(
		input InputContract,
		output OutputContract,
		throwable error,
		handler ThrowableCaughtHandlerContract,
	) OutputContract
}

type ProcessExitingMiddlewareContract interface {
	// ProcessExiting runs the middleware.
	ProcessExiting(input InputContract, output OutputContract, handler ProcessExitingHandlerContract)
}

type HandlerContract interface {
	// Add appends each middleware, by binding key, after the ones the handler
	// holds.
	Add(middleware ...string)
}

type InputReceivedHandlerContract interface {
	HandlerContract

	// InputReceived runs each middleware that the handler holds.
	InputReceived(input InputContract) InputReceivedResultContract
}

type RouteMatchedHandlerContract interface {
	HandlerContract

	// RouteMatched runs each middleware that the handler holds.
	RouteMatched(input InputContract, route RouteContract) RouteMatchedResultContract
}

type RouteNotMatchedHandlerContract interface {
	HandlerContract

	// RouteNotMatched runs each middleware that the handler holds.
	RouteNotMatched(input InputContract, output OutputContract) OutputContract
}

type RouteDispatchedHandlerContract interface {
	HandlerContract

	// RouteDispatched runs each middleware that the handler holds.
	RouteDispatched(input InputContract, output OutputContract, route RouteContract) OutputContract
}

type ThrowableCaughtHandlerContract interface {
	HandlerContract

	// ThrowableCaught runs each middleware that the handler holds.
	ThrowableCaught(input InputContract, output OutputContract, throwable error) OutputContract
}

type ProcessExitingHandlerContract interface {
	HandlerContract

	// ProcessExiting runs each middleware that the handler holds.
	ProcessExiting(input InputContract, output OutputContract)
}
