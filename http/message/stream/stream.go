/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package stream holds the body of a message, as a stream that a reader reads
// and a writer writes.
package stream

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

// The point that a seek measures its offset from. The values are the ones that
// every port uses, and they match `io.SeekStart`, `io.SeekCurrent`, and
// `io.SeekEnd`.
const (
	// WhenceStart measures the offset from the start of the stream.
	WhenceStart = 0

	// WhenceCurrent measures the offset from the pointer.
	WhenceCurrent = 1

	// WhenceEnd measures the offset from the end of the stream.
	WhenceEnd = 2
)

// The keys that the metadata of a stream carries.
const (
	// MetadataKeySeekable reports whether a caller seeks in the stream.
	MetadataKeySeekable = "seekable"

	// MetadataKeyMode is the mode that the stream opened in.
	MetadataKeyMode = "mode"
)

// Stream is the body of a message, held in memory.
//
// The other ports hold the content in a buffer of bytes, and so does this one. A
// stream over a file or a network connection is an adapter that this package
// gains later.
type Stream struct {
	buffer   []byte
	position int
	closed   bool
	mode     constant.Mode
}

// NewStream builds a stream over the content, in the mode.
//
// A stream that a writer writes and no reader reads starts at the end of the
// content, so a write appends rather than overwrites.
func NewStream(content string, mode constant.Mode) *Stream {
	built := &Stream{
		buffer: []byte(content),
		mode:   mode,
	}

	if mode.IsWritable() && !mode.IsReadable() {
		built.position = len(built.buffer)
	}

	return built
}

// String returns the whole stream as a string, and an empty string where no
// reader reads it.
//
// The other ports rewind the stream and read it whole. A failure here has no
// return to carry it, so the string is empty instead.
func (s *Stream) String() string {
	if !s.IsReadable() {
		return ""
	}

	// A readable stream is open, so neither the rewind nor the read can fail:
	// each one reports a failure only for a stream that is closed, or that no
	// reader reads, and the guard above covers both.
	_ = s.Rewind()

	contents, _ := s.GetContents()

	return contents
}

// Close closes the stream and drops what it held.
func (s *Stream) Close() error {
	s.closed = true
	s.buffer = []byte{}
	s.position = 0

	return nil
}

// Detach removes the underlying buffer and returns it. It returns nil where the
// stream is closed already.
func (s *Stream) Detach() []byte {
	if s.closed {
		return nil
	}

	detached := s.buffer

	s.closed = true
	s.buffer = []byte{}
	s.position = 0

	return detached
}

// GetSize returns the size of the stream in bytes.
func (s *Stream) GetSize() int {
	return len(s.buffer)
}

// Tell returns the position of the pointer. It reports a failure where the
// stream is closed.
func (s *Stream) Tell() (int, error) {
	if s.closed {
		return 0, exception.NewHttpStreamStreamTellError()
	}

	return s.position, nil
}

// IsEof reports whether the pointer is at the end of the stream.
func (s *Stream) IsEof() bool {
	return s.closed || s.position >= len(s.buffer)
}

// IsSeekable reports whether a caller seeks in the stream.
func (s *Stream) IsSeekable() bool {
	return !s.closed
}

// Seek moves the pointer to the offset, from the point that whence names. It
// reports a failure where the position falls outside the stream.
func (s *Stream) Seek(offset int, whence int) error {
	if !s.IsSeekable() {
		return exception.NewHttpStreamUnseekableStreamError()
	}

	position := s.getPositionFrom(offset, whence)

	if position < 0 || position > len(s.buffer) {
		return exception.NewHttpStreamStreamSeekError()
	}

	s.position = position

	return nil
}

// Rewind moves the pointer to the start of the stream.
func (s *Stream) Rewind() error {
	return s.Seek(0, WhenceStart)
}

// IsWritable reports whether a writer writes to the stream.
func (s *Stream) IsWritable() bool {
	return !s.closed && s.mode.IsWritable()
}

// Write writes the string at the pointer and returns the number of bytes that it
// wrote. It overwrites what the pointer covers, and appends past the end.
func (s *Stream) Write(value string) (int, error) {
	if !s.IsWritable() {
		return 0, exception.NewHttpStreamUnwritableStreamError()
	}

	chunk := []byte(value)

	written := make([]byte, 0, len(s.buffer)+len(chunk))
	written = append(written, s.buffer[:s.position]...)
	written = append(written, chunk...)

	if s.position+len(chunk) < len(s.buffer) {
		written = append(written, s.buffer[s.position+len(chunk):]...)
	}

	s.buffer = written
	s.position += len(chunk)

	return len(chunk), nil
}

// IsReadable reports whether a reader reads the stream.
func (s *Stream) IsReadable() bool {
	return !s.closed && s.mode.IsReadable()
}

// Read reads the number of bytes from the pointer and returns them. It reports a
// failure where the length is under zero, and where no reader reads the stream.
func (s *Stream) Read(length int) (string, error) {
	if length < 0 {
		return "", exception.NewHttpStreamInvalidLengthError(length)
	}

	if !s.IsReadable() {
		return "", exception.NewHttpStreamUnreadableStreamError()
	}

	end := min(s.position+length, len(s.buffer))

	chunk := s.buffer[s.position:end]
	s.position = end

	return string(chunk), nil
}

// GetContents returns what is left of the stream, from the pointer on.
func (s *Stream) GetContents() (string, error) {
	if !s.IsReadable() {
		return "", exception.NewHttpStreamUnreadableStreamError()
	}

	chunk := s.buffer[s.position:]
	s.position = len(s.buffer)

	return string(chunk), nil
}

// GetMetadata returns every metadata item of the stream, and an empty map where
// the stream is closed.
func (s *Stream) GetMetadata() map[string]any {
	if s.closed {
		return map[string]any{}
	}

	return map[string]any{
		MetadataKeySeekable: s.IsSeekable(),
		MetadataKeyMode:     s.mode,
	}
}

// GetMetadataItem returns one metadata item, and nil where the stream holds no
// item under the key.
func (s *Stream) GetMetadataItem(key string) any {
	return s.GetMetadata()[key]
}

// getPositionFrom returns the position that the offset names, measured from the
// point that whence names. An unknown whence measures from the end, which is
// what the other ports do.
func (s *Stream) getPositionFrom(offset int, whence int) int {
	switch whence {
	case WhenceStart:
		return offset
	case WhenceCurrent:
		return s.position + offset
	default:
		return len(s.buffer) + offset
	}
}
