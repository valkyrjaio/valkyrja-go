/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"maps"
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
)

// EventData is the event component's state. `sindri` generates one of these for
// the whole application, and the collection loads it at boot.
//
// The listener names of one event are a slice rather than a map, because the
// dispatcher runs the listeners in the order that the collection recorded them.
// The other ports rely on an ordered map for this — PHP on an array, Java on a
// `LinkedHashMap` — and a Go map has no order at all.
type EventData struct {
	events    map[string][]string
	listeners map[string]contract.ListenerFactory
}

// NewEventData builds the state from each map. It copies every map and every
// slice, and it accepts nil for a map that carries nothing.
func NewEventData(
	events map[string][]string,
	listeners map[string]contract.ListenerFactory,
) *EventData {
	return &EventData{
		events:    copyEvents(events),
		listeners: copyMap(listeners),
	}
}

// GetEvents returns a copy of each event identifier with the names of its
// listeners.
func (d *EventData) GetEvents() map[string][]string {
	return copyEvents(d.events)
}

// GetListeners returns a copy of each listener factory.
func (d *EventData) GetListeners() map[string]contract.ListenerFactory {
	return copyMap(d.listeners)
}

// copyEvents copies the map and each slice in it, so a later write to either one
// cannot reach the data.
func copyEvents(source map[string][]string) map[string][]string {
	target := make(map[string][]string, len(source))

	for eventID, listenerIDs := range source {
		target[eventID] = slices.Clone(listenerIDs)
	}

	return target
}

// copyMap returns a copy of the map, and an empty map where the map is nil. The
// copy is never nil, so a caller writes to it without a guard.
func copyMap[K comparable, V any](source map[K]V) map[K]V {
	target := make(map[K]V, len(source))

	maps.Copy(target, source)

	return target
}
