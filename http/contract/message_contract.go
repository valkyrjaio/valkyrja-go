/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds every contract of the HTTP component.
//
// The component keeps one `contract` package, for the reason that the container
// component keeps one: Go resolves an import cycle at the package level. HTTP is
// the case that forces it — `RouteContract` names each middleware contract, and
// each middleware contract names `RouteContract` back.
//
// Each `With` method returns a copy and leaves the receiver unchanged. The other
// ports return `this` or `static`; Go has no such return type, so each one
// returns the contract.
package contract

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

type MessageContract interface {
	// GetProtocolVersion returns the HTTP protocol version of the message.
	GetProtocolVersion() constant.ProtocolVersion

	// WithProtocolVersion returns a copy of the message for another protocol
	// version.
	WithProtocolVersion(version constant.ProtocolVersion) MessageContract

	// GetHeaders returns the headers of the message.
	GetHeaders() HeaderCollectionContract

	// WithHeaders returns a copy of the message with other headers.
	WithHeaders(headers HeaderCollectionContract) MessageContract

	// GetBody returns the body of the message.
	GetBody() StreamContract

	// WithBody returns a copy of the message with another body.
	WithBody(body StreamContract) MessageContract
}

//nolint:interfacebloat // Parity with the PHP reference implementation.
type StreamContract interface {
	// String returns the whole stream as a string.
	String() string

	// Close closes the stream.
	Close() error

	// Detach removes the underlying stream and returns it.
	Detach() []byte

	// GetSize returns the size of the stream in bytes.
	GetSize() int

	// Tell returns the position of the read and write pointer.
	Tell() (int, error)

	// IsEof reports whether the pointer is at the end of the stream.
	IsEof() bool

	// IsSeekable reports whether a caller seeks in the stream.
	IsSeekable() bool

	// Seek moves the pointer to the offset, from the point that whence names.
	Seek(offset int, whence int) error

	// Rewind moves the pointer to the start of the stream.
	Rewind() error

	// IsWritable reports whether a writer writes to the stream.
	IsWritable() bool

	// Write writes the string and returns the number of bytes that it wrote.
	Write(value string) (int, error)

	// IsReadable reports whether a reader reads the stream.
	IsReadable() bool

	// Read reads the number of bytes and returns them.
	Read(length int) (string, error)

	// GetContents returns what is left of the stream, from the pointer on.
	GetContents() (string, error)

	// GetMetadata returns every metadata item of the stream.
	GetMetadata() map[string]any

	// GetMetadataItem returns one metadata item, and nil where the stream holds
	// no item under the key.
	GetMetadataItem(key string) any
}

//nolint:interfacebloat // Parity with the PHP reference implementation.
type UriContract interface {
	// GetScheme returns the scheme of the URI.
	GetScheme() constant.Scheme

	// IsSecure reports whether the scheme is `https`.
	IsSecure() bool

	// GetAuthority returns the authority of the URI.
	GetAuthority() string

	// GetUsername returns the user name of the URI.
	GetUsername() string

	// GetPassword returns the password of the URI.
	GetPassword() string

	// GetUserInfo returns the user name and the password, separated by a colon.
	GetUserInfo() string

	// GetHost returns the host of the URI.
	GetHost() string

	// HasPort reports whether the URI names a port.
	HasPort() bool

	// GetPort returns the port of the URI.
	GetPort() int

	// GetHostPort returns the host and the port, separated by a colon.
	GetHostPort() string

	// GetSchemeHostPort returns the scheme, the host, and the port.
	GetSchemeHostPort() string

	// GetPath returns the path of the URI.
	GetPath() string

	// GetQuery returns the query string of the URI.
	GetQuery() string

	// GetFragment returns the fragment of the URI.
	GetFragment() string

	// WithScheme returns a copy of the URI for another scheme.
	WithScheme(scheme constant.Scheme) UriContract

	// WithUsername returns a copy of the URI for another user name.
	WithUsername(username string) UriContract

	// WithPassword returns a copy of the URI for another password.
	WithPassword(password string) UriContract

	// WithUserInfo returns a copy of the URI for another user name and
	// password.
	WithUserInfo(user string, password string) UriContract

	// WithHost returns a copy of the URI for another host.
	WithHost(host string) UriContract

	// WithPort returns a copy of the URI for another port.
	WithPort(port int) UriContract

	// WithPath returns a copy of the URI for another path.
	WithPath(path string) UriContract

	// WithQuery returns a copy of the URI for another query string.
	WithQuery(query string) UriContract

	// WithFragment returns a copy of the URI for another fragment.
	WithFragment(fragment string) UriContract

	// String returns the whole URI as a string.
	String() string
}
