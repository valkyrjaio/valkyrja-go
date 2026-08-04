/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package fixtures holds the reusable doubles that the HTTP middleware tests
// build on.
package fixtures

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/middleware/data"
)

// RecordingMiddlewareFixture is a middleware of every stage. It records that it
// ran, and it passes the run to the handler that it receives.
//
// One fixture serves each stage, because a test asserts the same two things at
// every stage: the middleware ran, and the handler continued.
type RecordingMiddlewareFixture struct {
	Name   string
	Record *[]string
}

// RequestReceived records the run and passes it on.
func (m *RecordingMiddlewareFixture) RequestReceived(
	request contract.ServerRequestContract,
	handler contract.RequestReceivedHandlerContract,
) contract.RequestReceivedResultContract {
	m.record()

	return handler.RequestReceived(request)
}

// RouteMatched records the run and passes it on.
func (m *RecordingMiddlewareFixture) RouteMatched(
	request contract.ServerRequestContract,
	route contract.RouteContract,
	handler contract.RouteMatchedHandlerContract,
) contract.RouteMatchedResultContract {
	m.record()

	return handler.RouteMatched(request, route)
}

// RouteNotMatched records the run and passes it on.
func (m *RecordingMiddlewareFixture) RouteNotMatched(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
	handler contract.RouteNotMatchedHandlerContract,
) contract.ResponseContract {
	m.record()

	return handler.RouteNotMatched(request, response)
}

// RouteDispatched records the run and passes it on.
func (m *RecordingMiddlewareFixture) RouteDispatched(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
	route contract.RouteContract,
	handler contract.RouteDispatchedHandlerContract,
) contract.ResponseContract {
	m.record()

	return handler.RouteDispatched(request, response, route)
}

// ThrowableCaught records the run and passes it on.
func (m *RecordingMiddlewareFixture) ThrowableCaught(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
	throwable error,
	handler contract.ThrowableCaughtHandlerContract,
) contract.ResponseContract {
	m.record()

	return handler.ThrowableCaught(request, response, throwable)
}

// SendingResponse records the run and passes it on.
func (m *RecordingMiddlewareFixture) SendingResponse(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
	handler contract.SendingResponseHandlerContract,
) contract.ResponseContract {
	m.record()

	return handler.SendingResponse(request, response)
}

// ResponseSent records the run and passes it on.
func (m *RecordingMiddlewareFixture) ResponseSent(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
	handler contract.ResponseSentHandlerContract,
) {
	m.record()

	handler.ResponseSent(request, response)
}

// record notes that the middleware ran.
func (m *RecordingMiddlewareFixture) record() {
	*m.Record = append(*m.Record, m.Name)
}

// EndingMiddlewareFixture is a middleware that ends the run with the response it
// holds. It never calls the handler, so no middleware after it runs.
type EndingMiddlewareFixture struct {
	Response contract.ResponseContract
	Record   *[]string
}

// RequestReceived ends the run with the response.
func (m *EndingMiddlewareFixture) RequestReceived(
	_ contract.ServerRequestContract,
	_ contract.RequestReceivedHandlerContract,
) contract.RequestReceivedResultContract {
	*m.Record = append(*m.Record, "ending")

	return data.NewRequestReceivedResponse(m.Response)
}

// RouteMatched ends the run with the response.
func (m *EndingMiddlewareFixture) RouteMatched(
	_ contract.ServerRequestContract,
	_ contract.RouteContract,
	_ contract.RouteMatchedHandlerContract,
) contract.RouteMatchedResultContract {
	*m.Record = append(*m.Record, "ending")

	return data.NewRouteMatchedResponse(m.Response)
}

// RouteNotMatched ends the run with the response.
func (m *EndingMiddlewareFixture) RouteNotMatched(
	_ contract.ServerRequestContract,
	_ contract.ResponseContract,
	_ contract.RouteNotMatchedHandlerContract,
) contract.ResponseContract {
	*m.Record = append(*m.Record, "ending")

	return m.Response
}
