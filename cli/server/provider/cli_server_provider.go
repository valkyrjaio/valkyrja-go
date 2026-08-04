/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the CLI server sub-component's providers.
package provider

import (
	"os"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	interactionfactory "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/factory"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/handler"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/command"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/server/handler"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// CliServerCliRoutesProvider registers every command that the CLI server ships.
//
// The other ports read a command from an annotation on its class. Go has no
// annotation, so each command declares its own route, and this provider returns
// them as a literal slice that `sindri` reads from the source.
type CliServerCliRoutesProvider struct{}

// GetRoutes returns each command that the CLI server registers.
func (p *CliServerCliRoutesProvider) GetRoutes() []clicontract.RouteContract {
	return []clicontract.RouteContract{
		command.NewListCommand().GetRoute(),
		command.NewListBashCommand().GetRoute(),
		command.NewHelpCommand().GetRoute(),
		command.NewVersionCommand().GetRoute(),
	}
}

// CliServerServiceProvider publishes the bindings of the CLI server
// sub-component.
type CliServerServiceProvider struct{}

// Publishers returns a publisher for each binding key that the sub-component
// defers.
func (p *CliServerServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.InputHandlerContractServiceID: PublishInputHandler,
	}
}

// PublishInputHandler binds the server's entry point for one input.
func PublishInputHandler(container containercontract.ContainerContract) {
	config := getCliConfig(container)

	container.SetSingleton(constant.InputHandlerContractServiceID, handler.NewInputHandler(
		container,
		getRouter(container),
		middlewarehandler.NewInputReceivedHandler(container, getInputReceivedMiddleware(config)...),
		middlewarehandler.NewThrowableCaughtHandler(container, getThrowableCaughtMiddleware(config)...),
		middlewarehandler.NewProcessExitingHandler(container, getProcessExitingMiddleware(config)...),
		getOutputFactory(container),
		os.Exit,
	))
}

// getCliConfig returns the CLI configuration that the application publishes.
func getCliConfig(container containercontract.ContainerContract) applicationcontract.CliConfigContract {
	resolved, err := container.GetSingleton(applicationconstant.CliConfigContractServiceID)
	if err != nil {
		return nil
	}

	config, isConfig := resolved.(applicationcontract.CliConfigContract)
	if !isConfig {
		return nil
	}

	return config
}

// getRouter returns the router that the routing sub-component published.
func getRouter(container containercontract.ContainerContract) clicontract.RouterContract {
	resolved, err := container.GetSingleton(routingconstant.RouterContractServiceID)
	if err != nil {
		return nil
	}

	built, isRouter := resolved.(clicontract.RouterContract)
	if !isRouter {
		return nil
	}

	return built
}

// getOutputFactory returns the factory that the interaction sub-component
// published.
func getOutputFactory(container containercontract.ContainerContract) clicontract.OutputFactoryContract {
	resolved, err := container.GetSingleton(interactionconstant.OutputFactoryContractServiceID)
	if err != nil {
		return interactionfactory.NewOutputFactory(nil)
	}

	built, isFactory := resolved.(clicontract.OutputFactoryContract)
	if !isFactory {
		return interactionfactory.NewOutputFactory(nil)
	}

	return built
}

// getInputReceivedMiddleware returns the input-received middleware that the
// application names.
func getInputReceivedMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetInputReceivedMiddleware()
}

// getThrowableCaughtMiddleware returns the throwable-caught middleware that the
// application names.
func getThrowableCaughtMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetThrowableCaughtMiddleware()
}

// getProcessExitingMiddleware returns the process-exiting middleware that the
// application names.
func getProcessExitingMiddleware(config applicationcontract.CliConfigContract) []string {
	if config == nil {
		return nil
	}

	return config.GetProcessExitingMiddleware()
}
