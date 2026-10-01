/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package collection holds every command of the application.
package collection

import (
	"maps"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
)

type Collection struct {
	routes map[string]contract.RouteFactory
}

// NewCollection builds an empty collection.
func NewCollection() *Collection {
	return &Collection{routes: map[string]contract.RouteFactory{}}
}

// GetData returns the collection's state.
func (c *Collection) GetData() contract.CliRoutingDataContract {
	return data.NewCliRoutingData(maps.Clone(c.routes))
}

// SetFromData replaces the collection's state.
func (c *Collection) SetFromData(routingData contract.CliRoutingDataContract) {
	c.routes = maps.Clone(routingData.GetRoutes())
}

// Add files each command under its own name.
func (c *Collection) Add(commands ...contract.RouteContract) contract.RouteCollectionContract {
	for _, command := range commands {
		c.routes[command.GetName()] = factoryOf(command)
	}

	return c
}

// Get returns the command under the name, and nil where the collection holds
// none.
func (c *Collection) Get(name string) contract.RouteContract {
	factory, found := c.routes[name]
	if !found {
		return nil
	}

	return c.resolve(name, factory)
}

// Has reports whether the collection holds a command under the name.
func (c *Collection) Has(name string) bool {
	_, found := c.routes[name]

	return found
}

// All returns every command, keyed by its own name.
func (c *Collection) All() map[string]contract.RouteContract {
	routes := make(map[string]contract.RouteContract, len(c.routes))

	for name, factory := range c.routes {
		routes[name] = c.resolve(name, factory)
	}

	return routes
}

// resolve builds the command from its factory, and keeps the command, so a
// later read returns the same one without building it again.
func (c *Collection) resolve(name string, factory contract.RouteFactory) contract.RouteContract {
	route := factory()

	c.routes[name] = factoryOf(route)

	return route
}

// factoryOf returns a factory that returns the command.
func factoryOf(route contract.RouteContract) contract.RouteFactory {
	return func() contract.RouteContract {
		return route
	}
}
