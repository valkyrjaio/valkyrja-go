/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
)

const (
	questionText = "What is your name?"
	yesResponse  = "yes"
)

// errReadFailed is what a failing reader reports.
var errReadFailed = errors.New("the reader failed")

type errReader struct{}

// Read reports the failure.
func (r *errReader) Read(_ []byte) (int, error) {
	return 0, errReadFailed
}

type partialReader struct {
	written bool
}

// Read writes the line once, then reports the failure.
func (r *partialReader) Read(into []byte) (int, error) {
	if r.written {
		return 0, errReadFailed
	}

	r.written = true

	return copy(into, "Melech"), nil
}

func TestAnAnswerReadsTheResponseThatTheCallerGave(t *testing.T) {
	t.Parallel()

	built := message.NewAnswer("yes or no", yesResponse, "no")

	if built.HasBeenAnswered() || built.GetUserResponse() != "" {
		t.Error("an answer that nobody gave must carry no response, but carried one")
	}

	answered := built.WithUserResponse(yesResponse)

	if !answered.HasBeenAnswered() || answered.GetUserResponse() != yesResponse {
		t.Error("the answer must carry the response that the caller gave, but did not")
	}

	if built.HasBeenAnswered() {
		t.Error("WithUserResponse must leave the receiver unchanged, but did not")
	}
}

func TestAnAnswerFallsBackToItsDefaultResponse(t *testing.T) {
	t.Parallel()

	built := message.NewAnswer("yes or no", yesResponse, "no").WithDefaultResponse("no")

	if built.GetDefaultResponse() != "no" || built.GetUserResponse() != "no" {
		t.Error("an answer that nobody gave must report its default response, but did not")
	}

	if built.WithUserResponse(yesResponse).GetUserResponse() != yesResponse {
		t.Error("a response that the caller gave must replace the default, but did not")
	}
}

func TestAnAnswerAcceptsAResponseThatItsListNames(t *testing.T) {
	t.Parallel()

	built := message.NewAnswer("yes or no", yesResponse, "no")

	if len(built.GetAllowedResponses()) != 2 {
		t.Errorf("the answer must accept each response, but accepted: %v", built.GetAllowedResponses())
	}

	if !built.WithUserResponse(yesResponse).IsValidResponse() {
		t.Error("a response that the list names must be valid, but was not")
	}

	if built.WithUserResponse("maybe").IsValidResponse() {
		t.Error("a response that the list does not name must be invalid, but was valid")
	}

	replaced := built.WithAllowedResponses("maybe")

	if !replaced.WithUserResponse("maybe").IsValidResponse() {
		t.Error("WithAllowedResponses must hold the new list, but did not")
	}
}

func TestAnAnswerAcceptsAResponseThatItsCallableAccepts(t *testing.T) {
	t.Parallel()

	built := message.NewAnswer("a number").WithValidationCallable(func(response string) bool {
		return response == "42"
	})

	if !built.HasValidationCallable() || built.GetValidationCallable() == nil {
		t.Error("the answer must carry the callable, but carried none")
	}

	if !built.WithUserResponse("42").IsValidResponse() {
		t.Error("a response that the callable accepts must be valid, but was not")
	}

	if built.WithUserResponse("41").IsValidResponse() {
		t.Error("a response that the callable rejects must be invalid, but was valid")
	}

	if built.WithoutValidationCallable().HasValidationCallable() {
		t.Error("WithoutValidationCallable must remove the callable, but did not")
	}
}

func TestWithHasBeenAnsweredReturnsACopy(t *testing.T) {
	t.Parallel()

	built := message.NewAnswer("yes or no")

	if !built.WithHasBeenAnswered(true).HasBeenAnswered() || built.HasBeenAnswered() {
		t.Error("WithHasBeenAnswered must hold the new flag and leave the receiver unchanged, but did not")
	}
}

func TestEachAnswerMessageMethodKeepsTheAnswer(t *testing.T) {
	t.Parallel()

	built := message.NewAnswer("yes or no", yesResponse).WithUserResponse(yesResponse)

	// A promoted `With` would return a plain message and drop the response, so
	// each one must return an answer that still carries it.
	assertStillAnAnswer(t, "WithText", built.WithText("other"))
	assertStillAnAnswer(t, "WithFormatter", built.WithFormatter(format.NewErrorFormatter()))
	assertStillAnAnswer(t, "WithoutFormatter", built.WithoutFormatter())

	if built.WithText("other").GetText() != "other" {
		t.Error("WithText must hold the new text, but did not")
	}
}

func TestAQuestionReadsOneLine(t *testing.T) {
	t.Parallel()

	built := message.NewQuestionForReader(
		questionText,
		nil,
		message.NewAnswer("a name"),
		strings.NewReader("Melech\n"),
	)

	if built.Ask().GetUserResponse() != "Melech" {
		t.Errorf("the question must read the line, but read: %q", built.Ask().GetUserResponse())
	}
}

func TestAQuestionReadsALastLineThatCarriesNoBreak(t *testing.T) {
	t.Parallel()

	built := message.NewQuestionForReader(
		questionText,
		nil,
		message.NewAnswer("a name"),
		strings.NewReader("Melech"),
	)

	if built.Ask().GetUserResponse() != "Melech" {
		t.Error("a last line that carries no break must still be read, but was not")
	}
}

func TestAQuestionReturnsItsAnswerWhereItReadsNothing(t *testing.T) {
	t.Parallel()

	tests := map[string]message.Question{
		"a reader that reports a failure": *message.NewQuestionForReader(
			questionText, nil, message.NewAnswer("a name"), &errReader{},
		),
		"a line that carries only space": *message.NewQuestionForReader(
			questionText, nil, message.NewAnswer("a name"), strings.NewReader("   \n"),
		),
		"a line that a failure cut short": *message.NewQuestionForReader(
			questionText, nil, message.NewAnswer("a name"), &partialReader{},
		),
	}

	for name, built := range tests {
		if built.Ask().HasBeenAnswered() {
			t.Errorf("%s must return the answer as it is, but did not", name)
		}
	}
}

func TestEachQuestionWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	answer := message.NewAnswer("a name")
	built := message.NewQuestion(questionText, nil, answer)

	callable := func(output contract.OutputContract, _ contract.AnswerContract) contract.OutputContract {
		return output
	}

	if built.WithCallable(callable).GetCallable() == nil {
		t.Error("WithCallable must hold the new callable, but did not")
	}

	other := message.NewAnswer("another name")
	if built.WithAnswer(other).GetAnswer().GetText() != "another name" {
		t.Error("WithAnswer must hold the new answer, but did not")
	}

	if built.GetCallable() != nil || built.GetAnswer().GetText() != "a name" {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachQuestionMessageMethodKeepsTheQuestion(t *testing.T) {
	t.Parallel()

	built := message.NewQuestion(questionText, nil, message.NewAnswer("a name"))

	assertStillAQuestion(t, "WithText", built.WithText("other"))
	assertStillAQuestion(t, "WithFormatter", built.WithFormatter(format.NewErrorFormatter()))
	assertStillAQuestion(t, "WithoutFormatter", built.WithoutFormatter())

	if built.WithoutFormatter().GetFormattedText() != questionText {
		t.Error("WithoutFormatter must remove the formatter, but did not")
	}
}

func TestAProgressReportsHowFarACommandHasCome(t *testing.T) {
	t.Parallel()

	built := message.NewProgress("Building")

	if built.IsComplete() || built.GetPercentage() != 0 {
		t.Error("a progress that has not started must report zero, but did not")
	}

	if built.WithPercentage(50).GetPercentage() != 50 {
		t.Error("WithPercentage must hold the new percentage, but did not")
	}

	if built.WithPercentage(50).IsComplete() {
		t.Error("a progress below 100 must report that the command has not finished, but did not")
	}
}

func TestAProgressAndItsCompleteFlagNeverDisagree(t *testing.T) {
	t.Parallel()

	built := message.NewProgress("Building")

	if built.WithPercentage(100).IsComplete() != true {
		t.Error("a progress at 100 must report that the command finished, but did not")
	}

	if built.WithIsComplete(true).GetPercentage() != 100 {
		t.Error("a progress that reports the command finished must report 100, but did not")
	}

	if built.WithIsComplete(false).GetPercentage() != 0 {
		t.Error("a progress that reports the command did not finish must keep its percentage, but did not")
	}
}

func TestAPercentageOutsideTheRangeTakesTheNearestEnd(t *testing.T) {
	t.Parallel()

	built := message.NewProgress("Building")

	if built.WithPercentage(-10).GetPercentage() != 0 {
		t.Error("a percentage below zero must take zero, but did not")
	}

	if built.WithPercentage(110).GetPercentage() != 100 {
		t.Error("a percentage above 100 must take 100, but did not")
	}
}

func TestEachProgressMessageMethodKeepsTheProgress(t *testing.T) {
	t.Parallel()

	built := message.NewProgress("Building").WithPercentage(50)

	assertStillAProgress(t, "WithText", built.WithText("other"))
	assertStillAProgress(t, "WithFormatter", built.WithFormatter(format.NewErrorFormatter()))
	assertStillAProgress(t, "WithoutFormatter", built.WithoutFormatter())
}

// assertStillAnAnswer fails the test where the message is no longer an answer
// that carries its response.
func assertStillAnAnswer(t *testing.T, name string, built contract.MessageContract) {
	t.Helper()

	answer, isAnswer := built.(contract.AnswerContract)
	if !isAnswer {
		t.Fatalf("%s must return an answer, but returned a plain message", name)
	}

	if answer.GetUserResponse() != yesResponse {
		t.Errorf("%s must keep the response, but returned: %q", name, answer.GetUserResponse())
	}
}

// assertStillAQuestion fails the test where the message is no longer a question
// that carries its answer.
func assertStillAQuestion(t *testing.T, name string, built contract.MessageContract) {
	t.Helper()

	question, isQuestion := built.(contract.QuestionContract)
	if !isQuestion {
		t.Fatalf("%s must return a question, but returned a plain message", name)
	}

	if question.GetAnswer() == nil {
		t.Errorf("%s must keep the answer, but dropped it", name)
	}
}

// assertStillAProgress fails the test where the message is no longer a progress
// that carries its percentage.
func assertStillAProgress(t *testing.T, name string, built contract.MessageContract) {
	t.Helper()

	progress, isProgress := built.(contract.ProgressContract)
	if !isProgress {
		t.Fatalf("%s must return a progress, but returned a plain message", name)
	}

	if progress.GetPercentage() != 50 {
		t.Errorf("%s must keep the percentage, but returned: %d", name, progress.GetPercentage())
	}
}
