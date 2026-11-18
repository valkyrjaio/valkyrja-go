/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package validator runs a set of rules over a set of subjects.
package validator

import (
	"maps"
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/validation/contract"
)

// messageSeparator separates the name of a subject from the message of the rule
// that it failed.
const messageSeparator = ": "

type Validator struct {
	rules         map[string][]contract.RuleContract
	subjectOrder  []string
	errorMessages map[string]string
}

// NewValidator builds the validator over its rules, keyed by the name of the
// subject that each set applies to.
func NewValidator(rules map[string][]contract.RuleContract) *Validator {
	built := &Validator{errorMessages: map[string]string{}}

	built.SetRules(rules)

	return built
}

// ValidateRules reports whether every subject passes every rule of its own.
func (v *Validator) ValidateRules() bool {
	v.errorMessages = map[string]string{}

	for _, subject := range v.subjectOrder {
		v.validateSubject(subject)
	}

	return len(v.errorMessages) == 0
}

// SetRules replaces the rules, keyed by the name of the subject that each set
// applies to.
func (v *Validator) SetRules(rules map[string][]contract.RuleContract) {
	v.rules = maps.Clone(rules)

	if v.rules == nil {
		v.rules = map[string][]contract.RuleContract{}
	}

	v.subjectOrder = slices.Sorted(maps.Keys(v.rules))
}

// GetErrorMessages returns one message for each subject that failed.
func (v *Validator) GetErrorMessages() map[string]string {
	return maps.Clone(v.errorMessages)
}

// HasFirstErrorMessage reports whether a subject failed.
func (v *Validator) HasFirstErrorMessage() bool {
	return len(v.errorMessages) > 0
}

// GetFirstErrorMessage returns the message of the first subject that failed, and
// an empty string where none failed.
func (v *Validator) GetFirstErrorMessage() string {
	for _, subject := range v.subjectOrder {
		message, failed := v.errorMessages[subject]
		if failed {
			return message
		}
	}

	return ""
}

// validateSubject runs each rule of the subject, and records the message of the
// first one that it fails.
func (v *Validator) validateSubject(subject string) {
	for _, rule := range v.rules[subject] {
		err := rule.Validate()
		if err == nil {
			continue
		}

		v.errorMessages[subject] = subject + messageSeparator + err.Error()

		return
	}
}
