/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds what a middleware returns where it returns one of two
// things.
//
// The TypeScript port writes a union — `RouteContract | ResponseContract` — and
// Go has no union. The Java port answers with a record that holds both, and this
// package holds the Go spelling of that record.
package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

type RequestReceivedResult struct {
	request  contract.ServerRequestContract
	response contract.ResponseContract
}

// NewRequestReceivedRequest builds the result that continues the run.
func NewRequestReceivedRequest(request contract.ServerRequestContract) *RequestReceivedResult {
	return &RequestReceivedResult{request: request}
}

// NewRequestReceivedResponse builds the result that ends the run.
func NewRequestReceivedResponse(response contract.ResponseContract) *RequestReceivedResult {
	return &RequestReceivedResult{response: response}
}

// GetRequest returns the request that the next middleware receives.
func (r *RequestReceivedResult) GetRequest() contract.ServerRequestContract {
	return r.request
}

// GetResponse returns the response that ends the run, and nil where the run
// continues.
func (r *RequestReceivedResult) GetResponse() contract.ResponseContract {
	return r.response
}

// IsResponse reports whether the result ends the run with a response.
func (r *RequestReceivedResult) IsResponse() bool {
	return r.response != nil
}

type RouteMatchedResult struct {
	route    contract.RouteContract
	response contract.ResponseContract
}

// NewRouteMatchedRoute builds the result that continues the run.
func NewRouteMatchedRoute(route contract.RouteContract) *RouteMatchedResult {
	return &RouteMatchedResult{route: route}
}

// NewRouteMatchedResponse builds the result that ends the run.
func NewRouteMatchedResponse(response contract.ResponseContract) *RouteMatchedResult {
	return &RouteMatchedResult{response: response}
}

// GetRoute returns the route that the next middleware receives.
func (r *RouteMatchedResult) GetRoute() contract.RouteContract {
	return r.route
}

// GetResponse returns the response that ends the run, and nil where the run
// continues.
func (r *RouteMatchedResult) GetResponse() contract.ResponseContract {
	return r.response
}

// IsResponse reports whether the result ends the run with a response.
func (r *RouteMatchedResult) IsResponse() bool {
	return r.response != nil
}
