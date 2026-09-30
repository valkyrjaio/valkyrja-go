/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

type ModeTranslation string

// The ModeTranslation values that the framework knows.
const (
	ModeTranslationNone       ModeTranslation = ""
	ModeTranslationWindows    ModeTranslation = "t"
	ModeTranslationBinarySafe ModeTranslation = "b"
)
