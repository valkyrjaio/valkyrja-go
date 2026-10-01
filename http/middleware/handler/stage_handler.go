/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package handler

import (
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/middleware/data"
)

type RequestReceivedHandler struct {
	Handler
}

// NewRequestReceivedHandler builds the handler.
func NewRequestReceivedHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *RequestReceivedHandler {
	return &RequestReceivedHandler{Handler: NewHandler(container, middleware...)}
}

// RequestReceived runs the next middleware, and returns the request where the
// handler holds none.
func (h *RequestReceivedHandler) RequestReceived(
	request contract.ServerRequestContract,
) contract.RequestReceivedResultContract {
	middleware, isMiddleware := h.getNext().(contract.RequestReceivedMiddlewareContract)
	if !isMiddleware {
		return data.NewRequestReceivedRequest(request)
	}

	return middleware.RequestReceived(request, h)
}

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
	request contract.ServerRequestContract,
	route contract.RouteContract,
) contract.RouteMatchedResultContract {
	middleware, isMiddleware := h.getNext().(contract.RouteMatchedMiddlewareContract)
	if !isMiddleware {
		return data.NewRouteMatchedRoute(route)
	}

	return middleware.RouteMatched(request, route, h)
}

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

// RouteNotMatched runs the next middleware, and returns the response where the
// handler holds none.
func (h *RouteNotMatchedHandler) RouteNotMatched(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
) contract.ResponseContract {
	middleware, isMiddleware := h.getNext().(contract.RouteNotMatchedMiddlewareContract)
	if !isMiddleware {
		return response
	}

	return middleware.RouteNotMatched(request, response, h)
}

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

// RouteDispatched runs the next middleware, and returns the response where the
// handler holds none.
func (h *RouteDispatchedHandler) RouteDispatched(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
	route contract.RouteContract,
) contract.ResponseContract {
	middleware, isMiddleware := h.getNext().(contract.RouteDispatchedMiddlewareContract)
	if !isMiddleware {
		return response
	}

	return middleware.RouteDispatched(request, response, route, h)
}

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

// ThrowableCaught runs the next middleware, and returns the response where the
// handler holds none.
func (h *ThrowableCaughtHandler) ThrowableCaught(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
	throwable error,
) contract.ResponseContract {
	middleware, isMiddleware := h.getNext().(contract.ThrowableCaughtMiddlewareContract)
	if !isMiddleware {
		return response
	}

	return middleware.ThrowableCaught(request, response, throwable, h)
}

type SendingResponseHandler struct {
	Handler
}

// NewSendingResponseHandler builds the handler.
func NewSendingResponseHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *SendingResponseHandler {
	return &SendingResponseHandler{Handler: NewHandler(container, middleware...)}
}

// SendingResponse runs the next middleware, and returns the response where the
// handler holds none.
func (h *SendingResponseHandler) SendingResponse(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
) contract.ResponseContract {
	middleware, isMiddleware := h.getNext().(contract.SendingResponseMiddlewareContract)
	if !isMiddleware {
		return response
	}

	return middleware.SendingResponse(request, response, h)
}

type ResponseSentHandler struct {
	Handler
}

// NewResponseSentHandler builds the handler.
func NewResponseSentHandler(
	container containercontract.ContainerContract,
	middleware ...string,
) *ResponseSentHandler {
	return &ResponseSentHandler{Handler: NewHandler(container, middleware...)}
}

// ResponseSent runs the next middleware, and returns where the handler holds
// none.
func (h *ResponseSentHandler) ResponseSent(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
) {
	middleware, isMiddleware := h.getNext().(contract.ResponseSentMiddlewareContract)
	if !isMiddleware {
		return
	}

	middleware.ResponseSent(request, response, h)
}
