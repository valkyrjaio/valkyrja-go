/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
)

// bannerPadding is the space that a banner puts on each side of its text.
const bannerPadding = "    "

// NewErrorMessage builds a message that the terminal prints as an error.
func NewErrorMessage(text string) *Message {
	return NewFormattedMessage(text, format.NewErrorFormatter())
}

// NewSuccessMessage builds a message that the terminal prints as a success.
func NewSuccessMessage(text string) *Message {
	return NewFormattedMessage(text, format.NewSuccessFormatter())
}

// NewWarningMessage builds a message that the terminal prints as a warning.
func NewWarningMessage(text string) *Message {
	return NewFormattedMessage(text, format.NewWarningFormatter())
}

// NewNewLine builds a message that carries a line break alone.
func NewNewLine() *Message {
	return NewMessage("\n")
}

// Messages is several messages that a caller passes as one.
type Messages struct {
	messages []contract.MessageContract
}

// NewMessages builds one message from the messages.
func NewMessages(messages ...contract.MessageContract) *Messages {
	return &Messages{messages: messages}
}

// NewBanner builds a message that prints the message on a block of its own
// color.
//
// The banner pads the text on each side, and it prints a line of the same width
// above and below the text, so the color reads as a block rather than a line.
func NewBanner(message contract.MessageContract) *Messages {
	text := bannerPadding + message.GetText() + bannerPadding
	spaces := strings.Repeat(" ", len(text))

	return NewMessages(
		NewNewLine(),
		message.WithText(spaces),
		NewNewLine(),
		message.WithText(text),
		NewNewLine(),
		message.WithText(spaces),
		NewNewLine(),
	)
}

// GetText returns the text of each message, one after the other.
func (m *Messages) GetText() string {
	return m.join(contract.MessageContract.GetText)
}

// GetFormattedText returns the formatted text of each message, one after the
// other.
func (m *Messages) GetFormattedText() string {
	return m.join(contract.MessageContract.GetFormattedText)
}

// WithText returns a copy that carries the text as its one message.
//
// The other ports let a caller replace the text of a group. A group holds
// several messages, each with its own formatter, so one text replaces all of
// them.
func (m *Messages) WithText(text string) contract.MessageContract {
	return NewMessages(NewMessage(text))
}

// HasFormatter reports that a group carries no formatter of its own.
func (m *Messages) HasFormatter() bool {
	return false
}

// GetFormatter returns nil, because each message in the group carries its own
// formatter.
func (m *Messages) GetFormatter() contract.FormatterContract {
	return nil
}

// WithFormatter returns a copy in which every message carries the formatter.
func (m *Messages) WithFormatter(formatter contract.FormatterContract) contract.MessageContract {
	return NewMessages(m.mapMessages(func(held contract.MessageContract) contract.MessageContract {
		return held.WithFormatter(formatter)
	})...)
}

// WithoutFormatter returns a copy in which no message carries a formatter.
func (m *Messages) WithoutFormatter() contract.MessageContract {
	return NewMessages(m.mapMessages(contract.MessageContract.WithoutFormatter)...)
}

// join renders each message with the reader and returns the results as one
// string.
func (m *Messages) join(read func(contract.MessageContract) string) string {
	built := &strings.Builder{}

	for _, held := range m.messages {
		built.WriteString(read(held))
	}

	return built.String()
}

// mapMessages returns each message, passed through the change.
func (m *Messages) mapMessages(
	change func(contract.MessageContract) contract.MessageContract,
) []contract.MessageContract {
	changed := make([]contract.MessageContract, 0, len(m.messages))

	for _, held := range m.messages {
		changed = append(changed, change(held))
	}

	return changed
}
