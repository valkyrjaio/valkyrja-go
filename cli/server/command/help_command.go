/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package command

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	routingdata "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	serverconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/server/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// commandOptionName names the option that says which command to describe.
const commandOptionName = "command"

// textIndent opens each line that the help text prints under a heading.
const textIndent = "  "

// listIndent opens the description of one parameter.
const listIndent = "    "

type HelpCommand struct{}

// NewHelpCommand builds the command.
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

// GetRoute returns the command that the router files.
func (c *HelpCommand) GetRoute() contract.RouteContract {
	return routingdata.NewRoute(
		serverconstant.CommandNameHelp,
		"Help for a command",
		c.Run,
	).
		WithHelpText(c.GetHelpText).
		WithOptions(
			routingdata.NewOptionParameter(
				commandOptionName,
				"The name of the command to get help for",
			).
				WithValueDisplayName(commandOptionName).
				WithMode(routingconstant.OptionModeRequired),
		)
}

// GetHelpText returns the help text of the command.
func (c *HelpCommand) GetHelpText() contract.MessageContract {
	return message.NewMessage("A command to get help for a specific command.")
}

// Run prints the help text of the command that the caller names.
func (c *HelpCommand) Run(
	container containercontract.ContainerContract,
	route contract.RouteContract,
) contract.OutputContract {
	outputFactory := getOutputFactory(container)
	name := route.GetOptionValue(commandOptionName, nil)
	collection := getCollection(container)

	if collection == nil || !collection.Has(name) {
		return outputFactory.CreateOutput(
			interactionconstant.ExitCodeError,
			message.NewBanner(message.NewErrorMessage("Command `"+name+"` was not found.")),
		)
	}

	help := helpMessages(collection.Get(name))

	messages := make([]contract.MessageContract, 0, len(help)+1)
	messages = append(messages, newHeader(getCliConfig(container), route))
	messages = append(messages, help...)

	return outputFactory.CreateOutput(interactionconstant.ExitCodeSuccess, messages...)
}

// helpMessages returns every line that the help text prints for the command.
func helpMessages(route contract.RouteContract) []contract.MessageContract {
	messages := []contract.MessageContract{
		message.NewNewLine(),
		headingWithText("Name: ", route.GetName()),
		message.NewNewLine(),
		message.NewNewLine(),
		headingOverText("Description:", message.NewMessage(route.GetDescription())),
		message.NewNewLine(),
		message.NewNewLine(),
		headingOverText("Usage:", message.NewMessage(usageText(route))),
		message.NewNewLine(),
		message.NewNewLine(),
	}

	messages = append(messages, argumentsMessages(route)...)
	messages = append(messages, optionsMessages(route)...)

	if !route.HasHelpText() {
		return messages
	}

	return append(messages,
		headingOverText("Help:", route.GetHelpTextMessage()),
		message.NewNewLine(),
		message.NewNewLine(),
	)
}

// headingWithText returns a heading and the text that follows it on the same
// line.
func headingWithText(heading string, text string) contract.MessageContract {
	return message.NewMessages(
		message.NewFormattedMessage(heading, format.NewHighlightedTextFormatter()),
		message.NewMessage(text),
	)
}

// headingOverText returns a heading and the text that follows it, indented on
// the next line.
func headingOverText(heading string, text contract.MessageContract) contract.MessageContract {
	return message.NewMessages(
		message.NewFormattedMessage(heading, format.NewHighlightedTextFormatter()),
		message.NewNewLine(),
		text.WithText(textIndent+strings.ReplaceAll(text.GetText(), "\n", "\n"+textIndent)),
	)
}

// usageText returns the line that says how a caller runs the command.
func usageText(route contract.RouteContract) string {
	usage := &strings.Builder{}
	usage.WriteString(route.GetName())

	if route.HasOptions() {
		usage.WriteString(" [options]")
	}

	usage.WriteString(" [global options]")

	for _, argument := range route.GetArguments() {
		usage.WriteString(" [" + argument.GetName())

		if argument.GetValueMode() == routingconstant.ArgumentValueModeArray {
			usage.WriteString("...")
		}

		usage.WriteString("]")
	}

	return usage.String()
}

// argumentsMessages returns the lines that describe each argument of the
// command.
func argumentsMessages(route contract.RouteContract) []contract.MessageContract {
	if !route.HasArguments() {
		return nil
	}

	messages := []contract.MessageContract{
		message.NewFormattedMessage("Arguments:", format.NewHighlightedTextFormatter()),
		message.NewNewLine(),
	}

	for _, argument := range route.GetArguments() {
		messages = append(messages, parameterMessages(argument.GetName(), argument.GetDescription())...)
	}

	return messages
}

// optionsMessages returns the lines that describe each option of the command,
// and each option that every command accepts.
func optionsMessages(route contract.RouteContract) []contract.MessageContract {
	messages := make([]contract.MessageContract, 0)

	if route.HasOptions() {
		messages = append(messages,
			message.NewFormattedMessage("Options:", format.NewHighlightedTextFormatter()),
			message.NewNewLine(),
		)

		for _, option := range route.GetOptions() {
			messages = append(messages, optionMessages(option)...)
		}
	}

	messages = append(messages,
		message.NewFormattedMessage("Global Options:", format.NewHighlightedTextFormatter()),
		message.NewNewLine(),
	)

	for _, option := range globalOptions() {
		messages = append(messages, optionMessages(option)...)
	}

	return messages
}

// globalOptions returns each option that every command of the application
// accepts.
func globalOptions() []contract.OptionParameterContract {
	return []contract.OptionParameterContract{
		routingdata.NewQuietOptionParameter(),
		routingdata.NewSilentOptionParameter(),
		routingdata.NewNoInteractionOptionParameter(),
		routingdata.NewHelpOptionParameter(),
		routingdata.NewVersionOptionParameter(),
	}
}

// optionMessages returns the lines that describe one option.
func optionMessages(option contract.OptionParameterContract) []contract.MessageContract {
	messages := []contract.MessageContract{
		message.NewMessage(textIndent),
		message.NewFormattedMessage("--"+option.GetName(), format.NewFormatter(
			format.NewTextColorFormat(interactionconstant.TextColorMagenta),
		)),
	}

	if shortNames := option.GetShortNames(); len(shortNames) > 0 {
		messages = append(messages,
			message.NewMessage(", "),
			message.NewFormattedMessage("-"+strings.Join(shortNames, "|"), format.NewFormatter(
				format.NewTextColorFormat(interactionconstant.TextColorMagenta),
			)),
		)
	}

	if option.HasValueDisplayName() {
		messages = append(messages,
			message.NewMessage(" "),
			message.NewFormattedMessage(valueDisplayText(option), format.NewHighlightedTextFormatter()),
		)
	}

	messages = append(messages, message.NewNewLine(), message.NewMessage(listIndent),
		message.NewMessage(option.GetDescription()))

	if validValues := option.GetValidValues(); len(validValues) > 0 {
		messages = append(messages,
			message.NewMessage(" ["+strings.Join(validValues, ", ")+"]"),
		)
	}

	return append(messages, message.NewNewLine(), message.NewNewLine())
}

// valueDisplayText returns how the help text spells the value of the option.
func valueDisplayText(option contract.OptionParameterContract) string {
	text := ""

	if option.GetValueMode() == routingconstant.OptionValueModeArray {
		text = "..."
	}

	if option.GetMode() == routingconstant.OptionModeRequired {
		return text + "=" + option.GetValueDisplayName()
	}

	return text + "[=" + option.GetValueDisplayName() + "]"
}

// parameterMessages returns the lines that describe one argument.
func parameterMessages(name string, description string) []contract.MessageContract {
	return []contract.MessageContract{
		message.NewMessage(textIndent),
		message.NewFormattedMessage(name, format.NewFormatter(
			format.NewTextColorFormat(interactionconstant.TextColorMagenta),
		)),
		message.NewNewLine(),
		message.NewMessage(listIndent),
		message.NewMessage(description),
		message.NewNewLine(),
		message.NewNewLine(),
	}
}
