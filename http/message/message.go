/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package message holds what an HTTP request and an HTTP response have in
// common.
package message

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
)

type Message struct {
	protocolVersion constant.ProtocolVersion
	headers         contract.HeaderCollectionContract
	body            contract.StreamContract
}

// NewMessage builds the shared state. It takes the defaults of the other ports
// where an argument is nil: HTTP 1.1, an empty header collection, and an empty
// body.
func NewMessage(
	protocolVersion constant.ProtocolVersion,
	headers contract.HeaderCollectionContract,
	body contract.StreamContract,
) Message {
	if protocolVersion == "" {
		protocolVersion = constant.ProtocolVersionV11
	}

	if headers == nil {
		headers = header.NewHeaderCollection()
	}

	if body == nil {
		body = stream.NewStream("", constant.ModeReadWrite)
	}

	return Message{
		protocolVersion: protocolVersion,
		headers:         headers,
		body:            body,
	}
}

// GetProtocolVersion returns the HTTP protocol version of the message.
func (m *Message) GetProtocolVersion() constant.ProtocolVersion {
	return m.protocolVersion
}

// GetHeaders returns the headers of the message.
func (m *Message) GetHeaders() contract.HeaderCollectionContract {
	return m.headers
}

// GetBody returns the body of the message.
func (m *Message) GetBody() contract.StreamContract {
	return m.body
}

// SetProtocolVersion replaces the protocol version. A message type calls it on
// its own copy, inside its own `WithProtocolVersion`.
func (m *Message) SetProtocolVersion(version constant.ProtocolVersion) {
	m.protocolVersion = version
}

// SetHeaders replaces the headers. A message type calls it on its own copy,
// inside its own `WithHeaders`.
func (m *Message) SetHeaders(headers contract.HeaderCollectionContract) {
	m.headers = headers
}

// SetBody replaces the body and rewinds it, so a reader reads the body whole.
func (m *Message) SetBody(body contract.StreamContract) {
	m.body = body

	// A body that no caller seeks in stays where it is, which is what the other
	// ports do: each one rewinds and ignores the result.
	_ = body.Rewind()
}

// InjectHeader returns the headers with the header in them.
func InjectHeader(
	injected contract.HeaderContract,
	headers contract.HeaderCollectionContract,
	override bool,
) contract.HeaderCollectionContract {
	name := injected.GetNormalizedName()

	if override || !headers.Has(name) {
		return headers.WithHeader(injected)
	}

	// The guard above reports that the collection holds the header, and `Has`
	// and `Get` normalize the name the same way, so the read cannot fail.
	existing, _ := headers.Get(name)

	return headers.WithHeader(existing.WithAddedValues(injected.GetValues()...))
}
