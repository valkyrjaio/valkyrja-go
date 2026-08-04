/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package fixtures holds the reusable doubles that the CLI routing tests build
// on.
package fixtures

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// RecordingHandlerFixture is a command handler that records the route it ran
// for, and returns an output that names it.
type RecordingHandlerFixture struct {
	Ran *[]string
}

// Run records the run and returns an output that names the command.
func (h *RecordingHandlerFixture) Run(
	_ containercontract.ContainerContract,
	route contract.RouteContract,
) contract.OutputContract {
	*h.Ran = append(*h.Ran, route.GetName())

	return output.NewOutput(nil, message.NewMessage(route.GetName())).
		WithExitCode(constant.ExitCodeSuccess)
}

// PanickingHandlerFixture is a command handler that panics.
type PanickingHandlerFixture struct {
	Recovered any
}

// Run panics with the value that the fixture holds.
func (h *PanickingHandlerFixture) Run(
	_ containercontract.ContainerContract,
	_ contract.RouteContract,
) contract.OutputContract {
	panic(h.Recovered)
}
