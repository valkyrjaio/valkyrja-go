/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package collection holds every route of the application.
package collection

import (
	"maps"
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
)

type RouteCollection struct {
	routes  map[string]contract.RouteContract
	paths   map[constant.RequestMethod]map[string]string
	regexes map[constant.RequestMethod]map[string]string
}

// NewRouteCollection builds an empty collection.
func NewRouteCollection() *RouteCollection {
	return &RouteCollection{
		routes:  map[string]contract.RouteContract{},
		paths:   map[constant.RequestMethod]map[string]string{},
		regexes: map[constant.RequestMethod]map[string]string{},
	}
}

// GetData returns the collection's state.
func (c *RouteCollection) GetData() contract.HttpRoutingDataContract {
	return data.NewHttpRoutingData(c.routes, c.paths, c.regexes)
}

// SetFromData replaces the collection's state.
func (c *RouteCollection) SetFromData(routingData contract.HttpRoutingDataContract) {
	c.routes = routingData.GetRoutes()
	c.paths = routingData.GetPaths()
	c.regexes = routingData.GetRegexes()
}

// Add files the route under its own name, and under its path or its regular
// expression, for each request method that it matches.
func (c *RouteCollection) Add(route contract.RouteContract) {
	c.routes[route.GetName()] = route

	dynamic, isDynamic := route.(contract.DynamicRouteContract)

	for _, method := range c.getMethodsOf(route) {
		if isDynamic && dynamic.GetRegex() != "" {
			c.fileUnder(c.regexes, method, dynamic.GetRegex(), route.GetName())

			continue
		}

		c.fileUnder(c.paths, method, route.GetPath(), route.GetName())
	}
}

// HasPath reports whether a route matches the static path and the request
// method.
func (c *RouteCollection) HasPath(path string, method constant.RequestMethod) bool {
	_, found := c.paths[method][path]

	return found
}

// GetByPath returns the route at the static path, and nil where none matches.
func (c *RouteCollection) GetByPath(path string, method constant.RequestMethod) contract.RouteContract {
	return c.routes[c.paths[method][path]]
}

// HasRegex reports whether a route matches the regular expression and the
// request method.
func (c *RouteCollection) HasRegex(regex string, method constant.RequestMethod) bool {
	_, found := c.regexes[method][regex]

	return found
}

// GetByRegex returns the route at the regular expression, and nil where none
// matches.
func (c *RouteCollection) GetByRegex(
	regex string,
	method constant.RequestMethod,
) contract.DynamicRouteContract {
	route, isDynamic := c.routes[c.regexes[method][regex]].(contract.DynamicRouteContract)
	if !isDynamic || route.GetRegex() == "" {
		return nil
	}

	return route
}

// GetPaths returns the name of the route at each static path.
func (c *RouteCollection) GetPaths(method constant.RequestMethod) map[string]string {
	return maps.Clone(c.paths[method])
}

// GetRegexes returns the name of the route at each regular expression.
func (c *RouteCollection) GetRegexes(method constant.RequestMethod) map[string]string {
	return maps.Clone(c.regexes[method])
}

// HasName reports whether the collection holds a route under the name.
func (c *RouteCollection) HasName(name string) bool {
	_, found := c.routes[name]

	return found
}

// GetByName returns the route under the name, and nil where the collection holds
// none.
func (c *RouteCollection) GetByName(name string) contract.RouteContract {
	return c.routes[name]
}

// GetAll returns every route that matches the request method, keyed by its own
// name.
func (c *RouteCollection) GetAll(method constant.RequestMethod) map[string]contract.RouteContract {
	all := map[string]contract.RouteContract{}

	for _, name := range c.paths[method] {
		all[name] = c.routes[name]
	}

	for _, name := range c.regexes[method] {
		all[name] = c.routes[name]
	}

	return all
}

// getMethodsOf returns each request method that the route files under. A route
// that matches the any method files under every one.
func (c *RouteCollection) getMethodsOf(route contract.RouteContract) []constant.RequestMethod {
	if slices.Contains(route.GetRequestMethods(), constant.RequestMethodAny) {
		return constant.GetAllRequestMethods()
	}

	return route.GetRequestMethods()
}

// fileUnder records the name of the route under the key, for the request method.
func (c *RouteCollection) fileUnder(
	target map[constant.RequestMethod]map[string]string,
	method constant.RequestMethod,
	key string,
	name string,
) {
	if target[method] == nil {
		target[method] = map[string]string{}
	}

	target[method][key] = name
}
