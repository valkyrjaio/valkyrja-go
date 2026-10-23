/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package param_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/param"
)

const (
	firstKey  = "first"
	secondKey = "second"

	firstValue  = "one"
	secondValue = "two"
)

// newCollection builds a collection that holds two parameters.
func newCollection() *param.ParamCollection {
	return param.NewParamCollection(map[string]any{
		firstKey:  firstValue,
		secondKey: secondValue,
	})
}

func TestNewParamCollectionHoldsEachParameter(t *testing.T) {
	t.Parallel()

	collection := newCollection()

	if !collection.Has(firstKey) {
		t.Error("Has must be true for a parameter that the collection holds, but is false")
	}

	if collection.Get(firstKey) != firstValue {
		t.Errorf("Get must return the parameter, but returned: %v", collection.Get(firstKey))
	}
}

func TestNewParamCollectionAcceptsNil(t *testing.T) {
	t.Parallel()

	collection := param.NewParamCollection(nil)

	if len(collection.GetAll()) != 0 {
		t.Errorf("the collection must be empty, but holds: %v", collection.GetAll())
	}
}

func TestHasAndGetReportAnUnknownKey(t *testing.T) {
	t.Parallel()

	collection := newCollection()

	if collection.Has("unknown") {
		t.Error("Has must be false for an unknown key, but is true")
	}

	if collection.Get("unknown") != nil {
		t.Errorf("Get must be nil for an unknown key, but is: %v", collection.Get("unknown"))
	}
}

func TestNewParamCollectionCopiesTheSourceMap(t *testing.T) {
	t.Parallel()

	params := map[string]any{firstKey: firstValue}
	collection := param.NewParamCollection(params)

	params[secondKey] = secondValue

	if collection.Has(secondKey) {
		t.Error("the collection must not follow a later write to the source map, but did")
	}
}

func TestGetAllReturnsACopy(t *testing.T) {
	t.Parallel()

	collection := newCollection()

	delete(collection.GetAll(), firstKey)

	if !collection.Has(firstKey) {
		t.Error("GetAll must return a copy, but the delete reached the collection")
	}
}

func TestGetOnlyAndGetAllExceptSplitTheCollection(t *testing.T) {
	t.Parallel()

	collection := newCollection()

	only := collection.GetOnly(firstKey)
	except := collection.GetAllExcept(firstKey)

	if len(only) != 1 || only[firstKey] != firstValue {
		t.Errorf("GetOnly must return the named parameter, but returned: %v", only)
	}

	if len(except) != 1 || except[secondKey] != secondValue {
		t.Errorf("GetAllExcept must return the other parameter, but returned: %v", except)
	}
}

func TestWithHoldsTheParametersAndNothingElse(t *testing.T) {
	t.Parallel()

	collection := newCollection()

	replaced := collection.With(map[string]any{"third": "three"})

	if replaced.Has(firstKey) {
		t.Error("With must drop what the collection held, but did not")
	}

	if !replaced.Has("third") {
		t.Error("With must hold the new parameter, but did not")
	}

	if !collection.Has(firstKey) {
		t.Error("With must leave the receiver unchanged, but did not")
	}
}

func TestWithAddedKeepsWhatTheCollectionHolds(t *testing.T) {
	t.Parallel()

	collection := newCollection()

	added := collection.WithAdded(map[string]any{"third": "three"})

	if !added.Has(firstKey) || !added.Has("third") {
		t.Error("WithAdded must hold both the old and the new parameters, but did not")
	}

	if collection.Has("third") {
		t.Error("WithAdded must leave the receiver unchanged, but did not")
	}
}

func TestWithAddedReplacesAParameterOfTheSameKey(t *testing.T) {
	t.Parallel()

	added := newCollection().WithAdded(map[string]any{firstKey: "replaced"})

	if added.Get(firstKey) != "replaced" {
		t.Errorf("WithAdded must replace a parameter of the same key, but is: %v", added.Get(firstKey))
	}
}

func TestTheCollectionSatisfiesEachNamedContract(t *testing.T) {
	t.Parallel()

	collection := newCollection()

	var server contract.ServerParamCollectionContract = collection
	var query contract.QueryParamCollectionContract = collection
	var attributes contract.AttributeParamCollectionContract = collection

	if !server.Has(firstKey) || !query.Has(firstKey) || !attributes.Has(firstKey) {
		t.Error("each named contract must read the collection, but did not")
	}
}
