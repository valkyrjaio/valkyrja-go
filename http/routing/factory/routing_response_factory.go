/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package factory builds a response that sends the client to a named route.
package factory

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// RoutingResponseFactory builds a response that sends the client to a named
// route.
type RoutingResponseFactory struct {
	url             contract.UrlContract
	responseFactory contract.ResponseFactoryContract
}

// NewRoutingResponseFactory builds the factory over what reads the URL of a
// route, and what builds a response.
func NewRoutingResponseFactory(
	url contract.UrlContract,
	responseFactory contract.ResponseFactoryContract,
) *RoutingResponseFactory {
	return &RoutingResponseFactory{
		url:             url,
		responseFactory: responseFactory,
	}
}

// CreateRouteRedirectResponse builds a response that sends the client to the
// route.
func (f *RoutingResponseFactory) CreateRouteRedirectResponse(
	name string,
	data map[string]string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.RedirectResponseContract {
	return f.responseFactory.CreateRedirectResponse(f.url.GetUrl(name, data), statusCode, headers)
}
