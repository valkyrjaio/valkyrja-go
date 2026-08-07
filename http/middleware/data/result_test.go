/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/middleware/data"
)

func TestTheRequestReceivedResultCarriesTheRequest(t *testing.T) {
	t.Parallel()

	built := request.NewServerRequest(nil, "", nil, nil)

	result := data.NewRequestReceivedRequest(built)

	if result.IsResponse() {
		t.Error("a result that carries a request must not report a response, but did")
	}

	if result.GetRequest() != contract.ServerRequestContract(built) {
		t.Error("the result must carry the request, but did not")
	}

	if result.GetResponse() != nil {
		t.Error("the result must carry no response, but carried one")
	}
}

func TestTheRequestReceivedResultCarriesTheResponse(t *testing.T) {
	t.Parallel()

	sent := response.NewResponse(nil, 0, nil)

	result := data.NewRequestReceivedResponse(sent)

	if !result.IsResponse() {
		t.Error("a result that carries a response must report one, but did not")
	}

	if result.GetResponse() != contract.ResponseContract(sent) {
		t.Error("the result must carry the response, but did not")
	}

	if result.GetRequest() != nil {
		t.Error("the result must carry no request, but carried one")
	}
}

func TestTheRouteMatchedResultCarriesTheRoute(t *testing.T) {
	t.Parallel()

	result := data.NewRouteMatchedRoute(nil)

	if result.IsResponse() {
		t.Error("a result that carries a route must not report a response, but did")
	}

	if result.GetResponse() != nil {
		t.Error("the result must carry no response, but carried one")
	}

	if result.GetRoute() != nil {
		t.Error("the result must carry the route that it received, but did not")
	}
}

func TestTheRouteMatchedResultCarriesTheResponse(t *testing.T) {
	t.Parallel()

	sent := response.NewResponse(nil, 0, nil)

	result := data.NewRouteMatchedResponse(sent)

	if !result.IsResponse() {
		t.Error("a result that carries a response must report one, but did not")
	}

	if result.GetResponse() != contract.ResponseContract(sent) {
		t.Error("the result must carry the response, but did not")
	}

	if result.GetRoute() != nil {
		t.Error("the result must carry no route, but carried one")
	}
}

func TestEachResultSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var received contract.RequestReceivedResultContract = data.NewRequestReceivedRequest(nil)
	var matched contract.RouteMatchedResultContract = data.NewRouteMatchedRoute(nil)

	if received.IsResponse() || matched.IsResponse() {
		t.Error("each contract must report no response, but one did")
	}
}
