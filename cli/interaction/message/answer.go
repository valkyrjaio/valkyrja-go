/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message

import (
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

type Answer struct {
	Message

	defaultResponse    string
	allowedResponses   []string
	userResponse       string
	validationCallable contract.AnswerValidationFunc
	hasBeenAnswered    bool
}

// NewAnswer builds an answer that carries the text, and that accepts each
// response.
func NewAnswer(text string, allowedResponses ...string) *Answer {
	return &Answer{
		Message:          *NewMessage(text),
		allowedResponses: allowedResponses,
	}
}

// WithText returns a copy of the answer with another text.
func (a *Answer) WithText(text string) contract.MessageContract {
	copied := *a
	copied.text = text

	return &copied
}

// WithFormatter returns a copy of the answer with another formatter.
func (a *Answer) WithFormatter(formatter contract.FormatterContract) contract.MessageContract {
	copied := *a
	copied.formatter = formatter

	return &copied
}

// WithoutFormatter returns a copy of the answer with no formatter.
func (a *Answer) WithoutFormatter() contract.MessageContract {
	copied := *a
	copied.formatter = nil

	return &copied
}

// GetDefaultResponse returns the response that the answer uses where the caller
// gives none.
func (a *Answer) GetDefaultResponse() string {
	return a.defaultResponse
}

// WithDefaultResponse returns a copy of the answer with another default
// response.
func (a *Answer) WithDefaultResponse(defaultResponse string) contract.AnswerContract {
	copied := *a
	copied.defaultResponse = defaultResponse

	return &copied
}

// GetAllowedResponses returns each response that the answer accepts.
func (a *Answer) GetAllowedResponses() []string {
	return a.allowedResponses
}

// WithAllowedResponses returns a copy of the answer that accepts other
// responses.
func (a *Answer) WithAllowedResponses(allowedResponses ...string) contract.AnswerContract {
	copied := *a
	copied.allowedResponses = allowedResponses

	return &copied
}

// GetUserResponse returns what the caller typed, and the default response where
// the caller typed nothing.
func (a *Answer) GetUserResponse() string {
	if a.userResponse == "" {
		return a.defaultResponse
	}

	return a.userResponse
}

// WithUserResponse returns a copy of the answer with another response from the
// caller.
func (a *Answer) WithUserResponse(userResponse string) contract.AnswerContract {
	copied := *a
	copied.userResponse = userResponse
	copied.hasBeenAnswered = true

	return &copied
}

// HasValidationCallable reports whether the answer validates a response with a
// callable.
func (a *Answer) HasValidationCallable() bool {
	return a.validationCallable != nil
}

// GetValidationCallable returns what the answer validates a response with, and
// nil where it validates with none.
func (a *Answer) GetValidationCallable() contract.AnswerValidationFunc {
	return a.validationCallable
}

// WithValidationCallable returns a copy of the answer that validates with
// another callable.
func (a *Answer) WithValidationCallable(
	validationCallable contract.AnswerValidationFunc,
) contract.AnswerContract {
	copied := *a
	copied.validationCallable = validationCallable

	return &copied
}

// WithoutValidationCallable returns a copy of the answer that validates with no
// callable.
func (a *Answer) WithoutValidationCallable() contract.AnswerContract {
	copied := *a
	copied.validationCallable = nil

	return &copied
}

// HasBeenAnswered reports whether the caller answered already.
func (a *Answer) HasBeenAnswered() bool {
	return a.hasBeenAnswered
}

// WithHasBeenAnswered returns a copy of the answer with another answered flag.
func (a *Answer) WithHasBeenAnswered(hasBeenAnswered bool) contract.AnswerContract {
	copied := *a
	copied.hasBeenAnswered = hasBeenAnswered

	return &copied
}

// IsValidResponse reports whether what the caller typed is a response that the
// answer accepts.
func (a *Answer) IsValidResponse() bool {
	response := a.GetUserResponse()

	if slices.Contains(a.allowedResponses, response) {
		return true
	}

	return a.validationCallable != nil && a.validationCallable(response)
}

// An answer satisfies its contract, which the compiler checks at build time
// rather than at run time.
var _ contract.AnswerContract = (*Answer)(nil)
