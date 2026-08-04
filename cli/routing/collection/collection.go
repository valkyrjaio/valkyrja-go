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

// Collection holds every command of the application, keyed by name.
//
// A command name is unique, so a command that is added under a name that the
// collection holds already replaces the one that is there.
type Collection struct {
	routes map[string]contract.RouteContract
}

// NewCollection builds an empty collection.
func NewCollection() *Collection {
	return &Collection{routes: map[string]contract.RouteContract{}}
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
		c.routes[command.GetName()] = command
	}

	return c
}

// Get returns the command under the name, and nil where the collection holds
// none.
func (c *Collection) Get(name string) contract.RouteContract {
	return c.routes[name]
}

// Has reports whether the collection holds a command under the name.
func (c *Collection) Has(name string) bool {
	_, found := c.routes[name]

	return found
}

// All returns every command, keyed by its own name.
func (c *Collection) All() map[string]contract.RouteContract {
	return maps.Clone(c.routes)
}
