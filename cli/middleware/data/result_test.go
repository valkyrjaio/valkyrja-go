/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"testing"

	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/input"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/middleware/data"
	routingdata "github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

func TestAnInputReceivedResultCarriesOneSideAtATime(t *testing.T) {
	t.Parallel()

	continued := data.NewInputReceivedInput(input.NewInput("", "cache:clear"))
	ended := data.NewInputReceivedOutput(output.NewOutput(nil))

	if continued.IsOutput() || continued.GetOutput() != nil {
		t.Error("a result that continues the run must carry no output, but carried one")
	}

	if continued.GetInput().GetCommandName() != "cache:clear" {
		t.Error("a result that continues the run must carry the input, but did not")
	}

	if !ended.IsOutput() || ended.GetOutput() == nil {
		t.Error("a result that ends the run must carry an output, but carried none")
	}

	if ended.GetInput() != nil {
		t.Error("a result that ends the run must carry no input, but carried one")
	}
}

func TestARouteMatchedResultCarriesOneSideAtATime(t *testing.T) {
	t.Parallel()

	route := routingdata.NewRoute("cache:clear", "Clear the cache", func(
		_ containercontract.ContainerContract,
		_ clicontract.RouteContract,
	) clicontract.OutputContract {
		return nil
	})

	continued := data.NewRouteMatchedRoute(route)
	ended := data.NewRouteMatchedOutput(output.NewOutput(nil))

	if continued.IsOutput() || continued.GetOutput() != nil {
		t.Error("a result that continues the run must carry no output, but carried one")
	}

	if continued.GetRoute().GetName() != "cache:clear" {
		t.Error("a result that continues the run must carry the route, but did not")
	}

	if !ended.IsOutput() || ended.GetOutput() == nil {
		t.Error("a result that ends the run must carry an output, but carried none")
	}

	if ended.GetRoute() != nil {
		t.Error("a result that ends the run must carry no route, but carried one")
	}
}
