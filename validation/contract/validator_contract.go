/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds every contract of the validation component.
package contract

type RuleContract interface {
	// GetSubject returns what the rule validates.
	GetSubject() any

	// IsValid reports whether the subject passes the rule.
	IsValid() bool

	// Validate runs the rule. It reports a failure where the subject fails it.
	//
	// The other ports throw here. Go reports a failure with a returned error,
	// which the method naming rules require.
	Validate() error
}

type ValidatorContract interface {
	// ValidateRules reports whether every subject passes every rule of its
	// own.
	ValidateRules() bool

	// SetRules replaces the rules, keyed by the name of the subject that each
	// set applies to.
	SetRules(rules map[string][]RuleContract)

	// GetErrorMessages returns one message for each subject that failed.
	GetErrorMessages() map[string]string

	// HasFirstErrorMessage reports whether a subject failed.
	HasFirstErrorMessage() bool

	// GetFirstErrorMessage returns the message of the first subject that
	// failed.
	GetFirstErrorMessage() string
}
