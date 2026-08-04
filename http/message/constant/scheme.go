/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// Scheme is the scheme of a URI.
type Scheme string

// The Scheme values that the framework knows.
const (
	SchemeEmpty Scheme = ""
	SchemeHttp  Scheme = "http"
	SchemeHttps Scheme = "https"
)
