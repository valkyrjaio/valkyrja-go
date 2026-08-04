/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package fixtures holds the reusable doubles that the CLI middleware tests
// build on.
package fixtures

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/data"
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

// InputReceived records the run and passes it on.
func (m *RecordingMiddlewareFixture) InputReceived(
	input contract.InputContract,
	handler contract.InputReceivedHandlerContract,
) contract.InputReceivedResultContract {
	m.record()

	return handler.InputReceived(input)
}

// RouteMatched records the run and passes it on.
func (m *RecordingMiddlewareFixture) RouteMatched(
	input contract.InputContract,
	route contract.RouteContract,
	handler contract.RouteMatchedHandlerContract,
) contract.RouteMatchedResultContract {
	m.record()

	return handler.RouteMatched(input, route)
}

// RouteNotMatched records the run and passes it on.
func (m *RecordingMiddlewareFixture) RouteNotMatched(
	input contract.InputContract,
	output contract.OutputContract,
	handler contract.RouteNotMatchedHandlerContract,
) contract.OutputContract {
	m.record()

	return handler.RouteNotMatched(input, output)
}

// RouteDispatched records the run and passes it on.
func (m *RecordingMiddlewareFixture) RouteDispatched(
	input contract.InputContract,
	output contract.OutputContract,
	route contract.RouteContract,
	handler contract.RouteDispatchedHandlerContract,
) contract.OutputContract {
	m.record()

	return handler.RouteDispatched(input, output, route)
}

// ThrowableCaught records the run and passes it on.
func (m *RecordingMiddlewareFixture) ThrowableCaught(
	input contract.InputContract,
	output contract.OutputContract,
	throwable error,
	handler contract.ThrowableCaughtHandlerContract,
) contract.OutputContract {
	m.record()

	return handler.ThrowableCaught(input, output, throwable)
}

// ProcessExiting records the run and passes it on.
func (m *RecordingMiddlewareFixture) ProcessExiting(
	input contract.InputContract,
	output contract.OutputContract,
	handler contract.ProcessExitingHandlerContract,
) {
	m.record()

	handler.ProcessExiting(input, output)
}

// record notes that the middleware ran.
func (m *RecordingMiddlewareFixture) record() {
	*m.Record = append(*m.Record, m.Name)
}

// EndingMiddlewareFixture is a middleware that ends the run with its own output.
type EndingMiddlewareFixture struct {
	Output contract.OutputContract
	Record *[]string
}

// InputReceived ends the run with the output that the fixture holds.
func (m *EndingMiddlewareFixture) InputReceived(
	_ contract.InputContract,
	_ contract.InputReceivedHandlerContract,
) contract.InputReceivedResultContract {
	*m.Record = append(*m.Record, "ending")

	return data.NewInputReceivedOutput(m.Output)
}

// RouteMatched ends the run with the output that the fixture holds.
func (m *EndingMiddlewareFixture) RouteMatched(
	_ contract.InputContract,
	_ contract.RouteContract,
	_ contract.RouteMatchedHandlerContract,
) contract.RouteMatchedResultContract {
	*m.Record = append(*m.Record, "ending")

	return data.NewRouteMatchedOutput(m.Output)
}
