/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package structure_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/param"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
	"github.com/valkyrjaio/valkyrja-go/v26/http/structure"
)

const (
	structureName = "users.create"
	nameField     = "name"
	emailField    = "email"
	nameValue     = "Melech"
	renamedField  = "full_name"
	extraValue    = "one"
)

// newJsonRequest builds a request that carries the body as JSON.
func newJsonRequest(t *testing.T, body string) contract.ServerRequestContract {
	t.Helper()

	contentType, err := header.NewHeader(
		constant.HeaderNameContentType,
		value.NewValueFromValue(constant.ContentTypeValueApplicationJson),
	)
	if err != nil {
		t.Fatalf("the header must be valid, but reported: %v", err)
	}

	return request.NewJsonServerRequest(
		nil,
		constant.RequestMethodPost,
		stream.NewStream(body, constant.ModeRead),
		header.NewHeaderCollection().WithHeader(contentType),
	)
}

func TestEachRequestStructureReadsItsOwnCollection(t *testing.T) {
	t.Parallel()

	params := param.NewParamCollection(map[string]any{nameField: nameValue, "extra": extraValue})

	built := request.NewServerRequest(nil, constant.RequestMethodPost, nil, nil)

	tests := map[string]struct {
		structure *structure.RequestStructure
		request   contract.ServerRequestContract
	}{
		"a query structure": {
			structure: structure.NewQueryRequestStructure(structureName, nameField),
			request:   built.WithQueryParams(params),
		},
		"a parsed body structure": {
			structure: structure.NewParsedBodyRequestStructure(structureName, nameField),
			request:   built.WithParsedBody(params),
		},
		"a JSON structure": {
			structure: structure.NewJsonRequestStructure(structureName, nameField),
			request:   newJsonRequest(t, `{"name":"Melech","extra":"one"}`),
		},
	}

	for name, test := range tests {
		data := test.structure.GetDataFromRequest(test.request)

		if len(data) != 1 || data[nameField] != nameValue {
			t.Errorf("%s must read the field that it names, but read: %v", name, data)
		}

		if !test.structure.DetermineIfRequestContainsExtraData(test.request) {
			t.Errorf("%s must report the field that it does not name, but did not", name)
		}
	}
}

func TestARequestStructureReportsNoExtraDataWhereItNamesEveryField(t *testing.T) {
	t.Parallel()

	built := request.NewServerRequest(nil, constant.RequestMethodPost, nil, nil).
		WithQueryParams(param.NewParamCollection(map[string]any{nameField: nameValue}))

	structured := structure.NewQueryRequestStructure(structureName, nameField)

	if structured.DetermineIfRequestContainsExtraData(built) {
		t.Error("a structure that names every field must report no extra data, but did")
	}
}

func TestAJsonStructureReadsAnEmptyCollectionWhereTheRequestCarriesNoJson(t *testing.T) {
	t.Parallel()

	// A request that parsed no JSON body reads as empty rather than as a
	// failure, because every server request of this port carries the parsed
	// JSON.
	built := request.NewServerRequest(nil, constant.RequestMethodPost, nil, nil)

	structured := structure.NewJsonRequestStructure(structureName, nameField)

	if len(structured.GetDataFromRequest(built)) != 0 {
		t.Error("a request that carries no JSON must read as empty, but did not")
	}
}

func TestARequestStructureReadsWhatItWasBuiltWith(t *testing.T) {
	t.Parallel()

	built := structure.NewQueryRequestStructure(structureName, nameField, emailField)

	if built.GetName() != structureName {
		t.Errorf("the structure must read its name, but read: %q", built.GetName())
	}

	if len(built.GetFields()) != 2 {
		t.Errorf("the structure must name each field, but named: %v", built.GetFields())
	}

	if built.GetValue() != nil {
		t.Error("a structure that was built with no value must carry none, but carried one")
	}

	if built.WithValue(extraValue).GetValue() != extraValue {
		t.Error("WithValue must hold the new value, but did not")
	}

	if built.GetValue() != nil {
		t.Error("WithValue must leave the receiver unchanged, but did not")
	}
}

func TestAResponseStructureRenamesEachFieldThatItNames(t *testing.T) {
	t.Parallel()

	built := structure.NewResponseStructure(structureName, map[string]string{
		nameField:  renamedField,
		emailField: "email_address",
	})

	structured := built.GetStructuredData(map[string]any{
		nameField:  nameValue,
		emailField: "melech@example.com",
		"secret":   "never",
	}, true)

	if structured[renamedField] != nameValue || structured["email_address"] != "melech@example.com" {
		t.Errorf("the structure must rename each field, but returned: %v", structured)
	}

	if len(structured) != 2 {
		t.Errorf("a field that the structure does not name must never reach a client, but did: %v", structured)
	}
}

func TestAResponseStructureReportsAFieldThatTheDataDoesNotCarry(t *testing.T) {
	t.Parallel()

	built := structure.NewResponseStructure(structureName, map[string]string{nameField: renamedField})

	withAll := built.GetStructuredData(map[string]any{}, true)

	held, carried := withAll[renamedField]
	if !carried || held != nil {
		t.Errorf("a missing field must read as nil where includeAll is true, but read: %v", withAll)
	}

	withoutAll := built.GetStructuredData(map[string]any{}, false)

	if len(withoutAll) != 0 {
		t.Errorf("a missing field must be left out where includeAll is false, but was: %v", withoutAll)
	}
}

func TestAResponseStructureReadsWhatItWasBuiltWith(t *testing.T) {
	t.Parallel()

	built := structure.NewResponseStructure(structureName, map[string]string{nameField: renamedField})

	if built.GetName() != structureName || len(built.GetFields()) != 1 {
		t.Error("the structure must read its name and its fields, but did not")
	}

	if built.GetValue() != nil || built.WithValue(extraValue).GetValue() != extraValue {
		t.Error("WithValue must hold the new value and leave the receiver unchanged, but did not")
	}
}
