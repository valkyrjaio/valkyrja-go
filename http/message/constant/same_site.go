/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

type SameSite string

// The SameSite values that the framework knows.
const (
	SameSiteNone   SameSite = "none"
	SameSiteLax    SameSite = "lax"
	SameSiteStrict SameSite = "strict"
)
