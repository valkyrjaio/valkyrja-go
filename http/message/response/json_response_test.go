/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package response_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
)

const (
	callbackName = "handle"
	jsonKey      = "one"
	jsonValue    = "two"
)

// readBody returns the body of the response.
func readBody(t *testing.T, built contract.ResponseContract) string {
	t.Helper()

	body := built.GetBody()

	_ = body.Rewind()

	contents, err := body.GetContents()
	if err != nil {
		t.Fatalf("the body must be readable, but reported: %v", err)
	}

	return contents
}

func TestAJsonResponseCarriesItsDataAsJson(t *testing.T) {
	t.Parallel()

	built := response.NewJsonResponseFromData(map[string]any{jsonKey: jsonValue}, constant.StatusCodeOk, nil)

	if readBody(t, built) != `{"one":"two"}` {
		t.Errorf("the response must carry the data as JSON, but carried: %q", readBody(t, built))
	}

	if built.GetBodyAsJson()[jsonKey] != jsonValue {
		t.Errorf("the response must read its body back as JSON, but read: %v", built.GetBodyAsJson())
	}
}

func TestAJsonpResponseWrapsTheJsonInItsCallback(t *testing.T) {
	t.Parallel()

	built := response.NewJsonpResponseFromData(
		callbackName,
		map[string]any{jsonKey: jsonValue},
		constant.StatusCodeOk,
		nil,
	)

	if readBody(t, built) != callbackName+`({"one":"two"})` {
		t.Errorf("the response must wrap the JSON in the callback, but carried: %q", readBody(t, built))
	}

	if built.GetBodyAsJson()[jsonKey] != jsonValue {
		t.Errorf("the response must read the JSON out of the call, but read: %v", built.GetBodyAsJson())
	}
}

func TestTheJsonOfAResponseIsACopy(t *testing.T) {
	t.Parallel()

	built := response.NewJsonResponseFromData(map[string]any{jsonKey: jsonValue}, constant.StatusCodeOk, nil)

	read := built.GetBodyAsJson()
	delete(read, jsonKey)

	if built.GetBodyAsJson()[jsonKey] != jsonValue {
		t.Error("GetBodyAsJson must return a copy, so a change to it must not reach the response, but it did")
	}
}

func TestEachJsonResponseWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := response.NewJsonResponseFromData(map[string]any{jsonKey: jsonValue}, constant.StatusCodeOk, nil)

	replaced := built.WithJsonAsBody(map[string]any{"three": "four"})
	if replaced.GetBodyAsJson()["three"] != "four" {
		t.Error("WithJsonAsBody must carry the new data, but did not")
	}

	wrapped := built.WithCallback(callbackName)
	if !strings.HasPrefix(readBody(t, wrapped), callbackName+"(") {
		t.Errorf("WithCallback must wrap the JSON, but carried: %q", readBody(t, wrapped))
	}

	unwrapped := wrapped.WithoutCallback()
	if strings.Contains(readBody(t, unwrapped), callbackName) {
		t.Errorf("WithoutCallback must remove the call, but carried: %q", readBody(t, unwrapped))
	}

	if built.GetBodyAsJson()[jsonKey] != jsonValue {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestCreateFromDataReturnsAResponseThatCarriesTheData(t *testing.T) {
	t.Parallel()

	built := response.NewJsonResponseFromData(nil, constant.StatusCodeOk, nil)

	created := built.CreateFromData(map[string]any{jsonKey: jsonValue}, constant.StatusCodeCreated, nil)

	if created.GetStatusCode() != constant.StatusCodeCreated {
		t.Error("the response must carry the status code, but did not")
	}

	if created.GetBodyAsJson()[jsonKey] != jsonValue {
		t.Error("the response must carry the data, but did not")
	}
}
