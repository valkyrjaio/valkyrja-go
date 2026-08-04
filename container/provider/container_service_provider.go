/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the container component's providers.
package provider

import (
	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// ContainerServiceProvider publishes the container component's own bindings.
type ContainerServiceProvider struct{}

// Publishers returns a publisher for each binding key that the component
// defers.
func (p *ContainerServiceProvider) Publishers() map[string]contract.PublishFunc {
	return map[string]contract.PublishFunc{
		constant.ContainerDataServiceID: PublishContainerData,
	}
}

// PublishContainerData registers every service provider that the application
// names, and sets the container's own state as a singleton.
//
// The publisher is a package-level function rather than a method, which is one
// of the two forms that `sindri` reads.
func PublishContainerData(container contract.ContainerContract) {
	resolved, err := container.GetSingleton(applicationconstant.ApplicationContractServiceID)
	if err != nil {
		return
	}

	app, isApplication := resolved.(applicationcontract.ApplicationContract)
	if !isApplication {
		return
	}

	for _, serviceProvider := range app.GetContainerProviders() {
		// A provider with no publisher is the developer's error, and the
		// container reports it where the application registers the provider.
		_ = container.Register(serviceProvider)
	}

	container.SetSingleton(constant.ContainerDataServiceID, container.GetData())
}

// ContainerComponentProvider is the container component's top-level provider.
type ContainerComponentProvider struct{}

// GetComponentProviders returns each component that the container needs.
func (p *ContainerComponentProvider) GetComponentProviders(
	_ applicationcontract.ApplicationContract,
) []applicationcontract.ComponentProviderContract {
	return []applicationcontract.ComponentProviderContract{}
}

// GetContainerProviders returns each service provider of the container.
func (p *ContainerComponentProvider) GetContainerProviders(
	_ applicationcontract.ApplicationContract,
) []contract.ServiceProviderContract {
	return []contract.ServiceProviderContract{&ContainerServiceProvider{}}
}

// GetEventProviders returns each listener provider of the container.
func (p *ContainerComponentProvider) GetEventProviders(
	_ applicationcontract.ApplicationContract,
) []eventcontract.ListenerProviderContract {
	return []eventcontract.ListenerProviderContract{}
}

// GetCliProviders returns each CLI route provider of the container.
func (p *ContainerComponentProvider) GetCliProviders(
	_ applicationcontract.ApplicationContract,
) []clicontract.CliRouteProviderContract {
	return []clicontract.CliRouteProviderContract{}
}

// GetHttpProviders returns each HTTP route provider of the container.
func (p *ContainerComponentProvider) GetHttpProviders(
	_ applicationcontract.ApplicationContract,
) []httpcontract.HttpRouteProviderContract {
	return []httpcontract.HttpRouteProviderContract{}
}
