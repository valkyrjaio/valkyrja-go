/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package collection holds the listener collection.
package collection

import (
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/data"
)

type ListenerCollection struct {
	events    map[string][]string
	listeners map[string]contract.ListenerFactory
	order     []string
}

// NewListenerCollection builds an empty collection.
func NewListenerCollection() *ListenerCollection {
	return &ListenerCollection{
		events:    map[string][]string{},
		listeners: map[string]contract.ListenerFactory{},
		order:     []string{},
	}
}

// GetData returns the collection's state.
func (c *ListenerCollection) GetData() contract.EventDataContract {
	return data.NewEventData(c.events, c.listeners)
}

// SetFromData replaces the collection's state.
func (c *ListenerCollection) SetFromData(eventData contract.EventDataContract) {
	c.events = eventData.GetEvents()
	c.listeners = eventData.GetListeners()
	c.order = c.getOrderFromEvents()
}

// HasListener reports whether the collection holds the listener.
func (c *ListenerCollection) HasListener(listener contract.ListenerContract) bool {
	return c.HasListenerByID(listener.GetName())
}

// HasListenerByID reports whether the collection holds a listener under the
// name.
func (c *ListenerCollection) HasListenerByID(listenerID string) bool {
	_, found := c.listeners[listenerID]

	return found
}

// AddListener files the listener under its own name, for its own event.
func (c *ListenerCollection) AddListener(listener contract.ListenerContract) {
	listenerID := listener.GetName()
	eventID := listener.GetEventID()

	if !slices.Contains(c.events[eventID], listenerID) {
		c.events[eventID] = append(c.events[eventID], listenerID)
	}

	if !slices.Contains(c.order, listenerID) {
		c.order = append(c.order, listenerID)
	}

	c.listeners[listenerID] = func() contract.ListenerContract { return listener }
}

// RemoveListener removes the listener from its own event.
func (c *ListenerCollection) RemoveListener(listener contract.ListenerContract) {
	listenerID := listener.GetName()
	eventID := listener.GetEventID()

	c.events[eventID] = removeValue(c.events[eventID], listenerID)

	delete(c.listeners, listenerID)

	c.order = removeValue(c.order, listenerID)
}

// RemoveListenerByID removes the listener from every event.
func (c *ListenerCollection) RemoveListenerByID(listenerID string) {
	for eventID := range c.events {
		c.events[eventID] = removeValue(c.events[eventID], listenerID)
	}

	delete(c.listeners, listenerID)

	c.order = removeValue(c.order, listenerID)
}

// HasListenersForEvent reports whether the collection holds a listener for the
// event.
func (c *ListenerCollection) HasListenersForEvent(event contract.EventContract) bool {
	return c.HasListenersForEventByID(event.GetEventID())
}

// HasListenersForEventByID reports whether the collection holds a listener for
// the event identifier.
func (c *ListenerCollection) HasListenersForEventByID(eventID string) bool {
	return len(c.events[eventID]) != 0
}

// GetListenersForEvent returns each listener for the event.
func (c *ListenerCollection) GetListenersForEvent(event contract.EventContract) []contract.ListenerContract {
	return c.GetListenersForEventByID(event.GetEventID())
}

// GetListenersForEventByID returns each listener for the event identifier, in
// the order that the collection recorded them.
func (c *ListenerCollection) GetListenersForEventByID(eventID string) []contract.ListenerContract {
	listenerIDs, found := c.events[eventID]
	if !found {
		return []contract.ListenerContract{}
	}

	return c.getListenersByID(listenerIDs)
}

// SetListenersForEvent files each listener for the event.
func (c *ListenerCollection) SetListenersForEvent(
	event contract.EventContract,
	listeners ...contract.ListenerContract,
) {
	c.SetListenersForEventByID(event.GetEventID(), listeners...)
}

// SetListenersForEventByID files each listener for the event identifier.
func (c *ListenerCollection) SetListenersForEventByID(
	eventID string,
	listeners ...contract.ListenerContract,
) {
	for _, listener := range listeners {
		c.AddListener(listener.WithEventID(eventID))
	}
}

// RemoveListenersForEvent removes every listener for the event.
func (c *ListenerCollection) RemoveListenersForEvent(event contract.EventContract) {
	c.RemoveListenersForEventByID(event.GetEventID())
}

// RemoveListenersForEventByID removes every listener for the event identifier.
func (c *ListenerCollection) RemoveListenersForEventByID(eventID string) {
	for _, listener := range c.GetListenersForEventByID(eventID) {
		c.RemoveListener(listener)
	}

	delete(c.events, eventID)
}

// GetListeners returns every listener that the collection holds, in the order
// that the collection recorded them.
func (c *ListenerCollection) GetListeners() []contract.ListenerContract {
	return c.getListenersByID(c.order)
}

// GetEvents returns the identifier of each event that has a listener.
func (c *ListenerCollection) GetEvents() []string {
	eventIDs := make([]string, 0, len(c.events))

	for eventID := range c.events {
		eventIDs = append(eventIDs, eventID)
	}

	slices.Sort(eventIDs)

	return eventIDs
}

// GetEventsWithListeners returns each event identifier with its listeners.
func (c *ListenerCollection) GetEventsWithListeners() map[string][]contract.ListenerContract {
	eventsWithListeners := make(map[string][]contract.ListenerContract, len(c.events))

	for eventID := range c.events {
		eventsWithListeners[eventID] = c.GetListenersForEventByID(eventID)
	}

	return eventsWithListeners
}

// getListenersByID builds each listener that the names identify, and skips a
// name that the collection holds no factory for.
func (c *ListenerCollection) getListenersByID(listenerIDs []string) []contract.ListenerContract {
	listeners := make([]contract.ListenerContract, 0, len(listenerIDs))

	for _, listenerID := range listenerIDs {
		factory, found := c.listeners[listenerID]
		if !found {
			continue
		}

		listeners = append(listeners, factory())
	}

	return listeners
}

// getOrderFromEvents rebuilds the order that GetListeners reports, from the
// state that the collection loaded. The state records the order of one event's
// listeners, so the rebuilt order reads the events by identifier to stay the
// same across two runs.
func (c *ListenerCollection) getOrderFromEvents() []string {
	order := make([]string, 0, len(c.listeners))

	eventIDs := make([]string, 0, len(c.events))
	for eventID := range c.events {
		eventIDs = append(eventIDs, eventID)
	}

	slices.Sort(eventIDs)

	for _, eventID := range eventIDs {
		for _, listenerID := range c.events[eventID] {
			if !slices.Contains(order, listenerID) {
				order = append(order, listenerID)
			}
		}
	}

	return order
}

// removeValue returns the slice without the value, and keeps the order of what
// stays.
func removeValue(values []string, value string) []string {
	index := slices.Index(values, value)
	if index == -1 {
		return values
	}

	return slices.Delete(slices.Clone(values), index, index+1)
}
