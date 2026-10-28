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
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
)

func TestEachTypedResponseSetsItsContentType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		built *response.Response
		want  string
	}{
		"text": {response.NewTextResponse("the text", 0, nil), constant.ContentTypeValueTextPlainUtf8},
		"html": {response.NewHtmlResponse("<p></p>", 0, nil), constant.ContentTypeValueTextHtmlUtf8},
		"json": {
			response.NewJsonResponse(map[string]any{"a": 1}, 0, nil),
			constant.ContentTypeValueApplicationJson,
		},
		"jsonp": {
			response.NewJsonpResponse("callback", map[string]any{"a": 1}, 0, nil),
			constant.ContentTypeValueTextJavascript,
		},
	}

	for name, test := range tests {
		line := test.built.GetHeaders().GetHeaderLine(constant.HeaderNameContentType)

		if line != test.want {
			t.Errorf("the %s response must set %q, but set %q", name, test.want, line)
		}
	}
}

func TestATypedResponseReplacesAContentTypeThatTheHeadersCarry(t *testing.T) {
	t.Parallel()

	existing, err := header.NewHeader(
		constant.HeaderNameContentType,
		value.NewValueFromValue("text/html"),
	)
	if err != nil {
		t.Fatalf("NewHeader must build the header, but reported: %v", err)
	}

	built := response.NewTextResponse("the text", 0, header.NewHeaderCollection(existing))

	line := built.GetHeaders().GetHeaderLine(constant.HeaderNameContentType)

	if line != constant.ContentTypeValueTextPlainUtf8 {
		t.Errorf("the response must replace the content type, but it is: %q", line)
	}
}

func TestTheTextResponseCarriesItsText(t *testing.T) {
	t.Parallel()

	built := response.NewTextResponse("the text", constant.StatusCodeCreated, nil)

	if built.GetBody().String() != "the text" {
		t.Errorf("the body must carry the text, but is: %q", built.GetBody().String())
	}

	if built.GetStatusCode() != constant.StatusCodeCreated {
		t.Errorf("the status code must be the one given, but is: %d", built.GetStatusCode())
	}
}

func TestTheEmptyResponseCarriesNothing(t *testing.T) {
	t.Parallel()

	built := response.NewEmptyResponse(nil)

	if built.GetStatusCode() != constant.StatusCodeNoContent {
		t.Errorf("the status code must be 204, but is: %d", built.GetStatusCode())
	}

	if built.GetBody().GetSize() != 0 {
		t.Error("the body must be empty, but is not")
	}

	if built.GetBody().IsWritable() {
		t.Error("no writer must write the body, but one does")
	}
}

func TestTheJsonResponseRendersTheData(t *testing.T) {
	t.Parallel()

	built := response.NewJsonResponse(map[string]any{"name": "valkyrja"}, 0, nil)

	if built.GetBody().String() != `{"name":"valkyrja"}` {
		t.Errorf("the body must render the data, but is: %q", built.GetBody().String())
	}
}

func TestTheJsonResponseRendersAnEmptyObjectForNoData(t *testing.T) {
	t.Parallel()

	if response.NewJsonResponse(nil, 0, nil).GetBody().String() != "{}" {
		t.Error("the body must render an empty object for no data, but did not")
	}
}

func TestTheJsonResponseRendersAnEmptyObjectForDataThatNoEncoderRenders(t *testing.T) {
	t.Parallel()

	built := response.NewJsonResponse(map[string]any{"channel": make(chan int)}, 0, nil)

	if built.GetBody().String() != "{}" {
		t.Errorf("the body must render an empty object, but is: %q", built.GetBody().String())
	}
}

func TestTheJsonpResponseWrapsTheJsonInTheCallback(t *testing.T) {
	t.Parallel()

	built := response.NewJsonpResponse("handle", map[string]any{"a": float64(1)}, 0, nil)

	if built.GetBody().String() != `handle({"a":1})` {
		t.Errorf("the body must wrap the JSON in the callback, but is: %q", built.GetBody().String())
	}
}

func TestTheRedirectResponseSetsTheLocation(t *testing.T) {
	t.Parallel()

	target, err := uri.NewUri(constant.SchemeHttps, "", "", "example.com", 0, "/path", "", "")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	built := response.NewRedirectResponse(target, 0, nil)

	if built.GetStatusCode() != constant.StatusCodeFound {
		t.Errorf("the status code must default to 302, but is: %d", built.GetStatusCode())
	}

	line := built.GetHeaders().GetHeaderLine(constant.HeaderNameLocation)

	if !strings.Contains(line, "example.com") {
		t.Errorf("the location must name the URI, but is: %q", line)
	}
}

func TestTheRedirectResponseTakesTheStatusCodeThatItReceives(t *testing.T) {
	t.Parallel()

	built := response.NewRedirectResponse(nil, constant.StatusCodeMovedPermanently, nil)

	if built.GetStatusCode() != constant.StatusCodeMovedPermanently {
		t.Errorf("the status code must be the one given, but is: %d", built.GetStatusCode())
	}

	if built.GetHeaders().GetHeaderLine(constant.HeaderNameLocation) != "" {
		t.Error("a redirect with no URI must set an empty location, but did not")
	}
}

func TestEachTypedResponseSatisfiesTheResponseContract(t *testing.T) {
	t.Parallel()

	responses := []contract.ResponseContract{
		response.NewTextResponse("", 0, nil),
		response.NewHtmlResponse("", 0, nil),
		response.NewEmptyResponse(nil),
		response.NewJsonResponse(nil, 0, nil),
		response.NewRedirectResponse(nil, 0, nil),
	}

	for _, built := range responses {
		if built.GetStatusCode() == 0 {
			t.Error("each response must carry a status code, but one did not")
		}
	}
}
