/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds the container's state as a value that the framework stores
// and reloads.
package data

import (
	"maps"

	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

type ContainerData struct {
	aliases           map[string]string
	deferredCallbacks map[string]contract.PublishFunc
	services          map[string]contract.ServiceFactory
	singletons        map[string]string
}

// NewContainerData builds the state from each map. It copies every map, and it
// accepts nil for a map that carries nothing.
func NewContainerData(
	aliases map[string]string,
	deferredCallbacks map[string]contract.PublishFunc,
	services map[string]contract.ServiceFactory,
	singletons map[string]string,
) *ContainerData {
	return &ContainerData{
		aliases:           copyMap(aliases),
		deferredCallbacks: copyMap(deferredCallbacks),
		services:          copyMap(services),
		singletons:        copyMap(singletons),
	}
}

// GetAliases returns a copy of each alias.
func (d *ContainerData) GetAliases() map[string]string {
	return copyMap(d.aliases)
}

// GetDeferredCallbacks returns a copy of each deferred publisher.
func (d *ContainerData) GetDeferredCallbacks() map[string]contract.PublishFunc {
	return copyMap(d.deferredCallbacks)
}

// GetServices returns a copy of each service factory.
func (d *ContainerData) GetServices() map[string]contract.ServiceFactory {
	return copyMap(d.services)
}

// GetSingletons returns a copy of each singleton binding.
func (d *ContainerData) GetSingletons() map[string]string {
	return copyMap(d.singletons)
}

// copyMap returns a copy of the map, and an empty map where the map is nil.
// The copy is never nil, so a caller writes to it without a guard.
func copyMap[K comparable, V any](source map[K]V) map[K]V {
	target := make(map[K]V, len(source))

	maps.Copy(target, source)

	return target
}
