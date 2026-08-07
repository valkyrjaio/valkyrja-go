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
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/provider"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

// quietConfigFixture is an application configuration that also states the
type quietConfigFixture struct {
	*applicationdata.Config
}

// IsQuiet reports that every output writes less.
func (c *quietConfigFixture) IsQuiet() bool {
	return true
}

// IsInteractive reports that no output asks the caller a question.
func (c *quietConfigFixture) IsInteractive() bool {
	return false
}

// IsSilent reports that every output writes something.
func (c *quietConfigFixture) IsSilent() bool {
	return false
}

func TestTheProviderDefersEachOfItsBindings(t *testing.T) {
	t.Parallel()

	publishers := (&provider.CliInteractionServiceProvider{}).Publishers()

	for _, id := range []string{
		constant.CliInteractionConfigContractServiceID,
		constant.OutputFactoryContractServiceID,
	} {
		if publishers[id] == nil {
			t.Errorf("the provider must defer %q, but deferred none", id)
		}
	}
}

func TestTheProviderBindsTheFrameworkDefaultWhereTheApplicationStatesNone(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishCliInteractionConfig(container)

	config := readConfig(t, container)

	if !config.IsInteractive() || config.IsQuiet() || config.IsSilent() {
		t.Error("the framework default must be interactive and write every message, but was not")
	}
}

func TestTheProviderBindsTheApplicationConfigurationWhereItHoldsTheSettings(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(
		applicationconstant.ConfigContractServiceID,
		&quietConfigFixture{Config: applicationdata.NewConfig()},
	)

	provider.PublishCliInteractionConfig(container)

	if !readConfig(t, container).IsQuiet() {
		t.Error("the application's own configuration must be bound where it holds the settings, but was not")
	}
}

func TestTheProviderBindsTheFrameworkDefaultWhereTheConfigurationHoldsNoSettings(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(applicationconstant.ConfigContractServiceID, applicationdata.NewConfig())

	provider.PublishCliInteractionConfig(container)

	if !readConfig(t, container).IsInteractive() {
		t.Error("a configuration that holds no setting must fall back to the framework default, but did not")
	}
}

func TestTheProviderBindsTheOutputFactory(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(
		applicationconstant.ConfigContractServiceID,
		&quietConfigFixture{Config: applicationdata.NewConfig()},
	)

	provider.PublishOutputFactory(container)

	resolved, err := container.GetSingleton(constant.OutputFactoryContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the output factory, but reported: %v", err)
	}

	built, isFactory := resolved.(clicontract.OutputFactoryContract)
	if !isFactory {
		t.Fatal("the container must hold an output factory, but held another value")
	}

	if !built.CreateOutput(constant.ExitCodeSuccess).IsQuiet() {
		t.Error("the factory must carry the settings of the configuration, but did not")
	}
}

// readConfig returns the configuration that the provider bound.
func readConfig(
	t *testing.T,
	container containercontract.ContainerContract,
) clicontract.CliInteractionConfigContract {
	t.Helper()

	resolved, err := container.GetSingleton(constant.CliInteractionConfigContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the configuration, but reported: %v", err)
	}

	config, isConfig := resolved.(clicontract.CliInteractionConfigContract)
	if !isConfig {
		t.Fatal("the container must hold a CLI interaction configuration, but held another value")
	}

	return config
}
