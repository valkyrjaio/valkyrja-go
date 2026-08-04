/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package factory_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/factory"
)

const responseContent = "the content"

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

func TestTheFactoryBuildsEachKindOfResponse(t *testing.T) {
	t.Parallel()

	built := factory.NewResponseFactory()

	plain := built.CreateResponse(responseContent, constant.StatusCodeOk, nil)
	if readBody(t, plain) != responseContent {
		t.Error("the response must carry the content, but did not")
	}

	text := built.CreateTextResponse(responseContent, constant.StatusCodeOk, nil)
	if text.GetHeaders().GetHeaderLine(constant.HeaderNameContentType) == "" {
		t.Error("a text response must state its content type, but did not")
	}

	html := built.CreateHtmlResponse(responseContent, constant.StatusCodeOk, nil)
	if readBody(t, html) != responseContent {
		t.Error("an HTML response must carry the content, but did not")
	}

	empty := built.CreateEmptyResponse(nil)
	if empty.GetStatusCode() != constant.StatusCodeNoContent {
		t.Errorf("an empty response must carry status 204, but carried: %d", empty.GetStatusCode())
	}
}

func TestTheFactoryBuildsEachKindOfJsonResponse(t *testing.T) {
	t.Parallel()

	built := factory.NewResponseFactory()
	data := map[string]any{"one": "two"}

	asJson := built.CreateJsonResponse(data, constant.StatusCodeOk, nil)
	if asJson.GetBodyAsJson()["one"] != "two" {
		t.Error("the response must carry the data as JSON, but did not")
	}

	asJsonp := built.CreateJsonpResponse("handle", data, constant.StatusCodeOk, nil)
	if readBody(t, asJsonp) != `handle({"one":"two"})` {
		t.Errorf("the response must wrap the JSON in the callback, but carried: %q", readBody(t, asJsonp))
	}
}

func TestTheFactoryBuildsARedirectResponse(t *testing.T) {
	t.Parallel()

	built := factory.NewResponseFactory()

	redirect := built.CreateRedirectResponse("https://example.com/target", 0, nil)
	if redirect.GetUri().GetPath() != "/target" {
		t.Errorf("the response must send the client to the target, but used: %q", redirect.GetUri().GetPath())
	}

	// A target that no parser reads sends the client to the root, because a
	// redirect with no target sends the client nowhere.
	broken := built.CreateRedirectResponse("://example.com", 0, nil)
	if broken.GetUri().GetPath() != "/" {
		t.Errorf("a target that no parser reads must send the client to the root, but used: %q",
			broken.GetUri().GetPath())
	}
}

func TestTheFactorySatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.ResponseFactoryContract = factory.NewResponseFactory()

	if built.CreateResponse("", constant.StatusCodeOk, nil).GetStatusCode() != constant.StatusCodeOk {
		t.Error("the factory must satisfy its contract, but did not")
	}
}
