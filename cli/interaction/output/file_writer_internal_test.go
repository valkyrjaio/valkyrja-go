/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package output

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/throwable/exception"
)

// errWrite is what a file that refuses a write reports.
var errWrite = errors.New("the file took no bytes")

// errClose is what a file that refuses a close reports.
var errClose = errors.New("the file did not close")

// refusingFileFixture is a file that opens, and then refuses the write or the
// close that follows.
type refusingFileFixture struct {
	writeErr error
	closeErr error
}

// Write reports the failure that the fixture holds.
func (f *refusingFileFixture) Write(text []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}

	return len(text), nil
}

// Close reports the failure that the fixture holds.
func (f *refusingFileFixture) Close() error {
	return f.closeErr
}

// newRefusingFileWriter builds a writer whose file opens and then refuses.
func newRefusingFileWriter(file *refusingFileFixture) *FileWriter {
	return &FileWriter{
		filepath: "/does/not/matter.log",
		open: func(_ string, _ int, _ os.FileMode) (io.WriteCloser, error) {
			return file, nil
		},
	}
}

func TestAFileThatRefusesTheWriteReachesTheCaller(t *testing.T) {
	t.Parallel()

	writer := newRefusingFileWriter(&refusingFileFixture{writeErr: errWrite})

	_, err := writer.Write(nil, message.NewMessage("the message"))

	if _, isFailure := errors.AsType[*exception.CliInteractionFileWriteError](err); !isFailure {
		t.Fatalf("a file that refuses the write must reach the caller, but reported: %v", err)
	}

	if !errors.Is(err, errWrite) {
		t.Error("the failure must unwrap to what the file reported, but did not")
	}
}

func TestAFileThatRefusesTheCloseReachesTheCaller(t *testing.T) {
	t.Parallel()

	// A close that fails after a whole write still loses the message, because
	// the bytes sit in a buffer that never reaches the file.
	writer := newRefusingFileWriter(&refusingFileFixture{closeErr: errClose})

	_, err := writer.Write(nil, message.NewMessage("the message"))

	if !errors.Is(err, errClose) {
		t.Fatalf("a file that refuses the close must reach the caller, but reported: %v", err)
	}
}

func TestAFileThatRefusesBothReportsTheWrite(t *testing.T) {
	t.Parallel()

	writer := newRefusingFileWriter(&refusingFileFixture{writeErr: errWrite, closeErr: errClose})

	_, err := writer.Write(nil, message.NewMessage("the message"))

	if !errors.Is(err, errWrite) {
		t.Errorf("the write failure must reach the caller ahead of the close, but reported: %v", err)
	}
}
