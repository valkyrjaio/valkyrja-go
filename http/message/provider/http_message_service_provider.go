/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the HTTP message sub-component's providers.
package provider

import (
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/factory"
)

type HttpMessageServiceProvider struct{}

// Publishers returns a publisher for each binding key that the sub-component
// defers.
func (p *HttpMessageServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.ResponseFactoryContractServiceID: PublishResponseFactory,
	}
}

// PublishResponseFactory binds the factory that builds each kind of response.
func PublishResponseFactory(container containercontract.ContainerContract) {
	container.SetSingleton(constant.ResponseFactoryContractServiceID, factory.NewResponseFactory())
}
