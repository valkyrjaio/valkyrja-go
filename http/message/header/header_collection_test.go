/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package header_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

const jsonValue = "application/json"

func TestTheCollectionKeysOnTheNormalizedName(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(
		newHeader(t, constant.HeaderNameContentType, value.NewValueFromValue(htmlValue)),
	)

	if !collection.Has("CONTENT-TYPE") {
		t.Error("Has must read the name in any case, but did not")
	}

	found, err := collection.Get("content-type")
	if err != nil {
		t.Fatalf("Get must return the header, but reported: %v", err)
	}

	if found.GetHeaderLine() != htmlValue {
		t.Errorf("Get must return the header, but returned: %q", found.GetHeaderLine())
	}
}

func TestGetReportsAHeaderThatTheCollectionDoesNotHold(t *testing.T) {
	t.Parallel()

	_, err := header.NewHeaderCollection().Get(acceptName)

	target, found := errors.AsType[*exception.HttpHeaderInvalidHeaderNameError](err)
	if !found {
		t.Fatalf("Get must report an unknown header, but reported: %v", err)
	}

	if target.GetName() != acceptName {
		t.Errorf("the error must name the header, but names: %q", target.GetName())
	}
}

func TestGetHeaderLineIsEmptyForAnUnknownHeader(t *testing.T) {
	t.Parallel()

	if header.NewHeaderCollection().GetHeaderLine(acceptName) != "" {
		t.Error("GetHeaderLine must be empty for an unknown header, but is not")
	}
}

func TestGetAllReturnsACopy(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(newHeader(t, constant.HeaderNameContentType))

	delete(collection.GetAll(), "content-type")

	if !collection.Has(constant.HeaderNameContentType) {
		t.Error("GetAll must return a copy, but the delete reached the collection")
	}
}

func TestGetOnlyAndGetAllExceptSplitTheCollection(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(
		newHeader(t, constant.HeaderNameContentType),
		newHeader(t, acceptName),
	)

	only := collection.GetOnly("CONTENT-TYPE")
	except := collection.GetAllExcept("CONTENT-TYPE")

	if len(only) != 1 {
		t.Errorf("GetOnly must return the named header, but returned: %d", len(only))
	}

	if _, found := only["content-type"]; !found {
		t.Error("GetOnly must key on the normalized name, but did not")
	}

	if len(except) != 1 {
		t.Errorf("GetAllExcept must return the other header, but returned: %d", len(except))
	}

	if _, found := except["accept"]; !found {
		t.Error("GetAllExcept must return the header that the name does not identify, but did not")
	}
}

func TestWithHeaderReplacesAHeaderOfTheSameName(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(
		newHeader(t, acceptName, value.NewValueFromValue(htmlValue)),
	)

	replaced := collection.WithHeader(newHeader(t, acceptName, value.NewValueFromValue(jsonValue)))

	if replaced.GetHeaderLine(acceptName) != jsonValue {
		t.Errorf("WithHeader must replace the header, but is: %q", replaced.GetHeaderLine(acceptName))
	}

	if collection.GetHeaderLine(acceptName) != htmlValue {
		t.Error("WithHeader must leave the receiver unchanged, but did not")
	}
}

func TestWithoutHeaderRemovesTheHeader(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(newHeader(t, acceptName))

	if header.NewHeaderCollection(newHeader(t, acceptName)).WithoutHeader("ACCEPT").Has(acceptName) {
		t.Error("WithoutHeader must remove the header, but did not")
	}

	if !collection.Has(acceptName) {
		t.Error("WithoutHeader must leave the receiver unchanged, but did not")
	}
}

func TestWithHeadersHoldsTheHeadersAndNothingElse(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(newHeader(t, constant.HeaderNameContentType))

	replaced := collection.WithHeaders(newHeader(t, acceptName))

	if replaced.Has(constant.HeaderNameContentType) {
		t.Error("WithHeaders must drop what the collection held, but did not")
	}

	if !replaced.Has(acceptName) {
		t.Error("WithHeaders must hold the new header, but did not")
	}
}

func TestWithAddedHeadersMergesAHeaderOfTheSameName(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(
		newHeader(t, acceptName, value.NewValueFromValue(htmlValue)),
	)

	added := collection.WithAddedHeaders(newHeader(t, acceptName, value.NewValueFromValue(jsonValue)))

	if added.GetHeaderLine(acceptName) != acceptedTypes {
		t.Errorf("WithAddedHeaders must merge the values, but is: %q", added.GetHeaderLine(acceptName))
	}
}

func TestWithAddedHeadersHoldsAHeaderOfAnotherName(t *testing.T) {
	t.Parallel()

	collection := header.NewHeaderCollection(newHeader(t, constant.HeaderNameContentType))

	added := collection.WithAddedHeaders(newHeader(t, acceptName))

	if !added.Has(constant.HeaderNameContentType) || !added.Has(acceptName) {
		t.Error("WithAddedHeaders must hold both headers, but did not")
	}

	if collection.Has(acceptName) {
		t.Error("WithAddedHeaders must leave the receiver unchanged, but did not")
	}
}
