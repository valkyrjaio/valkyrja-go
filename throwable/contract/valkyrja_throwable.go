/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds the root throwable contract for the framework.
package contract

// ValkyrjaThrowable is the root contract for every error that the framework
// raises. Each component declares its own contract that embeds this one, and
// each concrete error satisfies the contract of its component.
//
// Go has no exceptions, so a throwable is an error value. The other ports name
// the same type `ValkyrjaThrowable` and spell the concrete types `*Exception`;
// this port spells them `*Error`, because `Exception` is foreign to Go.
type ValkyrjaThrowable interface {
	error

	// GetTraceCode returns a stable identifier for the site that raised the
	// error. Two errors raised at the same site share a trace code, and a
	// reader quotes the code to locate the site.
	GetTraceCode() string
}
