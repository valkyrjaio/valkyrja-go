/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

type ContainerInvalidPublishCallbackError struct {
	ContainerRuntimeError

	id string
}

// NewContainerInvalidPublishCallbackError builds the error for a binding key.
func NewContainerInvalidPublishCallbackError(id string) *ContainerInvalidPublishCallbackError {
	return &ContainerInvalidPublishCallbackError{
		ContainerRuntimeError: newContainerRuntimeError(
			id+" should have a valid callable",
			nil,
		),
		id: id,
	}
}

// GetID returns the binding key that carries no publisher.
func (e *ContainerInvalidPublishCallbackError) GetID() string {
	return e.id
}
