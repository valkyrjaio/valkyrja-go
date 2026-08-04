/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package logger

import (
	"encoding/json"
	"io"
	"slices"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/log/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/log/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/log/throwable/exception"
)

// levelSeparator separates the severity of a message from the message itself.
const levelSeparator = ": "

// contextSeparator separates a message from the context that goes with it.
const contextSeparator = " "

// StreamLogger writes each message to a stream.
//
// A message reads `SEVERITY: text`, with the context after it as JSON where the
// caller gives one. That shape is what every port writes, so a reader of one
// port's log reads another port's log.
type StreamLogger struct {
	writer io.Writer
}

// NewStreamLogger builds the logger over a stream.
func NewStreamLogger(writer io.Writer) *StreamLogger {
	return &StreamLogger{writer: writer}
}

// Throwable writes a failure and the message that goes with it.
func (l *StreamLogger) Throwable(throwable error, message string) {
	l.Error(message, map[string]any{"throwable": throwable.Error()})
}

// Debug writes a message that a developer reads while it debugs.
func (l *StreamLogger) Debug(message string, context map[string]any) {
	l.write(constant.LogLevelDebug, message, context)
}

// Info writes a message that reports what happened.
func (l *StreamLogger) Info(message string, context map[string]any) {
	l.write(constant.LogLevelInfo, message, context)
}

// Notice writes a message that reports something out of the ordinary.
func (l *StreamLogger) Notice(message string, context map[string]any) {
	l.write(constant.LogLevelNotice, message, context)
}

// Warning writes a message that reports a condition to correct.
func (l *StreamLogger) Warning(message string, context map[string]any) {
	l.write(constant.LogLevelWarning, message, context)
}

// Error writes a message that reports a failure.
func (l *StreamLogger) Error(message string, context map[string]any) {
	l.write(constant.LogLevelError, message, context)
}

// Critical writes a message that reports a failure of a component.
func (l *StreamLogger) Critical(message string, context map[string]any) {
	l.write(constant.LogLevelCritical, message, context)
}

// Alert writes a message that reports a failure to correct at once.
func (l *StreamLogger) Alert(message string, context map[string]any) {
	l.write(constant.LogLevelAlert, message, context)
}

// Emergency writes a message that reports the application is unusable.
func (l *StreamLogger) Emergency(message string, context map[string]any) {
	l.write(constant.LogLevelEmergency, message, context)
}

// Log writes the message at the severity, and reports a severity that the
// framework does not report.
func (l *StreamLogger) Log(level constant.LogLevel, message string, context map[string]any) error {
	err := validateLevel(level)
	if err != nil {
		return err
	}

	l.write(level, message, context)

	return nil
}

// write puts one message on the stream.
//
// A stream that reports a failure changes nothing that a caller can act on — a
// logger that reports its own failure has nowhere to report it — so the failure
// is dropped.
func (l *StreamLogger) write(level constant.LogLevel, message string, context map[string]any) {
	line := &strings.Builder{}

	line.WriteString(strings.ToUpper(string(level)))
	line.WriteString(levelSeparator)
	line.WriteString(message)

	if len(context) > 0 {
		line.WriteString(contextSeparator)
		line.WriteString(encodeContext(context))
	}

	line.WriteString("\n")

	_, _ = io.WriteString(l.writer, line.String())
}

// encodeContext renders the context as JSON, and an empty object where no
// encoder can render it.
func encodeContext(context map[string]any) string {
	encoded, err := json.Marshal(context)
	if err != nil {
		return "{}"
	}

	return string(encoded)
}

// validateLevel reports a failure where the severity is not one that the
// framework reports.
func validateLevel(level constant.LogLevel) error {
	known := []constant.LogLevel{
		constant.LogLevelDebug,
		constant.LogLevelInfo,
		constant.LogLevelNotice,
		constant.LogLevelWarning,
		constant.LogLevelError,
		constant.LogLevelCritical,
		constant.LogLevelAlert,
		constant.LogLevelEmergency,
	}

	if slices.Contains(known, level) {
		return nil
	}

	return exception.NewLogInvalidLogLevelError(string(level))
}

// A logger satisfies its contract, which the compiler checks at build time
// rather than at run time.
var _ contract.LoggerContract = (*StreamLogger)(nil)
