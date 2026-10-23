/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package request_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

const (
	exampleHost = "example.com"
	basicPath   = "/path"
)

// newUri builds a URI and fails the test where a part is invalid.
func newUri(t *testing.T, host string, port int, path string, query string) contract.UriContract {
	t.Helper()

	built, err := uri.NewUri(constant.SchemeHttp, "", "", host, port, path, query, "")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	return built
}

func TestNewRequestTakesEachDefault(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(nil, "", nil, nil)

	if built.GetMethod() != constant.RequestMethodGet {
		t.Errorf("the method must default to GET, but is: %q", built.GetMethod())
	}

	if built.GetRequestTarget() != "/" {
		t.Errorf("the request target must default to the root, but is: %q", built.GetRequestTarget())
	}

	if built.GetUri().GetHost() != "" {
		t.Error("the URI must default to an empty one, but did not")
	}
}

func TestNewRequestAddsTheHostHeaderFromTheUri(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 8080, "", ""), "", nil, nil)

	if built.GetHeaders().GetHeaderLine(constant.HeaderNameHost) != "example.com:8080" {
		t.Errorf("the host header must follow the URI, but is: %q",
			built.GetHeaders().GetHeaderLine(constant.HeaderNameHost))
	}
}

func TestNewRequestLeavesOutTheHostHeaderWithNoHost(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, "", 0, "", ""), "", nil, nil)

	if built.GetHeaders().Has(constant.HeaderNameHost) {
		t.Error("a URI with no host must add no host header, but did")
	}
}

func TestNewRequestKeepsAHostHeaderThatTheHeadersCarry(t *testing.T) {
	t.Parallel()

	hostHeader, err := header.NewHeader(constant.HeaderNameHost, value.NewValueFromValue("other.com"))
	if err != nil {
		t.Fatalf("NewHeader must build the header, but reported: %v", err)
	}

	built := request.NewRequest(
		newUri(t, exampleHost, 0, "", ""),
		"",
		nil,
		header.NewHeaderCollection(hostHeader),
	)

	if built.GetHeaders().GetHeaderLine(constant.HeaderNameHost) != "other.com" {
		t.Error("the request must keep a host header that the headers carry, but did not")
	}
}

func TestGetRequestTargetReadsThePathAndTheQuery(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 0, basicPath, "a=b"), "", nil, nil)

	if built.GetRequestTarget() != "/path?a=b" {
		t.Errorf("the request target must join the path and the query, but is: %q", built.GetRequestTarget())
	}
}

func TestGetRequestTargetReadsThePathAlone(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 0, basicPath, ""), "", nil, nil)

	if built.GetRequestTarget() != basicPath {
		t.Errorf("the request target must be the path, but is: %q", built.GetRequestTarget())
	}
}

func TestWithRequestTargetHoldsTheNewTarget(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 0, basicPath, ""), "", nil, nil)

	if built.WithRequestTarget("/other").GetRequestTarget() != "/other" {
		t.Error("WithRequestTarget must hold the new target, but did not")
	}

	if built.GetRequestTarget() != basicPath {
		t.Error("WithRequestTarget must leave the receiver unchanged, but did not")
	}
}

func TestWithRequestTargetKeepsTheTargetWhereTheNewOneCarriesWhitespace(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 0, basicPath, ""), "", nil, nil)

	if built.WithRequestTarget("/a target").GetRequestTarget() != basicPath {
		t.Error("WithRequestTarget must keep the target where the new one is invalid, but did not")
	}
}

func TestValidateRequestTargetReportsWhitespace(t *testing.T) {
	t.Parallel()

	target, found := errors.AsType[*exception.HttpRequestInvalidRequestTargetError](
		request.ValidateRequestTarget("/a target"),
	)
	if !found {
		t.Fatal("ValidateRequestTarget must report whitespace, but did not")
	}

	if target.GetRequestTarget() != "/a target" {
		t.Errorf("the error must carry the target, but carries: %q", target.GetRequestTarget())
	}

	if request.ValidateRequestTarget(basicPath) != nil {
		t.Error("ValidateRequestTarget must accept a target with no whitespace, but did not")
	}
}

func TestWithMethodHoldsTheNewMethod(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(nil, "", nil, nil)

	if built.WithMethod(constant.RequestMethodPost).GetMethod() != constant.RequestMethodPost {
		t.Error("WithMethod must hold the new method, but did not")
	}

	if built.GetMethod() != constant.RequestMethodGet {
		t.Error("WithMethod must leave the receiver unchanged, but did not")
	}
}

func TestWithUriMovesTheHostHeader(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 0, "", ""), "", nil, nil)

	changed := built.WithUri(newUri(t, "other.com", 0, "", ""), false)

	if changed.GetHeaders().GetHeaderLine(constant.HeaderNameHost) != "other.com" {
		t.Error("WithUri must move the host header, but did not")
	}
}

func TestWithUriKeepsTheHostHeaderWherePreserveHostIsTrue(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 0, "", ""), "", nil, nil)

	changed := built.WithUri(newUri(t, "other.com", 0, "", ""), true)

	if changed.GetHeaders().GetHeaderLine(constant.HeaderNameHost) != exampleHost {
		t.Error("WithUri must keep the host header where preserveHost is true, but did not")
	}
}

func TestWithUriKeepsTheHostHeaderWhereTheNewUriNamesNoHost(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(newUri(t, exampleHost, 0, "", ""), "", nil, nil)

	changed := built.WithUri(newUri(t, "", 0, basicPath, ""), false)

	if changed.GetHeaders().GetHeaderLine(constant.HeaderNameHost) != exampleHost {
		t.Error("WithUri must keep the host header where the new URI names no host, but did not")
	}

	if changed.GetUri().GetPath() != basicPath {
		t.Error("WithUri must hold the new URI, but did not")
	}
}

func TestEachMessageWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := request.NewRequest(nil, "", nil, nil)

	if built.WithProtocolVersion(constant.ProtocolVersionV2).GetProtocolVersion() != constant.ProtocolVersionV2 {
		t.Error("WithProtocolVersion must hold the new version, but did not")
	}

	if built.WithHeaders(header.NewHeaderCollection()).GetProtocolVersion() != constant.ProtocolVersionV11 {
		t.Error("WithHeaders must keep the other state, but did not")
	}

	if built.WithBody(built.GetBody()).GetBody() != built.GetBody() {
		t.Error("WithBody must hold the new body, but did not")
	}

	if built.GetProtocolVersion() != constant.ProtocolVersionV11 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestTheRequestSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.RequestContract = request.NewRequest(nil, "", nil, nil)

	if built.GetMethod() != constant.RequestMethodGet {
		t.Errorf("the contract must read the method, but read: %q", built.GetMethod())
	}
}
