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

type EventThrowable interface {
	throwablecontract.ValkyrjaThrowable

	// IsEventThrowable marks the error as one that the event component raised.
	// The mark is what separates this contract from the root contract, which Go
	// otherwise treats as the same type.
	IsEventThrowable() bool
}
