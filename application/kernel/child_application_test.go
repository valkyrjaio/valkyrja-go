/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package kernel_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/application/data"
	"github.com/valkyrjaio/valkyrja-go/v26/application/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/application/kernel"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

func TestTheChildHoldsItsOwnContainer(t *testing.T) {
	t.Parallel()

	parentContainer := manager.NewContainer(nil)
	childContainer := manager.NewContainer(nil)

	parent := kernel.NewValkyrja(parentContainer, data.NewConfig())
	child := kernel.NewChildApplication(parent, childContainer)

	if child.GetContainer() == parent.GetContainer() {
		t.Error("the child must hold its own container, but holds the parent's")
	}
}

func TestTheChildReadsEverythingElseFromItsParent(t *testing.T) {
	t.Parallel()

	component := &fixtures.ComponentProviderFixture{}
	component.ContainerProviders = append(component.ContainerProviders, &fixtures.ServiceProviderFixture{})
	component.EventProviders = append(component.EventProviders, &fixtures.ListenerProviderFixture{})
	component.CliProviders = append(component.CliProviders, &fixtures.CliRouteProviderFixture{})
	component.HttpProviders = append(component.HttpProviders, &fixtures.HttpRouteProviderFixture{})

	ran := 0
	config := data.NewConfig(component)
	config.DebugMode = true
	config.Environment = testEnvironment
	config.Version = testVersion
	config.Callbacks = []contract.ApplicationCallbackFunc{func(_ contract.ApplicationContract) { ran++ }}

	parent := kernel.NewValkyrja(manager.NewContainer(nil), config)
	child := kernel.NewChildApplication(parent, manager.NewContainer(nil))

	child.PublishProviderCallbacks()

	if ran != 1 {
		t.Errorf("the child must run the parent's callbacks, but ran: %d", ran)
	}

	if len(child.GetProviders()) != 1 {
		t.Errorf("the child must read the parent's providers, but read: %d", len(child.GetProviders()))
	}

	if len(child.GetContainerProviders()) != 1 {
		t.Error("the child must read the parent's service providers, but did not")
	}

	if len(child.GetEventProviders()) != 1 {
		t.Error("the child must read the parent's listener providers, but did not")
	}

	if len(child.GetCliProviders()) != 1 {
		t.Error("the child must read the parent's CLI providers, but did not")
	}

	if len(child.GetHttpProviders()) != 1 {
		t.Error("the child must read the parent's HTTP providers, but did not")
	}

	if !child.GetDebugMode() {
		t.Error("the child must read the parent's debug mode, but is false")
	}

	if child.GetEnvironment() != testEnvironment {
		t.Errorf("the child must read the parent's environment, but is: %s", child.GetEnvironment())
	}

	if child.GetVersion() != testVersion {
		t.Errorf("the child must read the parent's version, but is: %s", child.GetVersion())
	}
}
