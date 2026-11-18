/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package uri holds the URI of a request, and the rules that RFC 3986 sets for
// each of its parts.
package uri

import (
	"strings"
)

// The characters that each part of a URI carries without percent-encoding.
//
// The other ports hold these as atoms of a regular expression character class.
// Go reads each byte instead, which needs no `regexp` and states the set where a
// reader sees it.
const (
	// unreservedChars is the set that every part carries.
	//
	// See https://tools.ietf.org/html/rfc3986#section-2.3
	unreservedChars = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"_-.~"

	// subDelimiterChars is the set of sub-delimiters, which every part below
	// carries.
	//
	// See https://tools.ietf.org/html/rfc3986#section-2.2
	subDelimiterChars = "!$&'()*+,;="

	// userInfoChars adds the colon that separates the user name from the
	// password.
	//
	// See https://tools.ietf.org/html/rfc3986#section-3.2.1
	userInfoChars = unreservedChars + subDelimiterChars + ":"

	// hostChars carries no character beyond the common set.
	//
	// See https://tools.ietf.org/html/rfc3986#section-3.2.2
	hostChars = unreservedChars + subDelimiterChars

	// pathChars adds a colon, an at sign, and the segment separator.
	//
	// See https://tools.ietf.org/html/rfc3986#section-3.3
	pathChars = unreservedChars + subDelimiterChars + ":@/"

	// queryChars adds a question mark to the path set. The fragment carries the
	// same set.
	//
	// See https://tools.ietf.org/html/rfc3986#section-3.4
	queryChars = pathChars + "?"
)

// hexDigits holds the digits that a percent-encoded triplet is written with.
const hexDigits = "0123456789ABCDEF"

// encode percent-encodes each byte that the part does not carry.
func encode(value string, allowed string) string {
	encoded := &strings.Builder{}

	for index := 0; index < len(value); index++ {
		character := value[index]

		if character == '%' && isTripletAt(value, index) {
			encoded.WriteByte('%')
			encoded.WriteByte(toUpperHex(value[index+1]))
			encoded.WriteByte(toUpperHex(value[index+2]))

			index += 2

			continue
		}

		if strings.IndexByte(allowed, character) != -1 {
			encoded.WriteByte(character)

			continue
		}

		encoded.WriteByte('%')
		encoded.WriteByte(hexDigits[character>>4])
		encoded.WriteByte(hexDigits[character&0x0F])
	}

	return encoded.String()
}

// isTripletAt reports whether a valid percent-encoded triplet starts at the
// index.
func isTripletAt(value string, index int) bool {
	if index+2 >= len(value) {
		return false
	}

	return isHexDigit(value[index+1]) && isHexDigit(value[index+2])
}

// isHexDigit reports whether the byte is a hexadecimal digit.
func isHexDigit(character byte) bool {
	switch {
	case character >= '0' && character <= '9',
		character >= 'a' && character <= 'f',
		character >= 'A' && character <= 'F':
		return true
	default:
		return false
	}
}

// toUpperHex returns the hexadecimal digit in upper case.
func toUpperHex(character byte) byte {
	if character >= 'a' && character <= 'f' {
		return character - ('a' - 'A')
	}

	return character
}
