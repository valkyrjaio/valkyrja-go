/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package rule

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/valkyrjaio/valkyrja-go/v26/validation/constant"
)

// NewAlpha builds the rule that every character of the subject is a letter.
func NewAlpha(subject any) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringAlpha, isAlphabetic)
}

// NewLowercase builds the rule that the subject carries no uppercase letter.
func NewLowercase(subject any) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringLowercase, isLowercase)
}

// NewUppercase builds the rule that the subject carries no lowercase letter.
func NewUppercase(subject any) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringUppercase, isUppercase)
}

// NewContains builds the rule that the subject carries the needle.
func NewContains(subject any, needle string) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringContains, func(text string) bool {
		return strings.Contains(text, needle)
	})
}

// NewStartsWith builds the rule that the subject starts with the needle.
func NewStartsWith(subject any, needle string) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringStartsWith, func(text string) bool {
		return strings.HasPrefix(text, needle)
	})
}

// NewEndsWith builds the rule that the subject ends with the needle.
func NewEndsWith(subject any, needle string) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringEndsWith, func(text string) bool {
		return strings.HasSuffix(text, needle)
	})
}

// NewMax builds the rule that the subject carries at most the number of
// characters.
func NewMax(subject any, longest int) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringMax, func(text string) bool {
		return utf8.RuneCountInString(text) <= longest
	})
}

// NewMin builds the rule that the subject carries at least the number of
// characters.
func NewMin(subject any, shortest int) *Rule {
	return newStringRule(subject, constant.ErrorMessageStringMin, func(text string) bool {
		return utf8.RuneCountInString(text) >= shortest
	})
}

// NewRegex builds the rule that the subject matches the regular expression.
//
// Warning: Go's regular expressions are RE2, and a pattern carries no delimiter.
// A pattern that PHP writes as `/^\d+$/` is `^\d+$` here. A pattern that RE2
// rejects fails every subject, because a rule has no return to carry the
// failure.
func NewRegex(subject any, pattern string) *Rule {
	compiled, err := regexp.Compile(pattern)

	return newStringRule(subject, constant.ErrorMessageStringRegex, func(text string) bool {
		return err == nil && text != "" && compiled.MatchString(text)
	})
}

// newStringRule builds a rule that fails a subject that is not a string, and
// otherwise gives the string to the check.
func newStringRule(subject any, errorMessage string, check func(text string) bool) *Rule {
	return NewRule(subject, errorMessage, func(held any) bool {
		text, isString := held.(string)

		return isString && check(text)
	})
}
