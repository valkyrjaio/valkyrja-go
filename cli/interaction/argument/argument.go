/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package argument holds what the caller typed on the command line.
package argument

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
)

type Argument struct {
	value string
}

// NewArgument builds an argument that carries the value.
func NewArgument(value string) *Argument {
	return &Argument{value: value}
}

// GetValue returns the value of the argument.
func (a *Argument) GetValue() string {
	return a.value
}

// WithValue returns a copy of the argument with another value.
func (a *Argument) WithValue(value string) contract.ArgumentContract {
	copied := *a
	copied.value = value

	return &copied
}

type Option struct {
	name       string
	value      string
	hasValue   bool
	optionType constant.OptionType
}

// NewOption builds an option under a name, of a type. The option carries no
// value until a caller gives it one.
func NewOption(name string, optionType constant.OptionType) *Option {
	return &Option{
		name:       name,
		optionType: optionType,
	}
}

// GetName returns the name of the option.
func (o *Option) GetName() string {
	return o.name
}

// WithName returns a copy of the option under another name.
func (o *Option) WithName(name string) contract.OptionContract {
	copied := *o
	copied.name = name

	return &copied
}

// HasValue reports whether the caller gave the option a value.
func (o *Option) HasValue() bool {
	return o.hasValue
}

// GetValue returns the value of the option.
func (o *Option) GetValue() string {
	return o.value
}

// WithValue returns a copy of the option with another value.
func (o *Option) WithValue(value string) contract.OptionContract {
	copied := *o
	copied.value = value
	copied.hasValue = true

	return &copied
}

// WithoutValue returns a copy of the option with no value.
func (o *Option) WithoutValue() contract.OptionContract {
	copied := *o
	copied.value = ""
	copied.hasValue = false

	return &copied
}

// GetType returns the type of the option, which says how the caller spells it.
func (o *Option) GetType() constant.OptionType {
	return o.optionType
}

// WithType returns a copy of the option of another type.
func (o *Option) WithType(optionType constant.OptionType) contract.OptionContract {
	copied := *o
	copied.optionType = optionType

	return &copied
}
