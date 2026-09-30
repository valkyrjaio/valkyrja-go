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
	"os"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/throwable/exception"
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
func (o *Output) WriteMessages() (contract.OutputContract, error) {
	written := contract.OutputContract(o)

	for _, message := range o.unwrittenMessages {
		next, err := written.WriteMessage(message)
		if err != nil {
			return next, err
		}

		written = next
	}

	return written.WithMessages(), nil
}

// WriteMessage writes one message, and records it as written.
func (o *Output) WriteMessage(message contract.MessageContract) (contract.OutputContract, error) {
	copied := *o
	copied.writtenMessages = appendMessages(o.writtenMessages, []contract.MessageContract{message})

	if o.silent || (o.quiet && o.exitCode == constant.ExitCodeSuccess) {
		return &copied, nil
	}

	for _, writer := range o.writers {
		if !writer.ShouldWriteMessage(message) {
			continue
		}

		_, err := writer.Write(&copied, message)
		if err != nil {
			return &copied, err
		}
	}

	return &copied, nil
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
) (contract.OutputContract, error) {
	if w.writer == nil {
		return output, exception.NewCliInteractionUnwritableStreamError()
	}

	return output, writeText(w.writer, message.GetFormattedText()+"\n")
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
) (contract.OutputContract, error) {
	if w.writer == nil {
		return output, exception.NewCliInteractionUnwritableStreamError()
	}

	return output, writeText(w.writer, message.GetText()+"\n")
}

// filePermissions is what a new output file is created with. The owner reads and
// writes the file, and every other user reads it.
const filePermissions = 0o644

// fileOpener opens the file that a writer appends to.
//
// The field is a seam. A file that opens and then refuses a write, or refuses a
// close, is not something a test can ask a real file system for.
type fileOpener func(name string, flag int, perm os.FileMode) (io.WriteCloser, error)

type FileWriter struct {
	filepath string
	open     fileOpener
}

// NewFileWriter builds a writer over the file at the path.
func NewFileWriter(filepath string) *FileWriter {
	return &FileWriter{filepath: filepath, open: openFile}
}

// openFile opens the file with the flags and the permissions that a caller
// names.
// The application names the file, the same way it does in every other port.
//
//nolint:gosec // The path is the application's own, not a client's.
func openFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return os.OpenFile(name, flag, perm)
}

// ShouldWriteMessage reports that this writer writes every message.
func (w *FileWriter) ShouldWriteMessage(_ contract.MessageContract) bool {
	return true
}

// Write appends the formatted text of the message to the file, followed by a
// line break.
//
// The writer opens the file for each message, because an output writes one
// message at a time and a write that truncates keeps only the last one.
func (w *FileWriter) Write(
	output contract.OutputContract,
	message contract.MessageContract,
) (contract.OutputContract, error) {
	file, err := w.open(w.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermissions)
	if err != nil {
		return output, exception.NewCliInteractionFileWriteError(w.filepath, err)
	}

	writeErr := writeText(file, message.GetFormattedText()+"\n")

	closeErr := file.Close()
	if writeErr == nil && closeErr != nil {
		writeErr = closeErr
	}

	if writeErr != nil {
		return output, exception.NewCliInteractionFileWriteError(w.filepath, writeErr)
	}

	return output, nil
}

// writeText puts the text on the writer, and reports a failure where the writer
// takes less than the whole text.
func writeText(writer io.Writer, text string) error {
	written, err := io.WriteString(writer, text)
	if err != nil {
		return exception.NewCliInteractionStreamWriteError(written, len(text), err)
	}

	if written != len(text) {
		return exception.NewCliInteractionStreamWriteError(written, len(text), nil)
	}

	return nil
}
