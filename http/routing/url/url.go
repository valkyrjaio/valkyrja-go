/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package url builds the URL of a named route.
package url

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

type Url struct {
	collection contract.RouteCollectionContract
}

// NewUrl builds the URL generator over a route collection.
func NewUrl(collection contract.RouteCollectionContract) *Url {
	return &Url{collection: collection}
}

// GetUrl returns the URL of the route, with the data filled into each parameter
// of its path.
func (u *Url) GetUrl(name string, data map[string]string) string {
	route := u.collection.GetByName(name)
	if route == nil {
		return ""
	}

	path := route.GetPath()

	for key, value := range data {
		path = strings.ReplaceAll(path, "{"+key+"}", value)
		path = strings.ReplaceAll(path, "{"+key+"?}", value)
	}

	return path
}
