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

// Url builds the URL of a named route.
type Url struct {
	collection contract.RouteCollectionContract
}

// NewUrl builds the URL generator over a route collection.
func NewUrl(collection contract.RouteCollectionContract) *Url {
	return &Url{collection: collection}
}

// GetUrl returns the URL of the route, with the data filled into each parameter
// of its path.
//
// It returns an empty string where the collection holds no route under the name.
// A parameter that the data does not name stays in the path, so a reader sees
// which value the caller left out.
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
