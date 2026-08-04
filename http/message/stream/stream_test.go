/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package stream_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

const content = "the body"

// newReadWriteStream builds a stream that a reader reads and a writer writes.
func newReadWriteStream() *stream.Stream {
	return stream.NewStream(content, constant.ModeReadWrite)
}

func TestNewStreamHoldsTheContent(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	if built.GetSize() != len(content) {
		t.Errorf("GetSize must be the size of the content, but is: %d", built.GetSize())
	}

	if built.String() != content {
		t.Errorf("String must be the content, but is: %q", built.String())
	}
}

func TestAWriteOnlyStreamStartsAtTheEnd(t *testing.T) {
	t.Parallel()

	built := stream.NewStream(content, constant.ModeWrite)

	position, err := built.Tell()
	if err != nil {
		t.Fatalf("Tell must return the position, but reported: %v", err)
	}

	if position != len(content) {
		t.Errorf("a write-only stream must start at the end, but is at: %d", position)
	}
}

func TestAReadWriteStreamStartsAtTheBeginning(t *testing.T) {
	t.Parallel()

	position, err := newReadWriteStream().Tell()
	if err != nil {
		t.Fatalf("Tell must return the position, but reported: %v", err)
	}

	if position != 0 {
		t.Errorf("a read-write stream must start at the beginning, but is at: %d", position)
	}
}

func TestReadMovesThePointer(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	read, err := built.Read(3)
	if err != nil {
		t.Fatalf("Read must read the bytes, but reported: %v", err)
	}

	if read != "the" {
		t.Errorf("Read must return the bytes, but returned: %q", read)
	}

	contents, err := built.GetContents()
	if err != nil {
		t.Fatalf("GetContents must read what is left, but reported: %v", err)
	}

	if contents != " body" {
		t.Errorf("GetContents must return what is left, but returned: %q", contents)
	}
}

func TestReadStopsAtTheEndOfTheStream(t *testing.T) {
	t.Parallel()

	read, err := newReadWriteStream().Read(1000)
	if err != nil {
		t.Fatalf("Read must read the bytes, but reported: %v", err)
	}

	if read != content {
		t.Errorf("Read must stop at the end of the stream, but returned: %q", read)
	}
}

func TestReadReportsALengthUnderZero(t *testing.T) {
	t.Parallel()

	_, err := newReadWriteStream().Read(-1)

	target, found := errors.AsType[*exception.HttpStreamInvalidLengthError](err)
	if !found {
		t.Fatalf("Read must report an invalid length, but reported: %v", err)
	}

	if target.GetLength() != -1 {
		t.Errorf("the error must carry the length, but carries: %d", target.GetLength())
	}
}

func TestReadReportsAStreamThatNoReaderReads(t *testing.T) {
	t.Parallel()

	built := stream.NewStream(content, constant.ModeWrite)

	_, readErr := built.Read(1)

	if _, found := errors.AsType[*exception.HttpStreamUnreadableStreamError](readErr); !found {
		t.Errorf("Read must report an unreadable stream, but reported: %v", readErr)
	}

	_, err := built.GetContents()

	if _, found := errors.AsType[*exception.HttpStreamUnreadableStreamError](err); !found {
		t.Errorf("GetContents must report an unreadable stream, but reported: %v", err)
	}
}

func TestWriteOverwritesFromThePointer(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	written, err := built.Write("THE")
	if err != nil {
		t.Fatalf("Write must write the bytes, but reported: %v", err)
	}

	if written != 3 {
		t.Errorf("Write must return the number of bytes, but returned: %d", written)
	}

	if built.String() != "THE body" {
		t.Errorf("Write must overwrite from the pointer, but the stream is: %q", built.String())
	}
}

func TestWriteAppendsPastTheEnd(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	err := built.Seek(0, 2)
	if err != nil {
		t.Fatalf("Seek must move the pointer, but reported: %v", err)
	}

	_, err = built.Write(" and more")
	if err != nil {
		t.Fatalf("Write must write the bytes, but reported: %v", err)
	}

	if built.String() != "the body and more" {
		t.Errorf("Write must append past the end, but the stream is: %q", built.String())
	}
}

func TestWriteReportsAStreamThatNoWriterWrites(t *testing.T) {
	t.Parallel()

	_, err := stream.NewStream(content, constant.ModeRead).Write("more")

	if _, found := errors.AsType[*exception.HttpStreamUnwritableStreamError](err); !found {
		t.Errorf("Write must report an unwritable stream, but reported: %v", err)
	}
}

func TestSeekMeasuresFromEachPoint(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		offset   int
		whence   int
		expected int
	}{
		"from the start":   {3, stream.WhenceStart, 3},
		"from the pointer": {2, stream.WhenceCurrent, 2},
		"from the end":     {-2, stream.WhenceEnd, len(content) - 2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			built := newReadWriteStream()

			err := built.Seek(test.offset, test.whence)
			if err != nil {
				t.Fatalf("Seek must move the pointer, but reported: %v", err)
			}

			position, _ := built.Tell()

			if position != test.expected {
				t.Errorf("Seek %s must reach %d, but reached: %d", name, test.expected, position)
			}
		})
	}
}

func TestSeekReportsAPositionOutsideTheStream(t *testing.T) {
	t.Parallel()

	offsets := []int{-1, len(content) + 1}

	for _, offset := range offsets {
		err := newReadWriteStream().Seek(offset, stream.WhenceStart)

		if _, found := errors.AsType[*exception.HttpStreamStreamSeekError](err); !found {
			t.Errorf("Seek to %d must report a failure, but reported: %v", offset, err)
		}
	}
}

func TestRewindReturnsThePointerToTheStart(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	_, err := built.Read(3)
	if err != nil {
		t.Fatalf("Read must read the bytes, but reported: %v", err)
	}

	err = built.Rewind()
	if err != nil {
		t.Fatalf("Rewind must move the pointer, but reported: %v", err)
	}

	position, _ := built.Tell()

	if position != 0 {
		t.Errorf("Rewind must return the pointer to the start, but it is at: %d", position)
	}
}

func TestIsEofReportsThePointerAtTheEnd(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	if built.IsEof() {
		t.Error("IsEof must be false at the start, but is true")
	}

	_, err := built.GetContents()
	if err != nil {
		t.Fatalf("GetContents must read the stream, but reported: %v", err)
	}

	if !built.IsEof() {
		t.Error("IsEof must be true at the end, but is false")
	}
}

func TestCloseDropsWhatTheStreamHeld(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	err := built.Close()
	if err != nil {
		t.Fatalf("Close must close the stream, but reported: %v", err)
	}

	if built.GetSize() != 0 {
		t.Errorf("Close must drop the content, but the size is: %d", built.GetSize())
	}

	if !built.IsEof() {
		t.Error("IsEof must be true for a closed stream, but is false")
	}

	if built.IsSeekable() || built.IsReadable() || built.IsWritable() {
		t.Error("a closed stream must be neither seekable, readable, nor writable, but is")
	}
}

func TestTellReportsAClosedStream(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	err := built.Close()
	if err != nil {
		t.Fatalf("Close must close the stream, but reported: %v", err)
	}

	_, err = built.Tell()

	if _, found := errors.AsType[*exception.HttpStreamStreamTellError](err); !found {
		t.Errorf("Tell must report a closed stream, but reported: %v", err)
	}
}

func TestSeekReportsAClosedStream(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	err := built.Close()
	if err != nil {
		t.Fatalf("Close must close the stream, but reported: %v", err)
	}

	err = built.Seek(0, stream.WhenceStart)

	if _, found := errors.AsType[*exception.HttpStreamUnseekableStreamError](err); !found {
		t.Errorf("Seek must report a closed stream, but reported: %v", err)
	}
}

func TestDetachReturnsTheBufferOnce(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	if string(built.Detach()) != content {
		t.Error("Detach must return the buffer, but did not")
	}

	if built.Detach() != nil {
		t.Error("Detach must return nil for a stream that is detached already, but did not")
	}
}

func TestGetMetadataReportsTheModeAndTheSeekableFlag(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	metadata := built.GetMetadata()

	if metadata[stream.MetadataKeySeekable] != true {
		t.Errorf("the metadata must report the seekable flag, but holds: %v", metadata)
	}

	if metadata[stream.MetadataKeyMode] != constant.ModeReadWrite {
		t.Errorf("the metadata must report the mode, but holds: %v", metadata)
	}

	if built.GetMetadataItem(stream.MetadataKeyMode) != constant.ModeReadWrite {
		t.Error("GetMetadataItem must return the mode, but did not")
	}

	if built.GetMetadataItem("unknown") != nil {
		t.Error("GetMetadataItem must be nil for an unknown key, but is not")
	}
}

func TestGetMetadataIsEmptyForAClosedStream(t *testing.T) {
	t.Parallel()

	built := newReadWriteStream()

	err := built.Close()
	if err != nil {
		t.Fatalf("Close must close the stream, but reported: %v", err)
	}

	if len(built.GetMetadata()) != 0 {
		t.Errorf("the metadata of a closed stream must be empty, but holds: %v", built.GetMetadata())
	}
}

func TestStringIsEmptyWhereNoReaderReadsTheStream(t *testing.T) {
	t.Parallel()

	if stream.NewStream(content, constant.ModeWrite).String() != "" {
		t.Error("String must be empty where no reader reads the stream, but is not")
	}
}

func TestTheStreamSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.StreamContract = newReadWriteStream()

	if built.GetSize() != len(content) {
		t.Errorf("the contract must report the size, but reported: %d", built.GetSize())
	}
}
