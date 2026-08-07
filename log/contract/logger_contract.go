/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds every contract of the log component.
package contract

type LoggerContract interface {
	// Throwable writes a failure and the message that goes with it.
	Throwable(throwable error, message string)

	// Debug writes a message that a developer reads while it debugs.
	Debug(message string, context map[string]any)

	// Info writes a message that reports what happened.
	Info(message string, context map[string]any)

	// Notice writes a message that reports something out of the ordinary.
	Notice(message string, context map[string]any)

	// Warning writes a message that reports a condition to correct.
	Warning(message string, context map[string]any)

	// Error writes a message that reports a failure.
	Error(message string, context map[string]any)

	// Critical writes a message that reports a failure of a component.
	Critical(message string, context map[string]any)

	// Alert writes a message that reports a failure to correct at once.
	Alert(message string, context map[string]any)

	// Emergency writes a message that reports the application is unusable.
	Emergency(message string, context map[string]any)
}

type LogConfigContract interface {
	// GetDefaultLogger returns the binding key of the logger that the
	// application writes through.
	GetDefaultLogger() string
}

type LogStreamConfigContract interface {
	// GetStreamFilePath returns the file that the logger writes to, and an
	// empty string where it writes to the standard error of the process.
	GetStreamFilePath() string
}
