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
	interactionfactory "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	interactionfixtures "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/input"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
	middlewarefixtures "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/fixtures"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/handler"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/dispatcher"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/handler"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

const (
	routeName     = "cache:clear"
	middlewareKey = "valkyrja.tests.cli.EndingMiddleware"
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

	routing := dispatcher.NewRouter(
		container,
		built,
		interactionfactory.NewOutputFactory(nil),
		middlewarehandler.NewRouteMatchedHandler(container),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	)

	return handler.NewInputHandler(
		container,
		routing,
		middlewarehandler.NewInputReceivedHandler(container, inputReceived...),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
		interactionfactory.NewOutputFactory(nil),
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

func TestRunReportsAFailedWriteThroughTheFactoryOutput(t *testing.T) {
	t.Parallel()

	// The command writes to a stream that takes nothing, so the run reports the
	// failure through an output that the factory builds instead.
	reported := &strings.Builder{}
	exited := []int{}

	container := manager.NewContainer(nil)
	failing := interactionfactory.NewOutputFactoryForWriter(nil, &interactionfixtures.FailingWriterFixture{})

	built := collection.NewCollection()
	built.Add(data.NewRoute(routeName, "Clear the cache", writesToAFailingStream(failing)))

	routing := dispatcher.NewRouter(
		container, built, failing,
		middlewarehandler.NewRouteMatchedHandler(container),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	)

	handler.NewInputHandler(
		container,
		routing,
		middlewarehandler.NewInputReceivedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
		interactionfactory.NewOutputFactoryForWriter(nil, reported),
		func(code int) { exited = append(exited, code) },
	).Run(input.NewInput("", routeName))

	if !strings.Contains(reported.String(), "Cli Server Error:") {
		t.Errorf("a failed write must reach the caller, but the report is: %q", reported.String())
	}

	if len(exited) != 1 || exited[0] != int(interactionconstant.ExitCodeError) {
		t.Errorf("a failed write must exit with a failure, but exited with: %v", exited)
	}
}

func TestRunExitsWithAFailureWhereEveryDestinationRefusesTheWrite(t *testing.T) {
	t.Parallel()

	// Nothing can be written anywhere, so the exit code is the last diagnostic.
	exited := []int{}

	container := manager.NewContainer(nil)
	failing := interactionfactory.NewOutputFactoryForWriter(nil, &interactionfixtures.FailingWriterFixture{})

	built := collection.NewCollection()
	built.Add(data.NewRoute(routeName, "Clear the cache", writesToAFailingStream(failing)))

	routing := dispatcher.NewRouter(
		container, built, failing,
		middlewarehandler.NewRouteMatchedHandler(container),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	)

	handler.NewInputHandler(
		container, routing,
		middlewarehandler.NewInputReceivedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
		failing,
		func(code int) { exited = append(exited, code) },
	).Run(input.NewInput("", routeName))

	if len(exited) != 1 || exited[0] != int(interactionconstant.ExitCodeError) {
		t.Errorf("the exit code must report the failure, but exited with: %v", exited)
	}
}

// writesToAFailingStream returns a command that reports through the factory, so
// the run reaches the write that the factory's stream refuses.
func writesToAFailingStream(built contract.OutputFactoryContract) contract.CliHandlerFunc {
	return func(_ containercontract.ContainerContract, _ contract.RouteContract) contract.OutputContract {
		return built.CreateOutput(interactionconstant.ExitCodeSuccess, message.NewMessage("the command ran"))
	}
}

func TestRunFallsBackWhereAMiddlewareReturnsAFailingOutput(t *testing.T) {
	t.Parallel()

	// The throwable-caught middleware answers with an output whose destination
	// is the one that just failed, so the run falls back to the factory.
	reported := &strings.Builder{}
	exited := []int{}

	container := manager.NewContainer(nil)
	failing := interactionfactory.NewOutputFactoryForWriter(nil, &interactionfixtures.FailingWriterFixture{})
	good := interactionfactory.NewOutputFactoryForWriter(nil, reported)

	container.SetSingleton(replacingMiddlewareID, &replacingThrowableCaughtMiddlewareFixture{
		output: failing.CreateOutput(interactionconstant.ExitCodeError, message.NewMessage("unreachable")),
	})

	built := collection.NewCollection()
	built.Add(data.NewRoute(routeName, "Clear the cache", writesToAFailingStream(failing)))

	routing := dispatcher.NewRouter(
		container, built, good,
		middlewarehandler.NewRouteMatchedHandler(container),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	)

	handler.NewInputHandler(
		container, routing,
		middlewarehandler.NewInputReceivedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container, replacingMiddlewareID),
		middlewarehandler.NewProcessExitingHandler(container),
		good,
		func(code int) { exited = append(exited, code) },
	).Run(input.NewInput("", routeName))

	if !strings.Contains(reported.String(), "Cli Server Error:") {
		t.Errorf("the fallback must reach the caller, but the report is: %q", reported.String())
	}

	if len(exited) != 1 || exited[0] != int(interactionconstant.ExitCodeError) {
		t.Errorf("the run must exit with a failure, but exited with: %v", exited)
	}
}

// replacingMiddlewareID is the binding key of the middleware below.
const replacingMiddlewareID = "valkyrja.tests.cli.ReplacingThrowableCaughtMiddleware"

// replacingThrowableCaughtMiddlewareFixture answers with an output of its own.
type replacingThrowableCaughtMiddlewareFixture struct {
	output contract.OutputContract
}

// ThrowableCaught returns the output that the fixture holds.
func (m *replacingThrowableCaughtMiddlewareFixture) ThrowableCaught(
	_ contract.InputContract,
	_ contract.OutputContract,
	_ error,
	_ contract.ThrowableCaughtHandlerContract,
) contract.OutputContract {
	return m.output
}
