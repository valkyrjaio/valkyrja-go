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
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	messageconstant "github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	messagefactory "github.com/valkyrjaio/valkyrja-go/v26/http/message/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/provider"
)

const routePath = "/users/{id}"

// routeProviderFixture registers one route with a parameter in its path, so the
type routeProviderFixture struct{}

// GetRoutes returns the one route that the fixture registers.
func (p *routeProviderFixture) GetRoutes() []contract.RouteContract {
	return []contract.RouteContract{
		data.NewRoute(routePath, "users.show", func(
			_ containercontract.ContainerContract,
			_ contract.RouteContract,
		) contract.ResponseContract {
			return nil
		}).WithParameters(data.NewParameter("id", constant.RegexNum)),
	}
}

// newApplication builds an application that registers the route provider above.
func newApplication() applicationcontract.ApplicationContract {
	component := &applicationfixtures.ComponentProviderFixture{}
	component.HttpProviders = append(component.HttpProviders, &routeProviderFixture{})

	return kernel.NewValkyrja(manager.NewContainer(nil), applicationdata.NewConfig(component))
}

// newPublishedContainer builds a container that every binding of the
// sub-component is published into.
func newPublishedContainer(withApplication bool) containercontract.ContainerContract {
	container := manager.NewContainer(nil)

	if withApplication {
		container.SetSingleton(applicationconstant.ApplicationContractServiceID, newApplication())
		container.SetSingleton(applicationconstant.HttpConfigContractServiceID, applicationdata.NewHttpConfig())
	}

	provider.PublishProcessor(container)
	provider.PublishRouteCollection(container)
	provider.PublishMatcher(container)
	provider.PublishUrl(container)
	provider.PublishRoutingResponseFactory(container)
	provider.PublishRouter(container)

	return container
}

func TestTheProviderDefersEachOfItsBindings(t *testing.T) {
	t.Parallel()

	publishers := (&provider.HttpRoutingServiceProvider{}).Publishers()

	for _, id := range []string{
		constant.ProcessorContractServiceID,
		constant.RouteCollectionContractServiceID,
		constant.MatcherContractServiceID,
		constant.UrlContractServiceID,
		constant.RoutingResponseFactoryContractServiceID,
		constant.RouterContractServiceID,
	} {
		if publishers[id] == nil {
			t.Errorf("the provider must defer %q, but deferred none", id)
		}
	}
}

func TestTheProviderBindsEveryServiceOfTheSubComponent(t *testing.T) {
	t.Parallel()

	container := newPublishedContainer(true)

	for _, id := range []string{
		constant.ProcessorContractServiceID,
		constant.RouteCollectionContractServiceID,
		constant.MatcherContractServiceID,
		constant.UrlContractServiceID,
		constant.RoutingResponseFactoryContractServiceID,
		constant.RouterContractServiceID,
	} {
		_, err := container.GetSingleton(id)
		if err != nil {
			t.Errorf("the provider must bind %q, but reported: %v", id, err)
		}
	}
}

func TestTheProviderFilesEveryRouteThatARouteProviderRegisters(t *testing.T) {
	t.Parallel()

	container := newPublishedContainer(true)

	resolved, err := container.GetSingleton(constant.MatcherContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the matcher, but reported: %v", err)
	}

	matcher, isMatcher := resolved.(contract.MatcherContract)
	if !isMatcher {
		t.Fatal("the container must hold a matcher, but held another value")
	}

	// The processor read the route before the collection filed it, so the
	// matcher reads the regular expression that the processor built.
	if matcher.Match("/users/1", messageconstant.RequestMethodGet) == nil {
		t.Error("the provider must file every route that a route provider registers, but did not")
	}
}

func TestTheProviderBindsEveryServiceWhereTheApplicationPublishesNothing(t *testing.T) {
	t.Parallel()

	container := newPublishedContainer(false)

	_, err := container.GetSingleton(constant.RouterContractServiceID)
	if err != nil {
		t.Errorf("the provider must bind a router even so, but reported: %v", err)
	}
}

func TestTheProviderIgnoresAServiceOfAnotherTypeUnderItsBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// A binding key that holds a value of another type is the developer's
	// error. The provider builds what it can rather than ending the process.
	container.SetSingleton(applicationconstant.ApplicationContractServiceID, "not an application")
	container.SetSingleton(applicationconstant.HttpConfigContractServiceID, "not a config")
	container.SetSingleton(constant.ProcessorContractServiceID, "not a processor")
	container.SetSingleton(constant.RouteCollectionContractServiceID, "not a collection")
	container.SetSingleton(constant.MatcherContractServiceID, "not a matcher")
	container.SetSingleton(constant.UrlContractServiceID, "not a url")
	container.SetSingleton(messageconstant.ResponseFactoryContractServiceID, "not a factory")

	// Each publisher reads what another one bound, so none of them may rebind a
	// key that this test states, or the value of another type never reaches it.
	provider.PublishRoutingResponseFactory(container)
	provider.PublishRouter(container)

	if len(readPaths(container)) != 0 {
		t.Error("an application of another type must register no route, but registered some")
	}

	resolved, err := container.GetSingleton(constant.RouterContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind a router even so, but reported: %v", err)
	}

	if _, isRouter := resolved.(contract.RouterContract); !isRouter {
		t.Fatal("the container must hold a router, but held another value")
	}
}

// readPaths returns the name of the route at each static path that the
// collection of the container holds.
func readPaths(container containercontract.ContainerContract) map[string]string {
	provider.PublishRouteCollection(container)

	resolved, err := container.GetSingleton(constant.RouteCollectionContractServiceID)
	if err != nil {
		return nil
	}

	built, isCollection := resolved.(contract.RouteCollectionContract)
	if !isCollection {
		return nil
	}

	return built.GetPaths(messageconstant.RequestMethodGet)
}

func TestEachPublisherBuildsItsOwnServiceWhereNoOtherPublisherRan(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// Each publisher reads what another one bound. A container that nothing was
	// published into is what the framework builds a default from.
	provider.PublishRouteCollection(container)
	provider.PublishMatcher(container)
	provider.PublishUrl(container)
	provider.PublishRoutingResponseFactory(container)
	provider.PublishRouter(container)

	for _, id := range []string{
		constant.RouteCollectionContractServiceID,
		constant.MatcherContractServiceID,
		constant.UrlContractServiceID,
		constant.RoutingResponseFactoryContractServiceID,
		constant.RouterContractServiceID,
	} {
		_, err := container.GetSingleton(id)
		if err != nil {
			t.Errorf("the provider must bind %q even so, but reported: %v", id, err)
		}
	}
}

func TestTheProviderReadsTheResponseFactoryThatTheMessageSubComponentPublished(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(
		messageconstant.ResponseFactoryContractServiceID,
		messagefactory.NewResponseFactory(),
	)

	provider.PublishRoutingResponseFactory(container)

	_, err := container.GetSingleton(constant.RoutingResponseFactoryContractServiceID)
	if err != nil {
		t.Errorf("the provider must bind the routing response factory, but reported: %v", err)
	}
}

func TestTheProviderBuildsItsOwnMatcherWhereTheContainerHoldsAnotherType(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// The router reads the matcher, and a binding key that holds a value of
	// another type is the developer's error. The router gets a matcher even so.
	container.SetSingleton(constant.MatcherContractServiceID, "not a matcher")

	provider.PublishRouter(container)

	_, err := container.GetSingleton(constant.RouterContractServiceID)
	if err != nil {
		t.Errorf("the provider must bind a router even so, but reported: %v", err)
	}

	// The same holds where nothing published a matcher at all.
	empty := manager.NewContainer(nil)

	provider.PublishRouter(empty)

	_, emptyErr := empty.GetSingleton(constant.RouterContractServiceID)
	if emptyErr != nil {
		t.Errorf("the provider must bind a router even so, but reported: %v", emptyErr)
	}
}
