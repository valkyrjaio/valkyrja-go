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
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	routingdata "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	serverconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/server/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// namespaceSeparator separates a namespace from the rest of a command name.
const namespaceSeparator = ":"

// ListBashCommand prints every command of the application, for bash completion.
type ListBashCommand struct{}

// NewListBashCommand builds the command.
func NewListBashCommand() *ListBashCommand {
	return &ListBashCommand{}
}

// GetRoute returns the command that the router files.
func (c *ListBashCommand) GetRoute() contract.RouteContract {
	return routingdata.NewRoute(
		serverconstant.CommandNameListBash,
		"List all commands for bash completion",
		c.Run,
	).
		WithHelpText(c.GetHelpText).
		WithArguments(
			routingdata.NewArgumentParameter("applicationName", "The application name"),
			routingdata.NewArgumentParameter(namespaceOptionName, "An optional namespace to filter commands by"),
		)
}

// GetHelpText returns the help text of the command.
func (c *ListBashCommand) GetHelpText() contract.MessageContract {
	return message.NewMessage(
		"A command to list all the commands present within the Cli component for bash completion.",
	)
}

// Run prints the name of every command, separated by a space.
//
// A namespace that carries a colon prints the rest of each name, so bash
// completes the part that the caller has not typed yet.
func (c *ListBashCommand) Run(
	container containercontract.ContainerContract,
	route contract.RouteContract,
) contract.OutputContract {
	namespace := getArgumentFirstValue(route, namespaceOptionName)
	routes := sortedRoutes(getCollection(container), namespace)
	names := make([]string, 0, len(routes))

	for _, listed := range routes {
		names = append(names, nameForBash(listed.GetName(), namespace))
	}

	return getOutputFactory(container).CreateOutput(
		interactionconstant.ExitCodeSuccess,
		message.NewMessage(strings.Join(names, " ")),
	)
}

// getArgumentFirstValue returns the first value that the caller gave the
// argument, and an empty string where the caller gave none.
func getArgumentFirstValue(route contract.RouteContract, name string) string {
	argument := route.GetArgument(name)
	if argument == nil {
		return ""
	}

	return argument.GetFirstValue()
}

// nameForBash returns the part of the name that bash completes.
func nameForBash(name string, namespace string) string {
	colonAt := strings.Index(namespace, namespaceSeparator)
	if colonAt < 0 {
		return name
	}

	return name[colonAt+1:]
}
