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
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
)

const body = "the body"

func TestNewResponseTakesEachDefault(t *testing.T) {
	t.Parallel()

	built := response.NewResponse(nil, 0, nil)

	if built.GetStatusCode() != constant.StatusCodeOk {
		t.Errorf("the status code must default to 200, but is: %d", built.GetStatusCode())
	}

	if built.GetReasonPhrase() != "OK" {
		t.Errorf("the reason phrase must follow the status code, but is: %q", built.GetReasonPhrase())
	}

	if built.GetProtocolVersion() != constant.ProtocolVersionV11 {
		t.Errorf("the protocol version must default to 1.1, but is: %q", built.GetProtocolVersion())
	}

	if built.GetBody().GetSize() != 0 {
		t.Error("the body must default to empty, but did not")
	}

	if len(built.GetHeaders().GetAll()) != 0 {
		t.Error("the headers must default to empty, but did not")
	}
}

func TestNewResponseFromContentCarriesTheContent(t *testing.T) {
	t.Parallel()

	built := response.NewResponseFromContent(body, constant.StatusCodeCreated, nil)

	if built.GetBody().String() != body {
		t.Errorf("the body must carry the content, but is: %q", built.GetBody().String())
	}

	if built.GetStatusCode() != constant.StatusCodeCreated {
		t.Errorf("the status code must be the one given, but is: %d", built.GetStatusCode())
	}
}

func TestCreateBuildsAResponseFromTheContent(t *testing.T) {
	t.Parallel()

	built := response.NewResponse(nil, 0, nil).Create(body, constant.StatusCodeNotFound, nil)

	if built.GetBody().String() != body {
		t.Errorf("Create must carry the content, but the body is: %q", built.GetBody().String())
	}

	if built.GetStatusCode() != constant.StatusCodeNotFound {
		t.Errorf("Create must take the status code, but is: %d", built.GetStatusCode())
	}
}

func TestWithStatusCodeMovesTheReasonPhrase(t *testing.T) {
	t.Parallel()

	built := response.NewResponse(nil, 0, nil)

	changed := built.WithStatusCode(constant.StatusCodeNotFound)

	if changed.GetReasonPhrase() != "Not Found" {
		t.Errorf("the reason phrase must follow the status code, but is: %q", changed.GetReasonPhrase())
	}

	if built.GetStatusCode() != constant.StatusCodeOk {
		t.Error("WithStatusCode must leave the receiver unchanged, but did not")
	}
}

func TestWithReasonPhraseTakesTheStatusPhraseWhereItIsEmpty(t *testing.T) {
	t.Parallel()

	built := response.NewResponse(nil, constant.StatusCodeNotFound, nil)

	if built.WithReasonPhrase("Gone Missing").GetReasonPhrase() != "Gone Missing" {
		t.Error("WithReasonPhrase must hold the new phrase, but did not")
	}

	if built.WithReasonPhrase("").GetReasonPhrase() != "Not Found" {
		t.Error("an empty phrase must take the phrase of the status code, but did not")
	}
}

func TestEachMessageWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := response.NewResponse(nil, 0, nil)

	other := stream.NewStream(body, constant.ModeReadWrite)
	headers := header.NewHeaderCollection()

	if built.WithProtocolVersion(constant.ProtocolVersionV2).GetProtocolVersion() != constant.ProtocolVersionV2 {
		t.Error("WithProtocolVersion must hold the new version, but did not")
	}

	if built.WithBody(other).GetBody().String() != body {
		t.Error("WithBody must hold the new body, but did not")
	}

	if built.WithHeaders(headers).GetHeaders() != contract.HeaderCollectionContract(headers) {
		t.Error("WithHeaders must hold the new headers, but did not")
	}

	if built.GetProtocolVersion() != constant.ProtocolVersionV11 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestWithBodyRewindsTheBody(t *testing.T) {
	t.Parallel()

	other := stream.NewStream(body, constant.ModeReadWrite)

	_, err := other.Read(3)
	if err != nil {
		t.Fatalf("Read must read the bytes, but reported: %v", err)
	}

	built := response.NewResponse(nil, 0, nil).WithBody(other)

	contents, err := built.GetBody().GetContents()
	if err != nil {
		t.Fatalf("GetContents must read the body, but reported: %v", err)
	}

	if contents != body {
		t.Errorf("WithBody must rewind the body, but it reads: %q", contents)
	}
}

func TestWithCookieSetsTheCookieHeader(t *testing.T) {
	t.Parallel()

	built := response.NewResponse(nil, 0, nil).WithCookie(value.NewCookie("session", "abc123"))

	line := built.GetHeaders().GetHeaderLine(constant.HeaderNameSetCookie)

	if !strings.Contains(line, "session=abc123") {
		t.Errorf("WithCookie must set the cookie header, but it is: %q", line)
	}
}

func TestWithoutCookieSetsTheDeletedCookie(t *testing.T) {
	t.Parallel()

	built := response.NewResponse(nil, 0, nil).WithoutCookie(value.NewCookie("session", "abc123"))

	line := built.GetHeaders().GetHeaderLine(constant.HeaderNameSetCookie)

	if !strings.Contains(line, "session=delete") {
		t.Errorf("WithoutCookie must set the deleted cookie, but the header is: %q", line)
	}
}

func TestTheResponseSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.ResponseContract = response.NewResponse(nil, 0, nil)

	if built.GetStatusCode() != constant.StatusCodeOk {
		t.Errorf("the contract must read the status code, but read: %d", built.GetStatusCode())
	}
}
