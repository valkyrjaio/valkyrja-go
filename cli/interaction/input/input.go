/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package input holds what the caller typed.
package input

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/argument"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
)

// longPrefix opens an option that the caller spells with two dashes.
const longPrefix = "--"

// shortPrefix opens an option that the caller spells with one dash.
const shortPrefix = "-"

// valueSeparator separates the name of an option from its value.
const valueSeparator = "="

type Input struct {
	caller      string
	commandName string
	arguments   []contract.ArgumentContract
	options     []contract.OptionContract
}

// NewInput builds an input from the caller and the command name.
func NewInput(caller string, commandName string) *Input {
	return &Input{
		caller:      caller,
		commandName: commandName,
		arguments:   []contract.ArgumentContract{},
		options:     []contract.OptionContract{},
	}
}

// NewInputFromArgs builds an input by reading what the caller typed.
func NewInputFromArgs(args []string) *Input {
	built := NewInput("", "")

	if len(args) > 0 {
		built.caller = args[0]
	}

	if len(args) > 1 {
		built.commandName = args[1]
	}

	if len(args) < 2 {
		return built
	}

	for _, arg := range args[2:] {
		built.readArg(arg)
	}

	return built
}

// newOption builds the option that the text spells, with its value where the
// text carries one.
func newOption(text string, optionType constant.OptionType) contract.OptionContract {
	name, value, found := strings.Cut(text, valueSeparator)

	built := argument.NewOption(name, optionType)

	if !found {
		return built
	}

	return built.WithValue(value)
}

// GetCaller returns the path that the caller ran.
func (i *Input) GetCaller() string {
	return i.caller
}

// WithCaller returns a copy of the input for another caller.
func (i *Input) WithCaller(caller string) contract.InputContract {
	copied := *i
	copied.caller = caller

	return &copied
}

// GetCommandName returns the name of the command that the caller ran.
func (i *Input) GetCommandName() string {
	return i.commandName
}

// WithCommandName returns a copy of the input for another command.
func (i *Input) WithCommandName(commandName string) contract.InputContract {
	copied := *i
	copied.commandName = commandName

	return &copied
}

// GetArguments returns each positional argument, in the order that the caller
// typed them.
func (i *Input) GetArguments() []contract.ArgumentContract {
	return i.arguments
}

// WithArguments returns a copy of the input with other arguments.
func (i *Input) WithArguments(arguments ...contract.ArgumentContract) contract.InputContract {
	copied := *i
	copied.arguments = arguments

	return &copied
}

// WithAddedArgument returns a copy of the input with the argument appended.
func (i *Input) WithAddedArgument(added contract.ArgumentContract) contract.InputContract {
	combined := make([]contract.ArgumentContract, 0, len(i.arguments)+1)
	combined = append(combined, i.arguments...)
	combined = append(combined, added)

	copied := *i
	copied.arguments = combined

	return &copied
}

// WithoutArgument returns a copy of the input without the argument that carries
// the value.
func (i *Input) WithoutArgument(value string) contract.InputContract {
	kept := make([]contract.ArgumentContract, 0, len(i.arguments))

	for _, held := range i.arguments {
		if held.GetValue() == value {
			continue
		}

		kept = append(kept, held)
	}

	copied := *i
	copied.arguments = kept

	return &copied
}

// WithoutArguments returns a copy of the input with no argument.
func (i *Input) WithoutArguments() contract.InputContract {
	copied := *i
	copied.arguments = []contract.ArgumentContract{}

	return &copied
}

// GetOptions returns each option that the caller typed.
func (i *Input) GetOptions() []contract.OptionContract {
	return i.options
}

// GetOption returns each option under the name.
func (i *Input) GetOption(name string) []contract.OptionContract {
	found := make([]contract.OptionContract, 0, len(i.options))

	for _, held := range i.options {
		if held.GetName() != name {
			continue
		}

		found = append(found, held)
	}

	return found
}

// HasOption reports whether the caller typed the option.
func (i *Input) HasOption(name string) bool {
	return len(i.GetOption(name)) > 0
}

// WithOptions returns a copy of the input with other options.
func (i *Input) WithOptions(options ...contract.OptionContract) contract.InputContract {
	copied := *i
	copied.options = options

	return &copied
}

// WithAddedOption returns a copy of the input with the option appended.
func (i *Input) WithAddedOption(added contract.OptionContract) contract.InputContract {
	combined := make([]contract.OptionContract, 0, len(i.options)+1)
	combined = append(combined, i.options...)
	combined = append(combined, added)

	copied := *i
	copied.options = combined

	return &copied
}

// WithoutOption returns a copy of the input without the option.
func (i *Input) WithoutOption(name string) contract.InputContract {
	kept := make([]contract.OptionContract, 0, len(i.options))

	for _, held := range i.options {
		if held.GetName() == name {
			continue
		}

		kept = append(kept, held)
	}

	copied := *i
	copied.options = kept

	return &copied
}

// WithoutOptions returns a copy of the input with no option.
func (i *Input) WithoutOptions() contract.InputContract {
	copied := *i
	copied.options = []contract.OptionContract{}

	return &copied
}

// readArg reads one argument into the input, as an option or as a positional
// argument.
func (i *Input) readArg(arg string) {
	switch {
	case strings.HasPrefix(arg, longPrefix):
		i.options = append(i.options, newOption(strings.TrimPrefix(arg, longPrefix), constant.OptionTypeLong))
	case strings.HasPrefix(arg, shortPrefix):
		i.options = append(i.options, newOption(strings.TrimPrefix(arg, shortPrefix), constant.OptionTypeShort))
	default:
		i.arguments = append(i.arguments, argument.NewArgument(arg))
	}
}
