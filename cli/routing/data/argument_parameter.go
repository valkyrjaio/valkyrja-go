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
	"github.com/valkyrjaio/valkyrja-go/v26/cli/throwable/exception"
)

// ArgumentParameter is one positional parameter of a command.
type ArgumentParameter struct {
	parameter

	mode      constant.ArgumentMode
	valueMode constant.ArgumentValueMode
	arguments []contract.ArgumentContract
}

// NewArgumentParameter builds an optional argument that takes one value.
//
// A caller states another mode, another value mode, or a cast with the `With`
// methods, the way every other port does.
func NewArgumentParameter(name string, description string) *ArgumentParameter {
	return &ArgumentParameter{
		parameter: parameter{
			name:        name,
			description: description,
		},
		mode:      constant.ArgumentModeOptional,
		valueMode: constant.ArgumentValueModeDefault,
		arguments: []contract.ArgumentContract{},
	}
}

// WithName returns a copy of the argument under another name.
func (p *ArgumentParameter) WithName(name string) contract.ParameterContract {
	copied := *p
	copied.name = name

	return &copied
}

// WithCast returns a copy of the argument for another cast.
func (p *ArgumentParameter) WithCast(cast contract.CastFunc) contract.ParameterContract {
	copied := *p
	copied.cast = cast

	return &copied
}

// WithoutCast returns a copy of the argument with no cast.
func (p *ArgumentParameter) WithoutCast() contract.ParameterContract {
	copied := *p
	copied.cast = nil

	return &copied
}

// WithDescription returns a copy of the argument with another description.
func (p *ArgumentParameter) WithDescription(description string) contract.ParameterContract {
	copied := *p
	copied.description = description

	return &copied
}

// GetMode returns whether the command needs the argument.
func (p *ArgumentParameter) GetMode() constant.ArgumentMode {
	return p.mode
}

// WithMode returns a copy of the argument in another mode.
func (p *ArgumentParameter) WithMode(mode constant.ArgumentMode) contract.ArgumentParameterContract {
	copied := *p
	copied.mode = mode

	return &copied
}

// GetValueMode returns how many values the argument takes.
func (p *ArgumentParameter) GetValueMode() constant.ArgumentValueMode {
	return p.valueMode
}

// WithValueMode returns a copy of the argument in another value mode.
func (p *ArgumentParameter) WithValueMode(valueMode constant.ArgumentValueMode) contract.ArgumentParameterContract {
	copied := *p
	copied.valueMode = valueMode

	return &copied
}

// GetArguments returns each argument that the caller typed for this parameter.
func (p *ArgumentParameter) GetArguments() []contract.ArgumentContract {
	return p.arguments
}

// WithArguments returns a copy of the parameter with other arguments.
func (p *ArgumentParameter) WithArguments(arguments ...contract.ArgumentContract) contract.ArgumentParameterContract {
	copied := *p
	copied.arguments = arguments

	return &copied
}

// WithAddedArguments returns a copy of the parameter with the arguments
// appended.
func (p *ArgumentParameter) WithAddedArguments(
	arguments ...contract.ArgumentContract,
) contract.ArgumentParameterContract {
	combined := make([]contract.ArgumentContract, 0, len(p.arguments)+len(arguments))
	combined = append(combined, p.arguments...)
	combined = append(combined, arguments...)

	copied := *p
	copied.arguments = combined

	return &copied
}

// GetCastValues returns each value, cast to the type that the argument names.
func (p *ArgumentParameter) GetCastValues() ([]any, error) {
	return getCastValues(p.cast, p.getValues())
}

// HasFirstValue reports whether the caller gave the argument a value.
func (p *ArgumentParameter) HasFirstValue() bool {
	return len(p.arguments) > 0
}

// GetFirstValue returns the first value that the caller gave, and an empty
// string where the caller gave none.
func (p *ArgumentParameter) GetFirstValue() string {
	if !p.HasFirstValue() {
		return ""
	}

	return p.arguments[0].GetValue()
}

// AreValuesValid reports whether every value is one that the argument accepts.
//
// A required argument needs a value. An argument in the default value mode takes
// one value at most; an argument in the array value mode takes as many as the
// caller gives.
func (p *ArgumentParameter) AreValuesValid() bool {
	if p.mode == constant.ArgumentModeRequired && len(p.arguments) == 0 {
		return false
	}

	return p.valueMode != constant.ArgumentValueModeDefault || len(p.arguments) <= 1
}

// ValidateValues reports a failure where a value is one that the argument does
// not accept.
func (p *ArgumentParameter) ValidateValues() error {
	if !p.AreValuesValid() {
		return exception.NewCliRoutingArgumentValuesValidationError(p.name)
	}

	return nil
}

// getValues returns the value of each argument that the caller typed.
func (p *ArgumentParameter) getValues() []string {
	values := make([]string, 0, len(p.arguments))

	for _, argument := range p.arguments {
		values = append(values, argument.GetValue())
	}

	return values
}
