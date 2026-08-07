/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the container component's binding keys and its
// enumerations.
package constant

// The container component's binding keys.
//
// Go has no `::class` equivalent, so every binding key is a string literal.
// Each key names one thing that the container resolves.
const (
	// ContainerContractServiceID is the binding key for the container itself.
	ContainerContractServiceID = "valkyrja.container.manager.ContainerContract"

	// ContainerDataServiceID is the binding key for the container's data.
	ContainerDataServiceID = "valkyrja.container.data.ContainerData"
)
