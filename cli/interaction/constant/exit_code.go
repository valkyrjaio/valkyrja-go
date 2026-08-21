/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the CLI interaction component's enumerations.
//
// Go has no enum, so each enumeration of the other ports is a defined type over
// a `const` block. The taxonomy puts every one of them in the `constant` segment
// for that reason.
package constant

type ExitCode int

// The exit codes that the framework uses. The values from 64 to 78 are the
// `sysexits.h` codes, which every port keeps so a command reports the same
// failure in each one.
const (
	ExitCodeSuccess ExitCode = 0
	ExitCodeError   ExitCode = 1

	ExitCodeUsageError    ExitCode = 64
	ExitCodeDataError     ExitCode = 65
	ExitCodeNoInput       ExitCode = 67
	ExitCodeNoUser        ExitCode = 68
	ExitCodeUnavailable   ExitCode = 69
	ExitCodeSoftwareError ExitCode = 70
	ExitCodeOsError       ExitCode = 71
	ExitCodeOsFileError   ExitCode = 72
	ExitCodeCantCreate    ExitCode = 73
	ExitCodeIoError       ExitCode = 74
	ExitCodeTempFail      ExitCode = 75
	ExitCodeProtocolError ExitCode = 76
	ExitCodeNoPermission  ExitCode = 77
	ExitCodeConfigError   ExitCode = 78

	ExitCodeAutoExit ExitCode = 255
)

type OptionType string

// The option types that the framework knows.
const (
	// OptionTypeShort is an option that the caller spells with one dash.
	OptionTypeShort OptionType = "SHORT"

	// OptionTypeLong is an option that the caller spells with two dashes.
	OptionTypeLong OptionType = "LONG"
)
