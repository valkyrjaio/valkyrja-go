/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/argument"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/throwable/exception"
)

const (
	parameterName        = "count"
	parameterDescription = "How many to build"
	otherName            = "other"
)

// errCast reports that a value does not convert.
var errCast = errors.New("the value does not convert")

// castToUpper is a cast that a test gives a parameter.
func castToUpper(value string) (any, error) {
	if value == "" {
		return nil, errCast
	}

	return "cast:" + value, nil
}

func TestAnArgumentParameterReadsWhatItWasBuiltWith(t *testing.T) {
	t.Parallel()

	built := data.NewArgumentParameter(parameterName, parameterDescription)

	if built.GetName() != parameterName || built.GetDescription() != parameterDescription {
		t.Error("the parameter must read its name and its description, but did not")
	}

	if built.HasCast() || built.GetCast() != nil {
		t.Error("a parameter that was built with no cast must carry none, but carried one")
	}

	if built.GetMode() != constant.ArgumentModeOptional {
		t.Error("an argument must be optional until a caller states otherwise, but was not")
	}

	if built.GetValueMode() != constant.ArgumentValueModeDefault {
		t.Error("an argument must take one value until a caller states otherwise, but did not")
	}
}

func TestEachArgumentParameterWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := data.NewArgumentParameter(parameterName, parameterDescription)

	if built.WithName(otherName).GetName() != otherName {
		t.Error("WithName must hold the new name, but did not")
	}

	if built.WithDescription(otherName).GetDescription() != otherName {
		t.Error("WithDescription must hold the new description, but did not")
	}

	if !built.WithCast(castToUpper).HasCast() {
		t.Error("WithCast must hold the cast, but did not")
	}

	if built.WithCast(castToUpper).WithoutCast().HasCast() {
		t.Error("WithoutCast must remove the cast, but did not")
	}

	if built.WithMode(constant.ArgumentModeRequired).GetMode() != constant.ArgumentModeRequired {
		t.Error("WithMode must hold the new mode, but did not")
	}

	valueMode := constant.ArgumentValueModeArray
	if built.WithValueMode(valueMode).GetValueMode() != valueMode {
		t.Error("WithValueMode must hold the new value mode, but did not")
	}

	if built.HasCast() || built.GetName() != parameterName {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestAnArgumentParameterHoldsItsArguments(t *testing.T) {
	t.Parallel()

	built := data.NewArgumentParameter(parameterName, parameterDescription)

	if built.HasFirstValue() || built.GetFirstValue() != "" {
		t.Error("a parameter with no argument must report no first value, but reported one")
	}

	filled := built.WithArguments(argument.NewArgument("first"))
	added := filled.WithAddedArguments(argument.NewArgument("second"))

	if !filled.HasFirstValue() || filled.GetFirstValue() != "first" {
		t.Error("the parameter must report the first value that the caller gave, but did not")
	}

	if len(added.GetArguments()) != 2 {
		t.Errorf("WithAddedArguments must append the argument, but held: %d", len(added.GetArguments()))
	}

	if len(filled.GetArguments()) != 1 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestAnArgumentParameterCastsItsValues(t *testing.T) {
	t.Parallel()

	built := data.NewArgumentParameter(parameterName, parameterDescription).
		WithArguments(argument.NewArgument("first"))

	plain, err := built.GetCastValues()
	if err != nil || len(plain) != 1 || plain[0] != "first" {
		t.Errorf("a parameter with no cast must return the value as it is, but returned: %v (%v)", plain, err)
	}

	cast, err := asArgumentParameter(t, built.WithCast(castToUpper)).GetCastValues()
	if err != nil || len(cast) != 1 || cast[0] != "cast:first" {
		t.Errorf("a parameter with a cast must convert the value, but returned: %v (%v)", cast, err)
	}

	failing := asArgumentParameter(t, data.NewArgumentParameter(parameterName, parameterDescription).
		WithArguments(argument.NewArgument("")).
		WithCast(castToUpper))

	_, err = failing.GetCastValues()
	if !errors.Is(err, errCast) {
		t.Errorf("a cast that reports a failure must reach the caller, but returned: %v", err)
	}
}

func TestAnArgumentParameterValidatesItsValues(t *testing.T) {
	t.Parallel()

	two := []contract.ArgumentContract{argument.NewArgument("first"), argument.NewArgument("second")}

	tests := map[string]struct {
		parameter contract.ArgumentParameterContract
		valid     bool
	}{
		"an optional argument with no value": {
			parameter: data.NewArgumentParameter(parameterName, parameterDescription),
			valid:     true,
		},
		"a required argument with no value": {
			parameter: data.NewArgumentParameter(parameterName, parameterDescription).
				WithMode(constant.ArgumentModeRequired),
			valid: false,
		},
		"a required argument with a value": {
			parameter: data.NewArgumentParameter(parameterName, parameterDescription).
				WithMode(constant.ArgumentModeRequired).
				WithArguments(argument.NewArgument("first")),
			valid: true,
		},
		"a default argument with two values": {
			parameter: data.NewArgumentParameter(parameterName, parameterDescription).WithArguments(two...),
			valid:     false,
		},
		"an array argument with two values": {
			parameter: data.NewArgumentParameter(parameterName, parameterDescription).
				WithValueMode(constant.ArgumentValueModeArray).
				WithArguments(two...),
			valid: true,
		},
	}

	for name, test := range tests {
		if test.parameter.AreValuesValid() != test.valid {
			t.Errorf("%s must report valid=%t, but did not", name, test.valid)
		}

		err := test.parameter.ValidateValues()
		if (err == nil) != test.valid {
			t.Errorf("%s must report an error only where it is invalid, but returned: %v", name, err)
		}

		var validationError *exception.CliRoutingArgumentValuesValidationError
		if !test.valid && !errors.As(err, &validationError) {
			t.Errorf("%s must report a validation error, but reported: %v", name, err)
		}
	}
}

func TestAnOptionParameterReadsWhatItWasBuiltWith(t *testing.T) {
	t.Parallel()

	built := data.NewOptionParameter(parameterName, parameterDescription)

	if built.GetName() != parameterName || built.GetDescription() != parameterDescription {
		t.Error("the parameter must read its name and its description, but did not")
	}

	if built.GetMode() != constant.OptionModeOptional || built.GetValueMode() != constant.OptionValueModeDefault {
		t.Error("an option must be optional and take one value until a caller states otherwise, but did not")
	}

	if built.HasValueDisplayName() || built.HasDefaultValue() {
		t.Error("an option must name no value and hold no default until a caller states one, but did")
	}
}

func TestEachOptionParameterWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := data.NewOptionParameter(parameterName, parameterDescription)

	if built.WithName(otherName).GetName() != otherName {
		t.Error("WithName must hold the new name, but did not")
	}

	if built.WithDescription(otherName).GetDescription() != otherName {
		t.Error("WithDescription must hold the new description, but did not")
	}

	if !built.WithCast(castToUpper).HasCast() || built.WithCast(castToUpper).WithoutCast().HasCast() {
		t.Error("WithCast and WithoutCast must hold and remove the cast, but did not")
	}

	if built.GetName() != parameterName || built.HasCast() {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachOptionParameterModeWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := data.NewOptionParameter(parameterName, parameterDescription)

	if built.WithMode(constant.OptionModeRequired).GetMode() != constant.OptionModeRequired {
		t.Error("WithMode must hold the new mode, but did not")
	}

	if built.WithValueMode(constant.OptionValueModeNone).GetValueMode() != constant.OptionValueModeNone {
		t.Error("WithValueMode must hold the new value mode, but did not")
	}

	if built.WithValueDisplayName("name").GetValueDisplayName() != "name" {
		t.Error("WithValueDisplayName must hold the new display name, but did not")
	}

	if built.WithDefaultValue("value").GetDefaultValue() != "value" {
		t.Error("WithDefaultValue must hold the new default value, but did not")
	}

	if built.HasValueDisplayName() || built.HasDefaultValue() {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

// asArgumentParameter reads the parameter back as an argument parameter, and
// fails the test where it is another type.
func asArgumentParameter(
	t *testing.T,
	parameter contract.ParameterContract,
) contract.ArgumentParameterContract {
	t.Helper()

	typed, isArgument := parameter.(contract.ArgumentParameterContract)
	if !isArgument {
		t.Fatal("the parameter must stay an argument parameter, but did not")
	}

	return typed
}

// asOptionParameter reads the parameter back as an option parameter, and fails
// the test where it is another type.
func asOptionParameter(
	t *testing.T,
	parameter contract.ParameterContract,
) contract.OptionParameterContract {
	t.Helper()

	typed, isOption := parameter.(contract.OptionParameterContract)
	if !isOption {
		t.Fatal("the parameter must stay an option parameter, but did not")
	}

	return typed
}

func TestAnOptionParameterAppendsAShortNameOnlyOnce(t *testing.T) {
	t.Parallel()

	built := data.NewOptionParameter(parameterName, parameterDescription).WithShortNames("c")

	added := built.WithAddedShortNames("c", "n")

	if len(added.GetShortNames()) != 2 {
		t.Errorf("a short name that the option holds must not be appended twice, but held: %v",
			added.GetShortNames())
	}

	if len(built.GetShortNames()) != 1 {
		t.Error("WithAddedShortNames must leave the receiver unchanged, but did not")
	}
}

func TestAnOptionParameterAppendsAValidValueOnlyOnce(t *testing.T) {
	t.Parallel()

	built := data.NewOptionParameter(parameterName, parameterDescription).WithValidValues("one")

	added := built.WithAddedValidValues("one", "two")

	if len(added.GetValidValues()) != 2 {
		t.Errorf("a valid value that the option holds must not be appended twice, but held: %v",
			added.GetValidValues())
	}

	if len(built.GetValidValues()) != 1 {
		t.Error("WithAddedValidValues must leave the receiver unchanged, but did not")
	}
}
