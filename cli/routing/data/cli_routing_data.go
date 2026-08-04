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

// CliRoutingData is the CLI routing component's state, as a value that the
// framework stores and reloads.
//
// `sindri` generates one of these for the application, so the framework loads
// every command without reading a provider at boot.
type CliRoutingData struct {
	routes map[string]contract.RouteContract
}

// NewCliRoutingData builds the state from the commands, keyed by name.
func NewCliRoutingData(routes map[string]contract.RouteContract) *CliRoutingData {
	if routes == nil {
		routes = map[string]contract.RouteContract{}
	}

	return &CliRoutingData{routes: routes}
}

// GetRoutes returns each route, keyed by its own name.
func (d *CliRoutingData) GetRoutes() map[string]contract.RouteContract {
	return d.routes
}
