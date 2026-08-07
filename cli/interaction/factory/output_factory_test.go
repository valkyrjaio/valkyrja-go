/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package factory_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
)

const outputText = "Something happened."

func TestEachKindOfOutputCarriesTheExitCodeAndTheMessages(t *testing.T) {
	t.Parallel()

	built := factory.NewOutputFactory(nil)
	sample := message.NewMessage(outputText)

	outputs := map[string]interface {
		GetMessages() []interface{ GetText() string }
	}{}
	_ = outputs

	created := built.CreateOutput(constant.ExitCodeError, sample)
	if created.GetExitCode() != constant.ExitCodeError || len(created.GetMessages()) != 1 {
		t.Error("the output must carry the exit code and the message, but did not")
	}

	empty := built.CreateEmptyOutput(constant.ExitCodeSuccess, sample)
	if len(empty.GetWriters()) != 0 {
		t.Error("an empty output must write through no writer, but wrote through one")
	}

	plain := built.CreatePlainOutput(constant.ExitCodeSuccess, sample)
	if len(plain.GetWriters()) != 1 {
		t.Error("a plain output must write through one writer, but did not")
	}
}

func TestTheOutputCarriesTheSettingsOfTheConfiguration(t *testing.T) {
	t.Parallel()

	config := data.NewCliInteractionConfigFromValues(true, false, true)

	created := factory.NewOutputFactory(config).CreateOutput(constant.ExitCodeSuccess)

	if !created.IsQuiet() || created.IsInteractive() || !created.IsSilent() {
		t.Error("the output must carry each setting of the configuration, but did not")
	}
}

func TestAStreamOutputWritesToTheWriter(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	factory.NewOutputFactory(nil).
		CreateStreamOutput(written, constant.ExitCodeSuccess, message.NewMessage(outputText)).
		WriteMessages()

	if !strings.Contains(written.String(), outputText) {
		t.Errorf("the output must write its message to the writer, but wrote: %q", written.String())
	}
}

func TestTheDefaultOutputWritesToTheWriterThatTheFactoryHolds(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	factory.NewOutputFactoryForWriter(nil, written).
		CreateOutput(constant.ExitCodeSuccess, message.NewMessage(outputText)).
		WriteMessages()

	if !strings.Contains(written.String(), outputText) {
		t.Errorf("the output must write its message to the writer, but wrote: %q", written.String())
	}
}

func TestAPlainOutputAppliesNoFormat(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	factory.NewOutputFactoryForWriter(nil, written).
		CreatePlainOutput(constant.ExitCodeSuccess, message.NewErrorMessage(outputText)).
		WriteMessages()

	if written.String() != outputText+"\n" {
		t.Errorf("a plain output must apply no format, but wrote: %q", written.String())
	}
}

func TestAFileOutputWritesToTheFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "output.log")

	factory.NewOutputFactory(nil).
		CreateFileOutput(path, constant.ExitCodeSuccess, message.NewMessage(outputText)).
		WriteMessages()

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the output must create the file, but reported: %v", err)
	}

	if !strings.Contains(string(contents), outputText) {
		t.Errorf("the output must write its message to the file, but wrote: %q", contents)
	}
}

func TestAFileOutputThatCannotOpenTheFileWritesNothing(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing", "output.log")

	created := factory.NewOutputFactory(nil).CreateFileOutput(path, constant.ExitCodeError)

	if len(created.GetWriters()) != 0 {
		t.Error("an output that cannot open its file must write through no writer, but wrote through one")
	}

	if created.GetExitCode() != constant.ExitCodeError {
		t.Error("an output that cannot open its file must keep its exit code, but did not")
	}
}
