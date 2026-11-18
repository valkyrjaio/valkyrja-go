/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package rule

import (
	"github.com/valkyrjaio/valkyrja-go/v26/validation/constant"
)

// NewRequired builds the rule that the subject carries a value.
func NewRequired(subject any) *Rule {
	return NewRule(subject, constant.ErrorMessageRequired, func(held any) bool {
		return !isEmpty(held)
	})
}

// NewIsEmpty builds the rule that the subject carries no value.
func NewIsEmpty(subject any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsEmpty, isEmpty)
}

// NewNotEmpty builds the rule that the subject carries a value.
func NewNotEmpty(subject any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsNotEmpty, func(held any) bool {
		return !isEmpty(held)
	})
}

// NewIsBool builds the rule that the subject is a boolean.
func NewIsBool(subject any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsBool, func(held any) bool {
		_, isBool := held.(bool)

		return isBool
	})
}

// NewIsString builds the rule that the subject is a string.
func NewIsString(subject any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsString, func(held any) bool {
		_, isString := held.(string)

		return isString
	})
}

// NewIsNumeric builds the rule that the subject is a number, or a string that
// reads as one.
func NewIsNumeric(subject any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsNumeric, isNumeric)
}

// NewEqual builds the rule that the subject is the value.
func NewEqual(subject any, value any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsEqual, func(held any) bool {
		return held == value
	})
}

// NewNotEqual builds the rule that the subject is not the value.
func NewNotEqual(subject any, value any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsNotEqual, func(held any) bool {
		return held != value
	})
}

// NewEmail builds the rule that the subject is an email address.
func NewEmail(subject any) *Rule {
	return NewRule(subject, constant.ErrorMessageIsEmail, func(held any) bool {
		text, isString := held.(string)

		return isString && isEmail(text)
	})
}
