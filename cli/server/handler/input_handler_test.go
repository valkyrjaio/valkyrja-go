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
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/input"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
	middlewarefixtures "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/fixtures"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/handler"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/router"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/handler"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

const (
	routeName     = "cache:clear"
	middlewareKey = "Valkyrja.Tests.Cli.EndingMiddleware"
)

// newHandler builds the server's entry point over one command, and returns the
// container that it binds into.
func newHandler(
	t *testing.T,
	route contract.RouteContract,
	inputReceived []string,
	exiter func(code int),
) (*handler.InputHandler, containercontract.ContainerContract) {
	t.Helper()

	container := manager.NewContainer(nil)
	built := collection.NewCollection()

	if route != nil {
		built.Add(route)
	}

	dispatcher := router.NewRouter(
		container,
		built,
		factory.NewOutputFactory(nil),
		middlewarehandler.NewRouteMatchedHandler(container),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	)

	return handler.NewInputHandler(
		container,
		dispatcher,
		middlewarehandler.NewInputReceivedHandler(container, inputReceived...),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
		factory.NewOutputFactory(nil),
		exiter,
	), container
}

func TestTheHandlerReturnsTheOutputOfTheCommand(t *testing.T) {
	t.Parallel()

	ran := []string{}
	commandHandler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	built, container := newHandler(t, data.NewRoute(routeName, "Clear the cache", commandHandler.Run), nil, nil)

	result := built.Handle(input.NewInput("", routeName))

	if result.GetExitCode() != interactionconstant.ExitCodeSuccess {
		t.Error("the handler must return the output of the command, but did not")
	}

	_, outputErr := container.GetSingleton(interactionconstant.OutputContractServiceID)
	if outputErr != nil {
		t.Error("the handler must bind the output in the container, but did not")
	}

	_, inputErr := container.GetSingleton(interactionconstant.InputContractServiceID)
	if inputErr != nil {
		t.Error("the handler must bind the input in the container, but did not")
	}
}

func TestTheHandlerTurnsAPanicIntoAnOutput(t *testing.T) {
	t.Parallel()

	tests := map[string]any{
		"a panic that carries an error": errors.New("the command failed"),
		"a panic that carries a string": "the command failed",
	}

	for name, recovered := range tests {
		commandHandler := &fixtures.PanickingHandlerFixture{Recovered: recovered}

		built, _ := newHandler(t, data.NewRoute(routeName, "Clear the cache", commandHandler.Run), nil, nil)

		result := built.Handle(input.NewInput("", routeName))

		if result.GetExitCode() != interactionconstant.ExitCodeError {
			t.Errorf("%s must report a failure, but did not", name)
		}

		if !strings.Contains(textOf(result), "the command failed") {
			t.Errorf("%s must reach the caller, but the output is: %q", name, textOf(result))
		}
	}
}

func TestAnInputReceivedMiddlewareThatEndsTheRunStopsTheCommand(t *testing.T) {
	t.Parallel()

	ran := []string{}
	commandHandler := &fixtures.RecordingHandlerFixture{Ran: &ran}
	record := []string{}
	ending := output.NewOutput(nil).WithExitCode(interactionconstant.ExitCodeUsageError)

	built, container := newHandler(
		t,
		data.NewRoute(routeName, "Clear the cache", commandHandler.Run),
		[]string{middlewareKey},
		nil,
	)

	container.Bind(middlewareKey, func(_ containercontract.ContainerContract, _ []any) any {
		return &middlewarefixtures.EndingMiddlewareFixture{Output: ending, Record: &record}
	})

	result := built.Handle(input.NewInput("", routeName))

	if result.GetExitCode() != interactionconstant.ExitCodeUsageError {
		t.Error("a middleware that ends the run must return its own output, but did not")
	}

	if len(ran) != 0 {
		t.Error("a command whose middleware ended the run must not run, but ran")
	}
}

func TestRunWritesTheOutputAndExitsWithItsCode(t *testing.T) {
	t.Parallel()

	ran := []string{}
	commandHandler := &fixtures.RecordingHandlerFixture{Ran: &ran}
	exited := []int{}

	built, _ := newHandler(
		t,
		data.NewRoute(routeName, "Clear the cache", commandHandler.Run),
		nil,
		func(code int) { exited = append(exited, code) },
	)

	built.Run(input.NewInput("", routeName))

	if len(exited) != 1 || exited[0] != int(interactionconstant.ExitCodeSuccess) {
		t.Errorf("Run must exit with the code of the output, but exited with: %v", exited)
	}
}

func TestRunEndsNoProcessWhereItWasGivenNoExiter(t *testing.T) {
	t.Parallel()

	ran := []string{}
	commandHandler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	built, _ := newHandler(t, data.NewRoute(routeName, "Clear the cache", commandHandler.Run), nil, nil)

	built.Run(input.NewInput("", routeName))

	if len(ran) != 1 {
		t.Error("Run must handle the input, but did not")
	}
}

// textOf returns the text of every message that the output holds.
func textOf(built contract.OutputContract) string {
	text := &strings.Builder{}

	for _, held := range built.GetMessages() {
		text.WriteString(held.GetText())
	}

	return text.String()
}
