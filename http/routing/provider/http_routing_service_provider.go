/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the HTTP routing sub-component's providers.
package provider

import (
	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	messageconstant "github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	messagefactory "github.com/valkyrjaio/valkyrja-go/v26/http/message/factory"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/http/middleware/handler"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/dispatcher"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/matcher"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/processor"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/url"
)

// HttpRoutingServiceProvider publishes the bindings of the HTTP routing
// sub-component.
type HttpRoutingServiceProvider struct{}

// Publishers returns a publisher for each binding key that the sub-component
// defers.
func (p *HttpRoutingServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.ProcessorContractServiceID:              PublishProcessor,
		constant.RouteCollectionContractServiceID:        PublishRouteCollection,
		constant.MatcherContractServiceID:                PublishMatcher,
		constant.UrlContractServiceID:                    PublishUrl,
		constant.RoutingResponseFactoryContractServiceID: PublishRoutingResponseFactory,
		constant.RouterContractServiceID:                 PublishRouter,
	}
}

// PublishProcessor binds what turns a declared route into one that the matcher
// reads.
func PublishProcessor(container containercontract.ContainerContract) {
	container.SetSingleton(constant.ProcessorContractServiceID, processor.NewProcessor())
}

// PublishRouteCollection binds the collection, and files every route that a
// route provider of the application registers.
//
// The processor reads each route before the collection files it, because the
// matcher reads the regular expression that the processor builds.
func PublishRouteCollection(container containercontract.ContainerContract) {
	built := collection.NewRouteCollection()
	routeProcessor := getProcessor(container)

	for _, routeProvider := range getRouteProviders(container) {
		for _, route := range routeProvider.GetRoutes() {
			built.Add(routeProcessor.Route(route))
		}
	}

	container.SetSingleton(constant.RouteCollectionContractServiceID, built)
}

// PublishMatcher binds what matches a request to a route.
func PublishMatcher(container containercontract.ContainerContract) {
	container.SetSingleton(constant.MatcherContractServiceID, matcher.NewMatcher(getCollection(container)))
}

// PublishUrl binds what reads the URL of a named route.
func PublishUrl(container containercontract.ContainerContract) {
	container.SetSingleton(constant.UrlContractServiceID, url.NewUrl(getCollection(container)))
}

// PublishRoutingResponseFactory binds what builds a response that sends the
// client to a named route.
func PublishRoutingResponseFactory(container containercontract.ContainerContract) {
	container.SetSingleton(
		constant.RoutingResponseFactoryContractServiceID,
		factory.NewRoutingResponseFactory(getUrl(container), getResponseFactory(container)),
	)
}

// PublishRouter binds the router, with the middleware that the application names
// for each stage.
func PublishRouter(container containercontract.ContainerContract) {
	config := getHttpConfig(container)

	container.SetSingleton(constant.RouterContractServiceID, dispatcher.NewRouter(
		container,
		getMatcher(container),
		getResponseFactory(container),
		middlewarehandler.NewRouteMatchedHandler(container, getRouteMatchedMiddleware(config)...),
		middlewarehandler.NewRouteNotMatchedHandler(container, getRouteNotMatchedMiddleware(config)...),
		middlewarehandler.NewRouteDispatchedHandler(container, getRouteDispatchedMiddleware(config)...),
	))
}

// getRouteProviders returns each HTTP route provider that the application
// registers.
func getRouteProviders(container containercontract.ContainerContract) []contract.HttpRouteProviderContract {
	resolved, err := container.GetSingleton(applicationconstant.ApplicationContractServiceID)
	if err != nil {
		return nil
	}

	app, isApplication := resolved.(applicationcontract.ApplicationContract)
	if !isApplication {
		return nil
	}

	return app.GetHttpProviders()
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

// getProcessor returns the processor that the sub-component published.
func getProcessor(container containercontract.ContainerContract) contract.ProcessorContract {
	resolved, err := container.GetSingleton(constant.ProcessorContractServiceID)
	if err != nil {
		return processor.NewProcessor()
	}

	built, isProcessor := resolved.(contract.ProcessorContract)
	if !isProcessor {
		return processor.NewProcessor()
	}

	return built
}

// getCollection returns the collection that the sub-component published.
func getCollection(container containercontract.ContainerContract) contract.RouteCollectionContract {
	resolved, err := container.GetSingleton(constant.RouteCollectionContractServiceID)
	if err != nil {
		return collection.NewRouteCollection()
	}

	built, isCollection := resolved.(contract.RouteCollectionContract)
	if !isCollection {
		return collection.NewRouteCollection()
	}

	return built
}

// getMatcher returns the matcher that the sub-component published.
func getMatcher(container containercontract.ContainerContract) contract.MatcherContract {
	resolved, err := container.GetSingleton(constant.MatcherContractServiceID)
	if err != nil {
		return matcher.NewMatcher(getCollection(container))
	}

	built, isMatcher := resolved.(contract.MatcherContract)
	if !isMatcher {
		return matcher.NewMatcher(getCollection(container))
	}

	return built
}

// getUrl returns what the sub-component published to read the URL of a route.
func getUrl(container containercontract.ContainerContract) contract.UrlContract {
	resolved, err := container.GetSingleton(constant.UrlContractServiceID)
	if err != nil {
		return url.NewUrl(getCollection(container))
	}

	built, isUrl := resolved.(contract.UrlContract)
	if !isUrl {
		return url.NewUrl(getCollection(container))
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

// getRouteMatchedMiddleware returns the route-matched middleware that the
// application names.
func getRouteMatchedMiddleware(config applicationcontract.HttpConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetRouteMatchedMiddleware()
}

// getRouteNotMatchedMiddleware returns the route-not-matched middleware that the
// application names.
func getRouteNotMatchedMiddleware(config applicationcontract.HttpConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetRouteNotMatchedMiddleware()
}

// getRouteDispatchedMiddleware returns the route-dispatched middleware that the
// application names.
func getRouteDispatchedMiddleware(config applicationcontract.HttpConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetRouteDispatchedMiddleware()
}
