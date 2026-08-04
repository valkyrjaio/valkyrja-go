/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package command_test

import (
	"strings"
	"testing"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationdata "github.com/valkyrjaio/valkyrja-go/v26/application/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	interactionfactory "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	routingdata "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/command"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

const detailedCommandName = "build:all"

// newDetailedContainer builds a container that holds one command with every
// parameter shape that the help text prints, and a CLI configuration.
func newDetailedContainer() containercontract.ContainerContract {
	container := manager.NewContainer(nil)

	config := applicationdata.NewCliConfig()
	config.ApplicationName = "Sample"
	config.Version = "9.9.9"

	built := collection.NewCollection()
	built.Add(newDetailedRoute())

	container.SetSingleton(applicationconstant.CliConfigContractServiceID, config)
	container.SetSingleton(routingconstant.RouteCollectionContractServiceID, built)
	container.SetSingleton(
		interactionconstant.OutputFactoryContractServiceID,
		interactionfactory.NewOutputFactory(nil),
	)

	return container
}

// newDetailedRoute builds the command that the help text prints every shape of.
func newDetailedRoute() contract.RouteContract {
	return routingdata.NewRoute(detailedCommandName, "Build everything", func(
		_ containercontract.ContainerContract,
		_ contract.RouteContract,
	) contract.OutputContract {
		return nil
	}).
		WithArguments(
			routingdata.NewArgumentParameter("targets", "Each target").
				WithValueMode(routingconstant.ArgumentValueModeArray),
		).
		WithOptions(
			routingdata.NewOptionParameter("mode", "How to build").
				WithValueDisplayName("mode").
				WithMode(routingconstant.OptionModeRequired).
				WithValidValues("fast", "safe"),
			routingdata.NewOptionParameter("tag", "Each tag").
				WithValueDisplayName("tag").
				WithValueMode(routingconstant.OptionValueModeArray),
		)
}

func TestTheHelpTextPrintsEveryParameterShape(t *testing.T) {
	t.Parallel()

	route := withOption(t, command.NewHelpCommand().GetRoute(), "command", detailedCommandName)

	text := textOf(command.NewHelpCommand().Run(newDetailedContainer(), route))

	for _, part := range []string{
		"[targets...]", // an array argument prints its ellipsis
		"=mode",        // a required option prints its value with no brackets
		"...[=tag]",    // an array option prints its ellipsis and its brackets
		"[fast, safe]", // an option that names its valid values prints them
		"Sample",       // the header names the application
		"v9.9.9",       // the header names the version of the application
	} {
		if !strings.Contains(text, part) {
			t.Errorf("the help text must carry %q, but printed: %q", part, text)
		}
	}
}

func TestTheHelpTextPrintsNoHelpSectionForACommandThatBuildsNone(t *testing.T) {
	t.Parallel()

	container := newDetailedContainer()
	route := withOption(t, command.NewHelpCommand().GetRoute(), "command", detailedCommandName)

	text := textOf(command.NewHelpCommand().Run(container, route))

	if strings.Contains(text, "Help:") {
		t.Errorf("a command that builds no help text must print no help section, but printed: %q", text)
	}
}

func TestACommandReadsNoValueFromAParameterThatItDoesNotTake(t *testing.T) {
	t.Parallel()

	container := newDetailedContainer()

	// The list command reads a namespace option, and the detailed command
	// declares none, so the reader must return an empty value rather than panic.
	text := textOf(command.NewListCommand().Run(container, newDetailedRoute()))

	if !strings.Contains(text, detailedCommandName) {
		t.Errorf("the list must print every command where it reads no namespace, but printed: %q", text)
	}

	bash := textOf(command.NewListBashCommand().Run(container, newDetailedRoute()))

	if bash != detailedCommandName {
		t.Errorf("the bash list must print every name where it reads no namespace, but printed: %q", bash)
	}
}
