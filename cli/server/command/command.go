/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package command holds every command that the CLI server ships.
//
// The other ports declare a command as a class, and read its route from an
// annotation. Go has no annotation, so each command declares its own route, and
// the route provider of the server returns them.
package command

import (
	"os"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	containerconstant "github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// resolve returns the service under the binding key, and the zero value of the
// type where the container resolves nothing under it.
//
// A command reads a service that the CLI component publishes, so a container
// that resolves nothing means the application registered no CLI component. The
// command then reports what it can, rather than ending the process.
func resolve[T any](container containercontract.ContainerContract, id string) T {
	var empty T

	resolved, err := container.Get(id, nil, containerconstant.NewInstanceOrThrowException)
	if err != nil {
		return empty
	}

	typed, isType := resolved.(T)
	if !isType {
		return empty
	}

	return typed
}

// getOutputFactory returns the output factory that the application publishes,
// and a default one where the application publishes none.
func getOutputFactory(container containercontract.ContainerContract) contract.OutputFactoryContract {
	resolved := resolve[contract.OutputFactoryContract](
		container,
		interactionconstant.OutputFactoryContractServiceID,
	)
	if resolved == nil {
		return factory.NewOutputFactory(nil)
	}

	return resolved
}

// getCollection returns the command collection that the application publishes.
func getCollection(container containercontract.ContainerContract) contract.RouteCollectionContract {
	return resolve[contract.RouteCollectionContract](
		container,
		routingconstant.RouteCollectionContractServiceID,
	)
}

// getCliConfig returns the CLI configuration that the application publishes.
func getCliConfig(container containercontract.ContainerContract) applicationcontract.CliConfigContract {
	return resolve[applicationcontract.CliConfigContract](
		container,
		applicationconstant.CliConfigContractServiceID,
	)
}

// getAppName returns the name that the CLI prints for the application.
func getAppName(config applicationcontract.CliConfigContract) string {
	if config == nil {
		return "Valkyrja"
	}

	return config.GetApplicationName()
}

// getAppVersion returns the version of the application.
func getAppVersion(config applicationcontract.CliConfigContract) string {
	if config == nil {
		return ""
	}

	return config.GetVersion()
}

// newHeader builds the header that a command prints above its own output.
func newHeader(
	config applicationcontract.CliConfigContract,
	route contract.RouteContract,
) contract.MessageContract {
	return message.NewHeader(getAppName(config), getAppVersion(config), getProjectRoot(), route)
}

// getProjectRoot returns the directory that the caller ran the command in.
//
// A process that cannot report its own directory returns an empty string, which
// the header prints as an empty line.
func getProjectRoot() string {
	root, _ := os.Getwd()

	return root
}
