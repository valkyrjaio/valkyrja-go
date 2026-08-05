/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds every contract of the container component.
//
// The other ports put each contract next to the type that implements it, which
// gives `Manager/Contract` and `Provider/Contract`. Go resolves an import cycle
// at the package level, and `ContainerContract` and `ServiceProviderContract`
// name each other, so a package for each one cannot compile. One `contract`
// package per component is therefore the Go spelling of the same taxonomy.
package contract

import (
	"github.com/valkyrjaio/valkyrja-go/v26/container/constant"
)

// ServiceFactory builds a service. The container passes itself and the
// arguments that the caller gave to `Get`.
type ServiceFactory func(container ContainerContract, arguments []any) any

// PublishFunc registers the bindings that one binding key needs. A service
// provider returns one for each key that it defers.
type PublishFunc func(container ContainerContract)

// ContainerDataContract is the container's state, as a value that the framework
// stores and reloads.
//
// The other ports type `GetData` with the concrete `ContainerData`. In Go that
// is a cycle: `ContainerData` holds a `ServiceFactory`, and a `ServiceFactory`
// takes a `ContainerContract`. The contract therefore names an interface, and
// the `data` package implements it.
type ContainerDataContract interface {
	// GetAliases returns each alias, keyed by the alias, valued by the binding
	// key that it points to.
	GetAliases() map[string]string

	// GetDeferredCallbacks returns each deferred publisher, keyed by the
	// binding key that it publishes.
	GetDeferredCallbacks() map[string]PublishFunc

	// GetServices returns each service factory, keyed by the binding key.
	GetServices() map[string]ServiceFactory

	// GetSingletons returns each singleton binding, keyed by the binding key.
	GetSingletons() map[string]string
}

// ProvidersAwareContract is the part of the container that registers a service
// provider and publishes what the provider defers.
type ProvidersAwareContract interface {
	// Register records each publisher that the provider defers. It reports a
	// failure where a publisher is nil.
	Register(provider ServiceProviderContract) error

	// IsDeferred reports whether a provider defers the binding key.
	IsDeferred(id string) bool

	// IsPublished reports whether the binding key is published already.
	IsPublished(id string) bool

	// Publish runs the deferred publisher for the binding key. It does nothing
	// where no publisher defers the key.
	Publish(id string)
}

// ContainerContract is the framework's service container.
//
// Go has no default parameter, so each method takes every argument that the
// other ports make optional. A caller that binds no arguments passes nil.
//
// The contract mirrors the PHP reference implementation method for method.
// Splitting it to satisfy the method count would give this port an API that no
// other port has, so the count is suppressed instead.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type ContainerContract interface {
	ProvidersAwareContract

	// GetData returns the container's state.
	GetData() ContainerDataContract

	// SetFromData merges the state into the container.
	SetFromData(data ContainerDataContract)

	// Has reports whether the container resolves the binding key.
	Has(id string) bool

	// Bind records a factory for the binding key. The container calls the
	// factory for each resolution.
	Bind(id string, factory ServiceFactory) ContainerContract

	// BindAlias records an alias that points to a binding key.
	BindAlias(alias string, id string) ContainerContract

	// BindSingleton records a factory for the binding key. The container calls
	// the factory once and returns that instance for each later resolution.
	BindSingleton(id string, factory ServiceFactory) ContainerContract

	// SetSingleton records an instance for the binding key.
	SetSingleton(id string, singleton any) ContainerContract

	// IsAlias reports whether the binding key is an alias.
	IsAlias(id string) bool

	// IsService reports whether a factory is bound to the binding key.
	IsService(id string) bool

	// IsSingleton reports whether the binding key is a singleton binding or a
	// singleton instance.
	IsSingleton(id string) bool

	// IsSingletonBinding reports whether the binding key is bound as a
	// singleton.
	IsSingletonBinding(id string) bool

	// IsSingletonInstance reports whether an instance is set for the binding
	// key.
	IsSingletonInstance(id string) bool

	// Get resolves the binding key. It reports a failure where the container
	// resolves nothing.
	Get(id string, arguments []any, mode constant.InvalidReferenceMode) (any, error)

	// GetAliased resolves the binding key as an alias. It reports a failure
	// where the key is no alias.
	GetAliased(id string, arguments []any) (any, error)

	// GetService resolves the binding key as a service. It reports a failure
	// where no factory is bound to the key.
	GetService(id string, arguments []any) (any, error)

	// GetSingleton resolves the binding key as a singleton. It reports a
	// failure where the key is no singleton.
	GetSingleton(id string) (any, error)
}

// ServiceProviderContract publishes the container bindings of one component.
//
// A provider returns a literal map, and never a computed one, because `sindri`
// reads the map from the source rather than by running it.
type ServiceProviderContract interface {
	// Publishers returns each publisher, keyed by the binding key that it
	// publishes.
	Publishers() map[string]PublishFunc
}
