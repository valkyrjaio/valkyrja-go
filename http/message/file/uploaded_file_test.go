/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package file_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/file"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

const (
	fileContents = "the contents"
	firstKey     = "first"
	secondKey    = "second"
)

// newUploadedFile builds a file that uploaded with success.
func newUploadedFile() *file.UploadedFile {
	return file.NewUploadedFile(
		stream.NewStream(fileContents, constant.ModeReadWrite),
		len(fileContents),
		constant.UploadErrorOk,
		"report.txt",
		"text/plain",
	)
}

func TestTheUploadedFileReadsWhatItHolds(t *testing.T) {
	t.Parallel()

	uploaded := newUploadedFile()

	if uploaded.GetSize() != len(fileContents) || !uploaded.HasSize() {
		t.Error("the file must report its size, but did not")
	}

	if uploaded.GetClientFilename() != "report.txt" || !uploaded.HasClientFilename() {
		t.Error("the file must report the client file name, but did not")
	}

	if uploaded.GetClientMediaType() != "text/plain" || !uploaded.HasClientMediaType() {
		t.Error("the file must report the client media type, but did not")
	}

	if uploaded.GetError() != nil {
		t.Errorf("GetError must be nil for a file that uploaded, but is: %v", uploaded.GetError())
	}

	if uploaded.GetStream().String() != fileContents {
		t.Error("the stream must carry the contents, but did not")
	}
}

func TestTheUploadedFileReportsWhatTheClientDidNotState(t *testing.T) {
	t.Parallel()

	uploaded := file.NewUploadedFile(nil, -1, constant.UploadErrorOk, "", "")

	if uploaded.HasSize() {
		t.Error("HasSize must be false where the client stated no size, but is true")
	}

	if uploaded.HasClientFilename() || uploaded.HasClientMediaType() {
		t.Error("the file must report that the client stated no name and no media type, but did not")
	}
}

func TestGetErrorReportsAnUploadThatWentWrong(t *testing.T) {
	t.Parallel()

	uploaded := file.NewUploadedFile(nil, 0, constant.UploadErrorNoFile, "", "")

	target, found := errors.AsType[*exception.HttpUploadedFileUploadError](uploaded.GetError())
	if !found {
		t.Fatalf("GetError must report the upload error, but reported: %v", uploaded.GetError())
	}

	if target.GetUploadError() != constant.UploadErrorNoFile {
		t.Errorf("the error must carry the upload error, but carries: %d", target.GetUploadError())
	}

	if target.Error() != "No file was uploaded" {
		t.Errorf("the error must state what went wrong, but is: %q", target.Error())
	}
}

func TestAnUnknownUploadErrorStatesAGeneralFailure(t *testing.T) {
	t.Parallel()

	err := exception.NewHttpUploadedFileUploadError(constant.UploadError(99))

	if err.Error() != "The file upload failed" {
		t.Errorf("an unknown upload error must state a general failure, but is: %q", err.Error())
	}
}

func TestMoveToWritesTheFile(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "report.txt")

	uploaded := newUploadedFile()

	err := uploaded.MoveTo(target)
	if err != nil {
		t.Fatalf("MoveTo must write the file, but reported: %v", err)
	}

	written, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("the test must read the file, but reported: %v", err)
	}

	if string(written) != fileContents {
		t.Errorf("the file must carry the contents, but carries: %q", string(written))
	}
}

func TestMoveToReportsAFileThatMovedAlready(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "report.txt")

	uploaded := newUploadedFile()

	err := uploaded.MoveTo(target)
	if err != nil {
		t.Fatalf("MoveTo must write the file, but reported: %v", err)
	}

	err = uploaded.MoveTo(target)

	if _, found := errors.AsType[*exception.HttpUploadedFileAlreadyMovedError](err); !found {
		t.Errorf("MoveTo must report a file that moved already, but reported: %v", err)
	}
}

func TestMoveToReportsAnUploadThatWentWrong(t *testing.T) {
	t.Parallel()

	uploaded := file.NewUploadedFile(nil, 0, constant.UploadErrorPartial, "", "")

	err := uploaded.MoveTo(filepath.Join(t.TempDir(), "report.txt"))

	if _, found := errors.AsType[*exception.HttpUploadedFileUploadError](err); !found {
		t.Errorf("MoveTo must report the upload error, but reported: %v", err)
	}
}

func TestMoveToReportsAnEmptyTargetPath(t *testing.T) {
	t.Parallel()

	err := newUploadedFile().MoveTo("")

	target, found := errors.AsType[*exception.HttpUploadedFileInvalidDirectoryError](err)
	if !found {
		t.Fatalf("MoveTo must report an invalid directory, but reported: %v", err)
	}

	if target.GetTargetPath() != "" {
		t.Errorf("the error must carry the target path, but carries: %q", target.GetTargetPath())
	}
}

func TestMoveToReportsAStreamThatNoReaderReads(t *testing.T) {
	t.Parallel()

	uploaded := file.NewUploadedFile(
		stream.NewStream(fileContents, constant.ModeWrite),
		0,
		constant.UploadErrorOk,
		"",
		"",
	)

	err := uploaded.MoveTo(filepath.Join(t.TempDir(), "report.txt"))

	if _, found := errors.AsType[*exception.HttpStreamUnreadableStreamError](err); !found {
		t.Errorf("MoveTo must carry the stream failure, but reported: %v", err)
	}
}

func TestMoveToReportsAPathThatTheFileSystemRefuses(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "missing", "report.txt")

	err := newUploadedFile().MoveTo(target)

	moveFailure, found := errors.AsType[*exception.HttpUploadedFileMoveFailureError](err)
	if !found {
		t.Fatalf("MoveTo must report a move failure, but reported: %v", err)
	}

	if moveFailure.GetTargetPath() != target {
		t.Errorf("the error must carry the target path, but carries: %q", moveFailure.GetTargetPath())
	}
}

func TestTheCollectionHoldsEachFile(t *testing.T) {
	t.Parallel()

	collection := file.NewUploadedFileCollection(map[string]contract.UploadedFileContract{
		firstKey: newUploadedFile(),
	})

	if !collection.Has(firstKey) {
		t.Error("Has must be true for a file that the collection holds, but is false")
	}

	if collection.Get(firstKey) == nil {
		t.Error("Get must return the file, but returned nil")
	}

	if collection.Has(secondKey) || collection.Get(secondKey) != nil {
		t.Error("the collection must report an unknown key, but did not")
	}
}

func TestTheCollectionCopiesItsSourceMap(t *testing.T) {
	t.Parallel()

	files := map[string]contract.UploadedFileContract{firstKey: newUploadedFile()}
	collection := file.NewUploadedFileCollection(files)

	files[secondKey] = newUploadedFile()

	if collection.Has(secondKey) {
		t.Error("the collection must not follow a later write to the source map, but did")
	}

	delete(collection.GetAll(), firstKey)

	if !collection.Has(firstKey) {
		t.Error("GetAll must return a copy, but the delete reached the collection")
	}
}

func TestGetOnlyAndGetAllExceptSplitTheCollection(t *testing.T) {
	t.Parallel()

	collection := file.NewUploadedFileCollection(map[string]contract.UploadedFileContract{
		firstKey:  newUploadedFile(),
		secondKey: newUploadedFile(),
	})

	if len(collection.GetOnly(firstKey)) != 1 {
		t.Error("GetOnly must return the named file, but did not")
	}

	if len(collection.GetAllExcept(firstKey)) != 1 {
		t.Error("GetAllExcept must return the other file, but did not")
	}
}

func TestWithAndWithAddedReturnACopy(t *testing.T) {
	t.Parallel()

	collection := file.NewUploadedFileCollection(map[string]contract.UploadedFileContract{
		firstKey: newUploadedFile(),
	})

	replaced := collection.With(map[string]contract.UploadedFileContract{secondKey: newUploadedFile()})
	added := collection.WithAdded(map[string]contract.UploadedFileContract{secondKey: newUploadedFile()})

	if replaced.Has(firstKey) || !replaced.Has(secondKey) {
		t.Error("With must hold the new files and nothing else, but did not")
	}

	if !added.Has(firstKey) || !added.Has(secondKey) {
		t.Error("WithAdded must hold both sets of files, but did not")
	}

	if collection.Has(secondKey) {
		t.Error("each method must leave the receiver unchanged, but did not")
	}
}

func TestTheCollectionSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var collection contract.UploadedFileCollectionContract = file.NewUploadedFileCollection(nil)

	if len(collection.GetAll()) != 0 {
		t.Error("the contract must read the collection, but did not")
	}
}
