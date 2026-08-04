/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package command

import (
	"sort"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	routingdata "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	serverconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/server/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// namespaceOptionName names the option that filters the list.
const namespaceOptionName = "namespace"

// ListCommand prints every command of the application.
type ListCommand struct{}

// NewListCommand builds the command.
func NewListCommand() *ListCommand {
	return &ListCommand{}
}

// GetRoute returns the command that the router files.
func (c *ListCommand) GetRoute() contract.RouteContract {
	return routingdata.NewRoute(
		serverconstant.CommandNameList,
		"List all commands",
		c.Run,
	).
		WithHelpText(c.GetHelpText).
		WithOptions(
			routingdata.NewOptionParameter(
				namespaceOptionName,
				"An optional namespace to filter commands by",
			).
				WithShortNames("n").
				WithValueDisplayName(namespaceOptionName),
		)
}

// GetHelpText returns the help text of the command.
func (c *ListCommand) GetHelpText() contract.MessageContract {
	return message.NewMessage("A command to list all the commands present within the Cli component.")
}

// Run prints every command of the application, in the order of its name.
//
// The namespace option keeps the commands whose name starts with the namespace.
// A namespace that matches no command reports a failure.
func (c *ListCommand) Run(
	container containercontract.ContainerContract,
	route contract.RouteContract,
) contract.OutputContract {
	outputFactory := getOutputFactory(container)
	namespace := getOptionFirstValue(route, namespaceOptionName)
	routes := sortedRoutes(getCollection(container), namespace)

	if len(routes) == 0 {
		return outputFactory.CreateOutput(
			interactionconstant.ExitCodeError,
			message.NewBanner(message.NewErrorMessage(noRoutesText(namespace))),
		)
	}

	messages := []contract.MessageContract{
		newHeader(getCliConfig(container), route),
		message.NewFormattedMessage(headingText(namespace), format.NewHighlightedTextFormatter()),
		message.NewNewLine(),
	}

	for _, listed := range routes {
		messages = append(messages, routeMessages(listed)...)
	}

	return outputFactory.CreateOutput(
		interactionconstant.ExitCodeSuccess,
		append(messages, message.NewNewLine())...,
	)
}

// getOptionFirstValue returns the first value that the caller gave the option,
// and an empty string where the caller gave none.
func getOptionFirstValue(route contract.RouteContract, name string) string {
	option := route.GetOption(name)
	if option == nil {
		return ""
	}

	return option.GetFirstValue()
}

// sortedRoutes returns every command under the namespace, in the order of its
// name.
func sortedRoutes(collection contract.RouteCollectionContract, namespace string) []contract.RouteContract {
	if collection == nil {
		return nil
	}

	all := collection.All()
	kept := make([]contract.RouteContract, 0, len(all))

	for _, held := range all {
		if strings.HasPrefix(held.GetName(), namespace) {
			kept = append(kept, held)
		}
	}

	sort.Slice(kept, func(first int, second int) bool {
		return kept[first].GetName() < kept[second].GetName()
	})

	return kept
}

// noRoutesText returns what the command reports where it lists no command.
func noRoutesText(namespace string) string {
	if namespace == "" {
		return "No routes found."
	}

	return "Namespace `" + namespace + "` was not found."
}

// headingText returns the heading that the list prints above the commands.
func headingText(namespace string) string {
	if namespace == "" {
		return "Commands:"
	}

	return "Commands [" + namespace + "]:"
}

// routeMessages returns the lines that the list prints for one command.
func routeMessages(route contract.RouteContract) []contract.MessageContract {
	return []contract.MessageContract{
		message.NewMessage("  "),
		message.NewFormattedMessage(route.GetName(), format.NewFormatter(
			format.NewTextColorFormat(interactionconstant.TextColorMagenta),
		)),
		message.NewNewLine(),
		message.NewMessage("    - "),
		message.NewFormattedMessage(route.GetDescription(), format.NewHighlightedTextFormatter()),
		message.NewNewLine(),
	}
}
