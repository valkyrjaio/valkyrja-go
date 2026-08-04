/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

func TestEachBaseErrorSatisfiesTheComponentContract(t *testing.T) {
	t.Parallel()

	runtimeError := exception.NewHttpRuntimeError("the runtime failed", nil)
	invalidArgumentError := exception.NewHttpInvalidArgumentError("the argument is invalid", nil)

	throwables := []contract.HttpThrowable{&runtimeError, &invalidArgumentError}

	for _, throwable := range throwables {
		if !throwable.IsHttpThrowable() {
			t.Errorf("IsHttpThrowable must be true, but is false for: %s", throwable.Error())
		}

		if throwable.GetTraceCode() == "" {
			t.Errorf("GetTraceCode must return a code, but is empty for: %s", throwable.Error())
		}
	}
}

func TestEachHeaderErrorReportsWhatItCarries(t *testing.T) {
	t.Parallel()

	invalidName := exception.NewHttpHeaderInvalidNameError("Bad Name")
	invalidValue := exception.NewHttpHeaderInvalidValueError("bad\nvalue")
	unknownName := exception.NewHttpHeaderInvalidHeaderNameError("Accept")

	if invalidName.GetName() != "Bad Name" {
		t.Errorf("GetName must be the header name, but is: %q", invalidName.GetName())
	}

	if invalidName.Error() != `"Bad Name" is not valid header name` {
		t.Errorf("Error must name the header, but is: %q", invalidName.Error())
	}

	if invalidValue.GetValue() != "bad\nvalue" {
		t.Errorf("GetValue must be the header value, but is: %q", invalidValue.GetValue())
	}

	if unknownName.GetName() != "Accept" {
		t.Errorf("GetName must be the header name, but is: %q", unknownName.GetName())
	}

	if unknownName.Error() != "Header Accept does not exist" {
		t.Errorf("Error must name the header, but is: %q", unknownName.Error())
	}
}

func TestEachHeaderErrorSatisfiesTheComponentContract(t *testing.T) {
	t.Parallel()

	throwables := []contract.HttpThrowable{
		exception.NewHttpHeaderInvalidNameError("Bad Name"),
		exception.NewHttpHeaderInvalidValueError("bad\nvalue"),
		exception.NewHttpHeaderInvalidHeaderNameError("Accept"),
	}

	for _, throwable := range throwables {
		if !throwable.IsHttpThrowable() {
			t.Errorf("IsHttpThrowable must be true, but is false for: %s", throwable.Error())
		}
	}
}

func TestErrorsAsSeparatesTheHeaderErrors(t *testing.T) {
	t.Parallel()

	var err error = exception.NewHttpHeaderInvalidNameError("Bad Name")

	if _, found := errors.AsType[*exception.HttpHeaderInvalidValueError](err); found {
		t.Error("errors.AsType must not find the value error in a name error, but did")
	}

	if _, found := errors.AsType[*exception.HttpHeaderInvalidNameError](err); !found {
		t.Error("errors.AsType must find the name error, but did not")
	}
}

func TestEachStreamErrorStatesWhatWentWrong(t *testing.T) {
	t.Parallel()

	invalidLength := exception.NewHttpStreamInvalidLengthError(-1)

	if invalidLength.GetLength() != -1 {
		t.Errorf("GetLength must be the length, but is: %d", invalidLength.GetLength())
	}

	tests := map[string]contract.HttpThrowable{
		"Invalid length of -1 provided. Length must be greater than 0": invalidLength,
		"Stream is not readable":                    exception.NewHttpStreamUnreadableStreamError(),
		"Stream is not writable":                    exception.NewHttpStreamUnwritableStreamError(),
		"Stream is not seekable":                    exception.NewHttpStreamUnseekableStreamError(),
		"Position is outside of the stream":         exception.NewHttpStreamStreamSeekError(),
		"Could not read the position of the stream": exception.NewHttpStreamStreamTellError(),
	}

	for message, throwable := range tests {
		if throwable.Error() != message {
			t.Errorf("Error must be %q, but is %q", message, throwable.Error())
		}

		if !throwable.IsHttpThrowable() {
			t.Errorf("IsHttpThrowable must be true, but is false for: %s", throwable.Error())
		}
	}
}

func TestErrorsAsSeparatesTheStreamErrors(t *testing.T) {
	t.Parallel()

	var err error = exception.NewHttpStreamUnreadableStreamError()

	if _, found := errors.AsType[*exception.HttpStreamUnwritableStreamError](err); found {
		t.Error("errors.AsType must not find the unwritable error in an unreadable error, but did")
	}

	if _, found := errors.AsType[*exception.HttpStreamUnreadableStreamError](err); !found {
		t.Error("errors.AsType must find the unreadable error, but did not")
	}
}

func TestEachUriErrorStatesWhatWentWrong(t *testing.T) {
	t.Parallel()

	invalidPort := exception.NewHttpUriInvalidPortError(70000)
	invalidPath := exception.NewHttpUriInvalidPathError("/path?a=b")
	invalidQuery := exception.NewHttpUriInvalidQueryError("a=b#part")

	if invalidPort.GetPort() != 70000 {
		t.Errorf("GetPort must be the port, but is: %d", invalidPort.GetPort())
	}

	if invalidPath.GetPath() != "/path?a=b" {
		t.Errorf("GetPath must be the path, but is: %q", invalidPath.GetPath())
	}

	if invalidQuery.GetQuery() != "a=b#part" {
		t.Errorf("GetQuery must be the query string, but is: %q", invalidQuery.GetQuery())
	}

	tests := map[string]contract.HttpThrowable{
		"Invalid port `70000` specified; must be a valid TCP/UDP port":                 invalidPort,
		"Invalid path of `/path?a=b` provided; must not contain a query string":        invalidPath,
		"Invalid query string of `a=b#part` provided; must not contain a URI fragment": invalidQuery,
	}

	for message, throwable := range tests {
		if throwable.Error() != message {
			t.Errorf("Error must be %q, but is %q", message, throwable.Error())
		}

		if !throwable.IsHttpThrowable() {
			t.Errorf("IsHttpThrowable must be true, but is false for: %s", throwable.Error())
		}
	}
}

func TestErrorsAsSeparatesTheUriErrors(t *testing.T) {
	t.Parallel()

	var err error = exception.NewHttpUriInvalidPathError("/path?a=b")

	if _, found := errors.AsType[*exception.HttpUriInvalidQueryError](err); found {
		t.Error("errors.AsType must not find the query error in a path error, but did")
	}

	if _, found := errors.AsType[*exception.HttpUriInvalidPathError](err); !found {
		t.Error("errors.AsType must find the path error, but did not")
	}
}

func TestTheRequestTargetErrorCarriesTheTarget(t *testing.T) {
	t.Parallel()

	err := exception.NewHttpRequestInvalidRequestTargetError("/a target")

	if err.GetRequestTarget() != "/a target" {
		t.Errorf("GetRequestTarget must be the target, but is: %q", err.GetRequestTarget())
	}

	if err.Error() != "Invalid request target provided; cannot contain whitespace" {
		t.Errorf("Error must state the rule, but is: %q", err.Error())
	}

	if !contract.HttpThrowable(err).IsHttpThrowable() {
		t.Error("IsHttpThrowable must be true, but is false")
	}
}

func TestEachUploadedFileErrorStatesWhatWentWrong(t *testing.T) {
	t.Parallel()

	uploadError := exception.NewHttpUploadedFileUploadError(constant.UploadErrorIniSize)
	alreadyMoved := exception.NewHttpUploadedFileAlreadyMovedError()
	invalidDirectory := exception.NewHttpUploadedFileInvalidDirectoryError("/missing")
	moveFailure := exception.NewHttpUploadedFileMoveFailureError("/missing/report.txt")

	if uploadError.GetUploadError() != constant.UploadErrorIniSize {
		t.Errorf("GetUploadError must be the upload error, but is: %d", uploadError.GetUploadError())
	}

	if invalidDirectory.GetTargetPath() != "/missing" {
		t.Errorf("GetTargetPath must be the path, but is: %q", invalidDirectory.GetTargetPath())
	}

	if moveFailure.GetTargetPath() != "/missing/report.txt" {
		t.Errorf("GetTargetPath must be the path, but is: %q", moveFailure.GetTargetPath())
	}

	tests := map[string]contract.HttpThrowable{
		"The uploaded file exceeds the upload max filesize directive": uploadError,
		"Cannot retrieve stream after it has already moved":           alreadyMoved,
		"Invalid directory provided for the move operation":           invalidDirectory,
		"Error occurred while moving uploaded file":                   moveFailure,
	}

	for message, throwable := range tests {
		if throwable.Error() != message {
			t.Errorf("Error must be %q, but is %q", message, throwable.Error())
		}

		if !throwable.IsHttpThrowable() {
			t.Errorf("IsHttpThrowable must be true, but is false for: %s", throwable.Error())
		}
	}
}

func TestAnUnknownUploadErrorStatesAGeneralFailure(t *testing.T) {
	t.Parallel()

	err := exception.NewHttpUploadedFileUploadError(constant.UploadError(99))

	if err.Error() != "The file upload failed" {
		t.Errorf("an unknown upload error must state a general failure, but is: %q", err.Error())
	}
}
