/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package argument_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/argument"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
)

const (
	optionName = "verbose"
	firstValue = "first"
)

func TestTheArgumentHoldsItsValue(t *testing.T) {
	t.Parallel()

	built := argument.NewArgument(firstValue)

	if built.GetValue() != firstValue {
		t.Errorf("the argument must hold its value, but holds: %q", built.GetValue())
	}

	if built.WithValue("second").GetValue() != "second" {
		t.Error("WithValue must hold the new value, but did not")
	}

	if built.GetValue() != firstValue {
		t.Error("WithValue must leave the receiver unchanged, but did not")
	}
}

func TestTheOptionHoldsItsNameAndItsType(t *testing.T) {
	t.Parallel()

	built := argument.NewOption(optionName, constant.OptionTypeLong)

	if built.GetName() != optionName {
		t.Errorf("the option must hold its name, but holds: %q", built.GetName())
	}

	if built.GetType() != constant.OptionTypeLong {
		t.Errorf("the option must hold its type, but holds: %q", built.GetType())
	}

	if built.HasValue() || built.GetValue() != "" {
		t.Error("an option must carry no value until a caller gives it one, but carried one")
	}
}

func TestAnOptionWithAnEmptyValueStillHasOne(t *testing.T) {
	t.Parallel()

	built := argument.NewOption(optionName, constant.OptionTypeLong).WithValue("")

	if !built.HasValue() {
		t.Error("an option that the caller set to an empty value must have one, but did not")
	}

	if built.GetValue() != "" {
		t.Errorf("the value must be empty, but is: %q", built.GetValue())
	}
}

func TestEachOptionWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := argument.NewOption(optionName, constant.OptionTypeLong)

	if built.WithName("quiet").GetName() != "quiet" {
		t.Error("WithName must hold the new name, but did not")
	}

	if built.WithValue("yes").GetValue() != "yes" {
		t.Error("WithValue must hold the new value, but did not")
	}

	if built.WithType(constant.OptionTypeShort).GetType() != constant.OptionTypeShort {
		t.Error("WithType must hold the new type, but did not")
	}

	if built.GetName() != optionName || built.HasValue() {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestWithoutValueRemovesTheValue(t *testing.T) {
	t.Parallel()

	built := argument.NewOption(optionName, constant.OptionTypeLong).WithValue("yes")

	without := built.WithoutValue()

	if without.HasValue() || without.GetValue() != "" {
		t.Error("WithoutValue must remove the value, but did not")
	}

	if !built.HasValue() {
		t.Error("WithoutValue must leave the receiver unchanged, but did not")
	}
}

func TestEachTypeSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.ArgumentContract = argument.NewArgument(firstValue)
	var option contract.OptionContract = argument.NewOption(optionName, constant.OptionTypeLong)

	if built.GetValue() != firstValue || option.GetName() != optionName {
		t.Error("each contract must read what it holds, but did not")
	}
}
