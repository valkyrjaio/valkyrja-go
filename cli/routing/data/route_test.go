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

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/argument"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

const (
	routeName        = "cache:clear"
	routeDescription = "Clear the cache"
	middlewareKey    = "valkyrja.test.Middleware"
)

// noopHandler is a handler that a route test gives a route.
func noopHandler(_ containercontract.ContainerContract, _ contract.RouteContract) contract.OutputContract {
	return nil
}

// newTestRoute builds the route that each test in this file starts from.
func newTestRoute() contract.RouteContract {
	return data.NewRoute(routeName, routeDescription, noopHandler)
}

func TestARouteReadsWhatItWasBuiltWith(t *testing.T) {
	t.Parallel()

	built := newTestRoute()

	if built.GetName() != routeName || built.GetDescription() != routeDescription {
		t.Error("the route must read its name and its description, but did not")
	}

	if built.GetHandler() == nil {
		t.Error("the route must hold the handler, but held none")
	}

	if built.HasHelpText() || built.GetHelpText() != nil || built.GetHelpTextMessage() != nil {
		t.Error("a route that was built with no help text must build none, but built one")
	}

	if built.HasArguments() || built.HasOptions() {
		t.Error("a route that was built with no parameter must take none, but took one")
	}
}

func TestEachRouteNameAndHandlerWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := newTestRoute()

	if built.WithName(otherName).GetName() != otherName {
		t.Error("WithName must hold the new name, but did not")
	}

	if built.WithDescription(otherName).GetDescription() != otherName {
		t.Error("WithDescription must hold the new description, but did not")
	}

	replaced := built.WithHandler(func(
		_ containercontract.ContainerContract,
		_ contract.RouteContract,
	) contract.OutputContract {
		return nil
	})

	if replaced.GetHandler() == nil {
		t.Error("WithHandler must hold the new handler, but did not")
	}

	if built.GetName() != routeName {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestARouteBuildsItsHelpText(t *testing.T) {
	t.Parallel()

	built := newTestRoute().WithHelpText(func() contract.MessageContract {
		return message.NewMessage("Help for the command.")
	})

	if !built.HasHelpText() || built.GetHelpText() == nil {
		t.Error("a route that was given help text must build it, but did not")
	}

	if built.GetHelpTextMessage().GetText() != "Help for the command." {
		t.Error("the route must build the help text that it was given, but did not")
	}
}

func TestARouteHoldsItsArguments(t *testing.T) {
	t.Parallel()

	first := data.NewArgumentParameter("first", parameterDescription)
	second := data.NewArgumentParameter("second", parameterDescription)

	built := newTestRoute().WithArguments(first)
	added := built.WithAddedArguments(second)

	if !built.HasArgument("first") || built.GetArgument("first") == nil {
		t.Error("the route must report the argument that it holds, but did not")
	}

	if built.HasArgument("second") || built.GetArgument("second") != nil {
		t.Error("the route must report no argument under a name that it does not hold, but did")
	}

	if len(added.GetArguments()) != 2 || len(built.GetArguments()) != 1 {
		t.Error("WithAddedArguments must append the argument and leave the receiver unchanged, but did not")
	}
}

func TestARouteHoldsItsOptions(t *testing.T) {
	t.Parallel()

	first := data.NewOptionParameter("first", parameterDescription)
	second := data.NewOptionParameter("second", parameterDescription)

	built := newTestRoute().WithOptions(first)
	added := built.WithAddedOptions(second)

	if !built.HasOption("first") || built.GetOption("first") == nil {
		t.Error("the route must report the option that it holds, but did not")
	}

	if built.HasOption("second") || built.GetOption("second") != nil {
		t.Error("the route must report no option under a name that it does not hold, but did")
	}

	if len(added.GetOptions()) != 2 || len(built.GetOptions()) != 1 {
		t.Error("WithAddedOptions must append the option and leave the receiver unchanged, but did not")
	}
}

func TestEachRouteMiddlewareStageHoldsItsOwnMiddleware(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		with  func(contract.RouteContract) contract.RouteContract
		added func(contract.RouteContract) contract.RouteContract
		read  func(contract.RouteContract) []string
	}{
		"route matched": {
			with: func(r contract.RouteContract) contract.RouteContract {
				return r.WithRouteMatchedMiddleware(middlewareKey)
			},
			added: func(r contract.RouteContract) contract.RouteContract {
				return r.WithAddedRouteMatchedMiddleware(middlewareKey)
			},
			read: contract.RouteContract.GetRouteMatchedMiddleware,
		},
		"route dispatched": {
			with: func(r contract.RouteContract) contract.RouteContract {
				return r.WithRouteDispatchedMiddleware(middlewareKey)
			},
			added: func(r contract.RouteContract) contract.RouteContract {
				return r.WithAddedRouteDispatchedMiddleware(middlewareKey)
			},
			read: contract.RouteContract.GetRouteDispatchedMiddleware,
		},
		"throwable caught": {
			with: func(r contract.RouteContract) contract.RouteContract {
				return r.WithThrowableCaughtMiddleware(middlewareKey)
			},
			added: func(r contract.RouteContract) contract.RouteContract {
				return r.WithAddedThrowableCaughtMiddleware(middlewareKey)
			},
			read: contract.RouteContract.GetThrowableCaughtMiddleware,
		},
		"process exiting": {
			with: func(r contract.RouteContract) contract.RouteContract {
				return r.WithProcessExitingMiddleware(middlewareKey)
			},
			added: func(r contract.RouteContract) contract.RouteContract {
				return r.WithAddedProcessExitingMiddleware(middlewareKey)
			},
			read: contract.RouteContract.GetProcessExitingMiddleware,
		},
	}

	for name, test := range tests {
		built := test.with(newTestRoute())

		if len(test.read(built)) != 1 {
			t.Errorf("the %s stage must hold the middleware, but held: %v", name, test.read(built))
		}

		// A middleware that is registered twice runs twice, so the same key
		// appended again must appear again.
		if len(test.read(test.added(built))) != 2 {
			t.Errorf("the %s stage must append a middleware that it holds already, but did not", name)
		}

		if len(test.read(built)) != 1 {
			t.Errorf("the %s stage must leave the receiver unchanged, but did not", name)
		}
	}
}

func TestTheRoutingDataHoldsEachRoute(t *testing.T) {
	t.Parallel()

	built := data.NewCliRoutingData(map[string]contract.RouteContract{routeName: newTestRoute()})

	if len(built.GetRoutes()) != 1 {
		t.Errorf("the routing data must hold each route, but held: %d", len(built.GetRoutes()))
	}

	if len(data.NewCliRoutingData(nil).GetRoutes()) != 0 {
		t.Error("routing data that was built with no route must hold none, but held some")
	}
}

func TestEachGlobalOptionIsBuiltWithNoValue(t *testing.T) {
	t.Parallel()

	options := map[string]contract.OptionParameterContract{
		constant.OptionNameHelp:          data.NewHelpOptionParameter(),
		constant.OptionNameVersion:       data.NewVersionOptionParameter(),
		constant.OptionNameQuiet:         data.NewQuietOptionParameter(),
		constant.OptionNameSilent:        data.NewSilentOptionParameter(),
		constant.OptionNameNoInteraction: data.NewNoInteractionOptionParameter(),
	}

	for name, option := range options {
		if option.GetName() != name {
			t.Errorf("the global option must be named %q, but was named %q", name, option.GetName())
		}

		if option.GetValueMode() != constant.OptionValueModeNone {
			t.Errorf("the global option %q must take no value, but took one", name)
		}

		if len(option.GetShortNames()) != 1 {
			t.Errorf("the global option %q must carry one short name, but carried: %v", name, option.GetShortNames())
		}
	}
}

func TestARouteReportsWhatTheCallerProvided(t *testing.T) {
	t.Parallel()

	declared := data.NewOptionParameter("namespace", parameterDescription)

	provided, err := declared.WithOptions(
		argument.NewOption("namespace", interactionconstant.OptionTypeLong).WithValue("db:"),
	)
	if err != nil {
		t.Fatalf("WithOptions must accept one value for a default option, but returned: %v", err)
	}

	withDefault := data.NewOptionParameter("namespace", parameterDescription).WithDefaultValue("app:")

	bare := newTestRoute()
	all := "all"

	if bare.HasProvidedOption("namespace") || bare.GetOptionValue("namespace", &all) != "all" {
		t.Error("a route that declares no option must report nothing provided, but did not")
	}

	if bare.GetOptionValue("namespace", nil) != "" {
		t.Error("a route that declares no option and takes no default must give an empty string, but did not")
	}

	onlyDeclared := newTestRoute().WithOptions(declared)

	if onlyDeclared.HasProvidedOption("namespace") {
		t.Error("an option the caller left out must not report itself provided, but did")
	}

	filled := newTestRoute().WithOptions(provided)

	if !filled.HasProvidedOption("namespace") || filled.GetOptionValue("namespace", nil) != "db:" {
		t.Error("an option the caller gave must report its value, but did not")
	}

	declaring := newTestRoute().WithOptions(withDefault)

	if declaring.GetOptionValue("namespace", nil) != "app:" {
		t.Error("a nil default must reach the option's declared default value, but did not")
	}

	empty := ""

	if declaring.GetOptionValue("namespace", &empty) != "" {
		t.Error("an empty default given at the call site must win, but did not")
	}
}

func TestARouteReportsWhatTheCallerProvidedForAnArgument(t *testing.T) {
	t.Parallel()

	declared := data.NewArgumentParameter("namespace", parameterDescription)
	provided := declared.WithArguments(argument.NewArgument("db:"))

	bare := newTestRoute()

	if bare.HasProvidedArgument("namespace") || bare.GetArgumentValue("namespace", "all") != "all" {
		t.Error("a route that declares no argument must report nothing provided, but did not")
	}

	onlyDeclared := newTestRoute().WithArguments(declared)

	if onlyDeclared.HasProvidedArgument("namespace") || onlyDeclared.GetArgumentValue("namespace", "all") != "all" {
		t.Error("an argument the caller left out must not report itself provided, but did")
	}

	filled := newTestRoute().WithArguments(provided)

	if !filled.HasProvidedArgument("namespace") || filled.GetArgumentValue("namespace", "all") != "db:" {
		t.Error("an argument the caller gave must report its value, but did not")
	}
}
