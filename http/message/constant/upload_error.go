/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// UploadError is what went wrong with one uploaded file.
type UploadError int

// The UploadError values that the framework knows. The numbers are PHP's own
// upload error codes, which every port keeps so an uploaded file reports the
// same failure in each one. The number 5 is absent in PHP as well.
const (
	// UploadErrorOk reports that nothing went wrong.
	UploadErrorOk UploadError = 0

	// UploadErrorIniSize reports a file over the server's size limit.
	UploadErrorIniSize UploadError = 1

	// UploadErrorFormSize reports a file over the form's size limit.
	UploadErrorFormSize UploadError = 2

	// UploadErrorPartial reports a file that arrived in part.
	UploadErrorPartial UploadError = 3

	// UploadErrorNoFile reports that no file arrived.
	UploadErrorNoFile UploadError = 4

	// UploadErrorNoTmpDir reports a missing temporary directory.
	UploadErrorNoTmpDir UploadError = 6

	// UploadErrorCantWrite reports a file that the server cannot write.
	UploadErrorCantWrite UploadError = 7

	// UploadErrorExtension reports an extension that stopped the upload.
	UploadErrorExtension UploadError = 8
)
