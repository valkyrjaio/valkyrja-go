/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the CLI component's top-level provider.
package provider

import (
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionprovider "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/provider"
	routingprovider "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/provider"
	serverprovider "github.com/valkyrjaio/valkyrja-go/v26/cli/server/provider"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	containerprovider "github.com/valkyrjaio/valkyrja-go/v26/container/provider"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

type CliComponentProvider struct{}

// GetComponentProviders returns each component that the CLI needs.
func (p *CliComponentProvider) GetComponentProviders(
	_ applicationcontract.ApplicationContract,
) []applicationcontract.ComponentProviderContract {
	return []applicationcontract.ComponentProviderContract{
		&containerprovider.ContainerComponentProvider{},
	}
}

// GetContainerProviders returns each service provider of the CLI.
func (p *CliComponentProvider) GetContainerProviders(
	_ applicationcontract.ApplicationContract,
) []containercontract.ServiceProviderContract {
	return []containercontract.ServiceProviderContract{
		&interactionprovider.CliInteractionServiceProvider{},
		&routingprovider.CliRoutingServiceProvider{},
		&serverprovider.CliServerServiceProvider{},
	}
}

// GetEventProviders returns each listener provider of the CLI.
func (p *CliComponentProvider) GetEventProviders(
	_ applicationcontract.ApplicationContract,
) []eventcontract.ListenerProviderContract {
	return []eventcontract.ListenerProviderContract{}
}

// GetCliProviders returns each CLI route provider of the CLI.
func (p *CliComponentProvider) GetCliProviders(
	_ applicationcontract.ApplicationContract,
) []clicontract.CliRouteProviderContract {
	return []clicontract.CliRouteProviderContract{
		&serverprovider.CliServerCliRoutesProvider{},
	}
}

// GetHttpProviders returns each HTTP route provider of the CLI.
func (p *CliComponentProvider) GetHttpProviders(
	_ applicationcontract.ApplicationContract,
) []httpcontract.HttpRouteProviderContract {
	return []httpcontract.HttpRouteProviderContract{}
}
