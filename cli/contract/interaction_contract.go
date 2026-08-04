/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds every contract of the CLI component.
//
// The component keeps one `contract` package, for the reason that the container
// component keeps one: Go resolves an import cycle at the package level, and the
// contracts here name each other — a route names its middleware, and each
// middleware names the route back.
//
// Each `With` method returns a copy and leaves the receiver unchanged. The other
// ports return `this` or `static`; Go has no such return type, so each one
// returns the contract.
package contract

import (
	"io"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
)

// ArgumentContract is one positional argument of a command.
type ArgumentContract interface {
	// GetValue returns the value of the argument.
	GetValue() string

	// WithValue returns a copy of the argument with another value.
	WithValue(value string) ArgumentContract
}

// OptionContract is one option of a command.
type OptionContract interface {
	// GetName returns the name of the option.
	GetName() string

	// WithName returns a copy of the option under another name.
	WithName(name string) OptionContract

	// HasValue reports whether the caller gave the option a value.
	HasValue() bool

	// GetValue returns the value of the option.
	GetValue() string

	// WithValue returns a copy of the option with another value.
	WithValue(value string) OptionContract

	// WithoutValue returns a copy of the option with no value.
	WithoutValue() OptionContract

	// GetType returns the type of the option, which says how the caller spells
	// it.
	GetType() constant.OptionType

	// WithType returns a copy of the option of another type.
	WithType(optionType constant.OptionType) OptionContract
}

// FormatContract is one pair of terminal codes that wrap a piece of text.
type FormatContract interface {
	// GetSetCode returns the code that starts the format.
	GetSetCode() string

	// WithSetCode returns a copy of the format for another start code.
	WithSetCode(setCode string) FormatContract

	// GetUnsetCode returns the code that ends the format.
	GetUnsetCode() string

	// WithUnsetCode returns a copy of the format for another end code.
	WithUnsetCode(unsetCode string) FormatContract
}

// FormatterContract applies each of its formats to a piece of text.
type FormatterContract interface {
	// GetFormats returns each format that the formatter applies.
	GetFormats() []FormatContract

	// WithFormats returns a copy of the formatter with other formats.
	WithFormats(formats ...FormatContract) FormatterContract

	// FormatText returns the text, wrapped in each format.
	FormatText(text string) string
}

// MessageContract is one message that a command writes.
type MessageContract interface {
	// GetText returns the text of the message.
	GetText() string

	// GetFormattedText returns the text with the formatter applied.
	GetFormattedText() string

	// WithText returns a copy of the message with another text.
	WithText(text string) MessageContract

	// HasFormatter reports whether the message carries a formatter.
	HasFormatter() bool

	// GetFormatter returns the formatter of the message.
	GetFormatter() FormatterContract

	// WithFormatter returns a copy of the message with another formatter.
	WithFormatter(formatter FormatterContract) MessageContract

	// WithoutFormatter returns a copy of the message with no formatter.
	WithoutFormatter() MessageContract
}

// QuestionContract is a message that asks the caller for an answer.
type QuestionContract interface {
	MessageContract

	// GetCallable returns what the question runs once the caller answers.
	GetCallable() QuestionCallableFunc

	// WithCallable returns a copy of the question that runs another callable.
	WithCallable(callable QuestionCallableFunc) QuestionContract

	// GetAnswer returns the answer of the question.
	GetAnswer() AnswerContract

	// WithAnswer returns a copy of the question with another answer.
	WithAnswer(answer AnswerContract) QuestionContract

	// Ask asks the caller and returns what the caller answered.
	Ask() AnswerContract
}

// QuestionCallableFunc runs once the caller answers a question.
type QuestionCallableFunc func(output OutputContract, answer AnswerContract) OutputContract

// AnswerValidationFunc reports whether a response is one that the answer
// accepts.
type AnswerValidationFunc func(response string) bool

// AnswerContract is what the caller answered to one question.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type AnswerContract interface {
	MessageContract

	// GetDefaultResponse returns the response that the answer uses where the
	// caller gives none.
	GetDefaultResponse() string

	// WithDefaultResponse returns a copy of the answer with another default
	// response.
	WithDefaultResponse(defaultResponse string) AnswerContract

	// GetAllowedResponses returns each response that the answer accepts.
	GetAllowedResponses() []string

	// WithAllowedResponses returns a copy of the answer that accepts other
	// responses.
	WithAllowedResponses(allowedResponses ...string) AnswerContract

	// GetUserResponse returns what the caller typed.
	GetUserResponse() string

	// WithUserResponse returns a copy of the answer with another response from
	// the caller.
	WithUserResponse(userResponse string) AnswerContract

	// HasValidationCallable reports whether the answer validates a response
	// with a callable.
	HasValidationCallable() bool

	// GetValidationCallable returns what the answer validates a response with.
	GetValidationCallable() AnswerValidationFunc

	// WithValidationCallable returns a copy of the answer that validates with
	// another callable.
	WithValidationCallable(validationCallable AnswerValidationFunc) AnswerContract

	// WithoutValidationCallable returns a copy of the answer that validates
	// with no callable.
	WithoutValidationCallable() AnswerContract

	// HasBeenAnswered reports whether the caller answered already.
	HasBeenAnswered() bool

	// WithHasBeenAnswered returns a copy of the answer with another answered
	// flag.
	WithHasBeenAnswered(hasBeenAnswered bool) AnswerContract

	// IsValidResponse reports whether what the caller typed is a response that
	// the answer accepts.
	IsValidResponse() bool
}

// ProgressContract is a message that reports how far a command has come.
type ProgressContract interface {
	MessageContract

	// IsComplete reports whether the command finished.
	IsComplete() bool

	// WithIsComplete returns a copy of the progress with another complete
	// flag.
	WithIsComplete(isComplete bool) ProgressContract

	// GetPercentage returns how far the command has come, from 0 to 100.
	GetPercentage() int

	// WithPercentage returns a copy of the progress at another percentage.
	WithPercentage(percentage int) ProgressContract
}

// WriterContract writes one message to one destination.
type WriterContract interface {
	// ShouldWriteMessage reports whether this writer writes the message.
	ShouldWriteMessage(message MessageContract) bool

	// Write writes the message and returns the output that holds the result.
	Write(output OutputContract, message MessageContract) OutputContract
}

// InputContract is what the caller typed.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type InputContract interface {
	// GetCaller returns the path that the caller ran.
	GetCaller() string

	// WithCaller returns a copy of the input for another caller.
	WithCaller(caller string) InputContract

	// GetCommandName returns the name of the command that the caller ran.
	GetCommandName() string

	// WithCommandName returns a copy of the input for another command.
	WithCommandName(commandName string) InputContract

	// GetArguments returns each positional argument, in the order that the
	// caller typed them.
	GetArguments() []ArgumentContract

	// WithArguments returns a copy of the input with other arguments.
	WithArguments(arguments ...ArgumentContract) InputContract

	// WithAddedArgument returns a copy of the input with the argument appended.
	WithAddedArgument(argument ArgumentContract) InputContract

	// WithoutArgument returns a copy of the input without the argument that
	// carries the value.
	WithoutArgument(value string) InputContract

	// WithoutArguments returns a copy of the input with no argument.
	WithoutArguments() InputContract

	// GetOptions returns each option that the caller typed.
	GetOptions() []OptionContract

	// GetOption returns each option under the name. A caller repeats an option,
	// so the result is a list.
	GetOption(name string) []OptionContract

	// HasOption reports whether the caller typed the option.
	HasOption(name string) bool

	// WithOptions returns a copy of the input with other options.
	WithOptions(options ...OptionContract) InputContract

	// WithAddedOption returns a copy of the input with the option appended.
	WithAddedOption(option OptionContract) InputContract

	// WithoutOption returns a copy of the input without the option.
	WithoutOption(name string) InputContract

	// WithoutOptions returns a copy of the input with no option.
	WithoutOptions() InputContract
}

// OutputContract is what a command writes back.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type OutputContract interface {
	// GetMessages returns every message that the output holds.
	GetMessages() []MessageContract

	// GetWrittenMessages returns each message that a writer wrote already.
	GetWrittenMessages() []MessageContract

	// HasWrittenMessage reports whether a writer wrote a message already.
	HasWrittenMessage() bool

	// GetUnwrittenMessages returns each message that no writer wrote yet.
	GetUnwrittenMessages() []MessageContract

	// HasUnwrittenMessage reports whether a message is waiting to be written.
	HasUnwrittenMessage() bool

	// WithMessages returns a copy of the output with other messages.
	WithMessages(messages ...MessageContract) OutputContract

	// WithAddedMessages returns a copy of the output with the messages
	// appended.
	WithAddedMessages(messages ...MessageContract) OutputContract

	// WithAddedMessage returns a copy of the output with the message appended.
	WithAddedMessage(message MessageContract) OutputContract

	// WriteMessages writes every message that is waiting.
	WriteMessages() OutputContract

	// WriteMessage writes one message.
	WriteMessage(message MessageContract) OutputContract

	// GetWriters returns each writer that the output writes through.
	GetWriters() []WriterContract

	// WithWriters returns a copy of the output with other writers.
	WithWriters(writers ...WriterContract) OutputContract

	// IsInteractive reports whether the output asks the caller a question.
	IsInteractive() bool

	// WithIsInteractive returns a copy of the output with another interactive
	// flag.
	WithIsInteractive(isInteractive bool) OutputContract

	// IsQuiet reports whether the output writes less.
	IsQuiet() bool

	// WithIsQuiet returns a copy of the output with another quiet flag.
	WithIsQuiet(isQuiet bool) OutputContract

	// IsSilent reports whether the output writes nothing.
	IsSilent() bool

	// WithIsSilent returns a copy of the output with another silent flag.
	WithIsSilent(isSilent bool) OutputContract

	// GetExitCode returns the code that the process exits with.
	GetExitCode() constant.ExitCode

	// WithExitCode returns a copy of the output for another exit code.
	WithExitCode(exitCode constant.ExitCode) OutputContract
}

// PlainOutputContract is an output that applies no format.
//
// The contract is an alias, because it adds no method. Go compares an interface
// by its method set, so a second interface with the same methods is the same
// type.
type PlainOutputContract = OutputContract

// EmptyOutputContract is an output that writes nothing.
type EmptyOutputContract = OutputContract

// StreamOutputContract is an output that writes to a stream.
type StreamOutputContract = OutputContract

// FileOutputContract is an output that writes to a file.
type FileOutputContract = OutputContract

// OutputFactoryContract builds each kind of output.
//
// Go has no default parameter, so a caller passes the exit code that the other
// ports leave out.
type OutputFactoryContract interface {
	// CreateOutput builds an output that writes through the default writers.
	CreateOutput(exitCode constant.ExitCode, messages ...MessageContract) OutputContract

	// CreateEmptyOutput builds an output that writes nothing.
	CreateEmptyOutput(exitCode constant.ExitCode, messages ...MessageContract) EmptyOutputContract

	// CreatePlainOutput builds an output that applies no format.
	CreatePlainOutput(exitCode constant.ExitCode, messages ...MessageContract) PlainOutputContract

	// CreateFileOutput builds an output that writes to the file.
	CreateFileOutput(
		filepath string,
		exitCode constant.ExitCode,
		messages ...MessageContract,
	) FileOutputContract

	// CreateStreamOutput builds an output that writes to the writer.
	CreateStreamOutput(
		writer io.Writer,
		exitCode constant.ExitCode,
		messages ...MessageContract,
	) StreamOutputContract
}

// CliInteractionConfigContract holds the settings that apply to every output of
// the application.
type CliInteractionConfigContract interface {
	// IsQuiet reports whether an output writes less.
	IsQuiet() bool

	// IsInteractive reports whether an output asks the caller a question.
	IsInteractive() bool

	// IsSilent reports whether an output writes nothing.
	IsSilent() bool
}
