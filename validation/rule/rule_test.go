/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package rule_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/validation/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/rule"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/throwable/exception"
)

const subjectText = "the subject"

func TestARuleReadsItsSubjectAndItsMessage(t *testing.T) {
	t.Parallel()

	built := rule.NewRule(subjectText, constant.ErrorMessageRequired, func(_ any) bool {
		return true
	})

	if built.GetSubject() != subjectText {
		t.Errorf("the rule must read its subject, but read: %v", built.GetSubject())
	}

	if built.GetErrorMessage() != constant.ErrorMessageRequired {
		t.Errorf("the rule must read its message, but read: %q", built.GetErrorMessage())
	}
}

func TestARuleReportsAFailureOnlyWhereTheSubjectFailsIt(t *testing.T) {
	t.Parallel()

	passing := rule.NewRule(subjectText, constant.ErrorMessageRequired, func(_ any) bool {
		return true
	})

	failing := rule.NewRule(subjectText, constant.ErrorMessageRequired, func(_ any) bool {
		return false
	})

	if !passing.IsValid() || passing.Validate() != nil {
		t.Error("a subject that passes must report no failure, but reported one")
	}

	err := failing.Validate()

	if _, isFailure := errors.AsType[*exception.ValidationRuleFailureError](err); !isFailure {
		t.Errorf("a subject that fails must report a rule failure, but reported: %v", err)
	}

	if err.Error() != constant.ErrorMessageRequired {
		t.Errorf("the failure must carry the message of the rule, but carried: %q", err.Error())
	}
}

func TestEachIsRuleReadsItsSubject(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		rule  *rule.Rule
		valid bool
	}{
		"a required subject that carries a value": {rule: rule.NewRequired(subjectText), valid: true},
		"a required subject that carries none":    {rule: rule.NewRequired(""), valid: false},
		"an empty subject that carries none":      {rule: rule.NewIsEmpty(""), valid: true},
		"an empty subject that carries a value":   {rule: rule.NewIsEmpty(subjectText), valid: false},
		"a not-empty subject that carries one":    {rule: rule.NewNotEmpty(subjectText), valid: true},
		"a not-empty subject that carries none":   {rule: rule.NewNotEmpty(nil), valid: false},
		"a boolean subject":                       {rule: rule.NewIsBool(true), valid: true},
		"a subject that is not a boolean":         {rule: rule.NewIsBool(1), valid: false},
		"a string subject":                        {rule: rule.NewIsString(subjectText), valid: true},
		"a subject that is not a string":          {rule: rule.NewIsString(1), valid: false},
		"a numeric subject":                       {rule: rule.NewIsNumeric(1), valid: true},
		"a string that reads as a number":         {rule: rule.NewIsNumeric("1.5"), valid: true},
		"a string that reads as no number":        {rule: rule.NewIsNumeric(subjectText), valid: false},
		"a subject that is not numeric":           {rule: rule.NewIsNumeric(true), valid: false},
		"a subject that is the value":             {rule: rule.NewEqual(1, 1), valid: true},
		"a subject that is another value":         {rule: rule.NewEqual(1, 2), valid: false},
		"a subject that is not the value":         {rule: rule.NewNotEqual(1, 2), valid: true},
		"a subject that is the same value":        {rule: rule.NewNotEqual(1, 1), valid: false},
		"an email address":                        {rule: rule.NewEmail("melech@example.com"), valid: true},
		"an address with a display name":          {rule: rule.NewEmail("M <m@example.com>"), valid: false},
		"a string that is no address":             {rule: rule.NewEmail(subjectText), valid: false},
		"a subject that is not a string address":  {rule: rule.NewEmail(1), valid: false},
	}

	for name, test := range tests {
		if test.rule.IsValid() != test.valid {
			t.Errorf("%s must report valid=%t, but did not", name, test.valid)
		}
	}
}

func TestTheEmptyRuleReadsEachKindOfSubject(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		subject any
		empty   bool
	}{
		"nothing at all":     {subject: nil, empty: true},
		"an empty string":    {subject: "", empty: true},
		"a string":           {subject: subjectText, empty: false},
		"a false":            {subject: false, empty: true},
		"a true":             {subject: true, empty: false},
		"a zero":             {subject: 0, empty: true},
		"a number":           {subject: 1, empty: false},
		"a zero as an int64": {subject: int64(0), empty: true},
		"a zero as a float":  {subject: 0.0, empty: true},
		"an empty list":      {subject: []any{}, empty: true},
		"a list":             {subject: []any{1}, empty: false},
		"an empty map":       {subject: map[string]any{}, empty: true},
		"a map":              {subject: map[string]any{"one": 1}, empty: false},
		"another type":       {subject: struct{}{}, empty: false},
	}

	for name, test := range tests {
		if rule.NewIsEmpty(test.subject).IsValid() != test.empty {
			t.Errorf("%s must report empty=%t, but did not", name, test.empty)
		}
	}
}

func TestEachStringRuleReadsItsSubject(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		rule  *rule.Rule
		valid bool
	}{
		"a subject of letters":                {rule: rule.NewAlpha("abc"), valid: true},
		"a subject with a digit":              {rule: rule.NewAlpha("ab1"), valid: false},
		"an empty subject of letters":         {rule: rule.NewAlpha(""), valid: false},
		"a subject that is not a string":      {rule: rule.NewAlpha(1), valid: false},
		"a lowercase subject":                 {rule: rule.NewLowercase("abc"), valid: true},
		"a subject with an uppercase letter":  {rule: rule.NewLowercase("aBc"), valid: false},
		"an uppercase subject":                {rule: rule.NewUppercase("ABC"), valid: true},
		"a subject with a lowercase letter":   {rule: rule.NewUppercase("AbC"), valid: false},
		"a subject that carries the needle":   {rule: rule.NewContains("abcdef", "cd"), valid: true},
		"a subject that carries none":         {rule: rule.NewContains("abcdef", "xy"), valid: false},
		"a subject that starts with it":       {rule: rule.NewStartsWith("abcdef", "ab"), valid: true},
		"a subject that starts with another":  {rule: rule.NewStartsWith("abcdef", "cd"), valid: false},
		"a subject that ends with it":         {rule: rule.NewEndsWith("abcdef", "ef"), valid: true},
		"a subject that ends with another":    {rule: rule.NewEndsWith("abcdef", "cd"), valid: false},
		"a subject that is short enough":      {rule: rule.NewMax("abc", 3), valid: true},
		"a subject that is too long":          {rule: rule.NewMax("abcd", 3), valid: false},
		"a subject that is long enough":       {rule: rule.NewMin("abc", 3), valid: true},
		"a subject that is too short":         {rule: rule.NewMin("ab", 3), valid: false},
		"a subject that matches the pattern":  {rule: rule.NewRegex("123", `^\d+$`), valid: true},
		"a subject that matches another":      {rule: rule.NewRegex("abc", `^\d+$`), valid: false},
		"an empty subject against a pattern":  {rule: rule.NewRegex("", `^\d*$`), valid: false},
		"a pattern that no engine compiles":   {rule: rule.NewRegex("abc", `(?=abc)`), valid: false},
		"a subject that is not a text needle": {rule: rule.NewContains(1, "cd"), valid: false},
	}

	for name, test := range tests {
		if test.rule.IsValid() != test.valid {
			t.Errorf("%s must report valid=%t, but did not", name, test.valid)
		}
	}
}

func TestEachIntRuleReadsItsSubject(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		rule  *rule.Rule
		valid bool
	}{
		"a subject above the number":     {rule: rule.NewGreaterThan(2, 1), valid: true},
		"a subject at the number":        {rule: rule.NewGreaterThan(1, 1), valid: false},
		"a subject below the number":     {rule: rule.NewLessThan(1, 2), valid: true},
		"a subject at the upper number":  {rule: rule.NewLessThan(2, 2), valid: false},
		"a subject that is not a number": {rule: rule.NewGreaterThan("2", 1), valid: false},
		"another subject that is no int": {rule: rule.NewLessThan("1", 2), valid: false},
	}

	for name, test := range tests {
		if test.rule.IsValid() != test.valid {
			t.Errorf("%s must report valid=%t, but did not", name, test.valid)
		}
	}
}
