/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package dispatcher_test

import (
	"errors"
	"slices"
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	containerexception "github.com/valkyrjaio/valkyrja-go/v26/container/throwable/exception"
	"github.com/valkyrjaio/valkyrja-go/v26/event/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/event/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/data"
	"github.com/valkyrjaio/valkyrja-go/v26/event/dispatcher"
	"github.com/valkyrjaio/valkyrja-go/v26/event/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/event/throwable/exception"
)

const (
	firstListenerID  = "first"
	secondListenerID = "second"
)

// newRecordingListener builds a listener that returns its own name and appends
// the name to the record.
func newRecordingListener(eventID string, name string, record *[]string) contract.ListenerContract {
	return data.NewListener(eventID, name, func(_ containercontract.ContainerContract, _ map[string]any) any {
		*record = append(*record, name)

		return name
	})
}

func TestDispatchRunsEachListenerInOrder(t *testing.T) {
	t.Parallel()

	record := []string{}
	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newRecordingListener(fixtures.EventFixtureID, firstListenerID, &record))
	listenerCollection.AddListener(newRecordingListener(fixtures.EventFixtureID, secondListenerID, &record))

	eventDispatcher := dispatcher.NewEventDispatcher(listenerCollection, manager.NewContainer(nil))
	event := &fixtures.EventFixture{}

	if eventDispatcher.Dispatch(event) != contract.EventContract(event) {
		t.Error("Dispatch must return the event that it received, but did not")
	}

	if !slices.Equal(record, []string{firstListenerID, secondListenerID}) {
		t.Errorf("Dispatch must run each listener in order, but ran: %v", record)
	}
}

func TestDispatchRunsNothingForAnEventWithNoListener(t *testing.T) {
	t.Parallel()

	eventDispatcher := dispatcher.NewEventDispatcher(
		collection.NewListenerCollection(),
		manager.NewContainer(nil),
	)
	event := &fixtures.EventFixture{}

	if eventDispatcher.Dispatch(event) != contract.EventContract(event) {
		t.Error("Dispatch must return the event that it received, but did not")
	}
}

func TestDispatchIfHasListenersRunsOnlyWhereTheEventHasOne(t *testing.T) {
	t.Parallel()

	record := []string{}
	listenerCollection := collection.NewListenerCollection()
	eventDispatcher := dispatcher.NewEventDispatcher(listenerCollection, manager.NewContainer(nil))

	eventDispatcher.DispatchIfHasListeners(&fixtures.EventFixture{})

	if len(record) != 0 {
		t.Errorf("DispatchIfHasListeners must run nothing where the event has no listener, but ran: %v", record)
	}

	listenerCollection.AddListener(newRecordingListener(fixtures.EventFixtureID, firstListenerID, &record))

	eventDispatcher.DispatchIfHasListeners(&fixtures.EventFixture{})

	if !slices.Equal(record, []string{firstListenerID}) {
		t.Errorf("DispatchIfHasListeners must run the listener, but ran: %v", record)
	}
}

func TestDispatchListenersStopsWhereAStoppableEventStops(t *testing.T) {
	t.Parallel()

	record := []string{}
	stopping := data.NewListener(
		fixtures.StoppableEventFixtureID,
		"stopping",
		func(_ containercontract.ContainerContract, arguments map[string]any) any {
			record = append(record, "stopping")

			event, isStoppable := arguments[constant.EventArgumentKey].(*fixtures.StoppableEventFixture)
			if isStoppable {
				event.Stopped = true
			}

			return nil
		},
	)

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newRecordingListener(fixtures.StoppableEventFixtureID, "before", &record))
	listenerCollection.AddListener(stopping)
	listenerCollection.AddListener(newRecordingListener(fixtures.StoppableEventFixtureID, "after", &record))

	eventDispatcher := dispatcher.NewEventDispatcher(listenerCollection, manager.NewContainer(nil))

	eventDispatcher.Dispatch(&fixtures.StoppableEventFixture{})

	if !slices.Equal(record, []string{"before", "stopping"}) {
		t.Errorf("Dispatch must stop after the listener that stopped the propagation, but ran: %v", record)
	}
}

func TestDispatchRunsEveryListenerOfAStoppableEventThatDoesNotStop(t *testing.T) {
	t.Parallel()

	record := []string{}
	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newRecordingListener(fixtures.StoppableEventFixtureID, firstListenerID, &record))
	listenerCollection.AddListener(newRecordingListener(fixtures.StoppableEventFixtureID, secondListenerID, &record))

	eventDispatcher := dispatcher.NewEventDispatcher(listenerCollection, manager.NewContainer(nil))

	eventDispatcher.Dispatch(&fixtures.StoppableEventFixture{})

	if !slices.Equal(record, []string{firstListenerID, secondListenerID}) {
		t.Errorf("Dispatch must run each listener where nothing stops it, but ran: %v", record)
	}
}

func TestDispatchListenerCollectsWhatEachListenerReturned(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(
		data.NewListener(fixtures.DispatchCollectableEventFixtureID, firstListenerID,
			func(_ containercontract.ContainerContract, _ map[string]any) any { return firstListenerID }),
	)
	listenerCollection.AddListener(
		data.NewListener(fixtures.DispatchCollectableEventFixtureID, secondListenerID,
			func(_ containercontract.ContainerContract, _ map[string]any) any { return secondListenerID }),
	)

	eventDispatcher := dispatcher.NewEventDispatcher(listenerCollection, manager.NewContainer(nil))
	event := &fixtures.DispatchCollectableEventFixture{}

	eventDispatcher.Dispatch(event)

	if !slices.Equal(event.GetDispatches(), []any{firstListenerID, secondListenerID}) {
		t.Errorf("the event must collect what each listener returned, but collected: %v", event.GetDispatches())
	}
}

func TestDispatchListenerRunsNothingForAListenerWithNoHandler(t *testing.T) {
	t.Parallel()

	eventDispatcher := dispatcher.NewEventDispatcher(
		collection.NewListenerCollection(),
		manager.NewContainer(nil),
	)
	event := &fixtures.DispatchCollectableEventFixture{}

	returned := eventDispatcher.DispatchListener(
		event,
		data.NewListener(fixtures.DispatchCollectableEventFixtureID, firstListenerID, nil),
	)

	if returned != contract.EventContract(event) {
		t.Error("DispatchListener must return the event that it received, but did not")
	}

	if len(event.GetDispatches()) != 0 {
		t.Errorf("a listener with no handler must collect nothing, but collected: %v", event.GetDispatches())
	}
}

func TestDispatchListenerPassesTheEventToTheHandler(t *testing.T) {
	t.Parallel()

	var received any

	eventDispatcher := dispatcher.NewEventDispatcher(
		collection.NewListenerCollection(),
		manager.NewContainer(nil),
	)
	event := &fixtures.EventFixture{}

	eventDispatcher.DispatchListener(event, data.NewListener(
		fixtures.EventFixtureID,
		firstListenerID,
		func(_ containercontract.ContainerContract, arguments map[string]any) any {
			received = arguments[constant.EventArgumentKey]

			return nil
		},
	))

	if received != any(event) {
		t.Error("the handler must receive the event under the event argument key, but did not")
	}
}

func TestDispatchByIDBuildsTheEventFromTheContainer(t *testing.T) {
	t.Parallel()

	record := []string{}
	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newRecordingListener(fixtures.EventFixtureID, firstListenerID, &record))

	container := manager.NewContainer(nil)
	container.Bind(fixtures.EventFixtureID, func(_ containercontract.ContainerContract, _ []any) any {
		return &fixtures.EventFixture{}
	})

	eventDispatcher := dispatcher.NewEventDispatcher(listenerCollection, container)

	event, err := eventDispatcher.DispatchByID(fixtures.EventFixtureID, nil)
	if err != nil {
		t.Fatalf("DispatchByID must build and dispatch the event, but reported: %v", err)
	}

	if event.GetEventID() != fixtures.EventFixtureID {
		t.Errorf("DispatchByID must return the event that it built, but returned: %s", event.GetEventID())
	}

	if !slices.Equal(record, []string{firstListenerID}) {
		t.Errorf("DispatchByID must run the listener, but ran: %v", record)
	}
}

func TestDispatchByIDFillsAnArgumentsCapableEvent(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(
		fixtures.ArgumentsCapableEventFixtureID,
		func(_ containercontract.ContainerContract, _ []any) any {
			return &fixtures.ArgumentsCapableEventFixture{}
		},
	)

	eventDispatcher := dispatcher.NewEventDispatcher(collection.NewListenerCollection(), container)

	event, err := eventDispatcher.DispatchByID(fixtures.ArgumentsCapableEventFixtureID, []any{firstListenerID})
	if err != nil {
		t.Fatalf("DispatchByID must build and dispatch the event, but reported: %v", err)
	}

	capable, isCapable := event.(*fixtures.ArgumentsCapableEventFixture)
	if !isCapable {
		t.Fatalf("DispatchByID must return the event, but returned: %T", event)
	}

	if !slices.Equal(capable.Arguments, []any{firstListenerID}) {
		t.Errorf("the event must hold the arguments, but holds: %v", capable.Arguments)
	}
}

func TestDispatchByIDCarriesTheContainerFailure(t *testing.T) {
	t.Parallel()

	eventDispatcher := dispatcher.NewEventDispatcher(
		collection.NewListenerCollection(),
		manager.NewContainer(nil),
	)

	_, err := eventDispatcher.DispatchByID(fixtures.EventFixtureID, nil)

	if _, found := errors.AsType[*containerexception.ContainerInvalidReferenceError](err); !found {
		t.Errorf("DispatchByID must carry the container's failure, but reported: %v", err)
	}
}

func TestDispatchByIDReportsABindingKeyThatIsNoEvent(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind("NotAnEventID", func(_ containercontract.ContainerContract, _ []any) any {
		return &fixtures.NotAnEventFixture{}
	})

	eventDispatcher := dispatcher.NewEventDispatcher(collection.NewListenerCollection(), container)

	_, err := eventDispatcher.DispatchByID("NotAnEventID", nil)

	target, found := errors.AsType[*exception.EventInvalidEventError](err)
	if !found {
		t.Fatalf("DispatchByID must report an invalid event, but reported: %v", err)
	}

	if target.GetID() != "NotAnEventID" {
		t.Errorf("the error must name the binding key, but names: %s", target.GetID())
	}
}

func TestDispatchByIDIfHasListenersRunsOnlyWhereTheEventHasOne(t *testing.T) {
	t.Parallel()

	record := []string{}
	listenerCollection := collection.NewListenerCollection()

	container := manager.NewContainer(nil)
	container.Bind(fixtures.EventFixtureID, func(_ containercontract.ContainerContract, _ []any) any {
		return &fixtures.EventFixture{}
	})

	eventDispatcher := dispatcher.NewEventDispatcher(listenerCollection, container)

	event, err := eventDispatcher.DispatchByIDIfHasListeners(fixtures.EventFixtureID, nil)
	if err != nil {
		t.Fatalf("DispatchByIDIfHasListeners must build the event, but reported: %v", err)
	}

	if event.GetEventID() != fixtures.EventFixtureID {
		t.Errorf("DispatchByIDIfHasListeners must return the event, but returned: %s", event.GetEventID())
	}

	if len(record) != 0 {
		t.Errorf("it must run nothing where the event has no listener, but ran: %v", record)
	}

	listenerCollection.AddListener(newRecordingListener(fixtures.EventFixtureID, firstListenerID, &record))

	_, err = eventDispatcher.DispatchByIDIfHasListeners(fixtures.EventFixtureID, nil)
	if err != nil {
		t.Fatalf("DispatchByIDIfHasListeners must dispatch the event, but reported: %v", err)
	}

	if !slices.Equal(record, []string{firstListenerID}) {
		t.Errorf("it must run the listener, but ran: %v", record)
	}
}

func TestDispatchByIDIfHasListenersCarriesTheContainerFailure(t *testing.T) {
	t.Parallel()

	eventDispatcher := dispatcher.NewEventDispatcher(
		collection.NewListenerCollection(),
		manager.NewContainer(nil),
	)

	_, err := eventDispatcher.DispatchByIDIfHasListeners(fixtures.EventFixtureID, nil)

	if _, found := errors.AsType[*containerexception.ContainerInvalidReferenceError](err); !found {
		t.Errorf("DispatchByIDIfHasListeners must carry the container's failure, but reported: %v", err)
	}
}
