/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// uploadErrorMessages holds the message that each upload error states. The
// wording is the wording of every port.
var uploadErrorMessages = map[constant.UploadError]string{
	constant.UploadErrorOk:        "There is no error, the file uploaded with success",
	constant.UploadErrorIniSize:   "The uploaded file exceeds the upload max filesize directive",
	constant.UploadErrorFormSize:  "The uploaded file exceeds the max file size directive of the form",
	constant.UploadErrorPartial:   "The uploaded file was only partially uploaded",
	constant.UploadErrorNoFile:    "No file was uploaded",
	constant.UploadErrorNoTmpDir:  "Missing a temporary folder",
	constant.UploadErrorCantWrite: "Failed to write file to disk",
	constant.UploadErrorExtension: "A PHP extension stopped the file upload",
}

type HttpUploadedFileRuntimeError struct {
	HttpRuntimeError
}

type HttpUploadedFileInvalidArgumentError struct {
	HttpInvalidArgumentError
}

type HttpUploadedFileUploadError struct {
	HttpUploadedFileRuntimeError

	uploadError constant.UploadError
}

// NewHttpUploadedFileUploadError builds the error for an upload error.
func NewHttpUploadedFileUploadError(uploadError constant.UploadError) *HttpUploadedFileUploadError {
	message, found := uploadErrorMessages[uploadError]
	if !found {
		message = "The file upload failed"
	}

	return &HttpUploadedFileUploadError{
		HttpUploadedFileRuntimeError: HttpUploadedFileRuntimeError{
			HttpRuntimeError: NewHttpRuntimeError(message, nil),
		},
		uploadError: uploadError,
	}
}

// GetUploadError returns the upload error that the error reports.
func (e *HttpUploadedFileUploadError) GetUploadError() constant.UploadError {
	return e.uploadError
}

type HttpUploadedFileAlreadyMovedError struct {
	HttpUploadedFileRuntimeError
}

// NewHttpUploadedFileAlreadyMovedError builds the error.
func NewHttpUploadedFileAlreadyMovedError() *HttpUploadedFileAlreadyMovedError {
	return &HttpUploadedFileAlreadyMovedError{
		HttpUploadedFileRuntimeError: HttpUploadedFileRuntimeError{
			HttpRuntimeError: NewHttpRuntimeError("Cannot retrieve stream after it has already moved", nil),
		},
	}
}

type HttpUploadedFileInvalidDirectoryError struct {
	HttpUploadedFileInvalidArgumentError

	targetPath string
}

// NewHttpUploadedFileInvalidDirectoryError builds the error for a target path.
func NewHttpUploadedFileInvalidDirectoryError(targetPath string) *HttpUploadedFileInvalidDirectoryError {
	return &HttpUploadedFileInvalidDirectoryError{
		HttpUploadedFileInvalidArgumentError: HttpUploadedFileInvalidArgumentError{
			HttpInvalidArgumentError: NewHttpInvalidArgumentError(
				"Invalid directory provided for the move operation",
				nil,
			),
		},
		targetPath: targetPath,
	}
}

// GetTargetPath returns the target path that the error reports.
func (e *HttpUploadedFileInvalidDirectoryError) GetTargetPath() string {
	return e.targetPath
}

type HttpUploadedFileMoveFailureError struct {
	HttpUploadedFileRuntimeError

	targetPath string
}

// NewHttpUploadedFileMoveFailureError builds the error for a target path.
func NewHttpUploadedFileMoveFailureError(targetPath string) *HttpUploadedFileMoveFailureError {
	return &HttpUploadedFileMoveFailureError{
		HttpUploadedFileRuntimeError: HttpUploadedFileRuntimeError{
			HttpRuntimeError: NewHttpRuntimeError("Error occurred while moving uploaded file", nil),
		},
		targetPath: targetPath,
	}
}

// GetTargetPath returns the target path that the error reports.
func (e *HttpUploadedFileMoveFailureError) GetTargetPath() string {
	return e.targetPath
}
