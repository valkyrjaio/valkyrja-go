/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package header holds one header of an HTTP message, and the collection that
// holds every header of one message.
package header

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

// The ASCII codes that a header value carries or rejects.
const (
	asciiTab            = 9
	asciiLineFeed       = 10
	asciiCarriageReturn = 13
	asciiSpace          = 32
	asciiLowestPrinting = 32
	asciiDelete         = 127
	asciiHighest        = 255
)

// IsValidName reports whether the name is one that RFC 7230 allows for a header.
func IsValidName(name string) bool {
	if name == "" {
		return false
	}

	for _, character := range []byte(name) {
		if !isValidNameByte(character) {
			return false
		}
	}

	return true
}

// isValidNameByte reports whether the byte is one that a header name carries.
func isValidNameByte(character byte) bool {
	switch {
	case character >= 'a' && character <= 'z',
		character >= 'A' && character <= 'Z',
		character >= '0' && character <= '9':
		return true
	default:
		return strings.IndexByte("!#$%&'*+-.^_`|~", character) != -1
	}
}

// ValidateName reports a failure where the name is not one that a header
// carries.
func ValidateName(name string) error {
	if !IsValidName(name) {
		return exception.NewHttpHeaderInvalidNameError(name)
	}

	return nil
}

// IsValidValue reports whether the value is one that RFC 7230 allows for a
// header.
func IsValidValue(value string) bool {
	characters := []byte(value)

	for index := 0; index < len(characters); index++ {
		character := characters[index]

		if character == asciiCarriageReturn || character == asciiLineFeed {
			if !isValidLineFold(characters, index) {
				return false
			}

			// The fold is three bytes. The loop steps over the line feed, so
			// the next pass does not read it as a bare line feed of its own.
			index++

			continue
		}

		if isInvalidValueByte(character) {
			return false
		}
	}

	return true
}

// isValidLineFold reports whether the byte at the index starts a line fold, which
// is a carriage return, a line feed, and then a space or a tab.
func isValidLineFold(characters []byte, index int) bool {
	if characters[index] != asciiCarriageReturn {
		return false
	}

	if index+2 >= len(characters) {
		return false
	}

	if characters[index+1] != asciiLineFeed {
		return false
	}

	return characters[index+2] == asciiSpace || characters[index+2] == asciiTab
}

// isInvalidValueByte reports whether the byte is one that a header value never
// carries.
func isInvalidValueByte(character byte) bool {
	if character < asciiLowestPrinting &&
		character != asciiTab &&
		character != asciiLineFeed &&
		character != asciiCarriageReturn {
		return true
	}

	return character == asciiDelete || character == asciiHighest
}

// ValidateValue reports a failure where the value is not one that a header
// carries.
func ValidateValue(value string) error {
	if !IsValidValue(value) {
		return exception.NewHttpHeaderInvalidValueError(value)
	}

	return nil
}
