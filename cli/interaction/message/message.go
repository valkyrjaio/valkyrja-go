/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package message holds one message that a command writes.
package message

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

// Message is one message that a command writes.
type Message struct {
	text      string
	formatter contract.FormatterContract
}

// NewMessage builds a message that carries the text, with no formatter.
func NewMessage(text string) *Message {
	return &Message{text: text}
}

// NewFormattedMessage builds a message that the formatter wraps.
func NewFormattedMessage(text string, formatter contract.FormatterContract) *Message {
	return &Message{
		text:      text,
		formatter: formatter,
	}
}

// GetText returns the text of the message.
func (m *Message) GetText() string {
	return m.text
}

// GetFormattedText returns the text with the formatter applied, and the text as
// it is where the message carries no formatter.
func (m *Message) GetFormattedText() string {
	if m.formatter == nil {
		return m.text
	}

	return m.formatter.FormatText(m.text)
}

// WithText returns a copy of the message with another text.
func (m *Message) WithText(text string) contract.MessageContract {
	copied := *m
	copied.text = text

	return &copied
}

// HasFormatter reports whether the message carries a formatter.
func (m *Message) HasFormatter() bool {
	return m.formatter != nil
}

// GetFormatter returns the formatter of the message.
func (m *Message) GetFormatter() contract.FormatterContract {
	return m.formatter
}

// WithFormatter returns a copy of the message with another formatter.
func (m *Message) WithFormatter(formatter contract.FormatterContract) contract.MessageContract {
	copied := *m
	copied.formatter = formatter

	return &copied
}

// WithoutFormatter returns a copy of the message with no formatter.
func (m *Message) WithoutFormatter() contract.MessageContract {
	copied := *m
	copied.formatter = nil

	return &copied
}
