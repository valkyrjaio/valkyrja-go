/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the event component's binding keys.
package constant

// The event component's binding keys.
const (
	// EventDataServiceID is the binding key for the event component's data.
	EventDataServiceID = "valkyrja.event.data.EventData"

	// ListenerCollectionContractServiceID is the binding key for the listener
	// collection.
	ListenerCollectionContractServiceID = "valkyrja.event.collection.ListenerCollectionContract"

	// EventDispatcherContractServiceID is the binding key for the event
	// dispatcher.
	EventDispatcherContractServiceID = "valkyrja.event.dispatcher.EventDispatcherContract"
)

// EventArgumentKey is the key that the dispatcher files the event under, in the
// arguments that it passes to a listener handler.
const EventArgumentKey = "event"
