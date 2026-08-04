/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"maps"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// HttpRoutingData is the HTTP routing component's state. `sindri` generates one
// for the whole application, and the route collection loads it at boot.
//
// Each getter returns a copy, so a holder of the data cannot change the
// collection behind its back.
type HttpRoutingData struct {
	routes  map[string]contract.RouteContract
	paths   map[constant.RequestMethod]map[string]string
	regexes map[constant.RequestMethod]map[string]string
}

// NewHttpRoutingData builds the state from each map. It copies every map, and it
// accepts nil for a map that carries nothing.
func NewHttpRoutingData(
	routes map[string]contract.RouteContract,
	paths map[constant.RequestMethod]map[string]string,
	regexes map[constant.RequestMethod]map[string]string,
) *HttpRoutingData {
	return &HttpRoutingData{
		routes:  maps.Clone(routes),
		paths:   copyByMethod(paths),
		regexes: copyByMethod(regexes),
	}
}

// GetRoutes returns a copy of each route, keyed by its own name.
func (d *HttpRoutingData) GetRoutes() map[string]contract.RouteContract {
	routes := make(map[string]contract.RouteContract, len(d.routes))

	maps.Copy(routes, d.routes)

	return routes
}

// GetPaths returns a copy of the name of the route at each static path.
func (d *HttpRoutingData) GetPaths() map[constant.RequestMethod]map[string]string {
	return copyByMethod(d.paths)
}

// GetRegexes returns a copy of the name of the route at each regular expression.
func (d *HttpRoutingData) GetRegexes() map[constant.RequestMethod]map[string]string {
	return copyByMethod(d.regexes)
}

// copyByMethod copies the outer map and each inner map, so a later write to
// either one cannot reach the data.
func copyByMethod(
	source map[constant.RequestMethod]map[string]string,
) map[constant.RequestMethod]map[string]string {
	target := make(map[constant.RequestMethod]map[string]string, len(source))

	for method, byKey := range source {
		target[method] = maps.Clone(byKey)
	}

	return target
}
