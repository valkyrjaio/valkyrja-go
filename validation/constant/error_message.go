/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the validation component's messages and binding keys.
package constant

// The message that each rule reports where a subject fails it.
//
// A rule takes its own message from a caller, and these are the messages that
// the framework's own rules take.
const (
	ErrorMessageRequired = "This field is required."

	ErrorMessageIntGreaterThan = "This field value is too low."

	ErrorMessageIntLessThan = "This field value is too high."

	ErrorMessageIsEmail = "This field must be a valid email."

	ErrorMessageIsEqual = "This field must be the same."

	ErrorMessageIsBool = "This field must be a boolean."

	ErrorMessageIsEmpty = "This field must be empty."

	ErrorMessageIsNumeric = "This field must be numeric."

	ErrorMessageIsString = "This field must be a string."

	ErrorMessageIsNotEmpty = "This field must not be empty."

	ErrorMessageIsNotEqual = "This field must not be the same."

	ErrorMessageEntityExists = "This field must match an existing entity."

	ErrorMessageEntityNotExists = "This field must not match an existing entity."

	ErrorMessageStringAlpha = "This field must be alphanumeric."

	ErrorMessageStringContains = "This field must contain another string."

	ErrorMessageStringEndsWith = "This field must end with another string."

	ErrorMessageStringLowercase = "This field must be lowercase."

	ErrorMessageStringMax = "This field is too long."

	ErrorMessageStringMin = "This field is too short."

	ErrorMessageStringRegex = "This field must match the regex."

	ErrorMessageStringStartsWith = "This field must start with another string."

	ErrorMessageStringUppercase = "This field must be uppercase."
)

// ValidatorContractServiceID is the binding key for the validator.
const ValidatorContractServiceID = "valkyrja.validation.validator.ValidatorContract"
