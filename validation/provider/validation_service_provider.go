/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the validation component's providers.
package provider

import (
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	containerprovider "github.com/valkyrjaio/valkyrja-go/v26/container/provider"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/validator"
)

// ValidationServiceProvider publishes the bindings of the validation component.
type ValidationServiceProvider struct{}

// Publishers returns a publisher for each binding key that the component defers.
func (p *ValidationServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.ValidatorContractServiceID: PublishValidator,
	}
}

// PublishValidator binds a validator that holds no rule.
//
// A caller states the rules of one validation, so the validator is bound with
// none, and the caller sets them before it validates.
func PublishValidator(container containercontract.ContainerContract) {
	container.SetSingleton(constant.ValidatorContractServiceID, validator.NewValidator(nil))
}

// ValidationComponentProvider is the validation component's top-level provider.
type ValidationComponentProvider struct{}

// GetComponentProviders returns each component that validation needs.
func (p *ValidationComponentProvider) GetComponentProviders(
	_ applicationcontract.ApplicationContract,
) []applicationcontract.ComponentProviderContract {
	return []applicationcontract.ComponentProviderContract{
		&containerprovider.ContainerComponentProvider{},
	}
}

// GetContainerProviders returns each service provider of validation.
func (p *ValidationComponentProvider) GetContainerProviders(
	_ applicationcontract.ApplicationContract,
) []containercontract.ServiceProviderContract {
	return []containercontract.ServiceProviderContract{&ValidationServiceProvider{}}
}

// GetEventProviders returns each listener provider of validation.
func (p *ValidationComponentProvider) GetEventProviders(
	_ applicationcontract.ApplicationContract,
) []eventcontract.ListenerProviderContract {
	return []eventcontract.ListenerProviderContract{}
}

// GetCliProviders returns each CLI route provider of validation.
func (p *ValidationComponentProvider) GetCliProviders(
	_ applicationcontract.ApplicationContract,
) []clicontract.CliRouteProviderContract {
	return []clicontract.CliRouteProviderContract{}
}

// GetHttpProviders returns each HTTP route provider of validation.
func (p *ValidationComponentProvider) GetHttpProviders(
	_ applicationcontract.ApplicationContract,
) []httpcontract.HttpRouteProviderContract {
	return []httpcontract.HttpRouteProviderContract{}
}
