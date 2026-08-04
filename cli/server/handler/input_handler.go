/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package handler is the server's entry point for one input.
package handler

import (
	"fmt"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	interactionconstant "github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// InputHandler is the server's entry point for one input.
//
// It runs the input-received middleware, dispatches the router, and turns a
// failure into an output. The output then writes its messages, and the
// process-exiting middleware runs before the process ends.
type InputHandler struct {
	container containercontract.ContainerContract
	router    contract.RouterContract

	inputReceivedHandler   contract.InputReceivedHandlerContract
	throwableCaughtHandler contract.ThrowableCaughtHandlerContract
	processExitingHandler  contract.ProcessExitingHandlerContract

	outputFactory contract.OutputFactoryContract
	exiter        func(code int)
}

// NewInputHandler builds the handler over a container, a router, the middleware
// handler of each stage, and an output factory.
//
// The exiter ends the process. A caller that passes nil gets a handler that ends
// no process, which is what a test and a long-running worker both need.
func NewInputHandler(
	container containercontract.ContainerContract,
	router contract.RouterContract,
	inputReceivedHandler contract.InputReceivedHandlerContract,
	throwableCaughtHandler contract.ThrowableCaughtHandlerContract,
	processExitingHandler contract.ProcessExitingHandlerContract,
	outputFactory contract.OutputFactoryContract,
	exiter func(code int),
) *InputHandler {
	return &InputHandler{
		container:              container,
		router:                 router,
		inputReceivedHandler:   inputReceivedHandler,
		throwableCaughtHandler: throwableCaughtHandler,
		processExitingHandler:  processExitingHandler,
		outputFactory:          outputFactory,
		exiter:                 exiter,
	}
}

// Handle returns the output for the input.
//
// A command that panics reaches the caller as an output that reports the
// failure, rather than ending the process. The other ports catch a throw here;
// Go has no throw, so this recovers a panic instead.
func (h *InputHandler) Handle(input contract.InputContract) contract.OutputContract {
	output := h.dispatch(input)

	h.container.SetSingleton(interactionconstant.OutputContractServiceID, output)

	return output
}

// Exit runs what is left before the process ends.
func (h *InputHandler) Exit(input contract.InputContract, output contract.OutputContract) {
	h.processExitingHandler.ProcessExiting(input, output)
}

// Run handles the input, writes what the command reported, and exits.
func (h *InputHandler) Run(input contract.InputContract) {
	output := h.Handle(input).WriteMessages()

	h.Exit(input, output)

	if h.exiter == nil {
		return
	}

	h.exiter(int(output.GetExitCode()))
}

// dispatch runs the input through the router, and turns a failure into an
// output.
func (h *InputHandler) dispatch(input contract.InputContract) (output contract.OutputContract) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		throwable, isError := recovered.(error)
		if !isError {
			throwable = fmt.Errorf("%v", recovered)
		}

		output = h.throwableCaughtHandler.ThrowableCaught(
			input,
			h.createOutputFromThrowable(input, throwable),
			throwable,
		)
	}()

	h.container.SetSingleton(interactionconstant.InputContractServiceID, input)

	result := h.inputReceivedHandler.InputReceived(input)
	if result.IsOutput() {
		return result.GetOutput()
	}

	received := result.GetInput()

	h.container.SetSingleton(interactionconstant.InputContractServiceID, received)

	return h.router.Dispatch(received)
}

// createOutputFromThrowable builds the output that reports what went wrong.
func (h *InputHandler) createOutputFromThrowable(
	input contract.InputContract,
	throwable error,
) contract.OutputContract {
	return h.outputFactory.CreateOutput(
		interactionconstant.ExitCodeError,
		message.NewBanner(message.NewErrorMessage("Cli Server Error:")),
		message.NewNewLine(),
		message.NewErrorMessage("Command:"),
		message.NewMessage(" "+input.GetCommandName()),
		message.NewNewLine(),
		message.NewNewLine(),
		message.NewErrorMessage("Message:"),
		message.NewMessage(" "+throwable.Error()),
	)
}
