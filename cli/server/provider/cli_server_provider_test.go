/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package provider_test

import (
	"testing"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationdata "github.com/valkyrjaio/valkyrja-go/v26/application/data"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	interactionfactory "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/handler"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	routingdispatcher "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/dispatcher"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/provider"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

func TestTheRouteProviderRegistersEveryCommandOfTheServer(t *testing.T) {
	t.Parallel()

	routes := (&provider.CliServerCliRoutesProvider{}).GetRoutes()

	if len(routes) != 4 {
		t.Fatalf("the provider must register every command of the server, but registered: %d", len(routes))
	}

	names := map[string]bool{}
	for _, route := range routes {
		names[route.GetName()] = true
	}

	for _, name := range []string{
		constant.CommandNameList,
		constant.CommandNameListBash,
		constant.CommandNameHelp,
		constant.CommandNameVersion,
	} {
		if !names[name] {
			t.Errorf("the provider must register %q, but did not", name)
		}
	}
}

func TestTheServiceProviderDefersItsBinding(t *testing.T) {
	t.Parallel()

	publishers := (&provider.CliServerServiceProvider{}).Publishers()

	if publishers[constant.InputHandlerContractServiceID] == nil {
		t.Error("the provider must defer the input handler, but deferred none")
	}
}

func TestTheProviderBindsTheInputHandler(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	built := collection.NewCollection()

	container.SetSingleton(applicationconstant.CliConfigContractServiceID, applicationdata.NewCliConfig())
	container.SetSingleton(
		interactionconstant.OutputFactoryContractServiceID,
		interactionfactory.NewOutputFactory(nil),
	)
	container.SetSingleton(routingconstant.RouterContractServiceID, routingdispatcher.NewRouter(
		container,
		built,
		interactionfactory.NewOutputFactory(nil),
		middlewarehandler.NewRouteMatchedHandler(container),
		middlewarehandler.NewRouteNotMatchedHandler(container),
		middlewarehandler.NewRouteDispatchedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewProcessExitingHandler(container),
	))

	provider.PublishInputHandler(container)

	resolved, err := container.GetSingleton(constant.InputHandlerContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the input handler, but reported: %v", err)
	}

	if _, isHandler := resolved.(clicontract.InputHandlerContract); !isHandler {
		t.Fatal("the container must hold an input handler, but held another value")
	}
}

func TestTheProviderBindsTheInputHandlerWhereTheApplicationPublishesNothing(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishInputHandler(container)

	_, err := container.GetSingleton(constant.InputHandlerContractServiceID)
	if err != nil {
		t.Errorf("the provider must bind an input handler even so, but reported: %v", err)
	}
}

func TestTheProviderIgnoresAServiceOfAnotherTypeUnderItsBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// A binding key that holds a value of another type is the developer's
	// error. The provider builds what it can rather than ending the process.
	container.SetSingleton(applicationconstant.CliConfigContractServiceID, "not a config")
	container.SetSingleton(routingconstant.RouterContractServiceID, "not a router")
	container.SetSingleton(interactionconstant.OutputFactoryContractServiceID, "not a factory")

	provider.PublishInputHandler(container)

	_, err := container.GetSingleton(constant.InputHandlerContractServiceID)
	if err != nil {
		t.Errorf("the provider must bind an input handler even so, but reported: %v", err)
	}
}
