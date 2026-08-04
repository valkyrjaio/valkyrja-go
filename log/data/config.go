/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds the log component's configuration.
package data

// DefaultLogger is the binding key of the logger that an application writes
// through until it states another one.
const DefaultLogger = "Valkyrja.Log.Logger.StreamLogger"

// LogConfig holds the settings that apply to the whole log component.
//
// The component config holds the component-wide setting alone. Each adapter gets
// a config contract of its own, so an application never builds the
// configuration of an adapter that it does not use.
type LogConfig struct {
	DefaultLoggerServiceID string
}

// NewLogConfig builds the configuration that every field takes its default value
// in.
func NewLogConfig() *LogConfig {
	return &LogConfig{DefaultLoggerServiceID: DefaultLogger}
}

// GetDefaultLogger returns the binding key of the logger that the application
// writes through.
func (c *LogConfig) GetDefaultLogger() string {
	return c.DefaultLoggerServiceID
}

// LogStreamConfig holds the settings of the logger that writes to a stream.
//
// Every property carries the adapter name, so one application config class can
// hold the settings of several adapters at once.
type LogStreamConfig struct {
	StreamFilePath string
}

// NewLogStreamConfig builds the configuration that every field takes its default
// value in.
//
// An empty path writes to the standard error of the process, which is where a
// application with no log file reads its messages.
func NewLogStreamConfig() *LogStreamConfig {
	return &LogStreamConfig{}
}

// GetStreamFilePath returns the file that the logger writes to, and an empty
// string where it writes to the standard error of the process.
func (c *LogStreamConfig) GetStreamFilePath() string {
	return c.StreamFilePath
}
