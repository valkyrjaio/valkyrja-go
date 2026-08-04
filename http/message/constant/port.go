/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// The port numbers that a URI uses.
const (
	// PortMin is the lowest port number that a URI accepts.
	PortMin = 1

	// PortMax is the highest port number that a URI accepts.
	PortMax = 65535

	// PortHttp is the default port of the `http` scheme.
	PortHttp = 80

	// PortHttps is the default port of the `https` scheme.
	PortHttps = 443
)

// IsValidPort reports whether the port number is one that a URI accepts.
//
// The other ports declare this beside the constants, as a static method on the
// constant class. Go has no class, so it is a function in the same package.
func IsValidPort(port int) bool {
	return port >= PortMin && port <= PortMax
}
