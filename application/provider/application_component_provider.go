/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the application component's providers.
package provider

import (
	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	containerprovider "github.com/valkyrjaio/valkyrja-go/v26/container/provider"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// ApplicationComponentProvider is the application's top-level provider. It names
// the container, which every application needs.
//
// The other ports declare a CLI variant and an HTTP variant that override
// `getComponentProviders`. Go has no method override, so each variant is its own
// provider once the CLI and the HTTP components land.
type ApplicationComponentProvider struct{}

// GetComponentProviders returns each component that the application needs.
func (p *ApplicationComponentProvider) GetComponentProviders(
	_ contract.ApplicationContract,
) []contract.ComponentProviderContract {
	return []contract.ComponentProviderContract{&containerprovider.ContainerComponentProvider{}}
}

// GetContainerProviders returns each service provider of the application.
func (p *ApplicationComponentProvider) GetContainerProviders(
	_ contract.ApplicationContract,
) []containercontract.ServiceProviderContract {
	return []containercontract.ServiceProviderContract{}
}

// GetEventProviders returns each listener provider of the application.
func (p *ApplicationComponentProvider) GetEventProviders(
	_ contract.ApplicationContract,
) []eventcontract.ListenerProviderContract {
	return []eventcontract.ListenerProviderContract{}
}

// GetCliProviders returns each CLI route provider of the application.
func (p *ApplicationComponentProvider) GetCliProviders(
	_ contract.ApplicationContract,
) []clicontract.CliRouteProviderContract {
	return []clicontract.CliRouteProviderContract{}
}

// GetHttpProviders returns each HTTP route provider of the application.
func (p *ApplicationComponentProvider) GetHttpProviders(
	_ contract.ApplicationContract,
) []httpcontract.HttpRouteProviderContract {
	return []httpcontract.HttpRouteProviderContract{}
}
