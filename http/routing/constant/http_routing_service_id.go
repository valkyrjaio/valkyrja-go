/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// The HTTP routing component's binding keys.
const (
	// RouterContractServiceID is the binding key for the router.
	RouterContractServiceID = "valkyrja.http.routing.dispatcher.RouterContract"

	// RouteCollectionContractServiceID is the binding key for the route
	// collection.
	RouteCollectionContractServiceID = "valkyrja.http.routing.collection.RouteCollectionContract"

	// RouteContractServiceID is the binding key for the route that matched. The
	// router binds it before it runs the handler, so a handler reads the route.
	RouteContractServiceID = "valkyrja.http.routing.data.RouteContract"

	// HttpRoutingDataServiceID is the binding key for the routing data.
	HttpRoutingDataServiceID = "valkyrja.http.routing.data.HttpRoutingData"

	// MatcherContractServiceID is the binding key for the matcher.
	MatcherContractServiceID = "valkyrja.http.routing.matcher.MatcherContract"

	// ProcessorContractServiceID is the binding key for the processor.
	ProcessorContractServiceID = "valkyrja.http.routing.processor.ProcessorContract"

	// UrlContractServiceID is the binding key for the URL generator.
	UrlContractServiceID = "valkyrja.http.routing.url.UrlContract"
)

// RoutingResponseFactoryContractServiceID is the binding key for what builds a
// response that sends the client to a named route.
const RoutingResponseFactoryContractServiceID = "valkyrja.http.routing.factory.RoutingResponseFactoryContract"
