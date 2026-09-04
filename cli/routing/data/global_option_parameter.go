/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
)

// NewHelpOptionParameter builds the option that prints the help text of a
// command.
func NewHelpOptionParameter() contract.OptionParameterContract {
	return newGlobalOptionParameter(
		constant.OptionNameHelp,
		"Help with this command",
		constant.OptionShortNameHelp,
	)
}

// NewVersionOptionParameter builds the option that prints the version of the
// application.
func NewVersionOptionParameter() contract.OptionParameterContract {
	return newGlobalOptionParameter(
		constant.OptionNameVersion,
		"Application version",
		constant.OptionShortNameVersion,
	)
}

// NewQuietOptionParameter builds the option that reports errors alone.
func NewQuietOptionParameter() contract.OptionParameterContract {
	return newGlobalOptionParameter(
		constant.OptionNameQuiet,
		"Output is suppressed, except for errors.",
		constant.OptionShortNameQuiet,
	)
}

// NewSilentOptionParameter builds the option that reports nothing.
func NewSilentOptionParameter() contract.OptionParameterContract {
	return newGlobalOptionParameter(
		constant.OptionNameSilent,
		"All output is suppressed",
		constant.OptionShortNameSilent,
	)
}

// NewNoInteractionOptionParameter builds the option that asks the caller no
// question.
func NewNoInteractionOptionParameter() contract.OptionParameterContract {
	return newGlobalOptionParameter(
		constant.OptionNameNoInteraction,
		"No interactive questions are asked.",
		constant.OptionShortNameNoInteraction,
	)
}

// newGlobalOptionParameter builds an option that every command accepts, and that
// takes no value.
func newGlobalOptionParameter(name string, description string, shortName string) contract.OptionParameterContract {
	return NewOptionParameter(name, description).
		WithShortNames(shortName).
		WithValueMode(constant.OptionValueModeNone)
}
