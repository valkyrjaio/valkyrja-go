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
	"github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/container/provider"
)

func TestThePublishersNameTheContainerData(t *testing.T) {
	t.Parallel()

	publishers := (&provider.ContainerServiceProvider{}).Publishers()

	if len(publishers) != 1 {
		t.Fatalf("Publishers must defer one binding key, but defers: %d", len(publishers))
	}

	if _, found := publishers[constant.ContainerDataServiceID]; !found {
		t.Error("Publishers must defer the container data, but does not")
	}
}

func TestPublishContainerDataRegistersEachServiceProvider(t *testing.T) {
	t.Parallel()

	component := &applicationfixtures.ComponentProviderFixture{}
	component.ContainerProviders = append(component.ContainerProviders, &providerFixture{})

	container := manager.NewContainer(nil)
	app := kernel.NewValkyrja(container, applicationdata.NewConfig(component))
	container.SetSingleton(applicationconstant.ApplicationContractServiceID, app)

	provider.PublishContainerData(container)

	if !container.IsDeferred(providedID) {
		t.Error("the publisher must register each service provider, but did not")
	}

	if !container.IsSingletonInstance(constant.ContainerDataServiceID) {
		t.Error("the publisher must set the container data, but did not")
	}
}

func TestPublishContainerDataStopsWhereNoApplicationIsBound(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishContainerData(container)

	if container.IsSingletonInstance(constant.ContainerDataServiceID) {
		t.Error("the publisher must set nothing where no application is bound, but set the data")
	}
}

func TestPublishContainerDataStopsWhereTheBindingIsNoApplication(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(applicationconstant.ApplicationContractServiceID, &providerFixture{})

	provider.PublishContainerData(container)

	if container.IsSingletonInstance(constant.ContainerDataServiceID) {
		t.Error("the publisher must set nothing where the binding is no application, but set the data")
	}
}

func TestTheContainerProviderNamesItsOwnServiceProvider(t *testing.T) {
	t.Parallel()

	app := kernel.NewValkyrja(manager.NewContainer(nil), applicationdata.NewConfig())
	componentProvider := &provider.ContainerComponentProvider{}

	serviceProviders := componentProvider.GetContainerProviders(app)

	if len(serviceProviders) != 1 {
		t.Fatalf("the provider must name one service provider, but named: %d", len(serviceProviders))
	}

	if _, isOwn := serviceProviders[0].(*provider.ContainerServiceProvider); !isOwn {
		t.Errorf("the provider must name its own service provider, but named: %T", serviceProviders[0])
	}
}

func TestTheContainerProviderNamesNoOtherProvider(t *testing.T) {
	t.Parallel()

	app := kernel.NewValkyrja(manager.NewContainer(nil), applicationdata.NewConfig())
	componentProvider := &provider.ContainerComponentProvider{}

	if len(componentProvider.GetComponentProviders(app)) != 0 {
		t.Error("the provider must name no component, but named one")
	}

	if len(componentProvider.GetEventProviders(app)) != 0 {
		t.Error("the provider must name no listener provider, but named one")
	}

	if len(componentProvider.GetCliProviders(app)) != 0 {
		t.Error("the provider must name no CLI provider, but named one")
	}

	if len(componentProvider.GetHttpProviders(app)) != 0 {
		t.Error("the provider must name no HTTP provider, but named one")
	}
}

// providedID is the binding key that providerFixture defers.
const providedID = "valkyrja.tests.container.ProvidedID"

// providerFixture is a service provider that defers one binding key. It stands
type providerFixture struct{}

// Publishers returns a publisher for the binding key that the provider defers.
func (p *providerFixture) Publishers() map[string]contract.PublishFunc {
	return map[string]contract.PublishFunc{
		providedID: func(_ contract.ContainerContract) {},
	}
}

// The fixture also stands in for a binding that is no application, which the
// publisher must not read as one.
var _ applicationcontract.ComponentProviderContract = (*applicationfixtures.ComponentProviderFixture)(nil)
