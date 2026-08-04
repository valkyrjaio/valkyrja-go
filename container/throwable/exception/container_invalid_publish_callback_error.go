/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

// ContainerInvalidPublishCallbackError reports a service provider that defers a
// binding key without a publisher for it.
//
// The other ports raise this where the publisher is not callable. Go types the
// publisher, so the one value that reaches the container without a publisher is
// a nil function.
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
