/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package kernel_test

import (
	"os"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/application/data"
	"github.com/valkyrjaio/valkyrja-go/v26/application/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/application/kernel"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

const (
	testEnvironment = "testing"
	testVersion     = "26.1.2"
	testName        = "first"
)

// newApplication builds an application over an empty container and a config that
// holds the providers.
func newApplication(providers ...contract.ComponentProviderContract) *kernel.Valkyrja {
	return kernel.NewValkyrja(manager.NewContainer(nil), data.NewConfig(providers...))
}

func TestGetContainerReturnsTheContainer(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	app := kernel.NewValkyrja(container, data.NewConfig())

	if app.GetContainer() != contract.ApplicationContract(app).GetContainer() {
		t.Error("GetContainer must return one container, but returned two")
	}
}

func TestNewValkyrjaSetsTheTimezone(t *testing.T) {
	// The test writes the process environment, so it does not run in parallel
	// with a test that reads it.
	config := data.NewConfig()
	config.Timezone = "America/Phoenix"

	kernel.NewValkyrja(manager.NewContainer(nil), config)

	if os.Getenv("TZ") != "America/Phoenix" {
		t.Errorf("the application must set the timezone, but it is: %s", os.Getenv("TZ"))
	}
}

func TestGetProvidersWalksWhatEachProviderNames(t *testing.T) {
	t.Parallel()

	leaf := &fixtures.ComponentProviderFixture{}
	branch := &fixtures.ComponentProviderFixture{
		ComponentProviders: []contract.ComponentProviderContract{leaf},
	}

	app := newApplication(branch)

	providers := app.GetProviders()

	if len(providers) != 2 {
		t.Fatalf("GetProviders must return each provider, but returned: %d", len(providers))
	}

	if providers[0] != contract.ComponentProviderContract(leaf) {
		t.Error("GetProviders must return what a provider names before the provider, but did not")
	}
}

func TestGetProvidersCollectsAProviderOnce(t *testing.T) {
	t.Parallel()

	shared := &fixtures.ComponentProviderFixture{}
	first := &fixtures.ComponentProviderFixture{
		ComponentProviders: []contract.ComponentProviderContract{shared},
	}
	second := &fixtures.ComponentProviderFixture{
		ComponentProviders: []contract.ComponentProviderContract{shared},
	}

	app := newApplication(first, second)

	if len(app.GetProviders()) != 3 {
		t.Errorf("a provider that two providers name must be collected once, but got: %d",
			len(app.GetProviders()))
	}
}

func TestGetProvidersStopsOnAProviderThatNamesItself(t *testing.T) {
	t.Parallel()

	looping := &fixtures.ComponentProviderFixture{}
	looping.ComponentProviders = []contract.ComponentProviderContract{looping}

	app := newApplication(looping)

	if len(app.GetProviders()) != 1 {
		t.Errorf("a provider that names itself must be collected once, but got: %d",
			len(app.GetProviders()))
	}
}

func TestGetProvidersWalksOnlyOnce(t *testing.T) {
	t.Parallel()

	app := newApplication(&fixtures.ComponentProviderFixture{})

	first := app.GetProviders()
	second := app.GetProviders()

	if len(first) != len(second) {
		t.Error("GetProviders must return the same providers on a second call, but did not")
	}
}

func TestEachProviderKindIsCollectedOnce(t *testing.T) {
	t.Parallel()

	serviceProvider := &fixtures.ServiceProviderFixture{Name: testName}
	listenerProvider := &fixtures.ListenerProviderFixture{Name: testName}
	cliProvider := &fixtures.CliRouteProviderFixture{Name: testName}
	httpProvider := &fixtures.HttpRouteProviderFixture{Name: testName}

	shared := &fixtures.ComponentProviderFixture{}
	shared.ContainerProviders = append(shared.ContainerProviders, serviceProvider)
	shared.EventProviders = append(shared.EventProviders, listenerProvider)
	shared.CliProviders = append(shared.CliProviders, cliProvider)
	shared.HttpProviders = append(shared.HttpProviders, httpProvider)

	other := &fixtures.ComponentProviderFixture{}
	other.ContainerProviders = append(other.ContainerProviders, serviceProvider)
	other.EventProviders = append(other.EventProviders, listenerProvider)
	other.CliProviders = append(other.CliProviders, cliProvider)
	other.HttpProviders = append(other.HttpProviders, httpProvider)

	app := newApplication(shared, other)

	if len(app.GetContainerProviders()) != 1 {
		t.Errorf("a service provider that two components name must be collected once, but got: %d",
			len(app.GetContainerProviders()))
	}

	if len(app.GetEventProviders()) != 1 {
		t.Errorf("a listener provider that two components name must be collected once, but got: %d",
			len(app.GetEventProviders()))
	}

	if len(app.GetCliProviders()) != 1 {
		t.Errorf("a CLI provider that two components name must be collected once, but got: %d",
			len(app.GetCliProviders()))
	}

	if len(app.GetHttpProviders()) != 1 {
		t.Errorf("an HTTP provider that two components name must be collected once, but got: %d",
			len(app.GetHttpProviders()))
	}
}

func TestEachProviderKindWalksOnlyOnce(t *testing.T) {
	t.Parallel()

	component := &fixtures.ComponentProviderFixture{}
	component.ContainerProviders = append(component.ContainerProviders, &fixtures.ServiceProviderFixture{})
	component.EventProviders = append(component.EventProviders, &fixtures.ListenerProviderFixture{})
	component.CliProviders = append(component.CliProviders, &fixtures.CliRouteProviderFixture{})
	component.HttpProviders = append(component.HttpProviders, &fixtures.HttpRouteProviderFixture{})

	app := newApplication(component)

	// Walk once, then give the component a second provider of each kind. A
	// second walk would find them; the held result does not.
	app.GetContainerProviders()
	app.GetEventProviders()
	app.GetCliProviders()
	app.GetHttpProviders()

	component.ContainerProviders = append(component.ContainerProviders, &fixtures.ServiceProviderFixture{})
	component.EventProviders = append(component.EventProviders, &fixtures.ListenerProviderFixture{})
	component.CliProviders = append(component.CliProviders, &fixtures.CliRouteProviderFixture{})
	component.HttpProviders = append(component.HttpProviders, &fixtures.HttpRouteProviderFixture{})

	if len(app.GetContainerProviders()) != 1 {
		t.Errorf("GetContainerProviders must walk only once, but returned: %d",
			len(app.GetContainerProviders()))
	}

	if len(app.GetEventProviders()) != 1 {
		t.Errorf("GetEventProviders must walk only once, but returned: %d", len(app.GetEventProviders()))
	}

	if len(app.GetCliProviders()) != 1 {
		t.Errorf("GetCliProviders must walk only once, but returned: %d", len(app.GetCliProviders()))
	}

	if len(app.GetHttpProviders()) != 1 {
		t.Errorf("GetHttpProviders must walk only once, but returned: %d", len(app.GetHttpProviders()))
	}
}

func TestPublishProviderCallbacksRunsEachCallback(t *testing.T) {
	t.Parallel()

	ran := 0
	config := data.NewConfig()
	config.Callbacks = []contract.ApplicationCallbackFunc{
		func(_ contract.ApplicationContract) { ran++ },
		func(_ contract.ApplicationContract) { ran++ },
	}

	kernel.NewValkyrja(manager.NewContainer(nil), config).PublishProviderCallbacks()

	if ran != 2 {
		t.Errorf("PublishProviderCallbacks must run each callback, but ran: %d", ran)
	}
}

func TestTheApplicationReadsItsConfig(t *testing.T) {
	t.Parallel()

	config := data.NewConfig()
	config.DebugMode = true
	config.Environment = testEnvironment
	config.Version = testVersion

	app := kernel.NewValkyrja(manager.NewContainer(nil), config)

	if !app.GetDebugMode() {
		t.Error("GetDebugMode must read the config, but is false")
	}

	if app.GetEnvironment() != testEnvironment {
		t.Errorf("GetEnvironment must read the config, but is: %s", app.GetEnvironment())
	}

	if app.GetVersion() != testVersion {
		t.Errorf("GetVersion must read the config, but is: %s", app.GetVersion())
	}
}
