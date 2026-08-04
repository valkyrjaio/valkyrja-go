/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package param holds the parameters that arrived with a request.
package param

import (
	"maps"
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// ParamCollection holds a set of parameters that arrived with a request.
//
// The other ports declare this abstract, and each named collection extends it
// without adding a method. Go has no abstract type, and a Go interface with the
// same method set is the same type, so this port declares one collection and
// gives each name an alias in the contract package.
type ParamCollection struct {
	params map[string]any
}

// NewParamCollection builds a collection over the parameters. It copies the map,
// so a later write to the source cannot reach the collection.
func NewParamCollection(params map[string]any) *ParamCollection {
	return &ParamCollection{params: copyParams(params)}
}

// Has reports whether the collection holds a parameter under the key.
func (c *ParamCollection) Has(key string) bool {
	_, found := c.params[key]

	return found
}

// Get returns the parameter under the key, and nil where the collection holds
// none.
func (c *ParamCollection) Get(key string) any {
	return c.params[key]
}

// GetAll returns a copy of every parameter.
func (c *ParamCollection) GetAll() map[string]any {
	return copyParams(c.params)
}

// GetOnly returns the parameters that the keys identify.
func (c *ParamCollection) GetOnly(keys ...string) map[string]any {
	return c.filterByKey(keys, true)
}

// GetAllExcept returns every parameter that the keys do not identify.
func (c *ParamCollection) GetAllExcept(keys ...string) map[string]any {
	return c.filterByKey(keys, false)
}

// With returns a copy of the collection that holds the parameters and nothing
// else.
func (c *ParamCollection) With(params map[string]any) contract.ParamCollectionContract {
	return NewParamCollection(params)
}

// WithAdded returns a copy of the collection with the parameters added to the
// ones it holds. A parameter of a key that the collection holds is replaced.
func (c *ParamCollection) WithAdded(params map[string]any) contract.ParamCollectionContract {
	combined := copyParams(c.params)

	maps.Copy(combined, params)

	return &ParamCollection{params: combined}
}

// filterByKey returns the parameters whose key is in the keys, where keep is
// true, and the parameters whose key is not, where keep is false.
func (c *ParamCollection) filterByKey(keys []string, keep bool) map[string]any {
	filtered := map[string]any{}

	for key, param := range c.params {
		if slices.Contains(keys, key) == keep {
			filtered[key] = param
		}
	}

	return filtered
}

// copyParams returns a copy of the map, and an empty map where the map is nil.
func copyParams(source map[string]any) map[string]any {
	target := make(map[string]any, len(source))

	maps.Copy(target, source)

	return target
}
