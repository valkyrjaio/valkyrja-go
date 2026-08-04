/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the HTTP server sub-component's providers.
package provider

import (
	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	messageconstant "github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	messagefactory "github.com/valkyrjaio/valkyrja-go/v26/http/message/factory"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/http/middleware/handler"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/server/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/server/handler"
)

// HttpServerServiceProvider publishes the bindings of the HTTP server
// sub-component.
type HttpServerServiceProvider struct{}

// Publishers returns a publisher for each binding key that the sub-component
// defers.
func (p *HttpServerServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.RequestHandlerContractServiceID: PublishRequestHandler,
	}
}

// PublishRequestHandler binds the server's entry point for one request.
func PublishRequestHandler(container containercontract.ContainerContract) {
	config := getHttpConfig(container)

	container.SetSingleton(constant.RequestHandlerContractServiceID, handler.NewRequestHandler(
		container,
		getRouter(container),
		middlewarehandler.NewRequestReceivedHandler(container, getRequestReceivedMiddleware(config)...),
		middlewarehandler.NewThrowableCaughtHandler(container, getThrowableCaughtMiddleware(config)...),
		middlewarehandler.NewSendingResponseHandler(container, getSendingResponseMiddleware(config)...),
		middlewarehandler.NewResponseSentHandler(container, getResponseSentMiddleware(config)...),
		getResponseFactory(container),
		getDebugMode(config),
	))
}

// getHttpConfig returns the HTTP configuration that the application publishes.
func getHttpConfig(container containercontract.ContainerContract) applicationcontract.HttpConfigContract {
	resolved, err := container.GetSingleton(applicationconstant.HttpConfigContractServiceID)
	if err != nil {
		return nil
	}

	config, isConfig := resolved.(applicationcontract.HttpConfigContract)
	if !isConfig {
		return nil
	}

	return config
}

// getRouter returns the router that the routing sub-component published.
func getRouter(container containercontract.ContainerContract) contract.RouterContract {
	resolved, err := container.GetSingleton(routingconstant.RouterContractServiceID)
	if err != nil {
		return nil
	}

	built, isRouter := resolved.(contract.RouterContract)
	if !isRouter {
		return nil
	}

	return built
}

// getResponseFactory returns the factory that the message sub-component
// published.
func getResponseFactory(container containercontract.ContainerContract) contract.ResponseFactoryContract {
	resolved, err := container.GetSingleton(messageconstant.ResponseFactoryContractServiceID)
	if err != nil {
		return messagefactory.NewResponseFactory()
	}

	built, isFactory := resolved.(contract.ResponseFactoryContract)
	if !isFactory {
		return messagefactory.NewResponseFactory()
	}

	return built
}

// getDebugMode reports whether a response carries what went wrong.
func getDebugMode(config applicationcontract.HttpConfigContract) bool {
	return config != nil && config.GetDebugMode()
}

// getRequestReceivedMiddleware returns the request-received middleware that the
// application names.
func getRequestReceivedMiddleware(config applicationcontract.HttpConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetRequestReceivedMiddleware()
}

// getThrowableCaughtMiddleware returns the throwable-caught middleware that the
// application names.
func getThrowableCaughtMiddleware(config applicationcontract.HttpConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetThrowableCaughtMiddleware()
}

// getSendingResponseMiddleware returns the sending-response middleware that the
// application names.
func getSendingResponseMiddleware(config applicationcontract.HttpConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetSendingResponseMiddleware()
}

// getResponseSentMiddleware returns the response-sent middleware that the
// application names.
func getResponseSentMiddleware(config applicationcontract.HttpConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetResponseSentMiddleware()
}
