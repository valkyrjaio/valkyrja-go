/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the CLI routing component's enumerations.
package constant

// ArgumentMode says whether a command needs the argument.
type ArgumentMode string

// The argument modes that the framework knows.
const (
	ArgumentModeRequired ArgumentMode = "REQUIRED"
	ArgumentModeOptional ArgumentMode = "OPTIONAL"
)

// ArgumentValueMode says how many values the argument takes.
type ArgumentValueMode string

// The argument value modes that the framework knows.
const (
	ArgumentValueModeDefault ArgumentValueMode = "DEFAULT"
	ArgumentValueModeArray   ArgumentValueMode = "ARRAY"
)

// OptionMode says whether a command needs the option.
type OptionMode string

// The option modes that the framework knows.
const (
	OptionModeRequired OptionMode = "REQUIRED"
	OptionModeOptional OptionMode = "OPTIONAL"
)

// OptionValueMode says how many values the option takes.
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
