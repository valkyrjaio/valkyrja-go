/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package handler_test

import (
	"errors"
	"slices"
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/middleware/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/http/middleware/handler"
)

const (
	firstID  = "Valkyrja.Tests.Http.FirstMiddleware"
	secondID = "Valkyrja.Tests.Http.SecondMiddleware"
	endingID = "Valkyrja.Tests.Http.EndingMiddleware"
)

// newContainer builds a container that resolves a recording middleware under
// each binding key, and an ending middleware under its own.
func newContainer(record *[]string, ending contract.ResponseContract) containercontract.ContainerContract {
	container := manager.NewContainer(nil)

	for _, id := range []string{firstID, secondID} {
		container.Bind(id, func(_ containercontract.ContainerContract, _ []any) any {
			return &fixtures.RecordingMiddlewareFixture{Name: id, Record: record}
		})
	}

	container.Bind(endingID, func(_ containercontract.ContainerContract, _ []any) any {
		return &fixtures.EndingMiddlewareFixture{Response: ending, Record: record}
	})

	return container
}

func TestEachHandlerRunsItsMiddlewareInOrder(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	built := request.NewServerRequest(nil, "", nil, nil)
	sent := response.NewResponse(nil, 0, nil)

	handler.NewRequestReceivedHandler(container, firstID, secondID).RequestReceived(built)

	if !slices.Equal(record, []string{firstID, secondID}) {
		t.Errorf("the handler must run each middleware in order, but ran: %v", record)
	}

	record = record[:0]

	handler.NewRouteMatchedHandler(container, firstID, secondID).RouteMatched(built, nil)

	if !slices.Equal(record, []string{firstID, secondID}) {
		t.Errorf("the route-matched handler must run each middleware, but ran: %v", record)
	}

	record = record[:0]

	handler.NewRouteNotMatchedHandler(container, firstID).RouteNotMatched(built, sent)
	handler.NewRouteDispatchedHandler(container, firstID).RouteDispatched(built, sent, nil)
	handler.NewThrowableCaughtHandler(container, firstID).ThrowableCaught(built, sent, errors.New("failed"))
	handler.NewSendingResponseHandler(container, firstID).SendingResponse(built, sent)
	handler.NewResponseSentHandler(container, firstID).ResponseSent(built, sent)

	if len(record) != 5 {
		t.Errorf("each remaining handler must run its middleware, but the record is: %v", record)
	}
}

func TestAHandlerWithNoMiddlewareReturnsWhatItReceived(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	built := request.NewServerRequest(nil, "", nil, nil)
	sent := response.NewResponse(nil, 0, nil)

	result := handler.NewRequestReceivedHandler(container).RequestReceived(built)

	if result.IsResponse() {
		t.Error("a handler with no middleware must return the request, but returned a response")
	}

	if result.GetRequest() != contract.ServerRequestContract(built) {
		t.Error("a handler with no middleware must return the request that it received, but did not")
	}

	matched := handler.NewRouteMatchedHandler(container).RouteMatched(built, nil)

	if matched.IsResponse() {
		t.Error("a handler with no middleware must return the route, but returned a response")
	}

	if handler.NewRouteNotMatchedHandler(container).RouteNotMatched(built, sent) != contract.ResponseContract(sent) {
		t.Error("a handler with no middleware must return the response, but did not")
	}

	if len(record) != 0 {
		t.Errorf("a handler with no middleware must run nothing, but ran: %v", record)
	}
}

func TestAMiddlewareThatEndsTheRunStopsTheOnesAfterIt(t *testing.T) {
	t.Parallel()

	record := []string{}
	ending := response.NewResponse(nil, 0, nil)
	container := newContainer(&record, ending)

	built := request.NewServerRequest(nil, "", nil, nil)

	result := handler.NewRequestReceivedHandler(container, endingID, firstID).RequestReceived(built)

	if !result.IsResponse() {
		t.Fatal("the result must end the run with a response, but did not")
	}

	if result.GetResponse() != contract.ResponseContract(ending) {
		t.Error("the result must carry the response that the middleware returned, but did not")
	}

	if !slices.Equal(record, []string{"ending"}) {
		t.Errorf("no middleware after the one that ended the run may run, but the record is: %v", record)
	}
}

func TestARouteMatchedMiddlewareThatEndsTheRunReturnsAResponse(t *testing.T) {
	t.Parallel()

	record := []string{}
	ending := response.NewResponse(nil, 0, nil)
	container := newContainer(&record, ending)

	result := handler.NewRouteMatchedHandler(container, endingID).
		RouteMatched(request.NewServerRequest(nil, "", nil, nil), nil)

	if !result.IsResponse() || result.GetResponse() != contract.ResponseContract(ending) {
		t.Error("the result must end the run with the response, but did not")
	}

	if result.GetRoute() != nil {
		t.Error("the result must carry no route where it ends the run, but carried one")
	}
}

func TestAddAppendsMiddlewareAfterTheOnesTheHandlerHolds(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	built := handler.NewRequestReceivedHandler(container, firstID)
	built.Add(secondID)

	built.RequestReceived(request.NewServerRequest(nil, "", nil, nil))

	if !slices.Equal(record, []string{firstID, secondID}) {
		t.Errorf("Add must append the middleware, but the record is: %v", record)
	}
}

func TestAddNeverDedupesAMiddleware(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	built := handler.NewRequestReceivedHandler(container, firstID)
	built.Add(firstID)

	built.RequestReceived(request.NewServerRequest(nil, "", nil, nil))

	if !slices.Equal(record, []string{firstID, firstID}) {
		t.Errorf("a middleware that is added twice must run twice, but the record is: %v", record)
	}
}

func TestAHandlerSkipsAMiddlewareThatTheContainerCannotResolve(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	built := request.NewServerRequest(nil, "", nil, nil)

	result := handler.NewRequestReceivedHandler(container, "Valkyrja.Tests.Http.Unbound").RequestReceived(built)

	if result.IsResponse() {
		t.Error("a handler must skip a middleware that the container cannot resolve, but ended the run")
	}

	if len(record) != 0 {
		t.Errorf("no middleware must run, but the record is: %v", record)
	}
}

func TestAHandlerSkipsABindingThatIsNoMiddleware(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)
	container.Bind("Valkyrja.Tests.Http.NotMiddleware", func(_ containercontract.ContainerContract, _ []any) any {
		return "not a middleware"
	})

	result := handler.NewRequestReceivedHandler(container, "Valkyrja.Tests.Http.NotMiddleware").
		RequestReceived(request.NewServerRequest(nil, "", nil, nil))

	if result.IsResponse() {
		t.Error("a handler must skip a binding that is no middleware, but ended the run")
	}
}

func TestEachHandlerSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	handlers := []contract.HandlerContract{
		handler.NewRequestReceivedHandler(container),
		handler.NewRouteMatchedHandler(container),
		handler.NewRouteNotMatchedHandler(container),
		handler.NewRouteDispatchedHandler(container),
		handler.NewThrowableCaughtHandler(container),
		handler.NewSendingResponseHandler(container),
		handler.NewResponseSentHandler(container),
	}

	for _, built := range handlers {
		built.Add(firstID)
	}
}
