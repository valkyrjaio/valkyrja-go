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
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/throwable/exception"
)

// newOption builds one option that a caller typed, under the name of the
// parameter that each test in this file fills.
func newOption(value string) contract.OptionContract {
	built := argument.NewOption(parameterName, interactionconstant.OptionTypeLong)

	if value == "" {
		return built
	}

	return built.WithValue(value)
}

func TestAnOptionParameterHoldsItsOptions(t *testing.T) {
	t.Parallel()

	built := data.NewOptionParameter(parameterName, parameterDescription)

	if built.HasFirstValue() || built.GetFirstValue() != "" {
		t.Error("a parameter with no option must report no first value, but reported one")
	}

	filled, err := built.WithOptions(newOption("one"))
	if err != nil {
		t.Fatalf("WithOptions must take an option that carries a value, but reported: %v", err)
	}

	added, err := filled.WithAddedOptions(newOption("two"))
	if err != nil {
		t.Fatalf("WithAddedOptions must take an option that carries a value, but reported: %v", err)
	}

	if !filled.HasFirstValue() || filled.GetFirstValue() != "one" {
		t.Error("the parameter must report the first value that the caller gave, but did not")
	}

	if len(added.GetOptions()) != 2 || len(filled.GetOptions()) != 1 {
		t.Error("WithAddedOptions must append the option and leave the receiver unchanged, but did not")
	}
}

func TestAnOptionParameterThatTakesNoValueRejectsOne(t *testing.T) {
	t.Parallel()

	built := data.NewOptionParameter(parameterName, parameterDescription).
		WithValueMode(constant.OptionValueModeNone)

	_, withErr := built.WithOptions(newOption("one"))
	_, addedErr := built.WithAddedOptions(newOption("one"))

	var invalid *exception.CliRoutingInvalidOptionWithValueError

	if !errors.As(withErr, &invalid) || !errors.As(addedErr, &invalid) {
		t.Errorf("an option that carries a value must be rejected, but reported: %v and %v", withErr, addedErr)
	}

	_, noValueErr := built.WithOptions(newOption(""))
	if noValueErr != nil {
		t.Errorf("an option that carries no value must be taken, but reported: %v", noValueErr)
	}
}

func TestAnOptionParameterCastsItsValues(t *testing.T) {
	t.Parallel()

	filled, err := data.NewOptionParameter(parameterName, parameterDescription).
		WithOptions(newOption("one"))
	if err != nil {
		t.Fatalf("the parameter must take the option, but reported: %v", err)
	}

	plain, err := filled.GetCastValues()
	if err != nil || len(plain) != 1 || plain[0] != "one" {
		t.Errorf("a parameter with no cast must return the value as it is, but returned: %v (%v)", plain, err)
	}

	cast, err := asOptionParameter(t, filled.WithCast(castToUpper)).GetCastValues()
	if err != nil || len(cast) != 1 || cast[0] != "cast:one" {
		t.Errorf("a parameter with a cast must convert the value, but returned: %v (%v)", cast, err)
	}
}

func TestAnOptionParameterValidatesItsValues(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		parameter contract.OptionParameterContract
		valid     bool
	}{
		"an optional option with no value": {
			parameter: data.NewOptionParameter(parameterName, parameterDescription),
			valid:     true,
		},
		"a required option with no value": {
			parameter: data.NewOptionParameter(parameterName, parameterDescription).
				WithMode(constant.OptionModeRequired),
			valid: false,
		},
		"a required option with a value": {
			parameter: withOptions(t, data.NewOptionParameter(parameterName, parameterDescription).
				WithMode(constant.OptionModeRequired), newOption("one")),
			valid: true,
		},
		"a default option with two values": {
			parameter: withOptions(t, data.NewOptionParameter(parameterName, parameterDescription),
				newOption("one"), newOption("two")),
			valid: false,
		},
		"an array option with two values": {
			parameter: withOptions(t, data.NewOptionParameter(parameterName, parameterDescription).
				WithValueMode(constant.OptionValueModeArray),
				newOption("one"), newOption("two")),
			valid: true,
		},
		"a value that the option accepts": {
			parameter: withOptions(t, data.NewOptionParameter(parameterName, parameterDescription).
				WithValidValues("one"), newOption("one")),
			valid: true,
		},
		"a value that the option does not accept": {
			parameter: withOptions(t, data.NewOptionParameter(parameterName, parameterDescription).
				WithValidValues("one"), newOption("two")),
			valid: false,
		},
	}

	for name, test := range tests {
		if test.parameter.AreValuesValid() != test.valid {
			t.Errorf("%s must report valid=%t, but did not", name, test.valid)
		}

		err := test.parameter.ValidateValues()

		var validationError *exception.CliRoutingOptionValuesValidationError
		if !test.valid && !errors.As(err, &validationError) {
			t.Errorf("%s must report a validation error, but reported: %v", name, err)
		}

		if test.valid && err != nil {
			t.Errorf("%s must report no error, but reported: %v", name, err)
		}
	}
}

// withOptions fills the parameter and fails the test where the parameter rejects
// an option.
func withOptions(
	t *testing.T,
	parameter contract.OptionParameterContract,
	options ...contract.OptionContract,
) contract.OptionParameterContract {
	t.Helper()

	filled, err := parameter.WithOptions(options...)
	if err != nil {
		t.Fatalf("the parameter must take each option, but reported: %v", err)
	}

	return filled
}
