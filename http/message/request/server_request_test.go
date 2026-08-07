/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package request_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/param"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
)

const (
	paramKey  = "key"
	jsonValue = "two"
	jsonBody  = `{"one":"two"}`
)

func TestNewServerRequestTakesAnEmptyCollectionForEachPart(t *testing.T) {
	t.Parallel()

	built := request.NewServerRequest(nil, "", nil, nil)

	collections := map[string]contract.ParamCollectionContract{
		"the server parameters": built.GetServerParams(),
		"the cookies":           built.GetCookieParams(),
		"the query parameters":  built.GetQueryParams(),
		"the parsed body":       built.GetParsedBody(),
		"the attributes":        built.GetAttributes(),
	}

	for name, collection := range collections {
		if len(collection.GetAll()) != 0 {
			t.Errorf("%s must default to empty, but hold: %v", name, collection.GetAll())
		}
	}

	if built.GetUploadedFiles() != nil {
		t.Error("the uploaded files must default to none, but did not")
	}
}

func TestEachServerRequestWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := request.NewServerRequest(nil, "", nil, nil)
	params := param.NewParamCollection(map[string]any{paramKey: "value"})

	changed := built.
		WithServerParams(params).
		WithCookieParams(params).
		WithQueryParams(params).
		WithParsedBody(params).
		WithAttributes(params)

	server, isServer := changed.(*request.ServerRequest)
	if !isServer {
		t.Fatalf("each With method must return a server request, but returned: %T", changed)
	}

	collections := map[string]contract.ParamCollectionContract{
		"the server parameters": server.GetServerParams(),
		"the cookies":           server.GetCookieParams(),
		"the query parameters":  server.GetQueryParams(),
		"the parsed body":       server.GetParsedBody(),
		"the attributes":        server.GetAttributes(),
	}

	for name, collection := range collections {
		if !collection.Has(paramKey) {
			t.Errorf("%s must hold the new parameters, but did not", name)
		}
	}

	if built.GetServerParams().Has(paramKey) {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestWithUploadedFilesHoldsTheFiles(t *testing.T) {
	t.Parallel()

	built := request.NewServerRequest(nil, "", nil, nil)

	changed := built.WithUploadedFiles(nil)

	if changed.GetUploadedFiles() != nil {
		t.Error("WithUploadedFiles must hold what it receives, but did not")
	}
}

func TestIsXmlHttpRequestReadsTheRequestedWithHeader(t *testing.T) {
	t.Parallel()

	requestedWith, err := header.NewHeader(
		constant.HeaderNameXRequestedWith,
		value.NewValueFromValue("XMLHttpRequest"),
	)
	if err != nil {
		t.Fatalf("NewHeader must build the header, but reported: %v", err)
	}

	fromScript := request.NewServerRequest(nil, "", nil, header.NewHeaderCollection(requestedWith))
	fromBrowser := request.NewServerRequest(nil, "", nil, nil)

	if !fromScript.IsXmlHttpRequest() {
		t.Error("IsXmlHttpRequest must be true where the header names a script, but is false")
	}

	if fromBrowser.IsXmlHttpRequest() {
		t.Error("IsXmlHttpRequest must be false where the header is absent, but is true")
	}
}

func TestTheServerRequestSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.ServerRequestContract = request.NewServerRequest(nil, "", nil, nil)

	if built.GetMethod() != constant.RequestMethodGet {
		t.Errorf("the contract must read the method, but read: %q", built.GetMethod())
	}
}

// newJsonHeaders builds the headers that state a JSON content type.
func newJsonHeaders(t *testing.T, contentType string) contract.HeaderCollectionContract {
	t.Helper()

	if contentType == "" {
		return header.NewHeaderCollection()
	}

	built, err := header.NewHeader(constant.HeaderNameContentType, value.NewValueFromValue(contentType))
	if err != nil {
		t.Fatalf("the header must be valid, but reported: %v", err)
	}

	return header.NewHeaderCollection().WithHeader(built)
}

func TestAJsonServerRequestParsesItsBody(t *testing.T) {
	t.Parallel()

	built := request.NewJsonServerRequest(
		nil,
		constant.RequestMethodPost,
		stream.NewStream(jsonBody, constant.ModeRead),
		newJsonHeaders(t, constant.ContentTypeValueApplicationJson),
	)

	if built.GetParsedJson().Get("one") != jsonValue {
		t.Errorf("the request must parse its body, but parsed: %v", built.GetParsedJson().GetAll())
	}
}

func TestAJsonServerRequestParsesAContentTypeThatCarriesACharset(t *testing.T) {
	t.Parallel()

	built := request.NewJsonServerRequest(
		nil,
		constant.RequestMethodPost,
		stream.NewStream(jsonBody, constant.ModeRead),
		newJsonHeaders(t, constant.ContentTypeValueApplicationJson+"; charset=utf-8"),
	)

	if built.GetParsedJson().Get("one") != jsonValue {
		t.Error("a content type that carries a charset must still be read as JSON, but was not")
	}
}

func TestAJsonServerRequestParsesNothingWhereItCannot(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		body        string
		contentType string
	}{
		"a request that states no content type": {body: jsonBody, contentType: ""},
		"a content type that is not JSON":       {body: jsonBody, contentType: "text/plain"},
		"a body that carries nothing":           {body: "", contentType: constant.ContentTypeValueApplicationJson},
		"a body that no decoder reads":          {body: "not json", contentType: constant.ContentTypeValueApplicationJson},
	}

	for name, test := range tests {
		built := request.NewJsonServerRequest(
			nil,
			constant.RequestMethodPost,
			stream.NewStream(test.body, constant.ModeRead),
			newJsonHeaders(t, test.contentType),
		)

		if len(built.GetParsedJson().GetAll()) != 0 {
			t.Errorf("%s must parse nothing, but parsed: %v", name, built.GetParsedJson().GetAll())
		}
	}
}

func TestWithParsedJsonReturnsACopyThatIsStillAServerRequest(t *testing.T) {
	t.Parallel()

	built := request.NewJsonServerRequest(
		nil,
		constant.RequestMethodPost,
		stream.NewStream(jsonBody, constant.ModeRead),
		newJsonHeaders(t, constant.ContentTypeValueApplicationJson),
	)

	replaced := built.WithParsedJson(param.NewParamCollection(map[string]any{"three": "four"}))

	if replaced.GetParsedJson().Get("three") != "four" {
		t.Error("WithParsedJson must hold the new body, but did not")
	}

	if built.GetParsedJson().Get("one") != jsonValue {
		t.Error("WithParsedJson must leave the receiver unchanged, but did not")
	}

	// One struct carries both shapes, so a `With` method of the server request
	// keeps the parsed JSON rather than dropping it.
	kept := built.WithAttributes(param.NewParamCollection(map[string]any{paramKey: "value"}))

	if kept.GetParsedJson().Get("one") != jsonValue {
		t.Error("a server request With method must keep the parsed JSON, but dropped it")
	}
}
