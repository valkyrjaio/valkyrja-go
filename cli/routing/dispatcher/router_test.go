/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package dispatcher_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/argument"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/input"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
	middlewarefixtures "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/fixtures"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/handler"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/dispatcher"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/fixtures"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

const (
	routeName     = "cache:clear"
	middlewareKey = "valkyrja.tests.cli.Middleware"
)

// newRouter builds a router over the commands, and returns the container that it
// binds into.
func newRouter(
	commands ...contract.RouteContract,
) (*dispatcher.Router, containercontract.ContainerContract) {
	container := manager.NewContainer(nil)
	built := collection.NewCollection()
	built.Add(commands...)

	return dispatcher.NewRouter(
		container,
		built,
		factory.NewOutputFactory(nil),
		middlewarehandler.NewRouteMatchedHandler(container),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	), container
}

func TestTheRouterRunsTheCommandThatTheInputNames(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	built, container := newRouter(data.NewRoute(routeName, "Clear the cache", handler.Run))

	output := built.Dispatch(input.NewInput("", routeName))

	if len(ran) != 1 || ran[0] != routeName {
		t.Errorf("the router must run the command that the input names, but ran: %v", ran)
	}

	if output.GetExitCode() != interactionconstant.ExitCodeSuccess {
		t.Error("the router must return the output of the command, but did not")
	}

	_, err := container.GetSingleton(constant.RouteContractServiceID)
	if err != nil {
		t.Error("the router must bind the command in the container, but did not")
	}
}

func TestTheRouterReportsACommandThatTheApplicationDoesNotHold(t *testing.T) {
	t.Parallel()

	built, _ := newRouter()

	output := built.Dispatch(input.NewInput("", "missing"))

	if output.GetExitCode() != interactionconstant.ExitCodeError {
		t.Error("an input that names no command must report a failure, but did not")
	}

	if !strings.Contains(output.GetMessages()[0].GetText(), "missing") {
		t.Error("the output must name the command that was not found, but did not")
	}
}

func TestTheRouterFillsEachArgumentInTheOrderTheCallerTyped(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	route := data.NewRoute(routeName, "Clear the cache", handler.Run).
		WithArguments(
			data.NewArgumentParameter("first", "The first"),
			data.NewArgumentParameter("second", "The second"),
		)

	built, container := newRouter(route)

	built.Dispatch(input.NewInput("", routeName).
		WithArguments(argument.NewArgument("one"), argument.NewArgument("two")))

	matched := readRoute(t, container)

	if matched.GetArgument("first").GetFirstValue() != "one" {
		t.Error("the first argument must take the first value that the caller typed, but did not")
	}

	if matched.GetArgument("second").GetFirstValue() != "two" {
		t.Error("the second argument must take the second value that the caller typed, but did not")
	}
}

func TestAnArrayArgumentTakesEveryValueThatIsLeft(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	route := data.NewRoute(routeName, "Clear the cache", handler.Run).
		WithArguments(
			data.NewArgumentParameter("names", "Each name").
				WithValueMode(constant.ArgumentValueModeArray),
		)

	built, container := newRouter(route)

	built.Dispatch(input.NewInput("", routeName).
		WithArguments(argument.NewArgument("one"), argument.NewArgument("two")))

	if len(readRoute(t, container).GetArgument("names").GetArguments()) != 2 {
		t.Error("an array argument must take every value that is left, but did not")
	}
}

func TestTheRouterFillsAnOptionUnderItsLongNameAndItsShortName(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	route := data.NewRoute(routeName, "Clear the cache", handler.Run).
		WithOptions(
			data.NewOptionParameter("force", "Clear without asking").
				WithShortNames("f").
				WithValueMode(constant.OptionValueModeArray),
		)

	built, container := newRouter(route)

	built.Dispatch(input.NewInput("", routeName).WithOptions(
		newOption("force", "one"),
		newOption("f", "two"),
		newOption("other", "three"),
	))

	filled := readRoute(t, container).GetOption("force")

	if len(filled.GetOptions()) != 2 {
		t.Errorf("the option must take its long name and its short name, but took: %d", len(filled.GetOptions()))
	}
}

func TestTheRouterReportsAValueThatTheCommandDoesNotAccept(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	route := data.NewRoute(routeName, "Clear the cache", handler.Run).
		WithArguments(
			data.NewArgumentParameter("first", "The first").
				WithMode(constant.ArgumentModeRequired),
		)

	built, _ := newRouter(route)

	output := built.Dispatch(input.NewInput("", routeName))

	if output.GetExitCode() != interactionconstant.ExitCodeError {
		t.Error("an argument that the caller left out must report a failure, but did not")
	}

	if len(ran) != 0 {
		t.Error("a command whose parameters are invalid must not run, but ran")
	}
}

func TestTheRouterReportsAnOptionValueThatTheCommandDoesNotAccept(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	tests := map[string]contract.RouteContract{
		"an option that takes no value": data.NewRoute(routeName, "Clear the cache", handler.Run).
			WithOptions(data.NewOptionParameter("force", "Clear without asking").
				WithValueMode(constant.OptionValueModeNone)),
		"a value that the option rejects": data.NewRoute(routeName, "Clear the cache", handler.Run).
			WithOptions(data.NewOptionParameter("force", "Clear without asking").
				WithValidValues("yes")),
	}

	for name, route := range tests {
		built, _ := newRouter(route)

		output := built.Dispatch(input.NewInput("", routeName).WithOptions(newOption("force", "no")))

		if output.GetExitCode() != interactionconstant.ExitCodeError {
			t.Errorf("%s must report a failure, but did not", name)
		}
	}
}

func TestTheRouterFilesTheMiddlewareOfTheMatchedCommand(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	route := data.NewRoute(routeName, "Clear the cache", handler.Run).
		WithProcessExitingMiddleware(middlewareKey)

	built, _ := newRouter(route)

	built.Dispatch(input.NewInput("", routeName))

	if built.GetProcessExitingHandler() == nil {
		t.Error("the router must hold the process-exiting handler, but held none")
	}
}

func TestARouteMatchedMiddlewareThatEndsTheRunStopsTheCommand(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}
	ending := output.NewOutput(nil).WithExitCode(interactionconstant.ExitCodeUsageError)

	record := []string{}
	container := manager.NewContainer(nil)
	container.Bind(middlewareKey, func(_ containercontract.ContainerContract, _ []any) any {
		return &middlewarefixtures.EndingMiddlewareFixture{Output: ending, Record: &record}
	})

	built := collection.NewCollection()
	built.Add(data.NewRoute(routeName, "Clear the cache", handler.Run))

	dispatched := dispatcher.NewRouter(
		container,
		built,
		factory.NewOutputFactory(nil),
		middlewarehandler.NewRouteMatchedHandler(container, middlewareKey),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	).Dispatch(input.NewInput("", routeName))

	if dispatched.GetExitCode() != interactionconstant.ExitCodeUsageError {
		t.Error("a middleware that ends the run must return its own output, but did not")
	}

	if len(ran) != 0 {
		t.Error("a command whose middleware ended the run must not run, but ran")
	}
}

// newOption builds one option that a caller typed.
func newOption(name string, value string) contract.OptionContract {
	return argument.NewOption(name, interactionconstant.OptionTypeLong).WithValue(value)
}

// readRoute returns the command that the router bound in the container.
func readRoute(t *testing.T, container containercontract.ContainerContract) contract.RouteContract {
	t.Helper()

	resolved, err := container.GetSingleton(constant.RouteContractServiceID)
	if err != nil {
		t.Fatalf("the router must bind the command in the container, but reported: %v", err)
	}

	route, isRoute := resolved.(contract.RouteContract)
	if !isRoute {
		t.Fatal("the container must hold the command that the router matched, but held another value")
	}

	return route
}

func TestAnArrayArgumentTakesOnlyTheValuesThatAreLeft(t *testing.T) {
	t.Parallel()

	ran := []string{}
	handler := &fixtures.RecordingHandlerFixture{Ran: &ran}

	// Each parameter consumes what it took, so an array parameter that follows
	// another one must not take the value that the first one already holds.
	route := data.NewRoute(routeName, "Clear the cache", handler.Run).
		WithArguments(
			data.NewArgumentParameter("name", "The name"),
			data.NewArgumentParameter("files", "Each file").
				WithValueMode(constant.ArgumentValueModeArray),
		)

	built, container := newRouter(route)

	built.Dispatch(input.NewInput("", routeName).WithArguments(
		argument.NewArgument("a"), argument.NewArgument("b"), argument.NewArgument("c"),
	))

	matched := readRoute(t, container)

	if matched.GetArgument("name").GetFirstValue() != "a" {
		t.Errorf("the first argument must take the first value, but took: %q",
			matched.GetArgument("name").GetFirstValue())
	}

	files := matched.GetArgument("files").GetArguments()

	if len(files) != 2 || files[0].GetValue() != "b" || files[1].GetValue() != "c" {
		t.Errorf("the array argument must take only the values that are left, but took: %v",
			valuesOf(files))
	}
}

// valuesOf returns the value of each argument.
func valuesOf(arguments []contract.ArgumentContract) []string {
	values := make([]string, 0, len(arguments))

	for _, held := range arguments {
		values = append(values, held.GetValue())
	}

	return values
}
