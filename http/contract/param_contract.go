/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

// ParamCollectionContract holds a set of parameters that arrived with a request.
//
// Each named collection below is an alias rather than its own interface. Go
// compares an interface by its method set, so a second interface with the same
// methods is the same type and adds no safety at all — the `iface` linter
// reports it. The TypeScript port declares each one as an alias for the same
// reason.
//
// The TypeScript port makes the value type a type parameter, and narrows it per
// collection — a query parameter is a string or a list of strings, and an
// attribute is anything. A Go interface that carries a type parameter cannot be
// named without instantiating it, so this port erases the value to `any`, which
// is what the Java port does.
type ParamCollectionContract interface {
	// Has reports whether the collection holds a parameter under the key.
	Has(key string) bool

	// Get returns the parameter under the key, and nil where the collection
	// holds none.
	Get(key string) any

	// GetAll returns every parameter, keyed by its own key.
	GetAll() map[string]any

	// GetOnly returns the parameters that the keys identify.
	GetOnly(keys ...string) map[string]any

	// GetAllExcept returns every parameter that the keys do not identify.
	GetAllExcept(keys ...string) map[string]any

	// With returns a copy of the collection with other parameters.
	With(params map[string]any) ParamCollectionContract

	// WithAdded returns a copy of the collection with the parameters added to
	// the ones it holds.
	WithAdded(params map[string]any) ParamCollectionContract
}

// ServerParamCollectionContract holds the server parameters of a request.
type ServerParamCollectionContract = ParamCollectionContract

// CookieParamCollectionContract holds the cookie parameters of a request.
type CookieParamCollectionContract = ParamCollectionContract

// QueryParamCollectionContract holds the query parameters of a request.
type QueryParamCollectionContract = ParamCollectionContract

// ParsedBodyParamCollectionContract holds the parsed body of a request.
type ParsedBodyParamCollectionContract = ParamCollectionContract

// ParsedJsonParamCollectionContract holds the parsed JSON body of a request.
type ParsedJsonParamCollectionContract = ParamCollectionContract

// AttributeParamCollectionContract holds the attributes that the framework puts
// on a request as it passes through.
type AttributeParamCollectionContract = ParamCollectionContract

// UploadedFileContract is one file that arrived with a request.
type UploadedFileContract interface {
	// GetStream returns the contents of the file as a stream.
	GetStream() StreamContract

	// MoveTo moves the file to the target path.
	MoveTo(targetPath string) error

	// HasSize reports whether the client stated the size of the file.
	HasSize() bool

	// GetSize returns the size of the file in bytes.
	GetSize() int

	// GetError returns what went wrong with the upload, and nil where nothing
	// went wrong.
	GetError() error

	// HasClientFilename reports whether the client stated a file name.
	HasClientFilename() bool

	// GetClientFilename returns the file name that the client stated.
	GetClientFilename() string

	// HasClientMediaType reports whether the client stated a media type.
	HasClientMediaType() bool

	// GetClientMediaType returns the media type that the client stated.
	GetClientMediaType() string
}

// UploadedFileCollectionContract holds the files that arrived with one request.
type UploadedFileCollectionContract interface {
	// Has reports whether the collection holds a file under the key.
	Has(key string) bool

	// Get returns the file under the key.
	Get(key string) UploadedFileContract

	// GetAll returns every file, keyed by its own key.
	GetAll() map[string]UploadedFileContract

	// GetOnly returns the files that the keys identify.
	GetOnly(keys ...string) map[string]UploadedFileContract

	// GetAllExcept returns every file that the keys do not identify.
	GetAllExcept(keys ...string) map[string]UploadedFileContract

	// With returns a copy of the collection with other files.
	With(collection map[string]UploadedFileContract) UploadedFileCollectionContract

	// WithAdded returns a copy of the collection with the files added to the
	// ones it holds.
	WithAdded(collection map[string]UploadedFileContract) UploadedFileCollectionContract
}
