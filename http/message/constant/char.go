/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// The characters that each part of a URI may hold, written for a regular
// expression's character class.
const (
	CharUnreserved = `a-zA-Z0-9_\-\.~`

	CharSubDelims = `!\$&\'\(\)\*\+,;=`

	CharUserInfo = CharUnreserved + CharSubDelims + `:`

	CharHost = CharUnreserved + CharSubDelims

	CharPath = CharUnreserved + CharSubDelims + `:@\/`

	CharQuery = CharUnreserved + CharSubDelims + `:@\/\?`
)
