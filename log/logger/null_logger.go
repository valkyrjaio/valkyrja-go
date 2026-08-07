/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package logger holds every logger that the framework ships.
package logger

import (
	"github.com/valkyrjaio/valkyrja-go/v26/log/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/log/contract"
)

type NullLogger struct{}

// NewNullLogger builds the logger.
func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

// Throwable writes nothing.
func (l *NullLogger) Throwable(_ error, _ string) {}

// Debug writes nothing.
func (l *NullLogger) Debug(_ string, _ map[string]any) {}

// Info writes nothing.
func (l *NullLogger) Info(_ string, _ map[string]any) {}

// Notice writes nothing.
func (l *NullLogger) Notice(_ string, _ map[string]any) {}

// Warning writes nothing.
func (l *NullLogger) Warning(_ string, _ map[string]any) {}

// Error writes nothing.
func (l *NullLogger) Error(_ string, _ map[string]any) {}

// Critical writes nothing.
func (l *NullLogger) Critical(_ string, _ map[string]any) {}

// Alert writes nothing.
func (l *NullLogger) Alert(_ string, _ map[string]any) {}

// Emergency writes nothing.
func (l *NullLogger) Emergency(_ string, _ map[string]any) {}

// Log writes nothing, and reports a severity that the framework does not report.
func (l *NullLogger) Log(level constant.LogLevel, _ string, _ map[string]any) error {
	return validateLevel(level)
}

// A logger satisfies its contract, which the compiler checks at build time
// rather than at run time.
var _ contract.LoggerContract = (*NullLogger)(nil)
