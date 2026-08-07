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

	"github.com/valkyrjaio/valkyrja-go/v26/application/data"
	"github.com/valkyrjaio/valkyrja-go/v26/application/kernel"
	"github.com/valkyrjaio/valkyrja-go/v26/application/provider"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	containerprovider "github.com/valkyrjaio/valkyrja-go/v26/container/provider"
)

func TestTheApplicationProviderNamesTheContainer(t *testing.T) {
	t.Parallel()

	app := kernel.NewValkyrja(manager.NewContainer(nil), data.NewConfig())
	componentProvider := &provider.ApplicationComponentProvider{}

	components := componentProvider.GetComponentProviders(app)

	if len(components) != 1 {
		t.Fatalf("the provider must name one component, but named: %d", len(components))
	}

	if _, isContainer := components[0].(*containerprovider.ContainerComponentProvider); !isContainer {
		t.Errorf("the provider must name the container, but named: %T", components[0])
	}
}

func TestTheApplicationProviderNamesNoOtherProvider(t *testing.T) {
	t.Parallel()

	app := kernel.NewValkyrja(manager.NewContainer(nil), data.NewConfig())
	componentProvider := &provider.ApplicationComponentProvider{}

	if len(componentProvider.GetContainerProviders(app)) != 0 {
		t.Error("the provider must name no service provider, but named one")
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
