/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

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

type ServerParamCollectionContract = ParamCollectionContract

type CookieParamCollectionContract = ParamCollectionContract

type QueryParamCollectionContract = ParamCollectionContract

type ParsedBodyParamCollectionContract = ParamCollectionContract

type ParsedJsonParamCollectionContract = ParamCollectionContract

type AttributeParamCollectionContract = ParamCollectionContract

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
