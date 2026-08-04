/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package kernel

import (
	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// ChildApplication is an application that holds its own container and reads
// everything else from a parent application.
type ChildApplication struct {
	parent    contract.ApplicationContract
	container containercontract.ContainerContract
}

// NewChildApplication builds a child of the parent, over its own container.
func NewChildApplication(
	parent contract.ApplicationContract,
	container containercontract.ContainerContract,
) *ChildApplication {
	return &ChildApplication{
		parent:    parent,
		container: container,
	}
}

// GetContainer returns the container of the child.
func (a *ChildApplication) GetContainer() containercontract.ContainerContract {
	return a.container
}

// PublishProviderCallbacks runs the parent's callbacks.
func (a *ChildApplication) PublishProviderCallbacks() {
	a.parent.PublishProviderCallbacks()
}

// GetProviders returns the parent's component providers.
func (a *ChildApplication) GetProviders() []contract.ComponentProviderContract {
	return a.parent.GetProviders()
}

// GetContainerProviders returns the parent's service providers.
func (a *ChildApplication) GetContainerProviders() []containercontract.ServiceProviderContract {
	return a.parent.GetContainerProviders()
}

// GetEventProviders returns the parent's listener providers.
func (a *ChildApplication) GetEventProviders() []eventcontract.ListenerProviderContract {
	return a.parent.GetEventProviders()
}

// GetCliProviders returns the parent's CLI route providers.
func (a *ChildApplication) GetCliProviders() []clicontract.CliRouteProviderContract {
	return a.parent.GetCliProviders()
}

// GetHttpProviders returns the parent's HTTP route providers.
func (a *ChildApplication) GetHttpProviders() []httpcontract.HttpRouteProviderContract {
	return a.parent.GetHttpProviders()
}

// GetDebugMode reports whether the parent runs in debug mode.
func (a *ChildApplication) GetDebugMode() bool {
	return a.parent.GetDebugMode()
}

// GetEnvironment returns the environment that the parent runs in.
func (a *ChildApplication) GetEnvironment() string {
	return a.parent.GetEnvironment()
}

// GetVersion returns the version of the parent.
func (a *ChildApplication) GetVersion() string {
	return a.parent.GetVersion()
}
