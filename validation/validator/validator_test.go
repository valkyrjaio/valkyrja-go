/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package validator_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/validation/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/rule"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/validator"
)

const (
	emailField = "email"
	nameField  = "name"
)

func TestTheValidatorReportsThatEverySubjectPasses(t *testing.T) {
	t.Parallel()

	built := validator.NewValidator(map[string][]contract.RuleContract{
		emailField: {rule.NewEmail("melech@example.com")},
		nameField:  {rule.NewRequired("Melech"), rule.NewMin("Melech", 2)},
	})

	if !built.ValidateRules() {
		t.Errorf("every subject passes, so the validator must report so, but reported: %v",
			built.GetErrorMessages())
	}

	if built.HasFirstErrorMessage() || built.GetFirstErrorMessage() != "" {
		t.Error("a validation that nothing failed must report no message, but reported one")
	}
}

func TestTheValidatorReportsOneMessageForEachSubjectThatFails(t *testing.T) {
	t.Parallel()

	built := validator.NewValidator(map[string][]contract.RuleContract{
		emailField: {rule.NewEmail("not an address")},
		nameField:  {rule.NewRequired(""), rule.NewMin("", 2)},
	})

	if built.ValidateRules() {
		t.Error("a subject that fails must report so, but reported that it passed")
	}

	messages := built.GetErrorMessages()

	if len(messages) != 2 {
		t.Fatalf("the validator must report one message for each subject, but reported: %d", len(messages))
	}

	if messages[emailField] != emailField+": "+constant.ErrorMessageIsEmail {
		t.Errorf("the message must name the subject and the rule, but read: %q", messages[emailField])
	}

	// The name failed two rules, and a person reads one message per field, so
	// only the first one is reported.
	if messages[nameField] != nameField+": "+constant.ErrorMessageRequired {
		t.Errorf("the message must be the one of the first rule that failed, but read: %q", messages[nameField])
	}
}

func TestTheFirstMessageIsTheOneOfTheFirstSubjectInTheOrderOfItsName(t *testing.T) {
	t.Parallel()

	built := validator.NewValidator(map[string][]contract.RuleContract{
		nameField:  {rule.NewRequired("")},
		emailField: {rule.NewEmail("not an address")},
	})

	built.ValidateRules()

	if !built.HasFirstErrorMessage() {
		t.Fatal("a subject failed, so the validator must report a message, but reported none")
	}

	// Go's map has no order, so the validator walks its subjects in the order of
	// their names, and `email` comes before `name`.
	if built.GetFirstErrorMessage() != emailField+": "+constant.ErrorMessageIsEmail {
		t.Errorf("the first message must be the one of the first subject, but read: %q",
			built.GetFirstErrorMessage())
	}
}

func TestSetRulesReplacesTheRules(t *testing.T) {
	t.Parallel()

	built := validator.NewValidator(map[string][]contract.RuleContract{
		nameField: {rule.NewRequired("")},
	})

	built.SetRules(map[string][]contract.RuleContract{
		nameField: {rule.NewRequired("Melech")},
	})

	if !built.ValidateRules() {
		t.Error("the validator must run the rules that replaced the ones it held, but did not")
	}
}

func TestAValidatorThatHoldsNoRuleReportsThatEverythingPasses(t *testing.T) {
	t.Parallel()

	built := validator.NewValidator(nil)

	if !built.ValidateRules() {
		t.Error("a validator that holds no rule must report that everything passes, but did not")
	}
}

func TestValidatingAgainDropsTheMessagesOfTheRunBefore(t *testing.T) {
	t.Parallel()

	built := validator.NewValidator(map[string][]contract.RuleContract{
		nameField: {rule.NewRequired("")},
	})

	built.ValidateRules()
	built.SetRules(map[string][]contract.RuleContract{
		nameField: {rule.NewRequired("Melech")},
	})

	if !built.ValidateRules() || built.HasFirstErrorMessage() {
		t.Error("a later run must drop the messages of the run before, but kept them")
	}
}

func TestTheMessagesAreACopy(t *testing.T) {
	t.Parallel()

	built := validator.NewValidator(map[string][]contract.RuleContract{
		nameField: {rule.NewRequired("")},
	})

	built.ValidateRules()

	messages := built.GetErrorMessages()
	delete(messages, nameField)

	if len(built.GetErrorMessages()) != 1 {
		t.Error("GetErrorMessages must return a copy, so a change to it must not reach the validator, but it did")
	}
}

func TestTheValidatorSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.ValidatorContract = validator.NewValidator(nil)

	if !built.ValidateRules() {
		t.Error("the validator must satisfy its contract, but did not")
	}
}
