/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
)

const kindText = "Something happened."

func TestEachKindOfMessageCarriesItsOwnFormatter(t *testing.T) {
	t.Parallel()

	kinds := map[string]interface{ GetFormattedText() string }{
		"an error message":  message.NewErrorMessage(kindText),
		"a success message": message.NewSuccessMessage(kindText),
		"a warning message": message.NewWarningMessage(kindText),
	}

	for name, built := range kinds {
		formatted := built.GetFormattedText()

		if !strings.Contains(formatted, kindText) {
			t.Errorf("%s must carry its text, but carried: %q", name, formatted)
		}

		if formatted == kindText {
			t.Errorf("%s must apply a format, but applied none", name)
		}
	}
}

func TestANewLineCarriesALineBreakAlone(t *testing.T) {
	t.Parallel()

	if message.NewNewLine().GetText() != "\n" {
		t.Error("a new line must carry a line break alone, but did not")
	}
}

func TestAGroupJoinsTheTextOfEachMessage(t *testing.T) {
	t.Parallel()

	built := message.NewMessages(message.NewMessage("one"), message.NewMessage("two"))

	if built.GetText() != "onetwo" {
		t.Errorf("a group must join the text of each message, but returned: %q", built.GetText())
	}

	if built.GetFormattedText() != "onetwo" {
		t.Errorf("a group of plain messages must apply no format, but returned: %q", built.GetFormattedText())
	}
}

func TestAGroupCarriesNoFormatterOfItsOwn(t *testing.T) {
	t.Parallel()

	built := message.NewMessages(message.NewMessage("one"))

	if built.HasFormatter() || built.GetFormatter() != nil {
		t.Error("a group must carry no formatter of its own, but carried one")
	}

	formatted := built.WithFormatter(format.NewErrorFormatter())
	if formatted.GetFormattedText() == "one" {
		t.Error("WithFormatter must give the formatter to each message, but did not")
	}

	if formatted.WithoutFormatter().GetFormattedText() != "one" {
		t.Error("WithoutFormatter must remove the formatter from each message, but did not")
	}
}

func TestAGroupWithATextCarriesThatTextAlone(t *testing.T) {
	t.Parallel()

	built := message.NewMessages(message.NewMessage("one"), message.NewMessage("two"))

	if built.WithText("three").GetText() != "three" {
		t.Error("WithText must replace every message with the text, but did not")
	}
}

func TestABannerFramesItsMessage(t *testing.T) {
	t.Parallel()

	built := message.NewBanner(message.NewMessage("Hi"))

	lines := strings.Split(built.GetText(), "\n")

	// The banner opens and closes with a line break, so the first and the last
	// piece are empty and the three lines sit between them.
	if len(lines) != 5 {
		t.Fatalf("a banner must print three lines between four line breaks, but printed: %d", len(lines))
	}

	if !strings.Contains(lines[2], "Hi") {
		t.Errorf("the middle line of a banner must carry the text, but carried: %q", lines[2])
	}

	if len(lines[1]) != len(lines[2]) || len(lines[3]) != len(lines[2]) {
		t.Error("each line of a banner must be as wide as the text, but was not")
	}
}
