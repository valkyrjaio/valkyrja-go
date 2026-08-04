/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package router runs the command that an input matches.
package router

import (
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// Router runs the command that an input matches.
//
// The router matches the command by name, fills each parameter of the command
// from what the caller typed, and runs the handler. A middleware handler runs at
// each stage, and it ends the run where it returns an output.
type Router struct {
	container  containercontract.ContainerContract
	collection contract.RouteCollectionContract

	outputFactory contract.OutputFactoryContract

	routeMatchedHandler    contract.RouteMatchedHandlerContract
	routeNotMatchedHandler contract.RouteNotMatchedHandlerContract
	routeDispatchedHandler contract.RouteDispatchedHandlerContract
	throwableCaughtHandler contract.ThrowableCaughtHandlerContract
	processExitingHandler  contract.ProcessExitingHandlerContract
}

// NewRouter builds the router over a container, a collection, an output factory,
// and the middleware handler of each stage.
func NewRouter(
	container containercontract.ContainerContract,
	collection contract.RouteCollectionContract,
	outputFactory contract.OutputFactoryContract,
	routeMatchedHandler contract.RouteMatchedHandlerContract,
	routeNotMatchedHandler contract.RouteNotMatchedHandlerContract,
	routeDispatchedHandler contract.RouteDispatchedHandlerContract,
	throwableCaughtHandler contract.ThrowableCaughtHandlerContract,
	processExitingHandler contract.ProcessExitingHandlerContract,
) *Router {
	return &Router{
		container:              container,
		collection:             collection,
		outputFactory:          outputFactory,
		routeMatchedHandler:    routeMatchedHandler,
		routeNotMatchedHandler: routeNotMatchedHandler,
		routeDispatchedHandler: routeDispatchedHandler,
		throwableCaughtHandler: throwableCaughtHandler,
		processExitingHandler:  processExitingHandler,
	}
}

// Dispatch matches the input to a command and runs it.
//
// An input that names no command of the application reaches the
// route-not-matched middleware with an output that reports the failure.
func (r *Router) Dispatch(input contract.InputContract) contract.OutputContract {
	name := input.GetCommandName()

	if !r.collection.Has(name) {
		return r.routeNotMatchedHandler.RouteNotMatched(input, r.createNotFoundOutput(name))
	}

	return r.DispatchRoute(input, r.collection.Get(name))
}

// DispatchRoute runs the command for the input.
//
// A parameter that holds a value the command does not accept ends the run with
// an output that reports the failure, the way an unmatched command does.
func (r *Router) DispatchRoute(
	input contract.InputContract,
	route contract.RouteContract,
) contract.OutputContract {
	filled, err := r.addParametersToRoute(input, route)
	if err != nil {
		return r.throwableCaughtHandler.ThrowableCaught(input, r.createThrowableOutput(err), err)
	}

	r.routeMatched(filled)

	result := r.routeMatchedHandler.RouteMatched(input, filled)
	if result.IsOutput() {
		return result.GetOutput()
	}

	matched := result.GetRoute()

	r.container.SetSingleton(constant.RouteContractServiceID, matched)

	output := matched.GetHandler()(r.container, matched)

	return r.routeDispatchedHandler.RouteDispatched(input, output, matched)
}

// GetProcessExitingHandler returns the handler that runs before the process
// exits.
//
// The router collects the process-exiting middleware of the matched command, so
// the server reads the handler from the router rather than holding its own.
func (r *Router) GetProcessExitingHandler() contract.ProcessExitingHandlerContract {
	return r.processExitingHandler
}

// createNotFoundOutput builds the output that reports that the application has
// no command under the name.
func (r *Router) createNotFoundOutput(name string) contract.OutputContract {
	text := "Command `" + name + "` was not found."

	return r.outputFactory.CreateOutput(
		interactionconstant.ExitCodeError,
		message.NewBanner(message.NewErrorMessage(text)),
	)
}

// createThrowableOutput builds the output that reports what went wrong.
func (r *Router) createThrowableOutput(throwable error) contract.OutputContract {
	return r.outputFactory.CreateOutput(
		interactionconstant.ExitCodeError,
		message.NewBanner(message.NewErrorMessage(throwable.Error())),
	)
}

// routeMatched files the middleware of the command with each handler, and binds
// the command in the container.
func (r *Router) routeMatched(route contract.RouteContract) {
	r.routeMatchedHandler.Add(route.GetRouteMatchedMiddleware()...)
	r.routeDispatchedHandler.Add(route.GetRouteDispatchedMiddleware()...)
	r.throwableCaughtHandler.Add(route.GetThrowableCaughtMiddleware()...)
	r.processExitingHandler.Add(route.GetProcessExitingMiddleware()...)

	r.container.SetSingleton(constant.RouteContractServiceID, route)
}

// addParametersToRoute fills each parameter of the command from what the caller
// typed.
func (r *Router) addParametersToRoute(
	input contract.InputContract,
	route contract.RouteContract,
) (contract.RouteContract, error) {
	withArguments, err := addArgumentsToRoute(input, route)
	if err != nil {
		return nil, err
	}

	return addOptionsToRoute(input, withArguments)
}

// addArgumentsToRoute fills each argument parameter from what the caller typed.
//
// Warning: an argument parameter in the array value mode takes every argument
// that is left, so it must be the last one that the command declares. A
// parameter that follows it receives nothing.
func addArgumentsToRoute(
	input contract.InputContract,
	route contract.RouteContract,
) (contract.RouteContract, error) {
	typed := input.GetArguments()
	parameters := route.GetArguments()
	filled := make([]contract.ArgumentParameterContract, 0, len(parameters))

	for index, parameter := range parameters {
		given := argumentsForParameter(parameter, typed, index)

		if parameter.GetValueMode() == constant.ArgumentValueModeArray {
			typed = nil
		}

		withArguments := parameter.WithArguments(given...)

		err := withArguments.ValidateValues()
		if err != nil {
			return nil, err
		}

		filled = append(filled, withArguments)
	}

	return route.WithArguments(filled...), nil
}

// argumentsForParameter returns each argument that the parameter takes.
func argumentsForParameter(
	parameter contract.ArgumentParameterContract,
	typed []contract.ArgumentContract,
	index int,
) []contract.ArgumentContract {
	if parameter.GetValueMode() == constant.ArgumentValueModeArray {
		return typed
	}

	if index >= len(typed) {
		return nil
	}

	return []contract.ArgumentContract{typed[index]}
}

// addOptionsToRoute fills each option parameter from what the caller typed.
//
// An option reaches a parameter under its long name, and under each short name
// that the parameter declares.
func addOptionsToRoute(
	input contract.InputContract,
	route contract.RouteContract,
) (contract.RouteContract, error) {
	typed := input.GetOptions()
	parameters := route.GetOptions()
	filled := make([]contract.OptionParameterContract, 0, len(parameters))

	for _, parameter := range parameters {
		withOptions, err := parameter.WithOptions(optionsForParameter(parameter, typed)...)
		if err != nil {
			return nil, err
		}

		err = withOptions.ValidateValues()
		if err != nil {
			return nil, err
		}

		filled = append(filled, withOptions)
	}

	return route.WithOptions(filled...), nil
}

// optionsForParameter returns each option that the parameter takes.
func optionsForParameter(
	parameter contract.OptionParameterContract,
	typed []contract.OptionContract,
) []contract.OptionContract {
	given := make([]contract.OptionContract, 0, len(typed))

	for _, option := range typed {
		if option.GetName() == parameter.GetName() || slices.Contains(parameter.GetShortNames(), option.GetName()) {
			given = append(given, option)
		}
	}

	return given
}
