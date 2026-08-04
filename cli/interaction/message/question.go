/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
)

// Question is a message that asks the caller for an answer.
type Question struct {
	Message

	callable contract.QuestionCallableFunc
	answer   contract.AnswerContract
	reader   io.Reader
}

// NewQuestion builds a question that carries the text, and that reads the
// standard input of the process.
//
// A caller that reads somewhere else uses `NewQuestionForReader`, which is what
// a test does.
func NewQuestion(
	text string,
	callable contract.QuestionCallableFunc,
	answer contract.AnswerContract,
) *Question {
	return &Question{
		Message:  *NewFormattedMessage(text, format.NewQuestionFormatter()),
		callable: callable,
		answer:   answer,
		reader:   os.Stdin,
	}
}

// NewQuestionForReader builds a question that reads the reader.
func NewQuestionForReader(
	text string,
	callable contract.QuestionCallableFunc,
	answer contract.AnswerContract,
	reader io.Reader,
) *Question {
	built := NewQuestion(text, callable, answer)
	built.reader = reader

	return built
}

// WithText returns a copy of the question with another text.
//
// The method is declared here rather than promoted from the embedded message,
// because a promoted `With` copies only the embedded struct and would return a
// plain message, dropping every field of the question.
func (q *Question) WithText(text string) contract.MessageContract {
	copied := *q
	copied.text = text

	return &copied
}

// WithFormatter returns a copy of the question with another formatter.
func (q *Question) WithFormatter(formatter contract.FormatterContract) contract.MessageContract {
	copied := *q
	copied.formatter = formatter

	return &copied
}

// WithoutFormatter returns a copy of the question with no formatter.
func (q *Question) WithoutFormatter() contract.MessageContract {
	copied := *q
	copied.formatter = nil

	return &copied
}

// GetCallable returns what the question runs once the caller answers.
func (q *Question) GetCallable() contract.QuestionCallableFunc {
	return q.callable
}

// WithCallable returns a copy of the question that runs another callable.
func (q *Question) WithCallable(callable contract.QuestionCallableFunc) contract.QuestionContract {
	copied := *q
	copied.callable = callable

	return &copied
}

// GetAnswer returns the answer of the question.
func (q *Question) GetAnswer() contract.AnswerContract {
	return q.answer
}

// WithAnswer returns a copy of the question with another answer.
func (q *Question) WithAnswer(answer contract.AnswerContract) contract.QuestionContract {
	copied := *q
	copied.answer = answer

	return &copied
}

// Ask reads one line and returns the answer that carries it.
//
// A reader that reports a failure, and a line that carries nothing but
// whitespace, both return the answer as it is. The answer then reports its own
// default response, so a caller that gives no answer still reads one.
func (q *Question) Ask() contract.AnswerContract {
	line, err := bufio.NewReader(q.reader).ReadString('\n')

	response := strings.TrimSpace(line)

	// A reader at its end returns the last line together with `io.EOF`, so the
	// line is read before the failure is.
	if response == "" {
		return q.answer
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return q.answer
	}

	return q.answer.WithUserResponse(response)
}

// A question satisfies its contract, which the compiler checks at build time
// rather than at run time.
var _ contract.QuestionContract = (*Question)(nil)
