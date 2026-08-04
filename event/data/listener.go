/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds the event component's data objects.
package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
)

// Listener is one listener for one event.
//
// Each `With` method returns a copy and leaves the receiver unchanged, so a
// listener that the collection holds cannot change under it.
type Listener struct {
	eventID string
	name    string
	handler contract.ListenerHandlerFunc
}

// NewListener builds a listener for an event, under a unique name.
func NewListener(eventID string, name string, handler contract.ListenerHandlerFunc) *Listener {
	return &Listener{
		eventID: eventID,
		name:    name,
		handler: handler,
	}
}

// GetEventID returns the identifier of the event that the listener listens for.
func (l *Listener) GetEventID() string {
	return l.eventID
}

// WithEventID returns a copy of the listener for another event.
func (l *Listener) WithEventID(eventID string) contract.ListenerContract {
	copied := *l
	copied.eventID = eventID

	return &copied
}

// GetName returns the unique name of the listener.
func (l *Listener) GetName() string {
	return l.name
}

// WithName returns a copy of the listener under another name.
func (l *Listener) WithName(name string) contract.ListenerContract {
	copied := *l
	copied.name = name

	return &copied
}

// GetHandler returns what the listener runs.
func (l *Listener) GetHandler() contract.ListenerHandlerFunc {
	return l.handler
}

// WithHandler returns a copy of the listener that runs another handler.
func (l *Listener) WithHandler(handler contract.ListenerHandlerFunc) contract.ListenerContract {
	copied := *l
	copied.handler = handler

	return &copied
}
