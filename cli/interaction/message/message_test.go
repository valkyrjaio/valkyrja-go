/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
)

const messageText = "the message"

func TestTheMessageHoldsItsText(t *testing.T) {
	t.Parallel()

	built := message.NewMessage(messageText)

	if built.GetText() != messageText {
		t.Errorf("the message must hold its text, but holds: %q", built.GetText())
	}

	if built.HasFormatter() || built.GetFormatter() != nil {
		t.Error("a message must carry no formatter by default, but carried one")
	}

	if built.GetFormattedText() != messageText {
		t.Error("a message with no formatter must return its text as it is, but did not")
	}
}

func TestAFormattedMessageWrapsItsText(t *testing.T) {
	t.Parallel()

	built := message.NewFormattedMessage(
		messageText,
		format.NewFormatter(format.NewTextColorFormat(constant.TextColorRed)),
	)

	if !built.HasFormatter() {
		t.Error("the message must carry a formatter, but did not")
	}

	if built.GetFormattedText() != "\x1b[31mthe message\x1b[39m" {
		t.Errorf("the message must wrap its text, but is: %q", built.GetFormattedText())
	}

	if built.GetText() != messageText {
		t.Error("GetText must return the text as it is, but did not")
	}
}

func TestEachMessageWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := message.NewMessage(messageText)
	formatter := format.NewFormatter(format.NewTextColorFormat(constant.TextColorRed))

	if built.WithText("other").GetText() != "other" {
		t.Error("WithText must hold the new text, but did not")
	}

	if !built.WithFormatter(formatter).HasFormatter() {
		t.Error("WithFormatter must hold the new formatter, but did not")
	}

	if built.HasFormatter() {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestWithoutFormatterRemovesTheFormatter(t *testing.T) {
	t.Parallel()

	built := message.NewFormattedMessage(
		messageText,
		format.NewFormatter(format.NewTextColorFormat(constant.TextColorRed)),
	)

	without := built.WithoutFormatter()

	if without.HasFormatter() || without.GetFormattedText() != messageText {
		t.Error("WithoutFormatter must remove the formatter, but did not")
	}

	if !built.HasFormatter() {
		t.Error("WithoutFormatter must leave the receiver unchanged, but did not")
	}
}

func TestTheMessageSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.MessageContract = message.NewMessage(messageText)

	if built.GetText() != messageText {
		t.Error("the contract must read the text, but did not")
	}
}
