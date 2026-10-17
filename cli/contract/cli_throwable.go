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

type CliThrowable interface {
	throwablecontract.ValkyrjaThrowable

	// IsCliThrowable marks the error as one that the CLI component raised. The
	// mark is what separates this contract from the root contract, which Go
	// otherwise treats as the same type.
	IsCliThrowable() bool
}

type CliInteractionThrowable interface {
	CliThrowable

	// IsCliInteractionThrowable marks the error as one that the CLI interaction sub-component raised.
	IsCliInteractionThrowable() bool
}

type CliMiddlewareThrowable interface {
	CliThrowable

	// IsCliMiddlewareThrowable marks the error as one that the CLI middleware raised.
	IsCliMiddlewareThrowable() bool
}

type CliRoutingThrowable interface {
	CliThrowable

	// IsCliRoutingThrowable marks the error as one that the CLI routing sub-component raised.
	IsCliRoutingThrowable() bool
}

type CliServerThrowable interface {
	CliThrowable

	// IsCliServerThrowable marks the error as one that the CLI server raised.
	IsCliServerThrowable() bool
}
