/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the HTTP component's top-level provider.
package provider

import (
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	containerprovider "github.com/valkyrjaio/valkyrja-go/v26/container/provider"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	messageprovider "github.com/valkyrjaio/valkyrja-go/v26/http/message/provider"
	routingprovider "github.com/valkyrjaio/valkyrja-go/v26/http/routing/provider"
	serverprovider "github.com/valkyrjaio/valkyrja-go/v26/http/server/provider"
)

// HttpComponentProvider is the HTTP component's top-level provider.
type HttpComponentProvider struct{}

// GetComponentProviders returns each component that HTTP needs.
func (p *HttpComponentProvider) GetComponentProviders(
	_ applicationcontract.ApplicationContract,
) []applicationcontract.ComponentProviderContract {
	return []applicationcontract.ComponentProviderContract{
		&containerprovider.ContainerComponentProvider{},
	}
}

// GetContainerProviders returns each service provider of HTTP.
func (p *HttpComponentProvider) GetContainerProviders(
	_ applicationcontract.ApplicationContract,
) []containercontract.ServiceProviderContract {
	return []containercontract.ServiceProviderContract{
		&messageprovider.HttpMessageServiceProvider{},
		&routingprovider.HttpRoutingServiceProvider{},
		&serverprovider.HttpServerServiceProvider{},
	}
}

// GetEventProviders returns each listener provider of HTTP.
func (p *HttpComponentProvider) GetEventProviders(
	_ applicationcontract.ApplicationContract,
) []eventcontract.ListenerProviderContract {
	return []eventcontract.ListenerProviderContract{}
}

// GetCliProviders returns each CLI route provider of HTTP.
func (p *HttpComponentProvider) GetCliProviders(
	_ applicationcontract.ApplicationContract,
) []clicontract.CliRouteProviderContract {
	return []clicontract.CliRouteProviderContract{}
}

// GetHttpProviders returns each HTTP route provider of HTTP.
func (p *HttpComponentProvider) GetHttpProviders(
	_ applicationcontract.ApplicationContract,
) []httpcontract.HttpRouteProviderContract {
	return []httpcontract.HttpRouteProviderContract{}
}
