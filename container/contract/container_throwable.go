/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	throwablecontract "github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
)

// ContainerThrowable is the contract that every error of the container
// component satisfies. A caller that handles any container failure asserts
// against this contract.
type ContainerThrowable interface {
	throwablecontract.ValkyrjaThrowable

	// IsContainerThrowable marks the error as one that the container raised.
	// The mark is what separates this contract from the root contract, which
	// Go otherwise treats as the same type.
	IsContainerThrowable() bool
}
