/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package fixtures holds the reusable doubles that the application component's
// tests build on.
package fixtures

import (
	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// ComponentProviderFixture is a component provider that returns what the test
// gave it.
type ComponentProviderFixture struct {
	ComponentProviders []contract.ComponentProviderContract
	ContainerProviders []containercontract.ServiceProviderContract
	EventProviders     []eventcontract.ListenerProviderContract
	CliProviders       []clicontract.CliRouteProviderContract
	HttpProviders      []httpcontract.HttpRouteProviderContract
}

// GetComponentProviders returns each component that the provider names.
func (p *ComponentProviderFixture) GetComponentProviders(
	_ contract.ApplicationContract,
) []contract.ComponentProviderContract {
	return p.ComponentProviders
}

// GetContainerProviders returns each service provider that the provider names.
func (p *ComponentProviderFixture) GetContainerProviders(
	_ contract.ApplicationContract,
) []containercontract.ServiceProviderContract {
	return p.ContainerProviders
}

// GetEventProviders returns each listener provider that the provider names.
func (p *ComponentProviderFixture) GetEventProviders(
	_ contract.ApplicationContract,
) []eventcontract.ListenerProviderContract {
	return p.EventProviders
}

// GetCliProviders returns each CLI route provider that the provider names.
func (p *ComponentProviderFixture) GetCliProviders(
	_ contract.ApplicationContract,
) []clicontract.CliRouteProviderContract {
	return p.CliProviders
}

// GetHttpProviders returns each HTTP route provider that the provider names.
func (p *ComponentProviderFixture) GetHttpProviders(
	_ contract.ApplicationContract,
) []httpcontract.HttpRouteProviderContract {
	return p.HttpProviders
}

// ServiceProviderFixture is a service provider that publishes nothing.
type ServiceProviderFixture struct {
	Name string
}

// Publishers returns no publisher.
func (p *ServiceProviderFixture) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{}
}

// ListenerProviderFixture is a listener provider that registers no listener.
type ListenerProviderFixture struct {
	Name string
}

// GetListeners returns no listener.
func (p *ListenerProviderFixture) GetListeners() []eventcontract.ListenerContract {
	return []eventcontract.ListenerContract{}
}

// CliRouteProviderFixture is a CLI route provider that registers no command.
type CliRouteProviderFixture struct {
	Name string
}

// GetRoutes returns no command.
func (p *CliRouteProviderFixture) GetRoutes() []clicontract.RouteContract {
	return []clicontract.RouteContract{}
}

// HttpRouteProviderFixture is an HTTP route provider that registers no route.
type HttpRouteProviderFixture struct {
	Name string
}

// GetRoutes returns no route.
func (p *HttpRouteProviderFixture) GetRoutes() []httpcontract.RouteContract {
	return []httpcontract.RouteContract{}
}
