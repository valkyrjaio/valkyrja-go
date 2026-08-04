/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package command

import (
	"runtime"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	routingdata "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	serverconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/server/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// shortOptionName names the option that prints the version number alone.
const shortOptionName = "short"

// plainOptionName names the option that prints the version with no frame.
const plainOptionName = "plain"

// VersionCommand prints the version of the application.
type VersionCommand struct{}

// NewVersionCommand builds the command.
func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

// GetRoute returns the command that the router files.
func (c *VersionCommand) GetRoute() contract.RouteContract {
	return routingdata.NewRoute(
		serverconstant.CommandNameVersion,
		"Get the application version",
		c.Run,
	).
		WithHelpText(c.GetHelpText).
		WithOptions(
			routingdata.NewOptionParameter(shortOptionName, "Output the version number only").
				WithShortNames("s"),
			routingdata.NewOptionParameter(plainOptionName, "Output version info without the banner").
				WithShortNames("p"),
		)
}

// GetHelpText returns the help text of the command.
func (c *VersionCommand) GetHelpText() contract.MessageContract {
	return message.NewMessage("A command to show the application version and info.")
}

// Run prints the version of the application.
//
// The short option prints the version number alone, so a script reads it. The
// plain option prints the same lines as the frame does, with no frame.
func (c *VersionCommand) Run(
	container containercontract.ContainerContract,
	route contract.RouteContract,
) contract.OutputContract {
	outputFactory := getOutputFactory(container)
	config := getCliConfig(container)
	version := getAppVersion(config)

	if hasOptionValue(route, shortOptionName) {
		return outputFactory.CreateOutput(interactionconstant.ExitCodeSuccess, message.NewMessage(version))
	}

	if hasOptionValue(route, plainOptionName) {
		return outputFactory.CreateOutput(
			interactionconstant.ExitCodeSuccess,
			message.NewMessage(getAppName(config)+" v"+version),
			message.NewNewLine(),
			message.NewMessage(
				"Built on Valkyrja v"+applicationconstant.Version+
					" (date: "+applicationconstant.VersionBuildDateTime+")",
			),
			message.NewNewLine(),
			message.NewMessage("Running on "+runtime.Version()),
		)
	}

	return outputFactory.CreateOutput(interactionconstant.ExitCodeSuccess, newHeader(config, route))
}

// hasOptionValue reports whether the caller typed the option.
//
// The route holds the parameter whether the caller typed it or not, so the
// parameter reports whether a value reached it.
func hasOptionValue(route contract.RouteContract, name string) bool {
	option := route.GetOption(name)

	return option != nil && len(option.GetOptions()) > 0
}
