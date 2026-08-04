/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
)

func TestGetDefaultReturnsTheCodeThatEndsTheStyle(t *testing.T) {
	t.Parallel()

	tests := map[constant.Style]int{
		constant.StyleBold:       22,
		constant.StyleUnderscore: 24,
		constant.StyleBlink:      25,
		constant.StyleInverse:    27,
		constant.StyleConceal:    28,
	}

	for style, code := range tests {
		if style.GetDefault() != code {
			t.Errorf("GetDefault for %d must be %d, but is %d", style, code, style.GetDefault())
		}
	}
}

func TestGetDefaultIsZeroForAnUnknownStyle(t *testing.T) {
	t.Parallel()

	if constant.Style(99).GetDefault() != 0 {
		t.Errorf("GetDefault for an unknown style must be zero, but is: %d", constant.Style(99).GetDefault())
	}
}
