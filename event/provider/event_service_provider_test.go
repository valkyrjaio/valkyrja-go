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
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/event/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/data"
	"github.com/valkyrjaio/valkyrja-go/v26/event/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/event/provider"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

const (
	eventID      = "valkyrja.tests.event.Event"
	listenerName = "valkyrja.tests.event.Listener"
)

// applicationFixture is an application that holds the listener providers a test
// gives it.
type applicationFixture struct {
	eventProviders []contract.ListenerProviderContract
}

func (a *applicationFixture) GetContainer() containercontract.ContainerContract { return nil }

func (a *applicationFixture) PublishProviderCallbacks() {}

func (a *applicationFixture) GetProviders() []applicationcontract.ComponentProviderContract {
	return nil
}

func (a *applicationFixture) GetContainerProviders() []containercontract.ServiceProviderContract {
	return nil
}

func (a *applicationFixture) GetEventProviders() []contract.ListenerProviderContract {
	return a.eventProviders
}

func (a *applicationFixture) GetCliProviders() []clicontract.CliRouteProviderContract { return nil }

func (a *applicationFixture) GetHttpProviders() []httpcontract.HttpRouteProviderContract {
	return nil
}

func (a *applicationFixture) GetDebugMode() bool { return false }

func (a *applicationFixture) GetEnvironment() string { return "" }

func (a *applicationFixture) GetVersion() string { return "" }

// newListener builds a listener for the test event.
func newListener() contract.ListenerContract {
	return data.NewListener(eventID, listenerName, func(
		_ containercontract.ContainerContract,
		_ map[string]any,
	) any {
		return nil
	})
}

// newContainerWithApplication returns a container that holds the application.
func newContainerWithApplication(app any) containercontract.ContainerContract {
	container := manager.NewContainer(nil)
	container.SetSingleton(applicationconstant.ApplicationContractServiceID, app)

	return container
}

func TestThePublishersCoverEveryBindingKeyOfTheComponent(t *testing.T) {
	t.Parallel()

	publishers := (&provider.EventServiceProvider{}).Publishers()

	ids := []string{
		constant.ListenerCollectionContractServiceID,
		constant.EventDispatcherContractServiceID,
		constant.EventDataServiceID,
	}

	if len(publishers) != len(ids) {
		t.Errorf("the provider must defer %d keys, but defers: %d", len(ids), len(publishers))
	}

	for _, id := range ids {
		if publishers[id] == nil {
			t.Errorf("the provider must defer %q, but deferred none", id)
		}
	}
}

func TestPublishListenerCollectionFilesEveryListenerThatAProviderRegisters(t *testing.T) {
	t.Parallel()

	container := newContainerWithApplication(&applicationFixture{
		eventProviders: []contract.ListenerProviderContract{
			&fixtures.ListenerProviderFixture{Listeners: []contract.ListenerContract{newListener()}},
		},
	})

	provider.PublishListenerCollection(container)

	resolved, err := container.GetSingleton(constant.ListenerCollectionContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the collection, but reported: %v", err)
	}

	collected, isCollection := resolved.(contract.ListenerCollectionContract)
	if !isCollection {
		t.Fatal("the provider must bind a listener collection, but bound another type")
	}

	if !collected.HasListenerByID(listenerName) {
		t.Error("the provider must file the listener that the provider registers, but did not")
	}
}

func TestPublishListenerCollectionBindsAnEmptyCollectionWhereNoApplicationResolves(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishListenerCollection(container)

	resolved, err := container.GetSingleton(constant.ListenerCollectionContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the collection, but reported: %v", err)
	}

	collected, isCollection := resolved.(contract.ListenerCollectionContract)
	if !isCollection {
		t.Fatal("the provider must bind a listener collection, but bound another type")
	}

	if collected.HasListenerByID(listenerName) {
		t.Error("the provider must file no listener where no application resolves, but filed one")
	}
}

func TestPublishListenerCollectionBindsAnEmptyCollectionWhereTheApplicationIsAnotherType(t *testing.T) {
	t.Parallel()

	container := newContainerWithApplication(&fixtures.NotAnEventFixture{})

	provider.PublishListenerCollection(container)

	_, err := container.GetSingleton(constant.ListenerCollectionContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the collection, but reported: %v", err)
	}
}

func TestPublishEventDispatcherBindsTheDispatcherOverThePublishedCollection(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishListenerCollection(container)
	provider.PublishEventDispatcher(container)

	resolved, err := container.GetSingleton(constant.EventDispatcherContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the dispatcher, but reported: %v", err)
	}

	if _, isDispatcher := resolved.(contract.EventDispatcherContract); !isDispatcher {
		t.Error("the provider must bind an event dispatcher, but bound another type")
	}
}

func TestPublishEventDispatcherBindsADispatcherWhereNoCollectionResolves(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishEventDispatcher(container)

	_, err := container.GetSingleton(constant.EventDispatcherContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the dispatcher, but reported: %v", err)
	}
}

func TestPublishEventDispatcherBindsADispatcherWhereTheCollectionIsAnotherType(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(constant.ListenerCollectionContractServiceID, &fixtures.NotAnEventFixture{})

	provider.PublishEventDispatcher(container)

	_, err := container.GetSingleton(constant.EventDispatcherContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the dispatcher, but reported: %v", err)
	}
}

func TestPublishEventDataBindsTheStateOfTheCollection(t *testing.T) {
	t.Parallel()

	container := newContainerWithApplication(&applicationFixture{
		eventProviders: []contract.ListenerProviderContract{
			&fixtures.ListenerProviderFixture{Listeners: []contract.ListenerContract{newListener()}},
		},
	})

	provider.PublishListenerCollection(container)
	provider.PublishEventData(container)

	resolved, err := container.GetSingleton(constant.EventDataServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the data, but reported: %v", err)
	}

	eventData, isData := resolved.(contract.EventDataContract)
	if !isData {
		t.Fatal("the provider must bind the component's data, but bound another type")
	}

	if len(eventData.GetListeners()) == 0 {
		t.Error("the data must hold the listener that the collection filed, but held none")
	}
}

func TestTheComponentProviderRegistersTheServiceProviderAndNothingElse(t *testing.T) {
	t.Parallel()

	component := &provider.EventComponentProvider{}

	if len(component.GetContainerProviders(nil)) != 1 {
		t.Errorf("the component must register one service provider, but registered: %d",
			len(component.GetContainerProviders(nil)))
	}

	if len(component.GetComponentProviders(nil)) != 0 {
		t.Error("the component must need no other component, but named one")
	}

	if len(component.GetEventProviders(nil)) != 0 {
		t.Error("the component must register no listener provider, but registered one")
	}

	if len(component.GetCliProviders(nil)) != 0 {
		t.Error("the component must register no CLI route provider, but registered one")
	}

	if len(component.GetHttpProviders(nil)) != 0 {
		t.Error("the component must register no HTTP route provider, but registered one")
	}
}

func TestTheProvidersSatisfyTheirContracts(t *testing.T) {
	t.Parallel()

	var service containercontract.ServiceProviderContract = &provider.EventServiceProvider{}
	var component applicationcontract.ComponentProviderContract = &provider.EventComponentProvider{}

	if len(service.Publishers()) == 0 {
		t.Error("the service provider must defer a binding key, but deferred none")
	}

	if component.GetContainerProviders(nil) == nil {
		t.Error("the component provider must register a service provider, but registered none")
	}
}
