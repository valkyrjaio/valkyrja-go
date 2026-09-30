/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

type CliRoutingData struct {
	routes map[string]contract.RouteFactory
}

// NewCliRoutingData builds the state from a factory for each command, keyed by
// name.
func NewCliRoutingData(routes map[string]contract.RouteFactory) *CliRoutingData {
	if routes == nil {
		routes = map[string]contract.RouteFactory{}
	}

	return &CliRoutingData{routes: routes}
}

// GetRoutes returns a factory for each route, keyed by the route's own name.
func (d *CliRoutingData) GetRoutes() map[string]contract.RouteFactory {
	return d.routes
}
