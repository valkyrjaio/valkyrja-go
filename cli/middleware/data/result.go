/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds what a CLI middleware returns.
//
// The TypeScript port returns a union, and Go has no union. Each result is a
// struct here, and the contract that it satisfies reports which side it carries.
package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

type InputReceivedResult struct {
	input  contract.InputContract
	output contract.OutputContract
}

// NewInputReceivedInput builds the result that continues the run with the input.
func NewInputReceivedInput(input contract.InputContract) *InputReceivedResult {
	return &InputReceivedResult{input: input}
}

// NewInputReceivedOutput builds the result that ends the run with the output.
func NewInputReceivedOutput(output contract.OutputContract) *InputReceivedResult {
	return &InputReceivedResult{output: output}
}

// GetInput returns the input that the next middleware receives.
func (r *InputReceivedResult) GetInput() contract.InputContract {
	return r.input
}

// GetOutput returns the output that ends the run, and nil where the run
// continues.
func (r *InputReceivedResult) GetOutput() contract.OutputContract {
	return r.output
}

// IsOutput reports whether the result ends the run with an output.
func (r *InputReceivedResult) IsOutput() bool {
	return r.output != nil
}

type RouteMatchedResult struct {
	route  contract.RouteContract
	output contract.OutputContract
}

// NewRouteMatchedRoute builds the result that continues the run with the route.
func NewRouteMatchedRoute(route contract.RouteContract) *RouteMatchedResult {
	return &RouteMatchedResult{route: route}
}

// NewRouteMatchedOutput builds the result that ends the run with the output.
func NewRouteMatchedOutput(output contract.OutputContract) *RouteMatchedResult {
	return &RouteMatchedResult{output: output}
}

// GetRoute returns the route that the next middleware receives.
func (r *RouteMatchedResult) GetRoute() contract.RouteContract {
	return r.route
}

// GetOutput returns the output that ends the run, and nil where the run
// continues.
func (r *RouteMatchedResult) GetOutput() contract.OutputContract {
	return r.output
}

// IsOutput reports whether the result ends the run with an output.
func (r *RouteMatchedResult) IsOutput() bool {
	return r.output != nil
}
