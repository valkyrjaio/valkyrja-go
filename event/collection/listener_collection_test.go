/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package collection_test

import (
	"slices"
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/data"
	"github.com/valkyrjaio/valkyrja-go/v26/event/fixtures"
)

const (
	otherEventID     = "io.valkyrja.tests.event.OtherEventFixture"
	firstListenerID  = "first"
	secondListenerID = "second"
	thirdListenerID  = "third"
)

// newListener builds a listener under the name, for the event.
func newListener(eventID string, name string) contract.ListenerContract {
	return data.NewListener(eventID, name, func(_ containercontract.ContainerContract, _ map[string]any) any {
		return name
	})
}

// getNames returns the name of each listener, in the order that it received
// them.
func getNames(listeners []contract.ListenerContract) []string {
	names := make([]string, 0, len(listeners))

	for _, listener := range listeners {
		names = append(names, listener.GetName())
	}

	return names
}

func TestAddListenerFilesTheListenerForItsEvent(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listener := newListener(fixtures.EventFixtureID, firstListenerID)

	listenerCollection.AddListener(listener)

	if !listenerCollection.HasListener(listener) {
		t.Error("HasListener must be true after AddListener, but is false")
	}

	if !listenerCollection.HasListenerByID(firstListenerID) {
		t.Error("HasListenerByID must be true after AddListener, but is false")
	}

	if !listenerCollection.HasListenersForEventByID(fixtures.EventFixtureID) {
		t.Error("HasListenersForEventByID must be true after AddListener, but is false")
	}
}

func TestTheCollectionKeepsTheOrderThatItRecorded(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()

	for _, name := range []string{firstListenerID, secondListenerID, thirdListenerID, "fourth", "fifth"} {
		listenerCollection.AddListener(newListener(fixtures.EventFixtureID, name))
	}

	names := getNames(listenerCollection.GetListenersForEventByID(fixtures.EventFixtureID))

	if !slices.Equal(names, []string{firstListenerID, secondListenerID, thirdListenerID, "fourth", "fifth"}) {
		t.Errorf("the listeners must keep the order that the collection recorded, but are: %v", names)
	}
}

func TestAddListenerReplacesAListenerOfTheSameName(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()

	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, secondListenerID))
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))

	names := getNames(listenerCollection.GetListenersForEventByID(fixtures.EventFixtureID))

	if !slices.Equal(names, []string{firstListenerID, secondListenerID}) {
		t.Errorf("a listener of the same name must not be filed twice, but the listeners are: %v", names)
	}
}

func TestHasListenerIsFalseForAnUnknownListener(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()

	if listenerCollection.HasListener(newListener(fixtures.EventFixtureID, firstListenerID)) {
		t.Error("HasListener must be false for an unknown listener, but is true")
	}

	if listenerCollection.HasListenerByID(firstListenerID) {
		t.Error("HasListenerByID must be false for an unknown name, but is true")
	}
}

func TestHasListenersForEventReadsTheEventIdentifier(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	event := &fixtures.EventFixture{}

	if listenerCollection.HasListenersForEvent(event) {
		t.Error("HasListenersForEvent must be false before a listener is filed, but is true")
	}

	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))

	if !listenerCollection.HasListenersForEvent(event) {
		t.Error("HasListenersForEvent must be true after a listener is filed, but is false")
	}
}

func TestGetListenersForEventReadsTheEventIdentifier(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))

	names := getNames(listenerCollection.GetListenersForEvent(&fixtures.EventFixture{}))

	if !slices.Equal(names, []string{firstListenerID}) {
		t.Errorf("GetListenersForEvent must return the event's listeners, but returned: %v", names)
	}
}

func TestGetListenersForEventByIDIsEmptyForAnUnknownEvent(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()

	if len(listenerCollection.GetListenersForEventByID(otherEventID)) != 0 {
		t.Error("GetListenersForEventByID must be empty for an unknown event, but is not")
	}
}

func TestRemoveListenerRemovesItFromItsEvent(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	first := newListener(fixtures.EventFixtureID, firstListenerID)

	listenerCollection.AddListener(first)
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, secondListenerID))
	listenerCollection.RemoveListener(first)

	if listenerCollection.HasListener(first) {
		t.Error("HasListener must be false after RemoveListener, but is true")
	}

	names := getNames(listenerCollection.GetListenersForEventByID(fixtures.EventFixtureID))

	if !slices.Equal(names, []string{secondListenerID}) {
		t.Errorf("RemoveListener must keep the order of what stays, but the listeners are: %v", names)
	}
}

func TestRemoveListenerIgnoresAnUnknownListener(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))

	listenerCollection.RemoveListener(newListener(fixtures.EventFixtureID, secondListenerID))

	if !listenerCollection.HasListenerByID(firstListenerID) {
		t.Error("RemoveListener must leave an unrelated listener in place, but did not")
	}
}

func TestRemoveListenerByIDRemovesItFromEveryEvent(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, "shared"))
	listenerCollection.AddListener(newListener(otherEventID, "shared"))
	listenerCollection.AddListener(newListener(otherEventID, "other"))

	listenerCollection.RemoveListenerByID("shared")

	if listenerCollection.HasListenerByID("shared") {
		t.Error("HasListenerByID must be false after RemoveListenerByID, but is true")
	}

	if listenerCollection.HasListenersForEventByID(fixtures.EventFixtureID) {
		t.Error("the first event must hold no listener after RemoveListenerByID, but holds one")
	}

	names := getNames(listenerCollection.GetListenersForEventByID(otherEventID))

	if !slices.Equal(names, []string{"other"}) {
		t.Errorf("the second event must keep its other listener, but holds: %v", names)
	}
}

func TestSetListenersForEventFilesEachListenerUnderTheEvent(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()

	listenerCollection.SetListenersForEvent(
		&fixtures.EventFixture{},
		newListener(otherEventID, firstListenerID),
		newListener(otherEventID, secondListenerID),
	)

	names := getNames(listenerCollection.GetListenersForEventByID(fixtures.EventFixtureID))

	if !slices.Equal(names, []string{firstListenerID, secondListenerID}) {
		t.Errorf("SetListenersForEvent must refile each listener, but the listeners are: %v", names)
	}

	if listenerCollection.HasListenersForEventByID(otherEventID) {
		t.Error("SetListenersForEvent must not leave the listener on its old event, but did")
	}
}

func TestRemoveListenersForEventRemovesEveryListener(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, secondListenerID))
	listenerCollection.AddListener(newListener(otherEventID, thirdListenerID))

	listenerCollection.RemoveListenersForEvent(&fixtures.EventFixture{})

	if listenerCollection.HasListenersForEventByID(fixtures.EventFixtureID) {
		t.Error("the event must hold no listener after RemoveListenersForEvent, but holds one")
	}

	if listenerCollection.HasListenerByID(firstListenerID) || listenerCollection.HasListenerByID(secondListenerID) {
		t.Error("RemoveListenersForEvent must remove each listener, but did not")
	}

	if !listenerCollection.HasListenerByID(thirdListenerID) {
		t.Error("RemoveListenersForEvent must leave another event's listener in place, but did not")
	}
}

func TestGetListenersReturnsEveryListenerInOrder(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))
	listenerCollection.AddListener(newListener(otherEventID, secondListenerID))
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, thirdListenerID))

	names := getNames(listenerCollection.GetListeners())

	if !slices.Equal(names, []string{firstListenerID, secondListenerID, thirdListenerID}) {
		t.Errorf("GetListeners must keep the order that the collection recorded, but returned: %v", names)
	}
}

func TestGetEventsReturnsEachEventIdentifier(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))
	listenerCollection.AddListener(newListener(otherEventID, secondListenerID))

	eventIDs := listenerCollection.GetEvents()

	if !slices.Equal(eventIDs, []string{fixtures.EventFixtureID, otherEventID}) {
		t.Errorf("GetEvents must return each event identifier, but returned: %v", eventIDs)
	}
}

func TestGetEventsWithListenersPairsEachEventWithItsListeners(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))
	listenerCollection.AddListener(newListener(otherEventID, secondListenerID))

	eventsWithListeners := listenerCollection.GetEventsWithListeners()

	if len(eventsWithListeners) != 2 {
		t.Fatalf("GetEventsWithListeners must hold each event, but holds: %d", len(eventsWithListeners))
	}

	if !slices.Equal(getNames(eventsWithListeners[fixtures.EventFixtureID]), []string{firstListenerID}) {
		t.Errorf("the first event must hold its own listener, but holds: %v",
			getNames(eventsWithListeners[fixtures.EventFixtureID]))
	}
}

func TestGetDataReturnsWhatTheCollectionHolds(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))
	listenerCollection.AddListener(newListener(fixtures.EventFixtureID, secondListenerID))

	eventData := listenerCollection.GetData()

	if !slices.Equal(eventData.GetEvents()[fixtures.EventFixtureID], []string{firstListenerID, secondListenerID}) {
		t.Errorf("GetData must hold the listener names in order, but holds: %v", eventData.GetEvents())
	}

	if len(eventData.GetListeners()) != 2 {
		t.Errorf("GetData must hold each listener factory, but holds: %d", len(eventData.GetListeners()))
	}
}

func TestSetFromDataReplacesWhatTheCollectionHolds(t *testing.T) {
	t.Parallel()

	source := collection.NewListenerCollection()
	source.AddListener(newListener(fixtures.EventFixtureID, firstListenerID))
	source.AddListener(newListener(otherEventID, secondListenerID))

	target := collection.NewListenerCollection()
	target.AddListener(newListener(fixtures.EventFixtureID, "replaced"))

	target.SetFromData(source.GetData())

	if target.HasListenerByID("replaced") {
		t.Error("SetFromData must replace what the collection holds, but the old listener stayed")
	}

	names := getNames(target.GetListeners())

	if !slices.Equal(names, []string{firstListenerID, secondListenerID}) {
		t.Errorf("SetFromData must rebuild the order, but the listeners are: %v", names)
	}
}

func TestSetFromDataSkipsAListenerNameWithNoFactory(t *testing.T) {
	t.Parallel()

	listenerCollection := collection.NewListenerCollection()

	listenerCollection.SetFromData(data.NewEventData(
		map[string][]string{fixtures.EventFixtureID: {firstListenerID, "missing"}},
		map[string]contract.ListenerFactory{
			firstListenerID: func() contract.ListenerContract {
				return newListener(fixtures.EventFixtureID, firstListenerID)
			},
		},
	))

	names := getNames(listenerCollection.GetListenersForEventByID(fixtures.EventFixtureID))

	if !slices.Equal(names, []string{firstListenerID}) {
		t.Errorf("a name with no factory must be skipped, but the listeners are: %v", names)
	}
}
