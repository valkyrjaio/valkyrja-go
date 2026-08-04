/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package rule

import (
	"net/mail"
	"strconv"
	"unicode"
)

// isEmpty reports whether the subject carries no value.
//
// The other ports read PHP's `empty`, which is false for an empty string, a
// zero, a false, an empty list, and a null. Go has no such rule, so this states
// the same one for the types that a subject carries.
func isEmpty(subject any) bool {
	switch held := subject.(type) {
	case nil:
		return true
	case string:
		return held == ""
	case bool:
		return !held
	case int:
		return held == 0
	case int64:
		return held == 0
	case float64:
		return held == 0
	case []any:
		return len(held) == 0
	case map[string]any:
		return len(held) == 0
	default:
		return false
	}
}

// isNumeric reports whether the subject is a number, or a string that reads as
// one.
func isNumeric(subject any) bool {
	switch held := subject.(type) {
	case int, int64, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(held, 64)

		return err == nil
	default:
		return false
	}
}

// isEmail reports whether the text is an email address.
//
// Go's `net/mail` reads an address with a display name, such as
// `Melech <melech@example.com>`, and a subject that carries one is not an
// address on its own. The parsed address must therefore match the text.
func isEmail(text string) bool {
	parsed, err := mail.ParseAddress(text)

	return err == nil && parsed.Address == text
}

// isAlphabetic reports whether every character of the text is a letter, and the
// text carries at least one.
func isAlphabetic(text string) bool {
	if text == "" {
		return false
	}

	for _, held := range text {
		if !unicode.IsLetter(held) {
			return false
		}
	}

	return true
}

// isLowercase reports whether the text carries no uppercase letter.
func isLowercase(text string) bool {
	for _, held := range text {
		if unicode.IsUpper(held) {
			return false
		}
	}

	return true
}

// isUppercase reports whether the text carries no lowercase letter.
func isUppercase(text string) bool {
	for _, held := range text {
		if unicode.IsLower(held) {
			return false
		}
	}

	return true
}
