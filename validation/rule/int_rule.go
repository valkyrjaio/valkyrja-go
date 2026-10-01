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

// NewGreaterThan builds the rule that the subject is greater than the number.
func NewGreaterThan(subject any, lowest int) *Rule {
	return newIntRule(subject, constant.ErrorMessageIntGreaterThan, func(number int) bool {
		return number > lowest
	})
}

// NewLessThan builds the rule that the subject is less than the number.
func NewLessThan(subject any, highest int) *Rule {
	return newIntRule(subject, constant.ErrorMessageIntLessThan, func(number int) bool {
		return number < highest
	})
}

// newIntRule builds a rule that fails a subject that is not a whole number, and
// otherwise gives the number to the check.
func newIntRule(subject any, errorMessage string, check func(number int) bool) *Rule {
	return NewRule(subject, errorMessage, func(held any) bool {
		number, isInt := held.(int)

		return isInt && check(number)
	})
}
