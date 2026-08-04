/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package file holds the files that arrived with a request.
package file

import (
	"maps"
	"os"
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

// filePermissions is the mode that a moved file takes. It gives the owner read
// and write, and every other reader read.
const filePermissions = 0o644

// UploadedFile is one file that arrived with a request.
type UploadedFile struct {
	body            contract.StreamContract
	size            int
	uploadError     constant.UploadError
	clientFilename  string
	clientMediaType string
	moved           bool
}

// NewUploadedFile builds an uploaded file over its stream.
//
// A size below zero states that the client named none, which is what `HasSize`
// reports.
func NewUploadedFile(
	body contract.StreamContract,
	size int,
	uploadError constant.UploadError,
	clientFilename string,
	clientMediaType string,
) *UploadedFile {
	return &UploadedFile{
		body:            body,
		size:            size,
		uploadError:     uploadError,
		clientFilename:  clientFilename,
		clientMediaType: clientMediaType,
	}
}

// GetStream returns the contents of the file as a stream.
func (f *UploadedFile) GetStream() contract.StreamContract {
	return f.body
}

// MoveTo writes the file to the target path.
//
// It reports a failure where the upload went wrong, where the file moved
// already, or where the write fails.
func (f *UploadedFile) MoveTo(targetPath string) error {
	if f.uploadError != constant.UploadErrorOk {
		return exception.NewHttpUploadedFileUploadError(f.uploadError)
	}

	if f.moved {
		return exception.NewHttpUploadedFileAlreadyMovedError()
	}

	if targetPath == "" {
		return exception.NewHttpUploadedFileInvalidDirectoryError(targetPath)
	}

	contents, err := f.body.GetContents()
	if err != nil {
		return err
	}

	err = os.WriteFile(targetPath, []byte(contents), filePermissions)
	if err != nil {
		return exception.NewHttpUploadedFileMoveFailureError(targetPath)
	}

	f.moved = true

	return nil
}

// HasSize reports whether the client stated the size of the file.
func (f *UploadedFile) HasSize() bool {
	return f.size >= 0
}

// GetSize returns the size of the file in bytes.
func (f *UploadedFile) GetSize() int {
	return f.size
}

// GetError returns what went wrong with the upload, and nil where nothing went
// wrong.
func (f *UploadedFile) GetError() error {
	if f.uploadError == constant.UploadErrorOk {
		return nil
	}

	return exception.NewHttpUploadedFileUploadError(f.uploadError)
}

// HasClientFilename reports whether the client stated a file name.
func (f *UploadedFile) HasClientFilename() bool {
	return f.clientFilename != ""
}

// GetClientFilename returns the file name that the client stated.
func (f *UploadedFile) GetClientFilename() string {
	return f.clientFilename
}

// HasClientMediaType reports whether the client stated a media type.
func (f *UploadedFile) HasClientMediaType() bool {
	return f.clientMediaType != ""
}

// GetClientMediaType returns the media type that the client stated.
func (f *UploadedFile) GetClientMediaType() string {
	return f.clientMediaType
}

// UploadedFileCollection holds the files that arrived with one request.
type UploadedFileCollection struct {
	files map[string]contract.UploadedFileContract
}

// NewUploadedFileCollection builds a collection over the files. It copies the
// map, so a later write to the source cannot reach the collection.
func NewUploadedFileCollection(files map[string]contract.UploadedFileContract) *UploadedFileCollection {
	return &UploadedFileCollection{files: copyFiles(files)}
}

// Has reports whether the collection holds a file under the key.
func (c *UploadedFileCollection) Has(key string) bool {
	_, found := c.files[key]

	return found
}

// Get returns the file under the key, and nil where the collection holds none.
func (c *UploadedFileCollection) Get(key string) contract.UploadedFileContract {
	return c.files[key]
}

// GetAll returns a copy of every file.
func (c *UploadedFileCollection) GetAll() map[string]contract.UploadedFileContract {
	return copyFiles(c.files)
}

// GetOnly returns the files that the keys identify.
func (c *UploadedFileCollection) GetOnly(keys ...string) map[string]contract.UploadedFileContract {
	return c.filterByKey(keys, true)
}

// GetAllExcept returns every file that the keys do not identify.
func (c *UploadedFileCollection) GetAllExcept(keys ...string) map[string]contract.UploadedFileContract {
	return c.filterByKey(keys, false)
}

// With returns a copy of the collection that holds the files and nothing else.
func (c *UploadedFileCollection) With(
	collection map[string]contract.UploadedFileContract,
) contract.UploadedFileCollectionContract {
	return NewUploadedFileCollection(collection)
}

// WithAdded returns a copy of the collection with the files added to the ones it
// holds. A file of a key that the collection holds is replaced.
func (c *UploadedFileCollection) WithAdded(
	collection map[string]contract.UploadedFileContract,
) contract.UploadedFileCollectionContract {
	combined := copyFiles(c.files)

	maps.Copy(combined, collection)

	return &UploadedFileCollection{files: combined}
}

// filterByKey returns the files whose key is in the keys, where keep is true,
// and the files whose key is not, where keep is false.
func (c *UploadedFileCollection) filterByKey(
	keys []string,
	keep bool,
) map[string]contract.UploadedFileContract {
	filtered := map[string]contract.UploadedFileContract{}

	for key, uploaded := range c.files {
		if slices.Contains(keys, key) == keep {
			filtered[key] = uploaded
		}
	}

	return filtered
}

// copyFiles returns a copy of the map, and an empty map where the map is nil.
func copyFiles(source map[string]contract.UploadedFileContract) map[string]contract.UploadedFileContract {
	target := make(map[string]contract.UploadedFileContract, len(source))

	maps.Copy(target, source)

	return target
}
