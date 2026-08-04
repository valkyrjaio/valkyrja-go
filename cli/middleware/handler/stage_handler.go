/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package handler

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/data"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// InputReceivedHandler runs the input-received middleware.
type InputReceivedHandler struct {
	Handler
}

// NewInputReceivedHandler builds the handler.
func NewInputReceivedHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *InputReceivedHandler {
	return &InputReceivedHandler{Handler: NewHandler(container, middleware...)}
}

// InputReceived runs the next middleware, and returns the input where the
// handler holds none.
func (h *InputReceivedHandler) InputReceived(input contract.InputContract) contract.InputReceivedResultContract {
	middleware, isMiddleware := h.getNext().(contract.InputReceivedMiddlewareContract)
	if !isMiddleware {
		return data.NewInputReceivedInput(input)
	}

	return middleware.InputReceived(input, h)
}

// RouteMatchedHandler runs the route-matched middleware.
type RouteMatchedHandler struct {
	Handler
}

// NewRouteMatchedHandler builds the handler.
func NewRouteMatchedHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *RouteMatchedHandler {
	return &RouteMatchedHandler{Handler: NewHandler(container, middleware...)}
}

// RouteMatched runs the next middleware, and returns the route where the handler
// holds none.
func (h *RouteMatchedHandler) RouteMatched(
	input contract.InputContract,
	route contract.RouteContract,
) contract.RouteMatchedResultContract {
	middleware, isMiddleware := h.getNext().(contract.RouteMatchedMiddlewareContract)
	if !isMiddleware {
		return data.NewRouteMatchedRoute(route)
	}

	return middleware.RouteMatched(input, route, h)
}

// RouteNotMatchedHandler runs the route-not-matched middleware.
type RouteNotMatchedHandler struct {
	Handler
}

// NewRouteNotMatchedHandler builds the handler.
func NewRouteNotMatchedHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *RouteNotMatchedHandler {
	return &RouteNotMatchedHandler{Handler: NewHandler(container, middleware...)}
}

// RouteNotMatched runs the next middleware, and returns the output where the
// handler holds none.
func (h *RouteNotMatchedHandler) RouteNotMatched(
	input contract.InputContract,
	output contract.OutputContract,
) contract.OutputContract {
	middleware, isMiddleware := h.getNext().(contract.RouteNotMatchedMiddlewareContract)
	if !isMiddleware {
		return output
	}

	return middleware.RouteNotMatched(input, output, h)
}

// RouteDispatchedHandler runs the route-dispatched middleware.
type RouteDispatchedHandler struct {
	Handler
}

// NewRouteDispatchedHandler builds the handler.
func NewRouteDispatchedHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *RouteDispatchedHandler {
	return &RouteDispatchedHandler{Handler: NewHandler(container, middleware...)}
}

// RouteDispatched runs the next middleware, and returns the output where the
// handler holds none.
func (h *RouteDispatchedHandler) RouteDispatched(
	input contract.InputContract,
	output contract.OutputContract,
	route contract.RouteContract,
) contract.OutputContract {
	middleware, isMiddleware := h.getNext().(contract.RouteDispatchedMiddlewareContract)
	if !isMiddleware {
		return output
	}

	return middleware.RouteDispatched(input, output, route, h)
}

// ThrowableCaughtHandler runs the throwable-caught middleware.
type ThrowableCaughtHandler struct {
	Handler
}

// NewThrowableCaughtHandler builds the handler.
func NewThrowableCaughtHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *ThrowableCaughtHandler {
	return &ThrowableCaughtHandler{Handler: NewHandler(container, middleware...)}
}

// ThrowableCaught runs the next middleware, and returns the output where the
// handler holds none.
func (h *ThrowableCaughtHandler) ThrowableCaught(
	input contract.InputContract,
	output contract.OutputContract,
	throwable error,
) contract.OutputContract {
	middleware, isMiddleware := h.getNext().(contract.ThrowableCaughtMiddlewareContract)
	if !isMiddleware {
		return output
	}

	return middleware.ThrowableCaught(input, output, throwable, h)
}

// ProcessExitingHandler runs the process-exiting middleware.
type ProcessExitingHandler struct {
	Handler
}

// NewProcessExitingHandler builds the handler.
func NewProcessExitingHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *ProcessExitingHandler {
	return &ProcessExitingHandler{Handler: NewHandler(container, middleware...)}
}

// ProcessExiting runs the next middleware, and returns where the handler holds
// none.
func (h *ProcessExitingHandler) ProcessExiting(input contract.InputContract, output contract.OutputContract) {
	middleware, isMiddleware := h.getNext().(contract.ProcessExitingMiddlewareContract)
	if !isMiddleware {
		return
	}

	middleware.ProcessExiting(input, output, h)
}
