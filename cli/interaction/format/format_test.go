/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package format_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
)

const plainText = "the text"

func TestNewFormatHoldsEachCode(t *testing.T) {
	t.Parallel()

	built := format.NewFormat("1", "22")

	if built.GetSetCode() != "1" || built.GetUnsetCode() != "22" {
		t.Error("the format must hold each code, but did not")
	}
}

func TestEachFormatWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := format.NewFormat("1", "22")

	if built.WithSetCode("4").GetSetCode() != "4" {
		t.Error("WithSetCode must hold the new code, but did not")
	}

	if built.WithUnsetCode("24").GetUnsetCode() != "24" {
		t.Error("WithUnsetCode must hold the new code, but did not")
	}

	if built.GetSetCode() != "1" || built.GetUnsetCode() != "22" {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachNamedFormatCarriesItsCodes(t *testing.T) {
	t.Parallel()

	text := format.NewTextColorFormat(constant.TextColorRed)
	background := format.NewBackgroundColorFormat(constant.BackgroundColorBlue)
	style := format.NewStyleFormat(constant.StyleBold)

	if text.GetSetCode() != "31" || text.GetUnsetCode() != "39" {
		t.Errorf("the text color must carry its codes, but carries: %q and %q",
			text.GetSetCode(), text.GetUnsetCode())
	}

	if background.GetSetCode() != "44" || background.GetUnsetCode() != "49" {
		t.Errorf("the background color must carry its codes, but carries: %q and %q",
			background.GetSetCode(), background.GetUnsetCode())
	}

	if style.GetSetCode() != "1" || style.GetUnsetCode() != "22" {
		t.Errorf("the style must carry its codes, but carries: %q and %q",
			style.GetSetCode(), style.GetUnsetCode())
	}
}

func TestFormatTextWrapsTheTextInEachFormat(t *testing.T) {
	t.Parallel()

	built := format.NewFormatter(
		format.NewTextColorFormat(constant.TextColorRed),
		format.NewStyleFormat(constant.StyleBold),
	)

	if built.FormatText(plainText) != "\x1b[31;1mthe text\x1b[39;22m" {
		t.Errorf("FormatText must wrap the text in each format, but is: %q", built.FormatText(plainText))
	}
}

func TestFormatTextReturnsTheTextWhereTheFormatterHoldsNoFormat(t *testing.T) {
	t.Parallel()

	if format.NewFormatter().FormatText(plainText) != plainText {
		t.Error("a formatter with no format must return the text as it is, but did not")
	}
}

func TestWithFormatsReturnsACopy(t *testing.T) {
	t.Parallel()

	built := format.NewFormatter(format.NewTextColorFormat(constant.TextColorRed))

	replaced := built.WithFormats(format.NewStyleFormat(constant.StyleBold))

	if len(replaced.GetFormats()) != 1 || replaced.GetFormats()[0].GetSetCode() != "1" {
		t.Error("WithFormats must replace the formats, but did not")
	}

	if built.GetFormats()[0].GetSetCode() != "31" {
		t.Error("WithFormats must leave the receiver unchanged, but did not")
	}
}

func TestEachTypeSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var applied contract.FormatContract = format.NewFormat("1", "22")
	var formatter contract.FormatterContract = format.NewFormatter(applied)

	if formatter.FormatText(plainText) == plainText {
		t.Error("the contracts must wrap the text, but did not")
	}
}

func TestEachNamedFormatterWrapsTheTextInItsOwnCodes(t *testing.T) {
	t.Parallel()

	formatters := map[string]*format.Formatter{
		"the error formatter":            format.NewErrorFormatter(),
		"the success formatter":          format.NewSuccessFormatter(),
		"the warning formatter":          format.NewWarningFormatter(),
		"the highlighted text formatter": format.NewHighlightedTextFormatter(),
		"the question formatter":         format.NewQuestionFormatter(),
	}

	seen := map[string]string{}

	for name, formatter := range formatters {
		formatted := formatter.FormatText("text")

		if !strings.Contains(formatted, "text") {
			t.Errorf("%s must keep the text, but returned: %q", name, formatted)
		}

		if formatted == "text" {
			t.Errorf("%s must apply a format, but applied none", name)
		}

		if other, isSeen := seen[formatted]; isSeen {
			t.Errorf("%s must differ from %s, but applied the same codes", name, other)
		}

		seen[formatted] = name
	}
}
