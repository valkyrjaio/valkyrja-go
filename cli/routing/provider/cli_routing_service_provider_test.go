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
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	applicationdata "github.com/valkyrjaio/valkyrja-go/v26/application/data"
	applicationfixtures "github.com/valkyrjaio/valkyrja-go/v26/application/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/application/kernel"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	interactionfactory "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/provider"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

const routeName = "cache:clear"

// routeProviderFixture registers one command.
type routeProviderFixture struct{}

// GetRoutes returns the one command that the fixture registers.
func (p *routeProviderFixture) GetRoutes() []clicontract.RouteContract {
	return []clicontract.RouteContract{
		data.NewRoute(routeName, "Clear the cache", func(
			_ containercontract.ContainerContract,
			_ clicontract.RouteContract,
		) clicontract.OutputContract {
			return nil
		}),
	}
}

func TestTheProviderDefersEachOfItsBindings(t *testing.T) {
	t.Parallel()

	publishers := (&provider.CliRoutingServiceProvider{}).Publishers()

	for _, id := range []string{
		constant.RouteCollectionContractServiceID,
		constant.RouterContractServiceID,
	} {
		if publishers[id] == nil {
			t.Errorf("the provider must defer %q, but deferred none", id)
		}
	}
}

func TestTheProviderFilesEveryCommandThatARouteProviderRegisters(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(applicationconstant.ApplicationContractServiceID, newApplication())

	provider.PublishRouteCollection(container)

	if !readCollection(t, container).Has(routeName) {
		t.Error("the provider must file every command that a route provider registers, but did not")
	}
}

func TestTheProviderFilesNoCommandWhereTheContainerHoldsNoApplication(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishRouteCollection(container)

	if len(readCollection(t, container).All()) != 0 {
		t.Error("a container with no application must file no command, but filed some")
	}
}

func TestTheProviderBindsTheRouter(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(applicationconstant.ApplicationContractServiceID, newApplication())
	container.SetSingleton(applicationconstant.CliConfigContractServiceID, applicationdata.NewCliConfig())
	container.SetSingleton(
		interactionconstant.OutputFactoryContractServiceID,
		interactionfactory.NewOutputFactory(nil),
	)

	provider.PublishRouteCollection(container)
	provider.PublishRouter(container)

	resolved, err := container.GetSingleton(constant.RouterContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the router, but reported: %v", err)
	}

	if _, isRouter := resolved.(clicontract.RouterContract); !isRouter {
		t.Fatal("the container must hold a router, but held another value")
	}
}

func TestTheProviderBindsTheRouterWhereTheApplicationPublishesNothing(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishRouter(container)

	_, err := container.GetSingleton(constant.RouterContractServiceID)
	if err != nil {
		t.Errorf("the provider must bind a router even so, but reported: %v", err)
	}
}

// newApplication builds an application that registers the route provider above.
func newApplication() applicationcontract.ApplicationContract {
	component := &applicationfixtures.ComponentProviderFixture{}
	component.CliProviders = append(component.CliProviders, &routeProviderFixture{})

	return kernel.NewValkyrja(manager.NewContainer(nil), applicationdata.NewConfig(component))
}

// readCollection returns the collection that the provider bound.
func readCollection(
	t *testing.T,
	container containercontract.ContainerContract,
) clicontract.RouteCollectionContract {
	t.Helper()

	resolved, err := container.GetSingleton(constant.RouteCollectionContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the collection, but reported: %v", err)
	}

	built, isCollection := resolved.(clicontract.RouteCollectionContract)
	if !isCollection {
		t.Fatal("the container must hold a collection, but held another value")
	}

	return built
}

func TestTheProviderIgnoresAServiceOfAnotherTypeUnderItsBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// A binding key that holds a value of another type is the developer's
	// error. The provider builds what it can rather than ending the process.
	container.SetSingleton(applicationconstant.ApplicationContractServiceID, "not an application")
	container.SetSingleton(applicationconstant.CliConfigContractServiceID, "not a config")
	container.SetSingleton(constant.RouteCollectionContractServiceID, "not a collection")
	container.SetSingleton(interactionconstant.OutputFactoryContractServiceID, "not a factory")

	provider.PublishRouter(container)

	resolved, err := container.GetSingleton(constant.RouterContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind a router even so, but reported: %v", err)
	}

	if _, isRouter := resolved.(clicontract.RouterContract); !isRouter {
		t.Fatal("the container must hold a router, but held another value")
	}

	provider.PublishRouteCollection(container)

	if len(readCollection(t, container).All()) != 0 {
		t.Error("an application of another type must register no command, but registered some")
	}
}
