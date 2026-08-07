/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package manager holds the container and the child container.
package manager

import (
	"maps"

	"github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/data"
	"github.com/valkyrjaio/valkyrja-go/v26/container/throwable/exception"
)

type containerInternal interface {
	contract.ContainerContract

	getSingletonWithoutChecks(id string) (any, bool, error)
	getServiceWithoutChecks(id string, arguments []any) (any, bool, error)
	getAliasedWithoutChecks(id string, arguments []any) (any, bool, error)
}

type Container struct {
	self              containerInternal
	aliases           map[string]string
	instances         map[string]any
	services          map[string]contract.ServiceFactory
	singletons        map[string]string
	deferredCallbacks map[string]contract.PublishFunc
	published         map[string]bool
}

// NewContainer builds a container and loads the state into it. It accepts nil
// for an empty container.
func NewContainer(containerData contract.ContainerDataContract) *Container {
	container := &Container{
		aliases:           map[string]string{},
		instances:         map[string]any{},
		services:          map[string]contract.ServiceFactory{},
		singletons:        map[string]string{},
		deferredCallbacks: map[string]contract.PublishFunc{},
		published:         map[string]bool{},
	}

	container.self = container

	if containerData != nil {
		container.SetFromData(containerData)
	}

	return container
}

// GetData returns the container's state.
func (c *Container) GetData() contract.ContainerDataContract {
	return data.NewContainerData(c.aliases, c.deferredCallbacks, c.services, c.singletons)
}

// SetFromData merges the state into the container.
func (c *Container) SetFromData(containerData contract.ContainerDataContract) {
	maps.Copy(c.aliases, containerData.GetAliases())
	maps.Copy(c.deferredCallbacks, containerData.GetDeferredCallbacks())
	maps.Copy(c.services, containerData.GetServices())
	maps.Copy(c.singletons, containerData.GetSingletons())
}

// Has reports whether the container resolves the binding key.
func (c *Container) Has(id string) bool {
	return c.self.IsDeferred(id) ||
		c.self.IsSingleton(id) ||
		c.self.IsService(id) ||
		c.self.IsAlias(id)
}

// Bind records a factory for the binding key.
func (c *Container) Bind(id string, factory contract.ServiceFactory) contract.ContainerContract {
	c.services[id] = factory
	c.published[id] = true

	return c.self
}

// BindAlias records an alias that points to a binding key.
func (c *Container) BindAlias(alias string, id string) contract.ContainerContract {
	c.aliases[alias] = id

	return c.self
}

// BindSingleton records a factory that the container calls once.
func (c *Container) BindSingleton(id string, factory contract.ServiceFactory) contract.ContainerContract {
	c.singletons[id] = id

	c.self.Bind(id, factory)

	return c.self
}

// SetSingleton records an instance for the binding key.
func (c *Container) SetSingleton(id string, singleton any) contract.ContainerContract {
	c.instances[id] = singleton
	c.published[id] = true

	return c.self
}

// IsAlias reports whether the binding key is an alias.
func (c *Container) IsAlias(id string) bool {
	_, found := c.aliases[id]

	return found
}

// IsService reports whether a factory is bound to the binding key.
func (c *Container) IsService(id string) bool {
	_, found := c.services[id]

	return found
}

// IsSingleton reports whether the binding key is a singleton binding or a
// singleton instance.
func (c *Container) IsSingleton(id string) bool {
	return c.self.IsSingletonBinding(id) || c.self.IsSingletonInstance(id)
}

// IsSingletonBinding reports whether the binding key is bound as a singleton.
func (c *Container) IsSingletonBinding(id string) bool {
	_, found := c.singletons[id]

	return found
}

// IsSingletonInstance reports whether an instance is set for the binding key.
func (c *Container) IsSingletonInstance(id string) bool {
	_, found := c.instances[id]

	return found
}

// Get resolves the binding key.
func (c *Container) Get(id string, arguments []any, mode constant.InvalidReferenceMode) (any, error) {
	c.publishUnpublishedProvided(id)

	singleton, found, err := c.self.getSingletonWithoutChecks(id)
	if err != nil {
		return nil, err
	}

	if found {
		return singleton, nil
	}

	service, found, err := c.self.getServiceWithoutChecks(id, arguments)
	if err != nil {
		return nil, err
	}

	if found {
		return service, nil
	}

	aliased, found, err := c.self.getAliasedWithoutChecks(id, arguments)
	if err != nil {
		return nil, err
	}

	if found {
		return aliased, nil
	}

	return nil, exception.NewContainerInvalidReferenceError(id)
}

// GetAliased resolves the binding key as an alias.
func (c *Container) GetAliased(id string, arguments []any) (any, error) {
	aliased, found, err := c.self.getAliasedWithoutChecks(id, arguments)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, exception.NewContainerInvalidReferenceError(id)
	}

	return aliased, nil
}

// GetService resolves the binding key as a service.
func (c *Container) GetService(id string, arguments []any) (any, error) {
	c.publishUnpublishedProvided(id)

	service, found, err := c.self.getServiceWithoutChecks(id, arguments)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, exception.NewContainerInvalidReferenceError(id)
	}

	return service, nil
}

// GetSingleton resolves the binding key as a singleton.
func (c *Container) GetSingleton(id string) (any, error) {
	c.publishUnpublishedProvided(id)

	singleton, found, err := c.self.getSingletonWithoutChecks(id)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, exception.NewContainerInvalidReferenceError(id)
	}

	return singleton, nil
}

// Register records each publisher that the provider defers.
func (c *Container) Register(provider contract.ServiceProviderContract) error {
	for id, publisher := range provider.Publishers() {
		if publisher == nil {
			return exception.NewContainerInvalidPublishCallbackError(id)
		}

		c.deferredCallbacks[id] = publisher
	}

	return nil
}

// IsDeferred reports whether a provider defers the binding key.
func (c *Container) IsDeferred(id string) bool {
	_, found := c.deferredCallbacks[id]

	return found
}

// IsPublished reports whether the binding key is published already.
func (c *Container) IsPublished(id string) bool {
	_, found := c.published[id]

	return found
}

// Publish runs the deferred publisher for the binding key.
func (c *Container) Publish(id string) {
	publisher, found := c.deferredCallbacks[id]
	if !found {
		return
	}

	publisher(c.self)

	c.published[id] = true
}

// getSingletonWithoutChecks resolves a singleton, and reports whether it
// resolved one. It records the instance the first time that it builds one.
func (c *Container) getSingletonWithoutChecks(id string) (any, bool, error) {
	instance, found := c.instances[id]
	if found {
		return instance, true, nil
	}

	if !c.self.IsSingletonBinding(id) {
		return nil, false, nil
	}

	singleton, found, err := c.self.getServiceWithoutChecks(id, nil)
	if err != nil {
		return nil, false, err
	}

	if found {
		c.instances[id] = singleton
	}

	return singleton, found, nil
}

// getServiceWithoutChecks calls the factory for the binding key, and reports
// whether a factory is bound to it. It carries an error, because a child
// container reads the factory through its parent.
func (c *Container) getServiceWithoutChecks(id string, arguments []any) (any, bool, error) {
	factory, found := c.services[id]
	if !found {
		return nil, false, nil
	}

	return factory(c.self, arguments), true, nil
}

// getAliasedWithoutChecks resolves what the alias points to, and reports
// whether the binding key is an alias.
func (c *Container) getAliasedWithoutChecks(id string, arguments []any) (any, bool, error) {
	aliased, found := c.aliases[id]
	if !found {
		return nil, false, nil
	}

	resolved, err := c.self.Get(aliased, arguments, constant.NewInstanceOrThrowException)
	if err != nil {
		return nil, false, err
	}

	return resolved, true, nil
}

// publishUnpublishedProvided publishes the binding key where a provider defers
// it and nothing published it yet.
func (c *Container) publishUnpublishedProvided(id string) {
	if c.self.IsDeferred(id) && !c.self.IsPublished(id) {
		c.self.Publish(id)
	}
}
