/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the CLI routing sub-component's providers.
package provider

import (
	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	interactionfactory "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/handler"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/router"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// CliRoutingServiceProvider publishes the bindings of the CLI routing
// sub-component.
type CliRoutingServiceProvider struct{}

// Publishers returns a publisher for each binding key that the sub-component
// defers.
func (p *CliRoutingServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.RouteCollectionContractServiceID: PublishRouteCollection,
		constant.RouterContractServiceID:          PublishRouter,
	}
}

// PublishRouteCollection binds the collection, and files every command that a
// route provider of the application registers.
func PublishRouteCollection(container containercontract.ContainerContract) {
	built := collection.NewCollection()

	for _, routeProvider := range getRouteProviders(container) {
		built.Add(routeProvider.GetRoutes()...)
	}

	container.SetSingleton(constant.RouteCollectionContractServiceID, built)
}

// PublishRouter binds the router, with the middleware that the application names
// for each stage.
func PublishRouter(container containercontract.ContainerContract) {
	config := getCliConfig(container)

	container.SetSingleton(constant.RouterContractServiceID, router.NewRouter(
		container,
		getCollection(container),
		getOutputFactory(container),
		middlewarehandler.NewRouteMatchedHandler(container, getRouteMatchedMiddleware(config)...),
		middlewarehandler.NewRouteNotMatchedHandler(container, getRouteNotMatchedMiddleware(config)...),
		middlewarehandler.NewRouteDispatchedHandler(container, getRouteDispatchedMiddleware(config)...),
		middlewarehandler.NewThrowableCaughtHandler(container, getThrowableCaughtMiddleware(config)...),
		middlewarehandler.NewProcessExitingHandler(container, getProcessExitingMiddleware(config)...),
	))
}

// getRouteProviders returns each CLI route provider that the application
// registers.
func getRouteProviders(container containercontract.ContainerContract) []clicontract.CliRouteProviderContract {
	resolved, err := container.GetSingleton(applicationconstant.ApplicationContractServiceID)
	if err != nil {
		return nil
	}

	app, isApplication := resolved.(applicationcontract.ApplicationContract)
	if !isApplication {
		return nil
	}

	return app.GetCliProviders()
}

// getCliConfig returns the CLI configuration that the application publishes.
func getCliConfig(container containercontract.ContainerContract) applicationcontract.CliConfigContract {
	resolved, err := container.GetSingleton(applicationconstant.CliConfigContractServiceID)
	if err != nil {
		return nil
	}

	config, isConfig := resolved.(applicationcontract.CliConfigContract)
	if !isConfig {
		return nil
	}

	return config
}

// getCollection returns the collection that the sub-component published.
func getCollection(container containercontract.ContainerContract) clicontract.RouteCollectionContract {
	resolved, err := container.GetSingleton(constant.RouteCollectionContractServiceID)
	if err != nil {
		return collection.NewCollection()
	}

	built, isCollection := resolved.(clicontract.RouteCollectionContract)
	if !isCollection {
		return collection.NewCollection()
	}

	return built
}

// getOutputFactory returns the factory that the interaction sub-component
// published.
func getOutputFactory(container containercontract.ContainerContract) clicontract.OutputFactoryContract {
	resolved, err := container.GetSingleton(interactionconstant.OutputFactoryContractServiceID)
	if err != nil {
		return interactionfactory.NewOutputFactory(nil)
	}

	built, isFactory := resolved.(clicontract.OutputFactoryContract)
	if !isFactory {
		return interactionfactory.NewOutputFactory(nil)
	}

	return built
}

// getRouteMatchedMiddleware returns the route-matched middleware that the
// application names.
func getRouteMatchedMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetRouteMatchedMiddleware()
}

// getRouteNotMatchedMiddleware returns the route-not-matched middleware that the
// application names.
func getRouteNotMatchedMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetRouteNotMatchedMiddleware()
}

// getRouteDispatchedMiddleware returns the route-dispatched middleware that the
// application names.
func getRouteDispatchedMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetRouteDispatchedMiddleware()
}

// getThrowableCaughtMiddleware returns the throwable-caught middleware that the
// application names.
func getThrowableCaughtMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetThrowableCaughtMiddleware()
}

// getProcessExitingMiddleware returns the process-exiting middleware that the
// application names.
func getProcessExitingMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetProcessExitingMiddleware()
}
