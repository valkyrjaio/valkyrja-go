/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package kernel holds the application and the child application.
package kernel

import (
	"os"

	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// timezoneEnvName is the environment variable that the Go runtime reads a
// timezone from, the same one that the TypeScript port writes.
const timezoneEnvName = "TZ"

// Valkyrja is the application.
//
// It walks the provider tree once and holds what it found, so a second call
// returns the same providers rather than walking again.
type Valkyrja struct {
	container containercontract.ContainerContract
	config    contract.ConfigContract

	providers          []contract.ComponentProviderContract
	serviceProviders   []containercontract.ServiceProviderContract
	eventProviders     []eventcontract.ListenerProviderContract
	cliRouteProviders  []clicontract.CliRouteProviderContract
	httpRouteProviders []httpcontract.HttpRouteProviderContract
}

// NewValkyrja builds the application over a container and a config, and sets the
// timezone that the config names.
func NewValkyrja(
	container containercontract.ContainerContract,
	config contract.ConfigContract,
) *Valkyrja {
	app := &Valkyrja{
		container: container,
		config:    config,
	}

	app.bootstrapTimezone()

	return app
}

// GetContainer returns the container of the application.
func (a *Valkyrja) GetContainer() containercontract.ContainerContract {
	return a.container
}

// PublishProviderCallbacks runs each callback that the config holds.
func (a *Valkyrja) PublishProviderCallbacks() {
	for _, callback := range a.config.GetCallbacks() {
		callback(a)
	}
}

// GetProviders returns every component provider, including the ones that a
// component provider names.
//
// A provider that names itself, or that two providers both name, is collected
// once. The other ports rely on a set for that; Go has no set, so the walk keeps
// a map of what it saw.
func (a *Valkyrja) GetProviders() []contract.ComponentProviderContract {
	if len(a.providers) > 0 {
		return a.providers
	}

	seen := map[contract.ComponentProviderContract]bool{}

	for _, provider := range a.config.GetProviders() {
		a.collectProviders(provider, seen)
	}

	return a.providers
}

// GetContainerProviders returns every service provider of every component.
func (a *Valkyrja) GetContainerProviders() []containercontract.ServiceProviderContract {
	if len(a.serviceProviders) > 0 {
		return a.serviceProviders
	}

	seen := map[containercontract.ServiceProviderContract]bool{}

	for _, provider := range a.GetProviders() {
		for _, serviceProvider := range provider.GetContainerProviders(a) {
			if seen[serviceProvider] {
				continue
			}

			seen[serviceProvider] = true

			a.serviceProviders = append(a.serviceProviders, serviceProvider)
		}
	}

	return a.serviceProviders
}

// GetEventProviders returns every listener provider of every component.
func (a *Valkyrja) GetEventProviders() []eventcontract.ListenerProviderContract {
	if len(a.eventProviders) > 0 {
		return a.eventProviders
	}

	seen := map[eventcontract.ListenerProviderContract]bool{}

	for _, provider := range a.GetProviders() {
		for _, eventProvider := range provider.GetEventProviders(a) {
			if seen[eventProvider] {
				continue
			}

			seen[eventProvider] = true

			a.eventProviders = append(a.eventProviders, eventProvider)
		}
	}

	return a.eventProviders
}

// GetCliProviders returns every CLI route provider of every component.
func (a *Valkyrja) GetCliProviders() []clicontract.CliRouteProviderContract {
	if len(a.cliRouteProviders) > 0 {
		return a.cliRouteProviders
	}

	seen := map[clicontract.CliRouteProviderContract]bool{}

	for _, provider := range a.GetProviders() {
		for _, cliProvider := range provider.GetCliProviders(a) {
			if seen[cliProvider] {
				continue
			}

			seen[cliProvider] = true

			a.cliRouteProviders = append(a.cliRouteProviders, cliProvider)
		}
	}

	return a.cliRouteProviders
}

// GetHttpProviders returns every HTTP route provider of every component.
func (a *Valkyrja) GetHttpProviders() []httpcontract.HttpRouteProviderContract {
	if len(a.httpRouteProviders) > 0 {
		return a.httpRouteProviders
	}

	seen := map[httpcontract.HttpRouteProviderContract]bool{}

	for _, provider := range a.GetProviders() {
		for _, httpProvider := range provider.GetHttpProviders(a) {
			if seen[httpProvider] {
				continue
			}

			seen[httpProvider] = true

			a.httpRouteProviders = append(a.httpRouteProviders, httpProvider)
		}
	}

	return a.httpRouteProviders
}

// GetDebugMode reports whether the application runs in debug mode.
func (a *Valkyrja) GetDebugMode() bool {
	return a.config.GetDebugMode()
}

// GetEnvironment returns the environment that the application runs in.
func (a *Valkyrja) GetEnvironment() string {
	return a.config.GetEnvironment()
}

// GetVersion returns the version of the application.
func (a *Valkyrja) GetVersion() string {
	return a.config.GetVersion()
}

// collectProviders walks what the provider names, depth first, and records the
// provider after the ones it names.
func (a *Valkyrja) collectProviders(
	provider contract.ComponentProviderContract,
	seen map[contract.ComponentProviderContract]bool,
) {
	if seen[provider] {
		return
	}

	seen[provider] = true

	for _, subProvider := range provider.GetComponentProviders(a) {
		a.collectProviders(subProvider, seen)
	}

	a.providers = append(a.providers, provider)
}

// bootstrapTimezone sets the timezone that the config names.
func (a *Valkyrja) bootstrapTimezone() {
	// The other ports assign the variable and read no result, so a failure
	// here changes nothing that the application can act on.
	_ = os.Setenv(timezoneEnvName, a.config.GetTimezone())
}
