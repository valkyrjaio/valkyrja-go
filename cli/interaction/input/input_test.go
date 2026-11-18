/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package input_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/argument"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/input"
)

const (
	callerPath  = "/usr/local/bin/valkyrja"
	commandName = "cache:clear"
	verboseName = "verbose"
	secondValue = "second"
)

func TestNewInputFromArgsReadsTheCallerAndTheCommand(t *testing.T) {
	t.Parallel()

	built := input.NewInputFromArgs([]string{callerPath, commandName})

	if built.GetCaller() != callerPath {
		t.Errorf("the input must read the caller, but read: %q", built.GetCaller())
	}

	if built.GetCommandName() != commandName {
		t.Errorf("the input must read the command name, but read: %q", built.GetCommandName())
	}
}

func TestNewInputFromArgsReadsEachArgumentAndOption(t *testing.T) {
	t.Parallel()

	built := input.NewInputFromArgs([]string{
		callerPath, commandName, "first", "--verbose", "-q", "--name=value", secondValue,
	})

	if len(built.GetArguments()) != 2 {
		t.Fatalf("the input must read each positional argument, but read: %d", len(built.GetArguments()))
	}

	if built.GetArguments()[0].GetValue() != "first" || built.GetArguments()[1].GetValue() != secondValue {
		t.Error("the input must keep the order of the arguments, but did not")
	}

	if len(built.GetOptions()) != 3 {
		t.Fatalf("the input must read each option, but read: %d", len(built.GetOptions()))
	}

	if !built.HasOption(verboseName) || !built.HasOption("q") || !built.HasOption("name") {
		t.Error("the input must read each option by name, but did not")
	}
}

func TestNewInputFromArgsReadsTheTypeOfEachOption(t *testing.T) {
	t.Parallel()

	built := input.NewInputFromArgs([]string{callerPath, commandName, "--verbose", "-q"})

	if built.GetOption(verboseName)[0].GetType() != constant.OptionTypeLong {
		t.Error("an option with two dashes must be long, but was not")
	}

	if built.GetOption("q")[0].GetType() != constant.OptionTypeShort {
		t.Error("an option with one dash must be short, but was not")
	}
}

func TestNewInputFromArgsReadsTheValueOfAnOption(t *testing.T) {
	t.Parallel()

	built := input.NewInputFromArgs([]string{callerPath, commandName, "--name=value", "--empty=", "--none"})

	if built.GetOption("name")[0].GetValue() != "value" {
		t.Error("an option must carry the value that the caller gave, but did not")
	}

	if !built.GetOption("empty")[0].HasValue() {
		t.Error("an option that the caller set to an empty value must have one, but did not")
	}

	if built.GetOption("none")[0].HasValue() {
		t.Error("an option with no value must have none, but had one")
	}
}

func TestNewInputFromArgsReadsAnArgumentListWithNoCommand(t *testing.T) {
	t.Parallel()

	tests := map[string][]string{
		"no argument at all": {},
		"the caller alone":   {callerPath},
	}

	for name, args := range tests {
		built := input.NewInputFromArgs(args)

		if built.GetCommandName() != "" {
			t.Errorf("%s must name no command, but named: %q", name, built.GetCommandName())
		}

		if len(built.GetArguments()) != 0 || len(built.GetOptions()) != 0 {
			t.Errorf("%s must read no argument and no option, but read some", name)
		}
	}
}

func TestGetOptionReturnsEachOptionUnderTheName(t *testing.T) {
	t.Parallel()

	built := input.NewInputFromArgs([]string{callerPath, commandName, "--tag=a", "--tag=b"})

	if len(built.GetOption("tag")) != 2 {
		t.Errorf("GetOption must return each option under the name, but returned: %d",
			len(built.GetOption("tag")))
	}

	if len(built.GetOption("missing")) != 0 || built.HasOption("missing") {
		t.Error("GetOption must be empty for an unknown name, but was not")
	}
}

func TestEachCallerAndCommandWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := input.NewInput(callerPath, commandName)

	if built.WithCaller("/other").GetCaller() != "/other" {
		t.Error("WithCaller must hold the new caller, but did not")
	}

	if built.WithCommandName("other").GetCommandName() != "other" {
		t.Error("WithCommandName must hold the new command, but did not")
	}

	if built.GetCaller() != callerPath || built.GetCommandName() != commandName {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachArgumentMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := input.NewInput(callerPath, commandName).
		WithArguments(argument.NewArgument("first"), argument.NewArgument(secondValue))

	added := built.WithAddedArgument(argument.NewArgument("third"))
	without := built.WithoutArgument("first")
	none := built.WithoutArguments()

	if len(added.GetArguments()) != 3 {
		t.Error("WithAddedArgument must append the argument, but did not")
	}

	if len(without.GetArguments()) != 1 || without.GetArguments()[0].GetValue() != secondValue {
		t.Error("WithoutArgument must remove the argument, but did not")
	}

	if len(none.GetArguments()) != 0 {
		t.Error("WithoutArguments must remove each argument, but did not")
	}

	if len(built.GetArguments()) != 2 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachOptionMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := input.NewInput(callerPath, commandName).WithOptions(
		argument.NewOption(verboseName, constant.OptionTypeLong),
		argument.NewOption("quiet", constant.OptionTypeLong),
	)

	added := built.WithAddedOption(argument.NewOption("silent", constant.OptionTypeLong))
	without := built.WithoutOption(verboseName)
	none := built.WithoutOptions()

	if len(added.GetOptions()) != 3 {
		t.Error("WithAddedOption must append the option, but did not")
	}

	if len(without.GetOptions()) != 1 || without.GetOptions()[0].GetName() != "quiet" {
		t.Error("WithoutOption must remove the option, but did not")
	}

	if len(none.GetOptions()) != 0 {
		t.Error("WithoutOptions must remove each option, but did not")
	}

	if len(built.GetOptions()) != 2 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestTheInputSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.InputContract = input.NewInput(callerPath, commandName)

	if built.GetCommandName() != commandName {
		t.Error("the contract must read the command name, but did not")
	}
}
