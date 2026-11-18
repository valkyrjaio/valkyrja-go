/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the CLI interaction sub-component's providers.
package provider

import (
	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

type CliInteractionServiceProvider struct{}

// Publishers returns a publisher for each binding key that the sub-component
// defers.
func (p *CliInteractionServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.CliInteractionConfigContractServiceID: PublishCliInteractionConfig,
		constant.OutputFactoryContractServiceID:        PublishOutputFactory,
	}
}

// PublishCliInteractionConfig binds the configuration that every output carries.
func PublishCliInteractionConfig(container containercontract.ContainerContract) {
	container.SetSingleton(constant.CliInteractionConfigContractServiceID, resolveConfig(container))
}

// PublishOutputFactory binds the factory that builds each kind of output.
func PublishOutputFactory(container containercontract.ContainerContract) {
	container.SetSingleton(
		constant.OutputFactoryContractServiceID,
		factory.NewOutputFactory(resolveConfig(container)),
	)
}

// resolveConfig returns the configuration that the application states, and the
// framework's default where the application states none.
func resolveConfig(container containercontract.ContainerContract) clicontract.CliInteractionConfigContract {
	resolved, err := container.GetSingleton(applicationconstant.ConfigContractServiceID)
	if err != nil {
		return data.NewCliInteractionConfig()
	}

	config, isConfig := resolved.(clicontract.CliInteractionConfigContract)
	if !isConfig {
		return data.NewCliInteractionConfig()
	}

	return config
}
