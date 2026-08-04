/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package manager

import (
	"maps"

	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// ChildContainer resolves what its own state holds, and falls back to a parent
// container for everything else.
//
// The child takes only the singleton bindings and the deferred publishers from
// the state. It reads each alias, each service, and each singleton instance
// through the parent.
type ChildContainer struct {
	*Container

	parent contract.ContainerContract
}

// NewChildContainer builds a child of the parent and loads the state into it.
// It accepts nil state for a child that holds nothing of its own.
func NewChildContainer(
	parent contract.ContainerContract,
	containerData contract.ContainerDataContract,
) *ChildContainer {
	child := &ChildContainer{
		Container: NewContainer(nil),
		parent:    parent,
	}

	child.self = child

	if containerData != nil {
		maps.Copy(child.singletons, containerData.GetSingletons())
		maps.Copy(child.deferredCallbacks, containerData.GetDeferredCallbacks())
	}

	return child
}

// IsAlias reports whether the child or the parent holds the alias.
func (c *ChildContainer) IsAlias(id string) bool {
	return c.Container.IsAlias(id) || c.parent.IsAlias(id)
}

// IsService reports whether the child or the parent holds the factory.
func (c *ChildContainer) IsService(id string) bool {
	return c.Container.IsService(id) || c.parent.IsService(id)
}

// IsSingletonInstance reports whether the child or the parent holds the
// instance.
func (c *ChildContainer) IsSingletonInstance(id string) bool {
	return c.Container.IsSingletonInstance(id) || c.parent.IsSingletonInstance(id)
}

// IsDeferred reports whether the child or the parent defers the binding key.
func (c *ChildContainer) IsDeferred(id string) bool {
	return c.Container.IsDeferred(id) || c.parent.IsDeferred(id)
}

// IsPublished reports whether the child or the parent published the binding
// key.
func (c *ChildContainer) IsPublished(id string) bool {
	return c.Container.IsPublished(id) || c.parent.IsPublished(id)
}

// getSingletonWithoutChecks reads the parent's instance where the child holds
// none of its own.
func (c *ChildContainer) getSingletonWithoutChecks(id string) (any, bool, error) {
	if !c.Container.IsSingletonInstance(id) && c.parent.IsSingletonInstance(id) {
		singleton, err := c.parent.GetSingleton(id)
		if err != nil {
			return nil, false, err
		}

		return singleton, true, nil
	}

	return c.Container.getSingletonWithoutChecks(id)
}

// getServiceWithoutChecks calls the parent's factory where the child holds none
// of its own.
func (c *ChildContainer) getServiceWithoutChecks(id string, arguments []any) (any, bool, error) {
	if !c.Container.IsService(id) && c.parent.IsService(id) {
		service, err := c.parent.GetService(id, arguments)
		if err != nil {
			return nil, false, err
		}

		return service, true, nil
	}

	return c.Container.getServiceWithoutChecks(id, arguments)
}

// getAliasedWithoutChecks reads the parent's alias where the child holds none
// of its own.
func (c *ChildContainer) getAliasedWithoutChecks(id string, arguments []any) (any, bool, error) {
	if !c.Container.IsAlias(id) && c.parent.IsAlias(id) {
		aliased, err := c.parent.GetAliased(id, arguments)
		if err != nil {
			return nil, false, err
		}

		return aliased, true, nil
	}

	return c.Container.getAliasedWithoutChecks(id, arguments)
}
