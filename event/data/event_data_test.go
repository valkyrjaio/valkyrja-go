/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"slices"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/data"
)

// newListenerFactory builds a factory for a listener under the name.
func newListenerFactory(name string) contract.ListenerFactory {
	return func() contract.ListenerContract {
		return data.NewListener(eventID, name, handlerReturning(name))
	}
}

func TestNewEventDataHoldsEachMap(t *testing.T) {
	t.Parallel()

	eventData := data.NewEventData(
		map[string][]string{eventID: {firstListenerID, secondListenerID}},
		map[string]contract.ListenerFactory{firstListenerID: newListenerFactory(firstListenerID)},
	)

	if !slices.Equal(eventData.GetEvents()[eventID], []string{firstListenerID, secondListenerID}) {
		t.Errorf("GetEvents must hold the listener names in order, but holds: %v", eventData.GetEvents())
	}

	if _, found := eventData.GetListeners()[firstListenerID]; !found {
		t.Error("GetListeners must hold the factory, but does not")
	}
}

func TestNewEventDataAcceptsNilForEachMap(t *testing.T) {
	t.Parallel()

	eventData := data.NewEventData(nil, nil)

	if len(eventData.GetEvents()) != 0 {
		t.Errorf("GetEvents must be empty, but holds: %v", eventData.GetEvents())
	}

	if len(eventData.GetListeners()) != 0 {
		t.Errorf("GetListeners must be empty, but holds: %v", eventData.GetListeners())
	}
}

func TestNewEventDataCopiesTheSliceOfEachEvent(t *testing.T) {
	t.Parallel()

	listenerIDs := []string{firstListenerID}
	eventData := data.NewEventData(map[string][]string{eventID: listenerIDs}, nil)

	listenerIDs[0] = "changed"

	if eventData.GetEvents()[eventID][0] != firstListenerID {
		t.Error("GetEvents must not follow a later write to the source slice, but did")
	}
}

func TestGetEventsReturnsACopyOfTheSlice(t *testing.T) {
	t.Parallel()

	eventData := data.NewEventData(map[string][]string{eventID: {firstListenerID}}, nil)

	eventData.GetEvents()[eventID][0] = "changed"
	delete(eventData.GetEvents(), eventID)

	if eventData.GetEvents()[eventID][0] != firstListenerID {
		t.Error("GetEvents must return a copy, but the write reached the data")
	}
}

func TestGetListenersReturnsACopy(t *testing.T) {
	t.Parallel()

	eventData := data.NewEventData(
		nil,
		map[string]contract.ListenerFactory{firstListenerID: newListenerFactory(firstListenerID)},
	)

	delete(eventData.GetListeners(), firstListenerID)

	if len(eventData.GetListeners()) != 1 {
		t.Error("GetListeners must return a copy, but the delete reached the data")
	}
}
