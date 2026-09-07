/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the log component's severities and binding keys.
package constant

type LogLevel string

// The severities that the framework reports.
const (
	LogLevelDebug LogLevel = "debug"

	LogLevelInfo LogLevel = "info"

	LogLevelNotice LogLevel = "notice"

	LogLevelWarning LogLevel = "warning"

	LogLevelError LogLevel = "error"

	LogLevelCritical LogLevel = "critical"

	LogLevelAlert LogLevel = "alert"

	LogLevelEmergency LogLevel = "emergency"
)

// The binding key of each service that the log component publishes.
const (
	LogConfigContractServiceID = "valkyrja.log.data.LogConfigContract"

	LoggerContractServiceID = "valkyrja.log.logger.LoggerContract"
)
