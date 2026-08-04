/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/throwable/exception"
)

// OptionParameter is one option parameter of a command.
type OptionParameter struct {
	parameter

	valueDisplayName string
	defaultValue     string
	shortNames       []string
	validValues      []string
	options          []contract.OptionContract
	mode             constant.OptionMode
	valueMode        constant.OptionValueMode
}

// NewOptionParameter builds an optional option that takes one value.
//
// A caller states another mode, another value mode, a short name, a valid value,
// or a cast with the `With` methods, the way every other port does.
func NewOptionParameter(name string, description string) *OptionParameter {
	return &OptionParameter{
		parameter: parameter{
			name:        name,
			description: description,
		},
		shortNames:  []string{},
		validValues: []string{},
		options:     []contract.OptionContract{},
		mode:        constant.OptionModeOptional,
		valueMode:   constant.OptionValueModeDefault,
	}
}

// WithName returns a copy of the option under another name.
func (p *OptionParameter) WithName(name string) contract.ParameterContract {
	copied := *p
	copied.name = name

	return &copied
}

// WithCast returns a copy of the option for another cast.
func (p *OptionParameter) WithCast(cast contract.CastFunc) contract.ParameterContract {
	copied := *p
	copied.cast = cast

	return &copied
}

// WithoutCast returns a copy of the option with no cast.
func (p *OptionParameter) WithoutCast() contract.ParameterContract {
	copied := *p
	copied.cast = nil

	return &copied
}

// WithDescription returns a copy of the option with another description.
func (p *OptionParameter) WithDescription(description string) contract.ParameterContract {
	copied := *p
	copied.description = description

	return &copied
}

// GetShortNames returns each one-letter name of the option.
func (p *OptionParameter) GetShortNames() []string {
	return p.shortNames
}

// WithShortNames returns a copy of the option with other short names.
func (p *OptionParameter) WithShortNames(shortNames ...string) contract.OptionParameterContract {
	copied := *p
	copied.shortNames = shortNames

	return &copied
}

// WithAddedShortNames returns a copy of the option with the short names
// appended.
//
// A short name that the option holds already is not appended twice, because a
// help text that names it twice reads as two options.
func (p *OptionParameter) WithAddedShortNames(shortNames ...string) contract.OptionParameterContract {
	copied := *p
	copied.shortNames = appendUnique(p.shortNames, shortNames)

	return &copied
}

// GetMode returns whether the command needs the option.
func (p *OptionParameter) GetMode() constant.OptionMode {
	return p.mode
}

// WithMode returns a copy of the option in another mode.
func (p *OptionParameter) WithMode(mode constant.OptionMode) contract.OptionParameterContract {
	copied := *p
	copied.mode = mode

	return &copied
}

// GetValueMode returns how many values the option takes.
func (p *OptionParameter) GetValueMode() constant.OptionValueMode {
	return p.valueMode
}

// WithValueMode returns a copy of the option in another value mode.
func (p *OptionParameter) WithValueMode(valueMode constant.OptionValueMode) contract.OptionParameterContract {
	copied := *p
	copied.valueMode = valueMode

	return &copied
}

// HasValueDisplayName reports whether the help text names the value.
func (p *OptionParameter) HasValueDisplayName() bool {
	return p.valueDisplayName != ""
}

// GetValueDisplayName returns the name that the help text gives the value.
func (p *OptionParameter) GetValueDisplayName() string {
	return p.valueDisplayName
}

// WithValueDisplayName returns a copy of the option with another display name.
func (p *OptionParameter) WithValueDisplayName(valueName string) contract.OptionParameterContract {
	copied := *p
	copied.valueDisplayName = valueName

	return &copied
}

// GetValidValues returns each value that the option accepts.
func (p *OptionParameter) GetValidValues() []string {
	return p.validValues
}

// WithValidValues returns a copy of the option that accepts other values.
func (p *OptionParameter) WithValidValues(validValues ...string) contract.OptionParameterContract {
	copied := *p
	copied.validValues = validValues

	return &copied
}

// WithAddedValidValues returns a copy of the option with the valid values
// appended.
func (p *OptionParameter) WithAddedValidValues(validValues ...string) contract.OptionParameterContract {
	copied := *p
	copied.validValues = appendUnique(p.validValues, validValues)

	return &copied
}

// HasDefaultValue reports whether the option has a value that the caller leaves
// out.
func (p *OptionParameter) HasDefaultValue() bool {
	return p.defaultValue != ""
}

// GetDefaultValue returns the value that the option uses where the caller gives
// none.
func (p *OptionParameter) GetDefaultValue() string {
	return p.defaultValue
}

// WithDefaultValue returns a copy of the option with another default value.
func (p *OptionParameter) WithDefaultValue(defaultValue string) contract.OptionParameterContract {
	copied := *p
	copied.defaultValue = defaultValue

	return &copied
}

// GetOptions returns each option that the caller typed for this parameter.
func (p *OptionParameter) GetOptions() []contract.OptionContract {
	return p.options
}

// WithOptions returns a copy of the parameter with other options.
//
// Warning: an option in the none value mode takes no value. An option that
// carries one reports a failure here, rather than at the moment the command
// reads it.
func (p *OptionParameter) WithOptions(
	options ...contract.OptionContract,
) (contract.OptionParameterContract, error) {
	err := p.validateOptions(options)
	if err != nil {
		return nil, err
	}

	copied := *p
	copied.options = options

	return &copied, nil
}

// WithAddedOptions returns a copy of the parameter with the options appended.
func (p *OptionParameter) WithAddedOptions(
	options ...contract.OptionContract,
) (contract.OptionParameterContract, error) {
	err := p.validateOptions(options)
	if err != nil {
		return nil, err
	}

	combined := make([]contract.OptionContract, 0, len(p.options)+len(options))
	combined = append(combined, p.options...)
	combined = append(combined, options...)

	copied := *p
	copied.options = combined

	return &copied, nil
}

// GetCastValues returns each value, cast to the type that the option names.
func (p *OptionParameter) GetCastValues() ([]any, error) {
	return getCastValues(p.cast, p.getValues())
}

// HasFirstValue reports whether the caller gave the option a value.
func (p *OptionParameter) HasFirstValue() bool {
	return len(p.options) > 0
}

// GetFirstValue returns the first value that the caller gave, and an empty
// string where the caller gave none.
func (p *OptionParameter) GetFirstValue() string {
	if !p.HasFirstValue() {
		return ""
	}

	return p.options[0].GetValue()
}

// AreValuesValid reports whether every value is one that the option accepts.
//
// A required option needs a value. An option in the default value mode takes one
// value at most. An option that names its valid values takes those values alone.
func (p *OptionParameter) AreValuesValid() bool {
	if len(p.validValues) > 0 && !p.areValuesInValidValues() {
		return false
	}

	if p.mode == constant.OptionModeRequired && len(p.options) == 0 {
		return false
	}

	return p.valueMode != constant.OptionValueModeDefault || len(p.options) <= 1
}

// ValidateValues reports a failure where a value is one that the option does not
// accept.
func (p *OptionParameter) ValidateValues() error {
	if !p.AreValuesValid() {
		return exception.NewCliRoutingOptionValuesValidationError(p.name)
	}

	return nil
}

// validateOptions reports a failure where an option carries a value that the
// parameter takes no value for.
func (p *OptionParameter) validateOptions(options []contract.OptionContract) error {
	if p.valueMode != constant.OptionValueModeNone {
		return nil
	}

	for _, option := range options {
		if option.HasValue() {
			return exception.NewCliRoutingInvalidOptionWithValueError(p.name)
		}
	}

	return nil
}

// areValuesInValidValues reports whether every value is one that the option
// names as valid.
func (p *OptionParameter) areValuesInValidValues() bool {
	for _, option := range p.options {
		if !slices.Contains(p.validValues, option.GetValue()) {
			return false
		}
	}

	return true
}

// getValues returns the value of each option that the caller typed.
func (p *OptionParameter) getValues() []string {
	values := make([]string, 0, len(p.options))

	for _, option := range p.options {
		values = append(values, option.GetValue())
	}

	return values
}
