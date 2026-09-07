/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds the root throwable contract for the framework.
package contract

type ValkyrjaThrowable interface {
	error

	// GetTraceCode returns a stable identifier for the site that raised the
	// error. Two errors raised at the same site share a trace code, and a
	// reader quotes the code to locate the site.
	GetTraceCode() string
}
