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

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/data"
)

const (
	eventID    = "valkyrja.tests.event.EventFixture"
	listenerID = "valkyrja.tests.event.ListenerFixture"

	firstListenerID  = "first"
	secondListenerID = "second"
)

// handlerReturning builds a handler that returns the value it is given.
func handlerReturning(value any) contract.ListenerHandlerFunc {
	return func(_ containercontract.ContainerContract, _ map[string]any) any {
		return value
	}
}

func TestNewListenerHoldsEachValue(t *testing.T) {
	t.Parallel()

	listener := data.NewListener(eventID, listenerID, handlerReturning(firstListenerID))

	if listener.GetEventID() != eventID {
		t.Errorf("GetEventID must be the event identifier, but is: %s", listener.GetEventID())
	}

	if listener.GetName() != listenerID {
		t.Errorf("GetName must be the listener name, but is: %s", listener.GetName())
	}

	if listener.GetHandler()(nil, nil) != firstListenerID {
		t.Error("GetHandler must return the handler that the listener holds, but did not")
	}
}

func TestEachWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	listener := data.NewListener(eventID, listenerID, handlerReturning(firstListenerID))

	withEventID := listener.WithEventID("OtherEventID")
	withName := listener.WithName("OtherName")
	withHandler := listener.WithHandler(handlerReturning(secondListenerID))

	if withEventID.GetEventID() != "OtherEventID" {
		t.Errorf("WithEventID must hold the new event identifier, but holds: %s", withEventID.GetEventID())
	}

	if withName.GetName() != "OtherName" {
		t.Errorf("WithName must hold the new name, but holds: %s", withName.GetName())
	}

	if withHandler.GetHandler()(nil, nil) != secondListenerID {
		t.Error("WithHandler must hold the new handler, but did not")
	}
}

func TestEachWithMethodLeavesTheReceiverUnchanged(t *testing.T) {
	t.Parallel()

	listener := data.NewListener(eventID, listenerID, handlerReturning(firstListenerID))

	listener.WithEventID("OtherEventID")
	listener.WithName("OtherName")
	listener.WithHandler(handlerReturning(secondListenerID))

	if listener.GetEventID() != eventID {
		t.Errorf("WithEventID must leave the receiver unchanged, but it holds: %s", listener.GetEventID())
	}

	if listener.GetName() != listenerID {
		t.Errorf("WithName must leave the receiver unchanged, but it holds: %s", listener.GetName())
	}

	if listener.GetHandler()(nil, nil) != firstListenerID {
		t.Error("WithHandler must leave the receiver unchanged, but it did not")
	}
}

func TestEachWithMethodKeepsTheOtherValues(t *testing.T) {
	t.Parallel()

	listener := data.NewListener(eventID, listenerID, handlerReturning(firstListenerID))

	withName := listener.WithName("OtherName")

	if withName.GetEventID() != eventID {
		t.Errorf("WithName must keep the event identifier, but holds: %s", withName.GetEventID())
	}

	if withName.GetHandler()(nil, nil) != firstListenerID {
		t.Error("WithName must keep the handler, but did not")
	}
}
