/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package collection_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

const routeName = "cache:clear"

// newRoute builds a command that a test files.
func newRoute(name string) contract.RouteContract {
	return data.NewRoute(name, "A command", func(
		_ containercontract.ContainerContract,
		_ contract.RouteContract,
	) contract.OutputContract {
		return nil
	})
}

func TestTheCollectionFilesEachCommandUnderItsName(t *testing.T) {
	t.Parallel()

	built := collection.NewCollection()

	built.Add(newRoute(routeName), newRoute("cache:warm"))

	if !built.Has(routeName) || built.Get(routeName).GetName() != routeName {
		t.Error("the collection must file the command under its own name, but did not")
	}

	if built.Has("missing") || built.Get("missing") != nil {
		t.Error("the collection must report no command under a name that it does not hold, but did")
	}

	if len(built.All()) != 2 {
		t.Errorf("the collection must hold each command, but held: %d", len(built.All()))
	}
}

func TestACommandReplacesTheOneUnderTheSameName(t *testing.T) {
	t.Parallel()

	built := collection.NewCollection()
	replacement := newRoute(routeName).WithDescription("Another description")

	built.Add(newRoute(routeName))
	built.Add(replacement)

	if len(built.All()) != 1 || built.Get(routeName).GetDescription() != "Another description" {
		t.Error("a command under a name that the collection holds must replace it, but did not")
	}
}

func TestTheCollectionReadsAndWritesItsState(t *testing.T) {
	t.Parallel()

	source := collection.NewCollection()
	source.Add(newRoute(routeName))

	target := collection.NewCollection()
	target.SetFromData(source.GetData())

	if !target.Has(routeName) {
		t.Error("the collection must take every command of the state, but did not")
	}

	source.Add(newRoute("cache:warm"))

	if target.Has("cache:warm") {
		t.Error("the state must be a copy, so a later command must not reach the other collection, but it did")
	}
}

func TestTheReturnedMapIsACopy(t *testing.T) {
	t.Parallel()

	built := collection.NewCollection()
	built.Add(newRoute(routeName))

	all := built.All()
	delete(all, routeName)

	if !built.Has(routeName) {
		t.Error("All must return a copy, so a change to it must not reach the collection, but it did")
	}
}
