/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the CLI routing component's enumerations.
package constant

type ArgumentMode string

// The argument modes that the framework knows.
const (
	ArgumentModeRequired ArgumentMode = "REQUIRED"
	ArgumentModeOptional ArgumentMode = "OPTIONAL"
)

type ArgumentValueMode string

// The argument value modes that the framework knows.
const (
	ArgumentValueModeDefault ArgumentValueMode = "DEFAULT"
	ArgumentValueModeArray   ArgumentValueMode = "ARRAY"
)

type OptionMode string

// The option modes that the framework knows.
const (
	OptionModeRequired OptionMode = "REQUIRED"
	OptionModeOptional OptionMode = "OPTIONAL"
)

type OptionValueMode string

// The option value modes that the framework knows.
const (
	OptionValueModeNone    OptionValueMode = "NONE"
	OptionValueModeDefault OptionValueMode = "DEFAULT"
	OptionValueModeArray   OptionValueMode = "ARRAY"
)

// The name of each option that every command of the application accepts.
const (
	OptionNameHelp = "help"

	OptionNameVersion = "version"

	OptionNameQuiet = "quiet"

	OptionNameSilent = "silent"

	OptionNameNoInteraction = "no-interaction"

	OptionNameToken = "token"
)

// The short name of each option that every command of the application accepts.
const (
	OptionShortNameHelp = "h"

	OptionShortNameVersion = "v"

	OptionShortNameQuiet = "q"

	OptionShortNameSilent = "s"

	OptionShortNameNoInteraction = "N"

	OptionShortNameToken = "t"
)
