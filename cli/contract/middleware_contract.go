/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

// InputReceivedResultContract is what an input-received middleware returns:
// either the input that the next middleware receives, or an output that ends the
// run.
//
// The TypeScript port returns `InputContract | OutputContract`, and Go has no
// union. The result is a contract of its own for the reason that the HTTP
// component's result contracts are.
type InputReceivedResultContract interface {
	// GetInput returns the input that the next middleware receives.
	GetInput() InputContract

	// GetOutput returns the output that ends the run, and nil where the run
	// continues.
	GetOutput() OutputContract

	// IsOutput reports whether the result ends the run with an output.
	IsOutput() bool
}

// RouteMatchedResultContract is what a route-matched middleware returns: either
// the route that the next middleware receives, or an output that ends the run.
type RouteMatchedResultContract interface {
	// GetRoute returns the route that the next middleware receives.
	GetRoute() RouteContract

	// GetOutput returns the output that ends the run, and nil where the run
	// continues.
	GetOutput() OutputContract

	// IsOutput reports whether the result ends the run with an output.
	IsOutput() bool
}

// InputReceivedMiddlewareContract runs when the server receives an input, before
// the router matches a command.
type InputReceivedMiddlewareContract interface {
	// InputReceived runs the middleware.
	InputReceived(input InputContract, handler InputReceivedHandlerContract) InputReceivedResultContract
}

// RouteMatchedMiddlewareContract runs when the router matches a command.
type RouteMatchedMiddlewareContract interface {
	// RouteMatched runs the middleware.
	RouteMatched(
		input InputContract,
		route RouteContract,
		handler RouteMatchedHandlerContract,
	) RouteMatchedResultContract
}

// RouteNotMatchedMiddlewareContract runs when the router matches no command.
type RouteNotMatchedMiddlewareContract interface {
	// RouteNotMatched runs the middleware.
	RouteNotMatched(
		input InputContract,
		output OutputContract,
		handler RouteNotMatchedHandlerContract,
	) OutputContract
}

// RouteDispatchedMiddlewareContract runs after the command handler returns.
type RouteDispatchedMiddlewareContract interface {
	// RouteDispatched runs the middleware.
	RouteDispatched(
		input InputContract,
		output OutputContract,
		route RouteContract,
		handler RouteDispatchedHandlerContract,
	) OutputContract
}

// ThrowableCaughtMiddlewareContract runs when something in the run reports a
// failure.
type ThrowableCaughtMiddlewareContract interface {
	// ThrowableCaught runs the middleware.
	ThrowableCaught(
		input InputContract,
		output OutputContract,
		throwable error,
		handler ThrowableCaughtHandlerContract,
	) OutputContract
}

// ProcessExitingMiddlewareContract runs before the process exits.
type ProcessExitingMiddlewareContract interface {
	// ProcessExiting runs the middleware.
	ProcessExiting(input InputContract, output OutputContract, handler ProcessExitingHandlerContract)
}

// HandlerContract holds the middleware of one stage and runs them in order.
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

// InputReceivedHandlerContract runs the input-received middleware.
type InputReceivedHandlerContract interface {
	HandlerContract

	// InputReceived runs each middleware that the handler holds.
	InputReceived(input InputContract) InputReceivedResultContract
}

// RouteMatchedHandlerContract runs the route-matched middleware.
type RouteMatchedHandlerContract interface {
	HandlerContract

	// RouteMatched runs each middleware that the handler holds.
	RouteMatched(input InputContract, route RouteContract) RouteMatchedResultContract
}

// RouteNotMatchedHandlerContract runs the route-not-matched middleware.
type RouteNotMatchedHandlerContract interface {
	HandlerContract

	// RouteNotMatched runs each middleware that the handler holds.
	RouteNotMatched(input InputContract, output OutputContract) OutputContract
}

// RouteDispatchedHandlerContract runs the route-dispatched middleware.
type RouteDispatchedHandlerContract interface {
	HandlerContract

	// RouteDispatched runs each middleware that the handler holds.
	RouteDispatched(input InputContract, output OutputContract, route RouteContract) OutputContract
}

// ThrowableCaughtHandlerContract runs the throwable-caught middleware.
type ThrowableCaughtHandlerContract interface {
	HandlerContract

	// ThrowableCaught runs each middleware that the handler holds.
	ThrowableCaught(input InputContract, output OutputContract, throwable error) OutputContract
}

// ProcessExitingHandlerContract runs the process-exiting middleware.
type ProcessExitingHandlerContract interface {
	HandlerContract

	// ProcessExiting runs each middleware that the handler holds.
	ProcessExiting(input InputContract, output OutputContract)
}
