/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package rule holds every rule that a subject either passes or fails.
//
// The other ports declare an abstract `Rule` that each rule extends, and each
// subclass overrides `isValid`. Go has no abstract type and no method override,
// so one struct holds the subject and the message, and a function decides
// whether the subject passes. Each constructor names that function.
package rule

import (
	"github.com/valkyrjaio/valkyrja-go/v26/validation/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/throwable/exception"
)

// CheckFunc reports whether the subject passes a rule.
type CheckFunc func(subject any) bool

// Rule is one rule that a subject either passes or fails.
type Rule struct {
	subject      any
	errorMessage string
	check        CheckFunc
}

// NewRule builds a rule over a subject, the message that a failure reports, and
// what decides whether the subject passes.
func NewRule(subject any, errorMessage string, check CheckFunc) *Rule {
	return &Rule{
		subject:      subject,
		errorMessage: errorMessage,
		check:        check,
	}
}

// GetSubject returns what the rule validates.
func (r *Rule) GetSubject() any {
	return r.subject
}

// GetErrorMessage returns the message that a failure reports.
func (r *Rule) GetErrorMessage() string {
	return r.errorMessage
}

// IsValid reports whether the subject passes the rule.
func (r *Rule) IsValid() bool {
	return r.check(r.subject)
}

// Validate reports a failure where the subject fails the rule.
func (r *Rule) Validate() error {
	if r.IsValid() {
		return nil
	}

	return exception.NewValidationRuleFailureError(r.errorMessage)
}

// A rule satisfies its contract, which the compiler checks at build time rather
// than at run time.
var _ contract.RuleContract = (*Rule)(nil)
