/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// InvalidReferenceMode says what the container does where it resolves nothing
// for a binding key.
//
// Go has no enum, so the other ports' enum is a defined type over a `const`
// block. The taxonomy puts it in the `constant` segment for that reason.
type InvalidReferenceMode int

const (
	// NewInstanceOrThrowException asks the container to build the type that the
	// key names, and to report a failure where it cannot.
	//
	// The PHP port builds the instance by reflection over the class name. Go
	// cannot construct a type from a string, so this port reports the failure,
	// which is what the TypeScript port does for the same reason. The mode
	// stays in the API, because an application selects it by name across the
	// ports.
	NewInstanceOrThrowException InvalidReferenceMode = iota

	// ThrowException asks the container to report a failure.
	ThrowException
)
