/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package output_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
)

const (
	firstText  = "the first"
	secondText = "the second"
)

// newOutput builds an output that writes to the builder.
func newOutput(written *strings.Builder, messages ...contract.MessageContract) contract.OutputContract {
	return output.NewOutput(
		[]contract.WriterContract{output.NewStreamWriter(written)},
		messages...,
	)
}

func TestNewOutputHoldsItsMessagesUnwritten(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := newOutput(written, message.NewMessage(firstText))

	if !built.HasUnwrittenMessage() || built.HasWrittenMessage() {
		t.Error("a new output must hold its messages unwritten, but did not")
	}

	if len(built.GetMessages()) != 1 {
		t.Errorf("GetMessages must return each message, but returned: %d", len(built.GetMessages()))
	}

	if !built.IsInteractive() || built.IsQuiet() || built.IsSilent() {
		t.Error("a new output must be interactive and neither quiet nor silent, but was not")
	}

	if built.GetExitCode() != constant.ExitCodeSuccess {
		t.Errorf("a new output must exit with success, but exits with: %d", built.GetExitCode())
	}
}

func TestWriteMessagesWritesEachMessageInOrder(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := newOutput(written, message.NewMessage(firstText), message.NewMessage(secondText)).
		WriteMessages()

	if written.String() != firstText+"\n"+secondText+"\n" {
		t.Errorf("WriteMessages must write each message in order, but wrote: %q", written.String())
	}

	if built.HasUnwrittenMessage() {
		t.Error("WriteMessages must leave no message unwritten, but did")
	}

	if len(built.GetWrittenMessages()) != 2 {
		t.Errorf("WriteMessages must record each message, but recorded: %d", len(built.GetWrittenMessages()))
	}
}

func TestWriteMessageRecordsTheMessage(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := newOutput(written).WriteMessage(message.NewMessage(firstText))

	if !built.HasWrittenMessage() {
		t.Error("WriteMessage must record the message, but did not")
	}

	if written.String() != firstText+"\n" {
		t.Errorf("WriteMessage must write the message, but wrote: %q", written.String())
	}
}

func TestASilentOutputWritesNothing(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := newOutput(written, message.NewMessage(firstText)).
		WithIsSilent(true).
		WriteMessages()

	if written.String() != "" {
		t.Errorf("a silent output must write nothing, but wrote: %q", written.String())
	}

	if !built.HasWrittenMessage() {
		t.Error("a silent output must still record the message, but did not")
	}
}

func TestAQuietOutputWritesNothingWhileItReportsSuccess(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	newOutput(written, message.NewMessage(firstText)).
		WithIsQuiet(true).
		WriteMessages()

	if written.String() != "" {
		t.Errorf("a quiet output must write nothing while it reports success, but wrote: %q", written.String())
	}
}

func TestAQuietOutputWritesWhereItReportsAFailure(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	newOutput(written, message.NewMessage(firstText)).
		WithIsQuiet(true).
		WithExitCode(constant.ExitCodeError).
		WriteMessages()

	if written.String() != firstText+"\n" {
		t.Errorf("a quiet output must write where it reports a failure, but wrote: %q", written.String())
	}
}

func TestEachMessageMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := newOutput(written, message.NewMessage(firstText))

	replaced := built.WithMessages(message.NewMessage(secondText))
	added := built.WithAddedMessages(message.NewMessage(secondText))
	addedOne := built.WithAddedMessage(message.NewMessage(secondText))

	if len(replaced.GetUnwrittenMessages()) != 1 ||
		replaced.GetUnwrittenMessages()[0].GetText() != secondText {
		t.Error("WithMessages must replace the messages, but did not")
	}

	if len(added.GetUnwrittenMessages()) != 2 || len(addedOne.GetUnwrittenMessages()) != 2 {
		t.Error("each added method must append the message, but did not")
	}

	if len(built.GetUnwrittenMessages()) != 1 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachFlagMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := newOutput(written)

	if built.WithIsInteractive(false).IsInteractive() {
		t.Error("WithIsInteractive must hold the new flag, but did not")
	}

	if !built.WithIsQuiet(true).IsQuiet() {
		t.Error("WithIsQuiet must hold the new flag, but did not")
	}

	if !built.WithIsSilent(true).IsSilent() {
		t.Error("WithIsSilent must hold the new flag, but did not")
	}

	if built.WithExitCode(constant.ExitCodeUsageError).GetExitCode() != constant.ExitCodeUsageError {
		t.Error("WithExitCode must hold the new code, but did not")
	}

	if !built.IsInteractive() || built.IsQuiet() || built.IsSilent() {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestWithWritersReturnsACopy(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}
	other := &strings.Builder{}

	built := newOutput(written)

	replaced := built.WithWriters(output.NewPlainWriter(other))

	replaced.WriteMessage(message.NewFormattedMessage(
		firstText,
		format.NewFormatter(format.NewTextColorFormat(constant.TextColorRed)),
	))

	if other.String() != firstText+"\n" {
		t.Errorf("the plain writer must write the text with no format, but wrote: %q", other.String())
	}

	if written.String() != "" {
		t.Error("WithWriters must replace the writers, but the old one wrote")
	}

	if len(built.GetWriters()) != 1 {
		t.Error("WithWriters must leave the receiver unchanged, but did not")
	}
}

func TestTheStreamWriterWritesTheFormattedText(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	newOutput(written).WriteMessage(message.NewFormattedMessage(
		firstText,
		format.NewFormatter(format.NewTextColorFormat(constant.TextColorRed)),
	))

	if written.String() != "\x1b[31mthe first\x1b[39m\n" {
		t.Errorf("the stream writer must write the formatted text, but wrote: %q", written.String())
	}
}

func TestAWriterThatSkipsAMessageWritesNothing(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := output.NewOutput([]contract.WriterContract{&skippingWriterFixture{}})

	built.WriteMessage(message.NewMessage(firstText))

	if written.String() != "" {
		t.Errorf("a writer that skips the message must write nothing, but wrote: %q", written.String())
	}
}

// skippingWriterFixture is a writer that writes no message at all.
type skippingWriterFixture struct{}

// ShouldWriteMessage reports that this writer writes no message.
func (w *skippingWriterFixture) ShouldWriteMessage(_ contract.MessageContract) bool {
	return false
}

// Write returns the output as it received it.
func (w *skippingWriterFixture) Write(
	output contract.OutputContract,
	_ contract.MessageContract,
) contract.OutputContract {
	return output
}

func TestTheOutputSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	built := newOutput(written)

	if built.GetExitCode() != constant.ExitCodeSuccess {
		t.Error("the contract must read the exit code, but did not")
	}
}
