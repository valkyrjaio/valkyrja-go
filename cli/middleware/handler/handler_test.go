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

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/input"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/handler"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

const (
	firstID   = "valkyrja.tests.cli.FirstMiddleware"
	secondID  = "valkyrja.tests.cli.SecondMiddleware"
	endingID  = "valkyrja.tests.cli.EndingMiddleware"
	missingID = "valkyrja.tests.cli.MissingMiddleware"
)

// newContainer builds a container that resolves a recording middleware under
// each binding key, and an ending middleware under its own.
func newContainer(record *[]string, ending contract.OutputContract) containercontract.ContainerContract {
	container := manager.NewContainer(nil)

	for _, id := range []string{firstID, secondID} {
		container.Bind(id, func(_ containercontract.ContainerContract, _ []any) any {
			return &fixtures.RecordingMiddlewareFixture{Name: id, Record: record}
		})
	}

	container.Bind(endingID, func(_ containercontract.ContainerContract, _ []any) any {
		return &fixtures.EndingMiddlewareFixture{Output: ending, Record: record}
	})

	return container
}

func TestEachHandlerRunsItsMiddlewareInOrder(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	typed := input.NewInput("", "cache:clear")
	written := output.NewOutput(nil)

	handler.NewInputReceivedHandler(container, firstID, secondID).InputReceived(typed)

	if !slices.Equal(record, []string{firstID, secondID}) {
		t.Errorf("the handler must run each middleware in order, but ran: %v", record)
	}

	record = record[:0]

	handler.NewRouteMatchedHandler(container, firstID, secondID).RouteMatched(typed, nil)

	if !slices.Equal(record, []string{firstID, secondID}) {
		t.Errorf("the route-matched handler must run each middleware, but ran: %v", record)
	}

	record = record[:0]

	handler.NewRouteNotMatchedHandler(container, firstID).RouteNotMatched(typed, written)
	handler.NewRouteDispatchedHandler(container, firstID).RouteDispatched(typed, written, nil)
	handler.NewThrowableCaughtHandler(container, firstID).ThrowableCaught(typed, written, errors.New("failed"))
	handler.NewProcessExitingHandler(container, firstID).ProcessExiting(typed, written)

	if len(record) != 4 {
		t.Errorf("each remaining handler must run its middleware, but the record is: %v", record)
	}
}

func TestAHandlerWithNoMiddlewareReturnsWhatItReceived(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	typed := input.NewInput("", "cache:clear")
	written := output.NewOutput(nil)

	result := handler.NewInputReceivedHandler(container).InputReceived(typed)
	if result.IsOutput() || result.GetInput() != contract.InputContract(typed) {
		t.Error("a handler with no middleware must return the input that it received, but did not")
	}

	matched := handler.NewRouteMatchedHandler(container).RouteMatched(typed, nil)
	if matched.IsOutput() {
		t.Error("a handler with no middleware must continue the run, but ended it")
	}

	if handler.NewRouteNotMatchedHandler(container).RouteNotMatched(typed, written) != written {
		t.Error("a handler with no middleware must return the output that it received, but did not")
	}

	if handler.NewRouteDispatchedHandler(container).RouteDispatched(typed, written, nil) != written {
		t.Error("a handler with no middleware must return the output that it received, but did not")
	}

	if handler.NewThrowableCaughtHandler(container).ThrowableCaught(typed, written, nil) != written {
		t.Error("a handler with no middleware must return the output that it received, but did not")
	}

	handler.NewProcessExitingHandler(container).ProcessExiting(typed, written)

	if len(record) != 0 {
		t.Errorf("a handler with no middleware must run none, but ran: %v", record)
	}
}

func TestAMiddlewareThatEndsTheRunStopsTheHandler(t *testing.T) {
	t.Parallel()

	record := []string{}
	ending := output.NewOutput(nil)
	container := newContainer(&record, ending)

	typed := input.NewInput("", "cache:clear")

	result := handler.NewInputReceivedHandler(container, endingID, firstID).InputReceived(typed)
	if !result.IsOutput() || result.GetOutput() != contract.OutputContract(ending) {
		t.Error("a middleware that ends the run must return its own output, but did not")
	}

	matched := handler.NewRouteMatchedHandler(container, endingID, firstID).RouteMatched(typed, nil)
	if !matched.IsOutput() {
		t.Error("a middleware that ends the run must end it, but did not")
	}

	if slices.Contains(record, firstID) {
		t.Errorf("a middleware after the one that ended the run must not run, but the record is: %v", record)
	}
}

func TestAHandlerSkipsAMiddlewareThatTheContainerCannotResolve(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	typed := input.NewInput("", "cache:clear")

	result := handler.NewInputReceivedHandler(container, missingID).InputReceived(typed)

	if result.IsOutput() {
		t.Error("a middleware that the container cannot resolve must be skipped, but ended the run")
	}
}

func TestAddAppendsEachMiddlewareAfterTheOnesTheHandlerHolds(t *testing.T) {
	t.Parallel()

	record := []string{}
	container := newContainer(&record, nil)

	built := handler.NewInputReceivedHandler(container, firstID)
	built.Add(secondID)

	built.InputReceived(input.NewInput("", "cache:clear"))

	if !slices.Equal(record, []string{firstID, secondID}) {
		t.Errorf("Add must append the middleware after the ones the handler holds, but ran: %v", record)
	}
}
