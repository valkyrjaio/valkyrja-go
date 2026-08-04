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
)

const paramKey = "key"

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
