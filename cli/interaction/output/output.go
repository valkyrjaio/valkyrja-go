/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package output holds what a command writes back.
package output

import (
	"io"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
)

type Output struct {
	unwrittenMessages []contract.MessageContract
	writtenMessages   []contract.MessageContract
	writers           []contract.WriterContract

	interactive bool
	quiet       bool
	silent      bool
	exitCode    constant.ExitCode
}

// NewOutput builds an output that writes through the writers.
func NewOutput(writers []contract.WriterContract, messages ...contract.MessageContract) *Output {
	return &Output{
		unwrittenMessages: messages,
		writtenMessages:   []contract.MessageContract{},
		writers:           writers,
		interactive:       true,
		exitCode:          constant.ExitCodeSuccess,
	}
}

// GetMessages returns every message that the output holds, written first.
func (o *Output) GetMessages() []contract.MessageContract {
	messages := make([]contract.MessageContract, 0, len(o.writtenMessages)+len(o.unwrittenMessages))
	messages = append(messages, o.writtenMessages...)
	messages = append(messages, o.unwrittenMessages...)

	return messages
}

// GetWrittenMessages returns each message that a writer wrote already.
func (o *Output) GetWrittenMessages() []contract.MessageContract {
	return o.writtenMessages
}

// HasWrittenMessage reports whether a writer wrote a message already.
func (o *Output) HasWrittenMessage() bool {
	return len(o.writtenMessages) > 0
}

// GetUnwrittenMessages returns each message that no writer wrote yet.
func (o *Output) GetUnwrittenMessages() []contract.MessageContract {
	return o.unwrittenMessages
}

// HasUnwrittenMessage reports whether a message is waiting to be written.
func (o *Output) HasUnwrittenMessage() bool {
	return len(o.unwrittenMessages) > 0
}

// WithMessages returns a copy of the output that holds the messages and nothing
// else.
func (o *Output) WithMessages(messages ...contract.MessageContract) contract.OutputContract {
	copied := *o
	copied.unwrittenMessages = messages

	return &copied
}

// WithAddedMessages returns a copy of the output with the messages appended.
func (o *Output) WithAddedMessages(messages ...contract.MessageContract) contract.OutputContract {
	copied := *o
	copied.unwrittenMessages = appendMessages(o.unwrittenMessages, messages)

	return &copied
}

// WithAddedMessage returns a copy of the output with the message appended.
func (o *Output) WithAddedMessage(message contract.MessageContract) contract.OutputContract {
	return o.WithAddedMessages(message)
}

// WriteMessages writes every message that is waiting.
func (o *Output) WriteMessages() contract.OutputContract {
	written := contract.OutputContract(o)

	for _, message := range o.unwrittenMessages {
		written = written.WriteMessage(message)
	}

	return written.WithMessages()
}

// WriteMessage writes one message, and records it as written.
func (o *Output) WriteMessage(message contract.MessageContract) contract.OutputContract {
	copied := *o
	copied.writtenMessages = appendMessages(o.writtenMessages, []contract.MessageContract{message})

	if o.silent || (o.quiet && o.exitCode == constant.ExitCodeSuccess) {
		return &copied
	}

	for _, writer := range o.writers {
		if !writer.ShouldWriteMessage(message) {
			continue
		}

		writer.Write(&copied, message)
	}

	return &copied
}

// GetWriters returns each writer that the output writes through.
func (o *Output) GetWriters() []contract.WriterContract {
	return o.writers
}

// WithWriters returns a copy of the output with other writers.
func (o *Output) WithWriters(writers ...contract.WriterContract) contract.OutputContract {
	copied := *o
	copied.writers = writers

	return &copied
}

// IsInteractive reports whether the output asks the caller a question.
func (o *Output) IsInteractive() bool {
	return o.interactive
}

// WithIsInteractive returns a copy of the output with another interactive flag.
func (o *Output) WithIsInteractive(isInteractive bool) contract.OutputContract {
	copied := *o
	copied.interactive = isInteractive

	return &copied
}

// IsQuiet reports whether the output writes less.
func (o *Output) IsQuiet() bool {
	return o.quiet
}

// WithIsQuiet returns a copy of the output with another quiet flag.
func (o *Output) WithIsQuiet(isQuiet bool) contract.OutputContract {
	copied := *o
	copied.quiet = isQuiet

	return &copied
}

// IsSilent reports whether the output writes nothing.
func (o *Output) IsSilent() bool {
	return o.silent
}

// WithIsSilent returns a copy of the output with another silent flag.
func (o *Output) WithIsSilent(isSilent bool) contract.OutputContract {
	copied := *o
	copied.silent = isSilent

	return &copied
}

// GetExitCode returns the code that the process exits with.
func (o *Output) GetExitCode() constant.ExitCode {
	return o.exitCode
}

// WithExitCode returns a copy of the output for another exit code.
func (o *Output) WithExitCode(exitCode constant.ExitCode) contract.OutputContract {
	copied := *o
	copied.exitCode = exitCode

	return &copied
}

// appendMessages returns the messages with the added ones after them, in a slice
// of its own, so a copy never shares a backing array with the output it came
// from.
func appendMessages(
	messages []contract.MessageContract,
	added []contract.MessageContract,
) []contract.MessageContract {
	combined := make([]contract.MessageContract, 0, len(messages)+len(added))
	combined = append(combined, messages...)
	combined = append(combined, added...)

	return combined
}

type StreamWriter struct {
	writer io.Writer
}

// NewStreamWriter builds a writer over a stream.
func NewStreamWriter(writer io.Writer) *StreamWriter {
	return &StreamWriter{writer: writer}
}

// ShouldWriteMessage reports that this writer writes every message.
func (w *StreamWriter) ShouldWriteMessage(_ contract.MessageContract) bool {
	return true
}

// Write writes the formatted text of the message, followed by a line break.
func (w *StreamWriter) Write(
	output contract.OutputContract,
	message contract.MessageContract,
) contract.OutputContract {
	_, _ = io.WriteString(w.writer, message.GetFormattedText()+"\n")

	return output
}

type PlainWriter struct {
	writer io.Writer
}

// NewPlainWriter builds a writer that applies no format.
func NewPlainWriter(writer io.Writer) *PlainWriter {
	return &PlainWriter{writer: writer}
}

// ShouldWriteMessage reports that this writer writes every message.
func (w *PlainWriter) ShouldWriteMessage(_ contract.MessageContract) bool {
	return true
}

// Write writes the text of the message, with no format, followed by a line
// break.
func (w *PlainWriter) Write(
	output contract.OutputContract,
	message contract.MessageContract,
) contract.OutputContract {
	_, _ = io.WriteString(w.writer, message.GetText()+"\n")

	return output
}
