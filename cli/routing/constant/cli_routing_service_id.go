/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// The binding key of each service that the CLI routing sub-component publishes.
//
// A language with no class token spells a binding key as a string, and the
// TypeScript port spells it the same way.
const (
	RouterContractServiceID = "valkyrja.cli.routing.dispatcher.RouterContract"

	RouteCollectionContractServiceID = "valkyrja.cli.routing.collection.RouteCollectionContract"

	RouteContractServiceID = "valkyrja.cli.routing.data.RouteContract"
)
