/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

type TextColor int

// The text colors that the framework writes.
const (
	TextColorBlack        TextColor = 30
	TextColorRed          TextColor = 31
	TextColorGreen        TextColor = 32
	TextColorYellow       TextColor = 33
	TextColorBlue         TextColor = 34
	TextColorMagenta      TextColor = 35
	TextColorCyan         TextColor = 36
	TextColorWhite        TextColor = 37
	TextColorDarkGray     TextColor = 90
	TextColorLightRed     TextColor = 91
	TextColorLightGreen   TextColor = 92
	TextColorLightYellow  TextColor = 93
	TextColorLightBlue    TextColor = 94
	TextColorLightMagenta TextColor = 95
	TextColorLightCyan    TextColor = 96
	TextColorLightWhite   TextColor = 97
)

// TextColorDefault is the code that returns the text to its own color.
const TextColorDefault = 39

type BackgroundColor int

// The background colors that the framework writes.
const (
	BackgroundColorBlack        BackgroundColor = 40
	BackgroundColorRed          BackgroundColor = 41
	BackgroundColorGreen        BackgroundColor = 42
	BackgroundColorYellow       BackgroundColor = 43
	BackgroundColorBlue         BackgroundColor = 44
	BackgroundColorMagenta      BackgroundColor = 45
	BackgroundColorCyan         BackgroundColor = 46
	BackgroundColorWhite        BackgroundColor = 47
	BackgroundColorDarkGray     BackgroundColor = 100
	BackgroundColorLightRed     BackgroundColor = 101
	BackgroundColorLightGreen   BackgroundColor = 102
	BackgroundColorLightYellow  BackgroundColor = 103
	BackgroundColorLightBlue    BackgroundColor = 104
	BackgroundColorLightMagenta BackgroundColor = 105
	BackgroundColorLightCyan    BackgroundColor = 106
	BackgroundColorLightWhite   BackgroundColor = 107
)

// BackgroundColorDefault is the code that returns the background to its own
// color.
const BackgroundColorDefault = 49

type Style int

// The styles that the framework writes.
const (
	StyleBold       Style = 1
	StyleUnderscore Style = 4
	StyleBlink      Style = 5
	StyleInverse    Style = 7
	StyleConceal    Style = 8
)

// styleDefaults holds the code that ends each style.
var styleDefaults = map[Style]int{
	StyleBold:       22,
	StyleUnderscore: 24,
	StyleBlink:      25,
	StyleInverse:    27,
	StyleConceal:    28,
}

// GetDefault returns the code that ends the style, and zero where the style is
// not one that the framework knows.
func (s Style) GetDefault() int {
	return styleDefaults[s]
}
