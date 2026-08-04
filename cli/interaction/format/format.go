/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package format holds the terminal codes that wrap a piece of text.
package format

import (
	"strconv"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
)

// escapeStart opens a terminal escape sequence.
const escapeStart = "\x1b["

// escapeEnd closes a terminal escape sequence.
const escapeEnd = "m"

// codeSeparator separates one terminal code from the next.
const codeSeparator = ";"

// Format is one pair of terminal codes that wrap a piece of text.
type Format struct {
	setCode   string
	unsetCode string
}

// NewFormat builds a format from the code that starts it and the code that ends
// it.
func NewFormat(setCode string, unsetCode string) *Format {
	return &Format{
		setCode:   setCode,
		unsetCode: unsetCode,
	}
}

// NewTextColorFormat builds the format that colors the text.
func NewTextColorFormat(color constant.TextColor) *Format {
	return NewFormat(strconv.Itoa(int(color)), strconv.Itoa(constant.TextColorDefault))
}

// NewBackgroundColorFormat builds the format that colors the background.
func NewBackgroundColorFormat(color constant.BackgroundColor) *Format {
	return NewFormat(strconv.Itoa(int(color)), strconv.Itoa(constant.BackgroundColorDefault))
}

// NewStyleFormat builds the format that styles the text.
func NewStyleFormat(style constant.Style) *Format {
	return NewFormat(strconv.Itoa(int(style)), strconv.Itoa(style.GetDefault()))
}

// GetSetCode returns the code that starts the format.
func (f *Format) GetSetCode() string {
	return f.setCode
}

// WithSetCode returns a copy of the format for another start code.
func (f *Format) WithSetCode(setCode string) contract.FormatContract {
	copied := *f
	copied.setCode = setCode

	return &copied
}

// GetUnsetCode returns the code that ends the format.
func (f *Format) GetUnsetCode() string {
	return f.unsetCode
}

// WithUnsetCode returns a copy of the format for another end code.
func (f *Format) WithUnsetCode(unsetCode string) contract.FormatContract {
	copied := *f
	copied.unsetCode = unsetCode

	return &copied
}

// Formatter applies each of its formats to a piece of text.
type Formatter struct {
	formats []contract.FormatContract
}

// NewFormatter builds a formatter over its formats.
func NewFormatter(formats ...contract.FormatContract) *Formatter {
	return &Formatter{formats: formats}
}

// NewErrorFormatter builds the formatter that an error message carries.
func NewErrorFormatter() *Formatter {
	return NewFormatter(
		NewTextColorFormat(constant.TextColorLightWhite),
		NewBackgroundColorFormat(constant.BackgroundColorRed),
	)
}

// NewSuccessFormatter builds the formatter that a success message carries.
func NewSuccessFormatter() *Formatter {
	return NewFormatter(
		NewTextColorFormat(constant.TextColorLightWhite),
		NewBackgroundColorFormat(constant.BackgroundColorGreen),
	)
}

// NewWarningFormatter builds the formatter that a warning message carries.
func NewWarningFormatter() *Formatter {
	return NewFormatter(
		NewTextColorFormat(constant.TextColorBlack),
		NewBackgroundColorFormat(constant.BackgroundColorYellow),
	)
}

// NewHighlightedTextFormatter builds the formatter that highlighted text
// carries.
func NewHighlightedTextFormatter() *Formatter {
	return NewFormatter(NewTextColorFormat(constant.TextColorCyan))
}

// NewQuestionFormatter builds the formatter that a question carries.
func NewQuestionFormatter() *Formatter {
	return NewFormatter(NewTextColorFormat(constant.TextColorLightBlue))
}

// GetFormats returns each format that the formatter applies.
func (f *Formatter) GetFormats() []contract.FormatContract {
	return f.formats
}

// WithFormats returns a copy of the formatter with other formats.
func (f *Formatter) WithFormats(formats ...contract.FormatContract) contract.FormatterContract {
	copied := *f
	copied.formats = formats

	return &copied
}

// FormatText returns the text, wrapped in each format.
//
// A formatter that holds no format returns the text as it received it, so a
// terminal that reads no escape sequence receives none.
func (f *Formatter) FormatText(text string) string {
	if len(f.formats) == 0 {
		return text
	}

	set := make([]string, 0, len(f.formats))
	unset := make([]string, 0, len(f.formats))

	for _, applied := range f.formats {
		set = append(set, applied.GetSetCode())
		unset = append(unset, applied.GetUnsetCode())
	}

	return escapeStart + strings.Join(set, codeSeparator) + escapeEnd +
		text +
		escapeStart + strings.Join(unset, codeSeparator) + escapeEnd
}
