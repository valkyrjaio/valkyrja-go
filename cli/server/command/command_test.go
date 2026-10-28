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
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/argument"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	interactionfactory "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	routingdata "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/command"
	serverconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/server/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

// newContainer builds a container that holds every command of the CLI server,
// and the output factory that a command writes through.
func newContainer() containercontract.ContainerContract {
	container := manager.NewContainer(nil)

	built := collection.NewCollection()
	built.Add(
		command.NewListCommand().GetRoute(),
		command.NewListBashCommand().GetRoute(),
		command.NewHelpCommand().GetRoute(),
		command.NewVersionCommand().GetRoute(),
	)

	container.SetSingleton(routingconstant.RouteCollectionContractServiceID, built)
	container.SetSingleton(
		interactionconstant.OutputFactoryContractServiceID,
		interactionfactory.NewOutputFactory(nil),
	)

	return container
}

// textOf returns the text of every message that the output holds.
func textOf(built contract.OutputContract) string {
	text := &strings.Builder{}

	for _, held := range built.GetMessages() {
		text.WriteString(held.GetText())
	}

	return text.String()
}

// withOption returns the route with the option filled, as the router fills it.
func withOption(t *testing.T, route contract.RouteContract, name string, value string) contract.RouteContract {
	t.Helper()

	filled, err := route.GetOption(name).WithOptions(
		argument.NewOption(name, interactionconstant.OptionTypeLong).WithValue(value),
	)
	if err != nil {
		t.Fatalf("the option must take the value, but reported: %v", err)
	}

	return route.WithOptions(filled)
}

func TestEachCommandDeclaresItsOwnRoute(t *testing.T) {
	t.Parallel()

	routes := map[string]contract.RouteContract{
		serverconstant.CommandNameList:     command.NewListCommand().GetRoute(),
		serverconstant.CommandNameListBash: command.NewListBashCommand().GetRoute(),
		serverconstant.CommandNameHelp:     command.NewHelpCommand().GetRoute(),
		serverconstant.CommandNameVersion:  command.NewVersionCommand().GetRoute(),
	}

	for name, route := range routes {
		if route.GetName() != name {
			t.Errorf("the command must be named %q, but was named %q", name, route.GetName())
		}

		if !route.HasHelpText() || route.GetHelpTextMessage().GetText() == "" {
			t.Errorf("the command %q must build its help text, but built none", name)
		}
	}
}

func TestTheListCommandPrintsEveryCommand(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := command.NewListCommand().GetRoute()

	text := textOf(command.NewListCommand().Run(container, route))

	for _, name := range []string{
		serverconstant.CommandNameList,
		serverconstant.CommandNameHelp,
		serverconstant.CommandNameVersion,
	} {
		if !strings.Contains(text, name) {
			t.Errorf("the list must name %q, but printed: %q", name, text)
		}
	}
}

func TestTheListCommandKeepsTheCommandsOfTheNamespace(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := withOption(t, command.NewListCommand().GetRoute(), "namespace", "list")

	text := textOf(command.NewListCommand().Run(container, route))

	if !strings.Contains(text, "[list]") {
		t.Errorf("the list must name the namespace in its heading, but printed: %q", text)
	}

	if strings.Contains(text, serverconstant.CommandNameVersion) {
		t.Errorf("the list must drop a command outside the namespace, but printed: %q", text)
	}
}

func TestTheListCommandReportsANamespaceThatMatchesNoCommand(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := withOption(t, command.NewListCommand().GetRoute(), "namespace", "missing")

	built := command.NewListCommand().Run(container, route)

	if built.GetExitCode() != interactionconstant.ExitCodeError {
		t.Error("a namespace that matches no command must report a failure, but did not")
	}

	if !strings.Contains(textOf(built), "missing") {
		t.Errorf("the output must name the namespace, but printed: %q", textOf(built))
	}
}

func TestTheListCommandReportsAnApplicationWithNoCommand(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	route := command.NewListCommand().GetRoute()

	built := command.NewListCommand().Run(container, route)

	if built.GetExitCode() != interactionconstant.ExitCodeError {
		t.Error("an application with no command must report a failure, but did not")
	}

	if !strings.Contains(textOf(built), "No routes found.") {
		t.Errorf("the output must report that it found no command, but printed: %q", textOf(built))
	}
}

func TestTheVersionCommandPrintsTheVersion(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := command.NewVersionCommand().GetRoute()

	text := textOf(command.NewVersionCommand().Run(container, route))

	if !strings.Contains(text, "Built on Valkyrja") && !strings.Contains(text, "Valkyrja") {
		t.Errorf("the version must name the framework, but printed: %q", text)
	}
}

func TestTheVersionCommandPrintsTheNumberAloneUnderTheShortOption(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := withOption(t, command.NewVersionCommand().GetRoute(), "short", "")

	text := textOf(command.NewVersionCommand().Run(container, route))

	if strings.Contains(text, "Built on Valkyrja") {
		t.Errorf("the short option must print the number alone, but printed: %q", text)
	}
}

func TestTheVersionCommandPrintsNoFrameUnderThePlainOption(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := withOption(t, command.NewVersionCommand().GetRoute(), "plain", "")

	text := textOf(command.NewVersionCommand().Run(container, route))

	if !strings.Contains(text, "Built on Valkyrja") {
		t.Errorf("the plain option must print the framework version, but printed: %q", text)
	}

	if strings.Contains(text, "╭──") {
		t.Errorf("the plain option must print no frame, but printed: %q", text)
	}
}

func TestTheVersionCommandReadsTheConfigurationOfTheApplication(t *testing.T) {
	t.Parallel()

	container := newContainer()
	container.SetSingleton(applicationconstant.CliConfigContractServiceID, nil)

	route := command.NewVersionCommand().GetRoute()

	if textOf(command.NewVersionCommand().Run(container, route)) == "" {
		t.Error("the command must print the version where the application publishes no configuration, but did not")
	}
}

func TestTheListBashCommandPrintsEachNameSeparatedByASpace(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := command.NewListBashCommand().GetRoute()

	text := textOf(command.NewListBashCommand().Run(container, route))

	if len(strings.Fields(text)) != 4 {
		t.Errorf("the command must print each name, but printed: %q", text)
	}
}

func TestTheListBashCommandPrintsTheRestOfTheNameUnderANamespace(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := command.NewListBashCommand().GetRoute()

	filled := route.GetArgument("namespace").WithArguments(argument.NewArgument("list:"))

	text := textOf(command.NewListBashCommand().Run(container, route.WithArguments(
		route.GetArgument("applicationName"),
		filled,
	)))

	if text != "bash" {
		t.Errorf("the command must print the rest of the name, but printed: %q", text)
	}
}

func TestTheListBashCommandStripsEveryNamespaceSegmentTheCallerTyped(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := command.NewListBashCommand().GetRoute()

	// The caller typed two segments, so the command prints what follows both.
	filled := route.GetArgument("namespace").WithArguments(argument.NewArgument("list:ba"))

	text := textOf(command.NewListBashCommand().Run(container, route.WithArguments(
		route.GetArgument("applicationName"),
		filled,
	)))

	if text != "sh" {
		t.Errorf("the command must strip every segment the caller typed, but printed: %q", text)
	}
}

func TestTheHelpCommandPrintsTheHelpTextOfTheCommand(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := withOption(t, command.NewHelpCommand().GetRoute(), "command", serverconstant.CommandNameList)

	text := textOf(command.NewHelpCommand().Run(container, route))

	for _, part := range []string{"Name: ", "Description:", "Usage:", "Options:", "Global Options:", "Help:"} {
		if !strings.Contains(text, part) {
			t.Errorf("the help text must carry %q, but printed: %q", part, text)
		}
	}
}

func TestTheHelpCommandPrintsEachArgumentOfTheCommand(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := withOption(t, command.NewHelpCommand().GetRoute(), "command", serverconstant.CommandNameListBash)

	text := textOf(command.NewHelpCommand().Run(container, route))

	if !strings.Contains(text, "Arguments:") || !strings.Contains(text, "applicationName") {
		t.Errorf("the help text must name each argument, but printed: %q", text)
	}
}

func TestTheHelpCommandReportsACommandThatTheApplicationDoesNotHold(t *testing.T) {
	t.Parallel()

	container := newContainer()
	route := withOption(t, command.NewHelpCommand().GetRoute(), "command", "missing")

	built := command.NewHelpCommand().Run(container, route)

	if built.GetExitCode() != interactionconstant.ExitCodeError {
		t.Error("a command that the application does not hold must report a failure, but did not")
	}
}

func TestTheHelpCommandReportsAnApplicationWithNoCommand(t *testing.T) {
	t.Parallel()

	route := withOption(t, command.NewHelpCommand().GetRoute(), "command", serverconstant.CommandNameList)

	built := command.NewHelpCommand().Run(manager.NewContainer(nil), route)

	if built.GetExitCode() != interactionconstant.ExitCodeError {
		t.Error("an application with no command must report a failure, but did not")
	}
}

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
