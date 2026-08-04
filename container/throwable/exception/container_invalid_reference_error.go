/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

import (
	"fmt"
)

// ContainerInvalidReferenceError reports a binding key that the container
// resolves nothing for.
type ContainerInvalidReferenceError struct {
	ContainerInvalidArgumentError

	id string
}

// NewContainerInvalidReferenceError builds the error for a binding key.
func NewContainerInvalidReferenceError(id string) *ContainerInvalidReferenceError {
	return &ContainerInvalidReferenceError{
		ContainerInvalidArgumentError: newContainerInvalidArgumentError(
			fmt.Sprintf("Service with `%s` not found", id),
			nil,
		),
		id: id,
	}
}

// GetID returns the binding key that the container resolved nothing for.
func (e *ContainerInvalidReferenceError) GetID() string {
	return e.id
}
