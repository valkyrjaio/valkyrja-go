/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package response_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
)

// newUri builds the URI that a redirect test sends the client to.
func newUri(t *testing.T, raw string) contract.UriContract {
	t.Helper()

	built, err := uri.NewUriFromString(raw)
	if err != nil {
		t.Fatalf("the parser must read %q, but reported: %v", raw, err)
	}

	return built
}

// newRequestWithReferer builds a request that names the referer.
func newRequestWithReferer(t *testing.T, referer string) contract.ServerRequestContract {
	t.Helper()

	var headers contract.HeaderCollectionContract = header.NewHeaderCollection()

	if referer != "" {
		built, err := header.NewHeader(constant.HeaderNameReferer, value.NewValueFromValue(referer))
		if err != nil {
			t.Fatalf("the header must be valid, but reported: %v", err)
		}

		headers = headers.WithHeader(built)
	}

	return request.NewServerRequest(
		newUri(t, "http://example.com/current"),
		constant.RequestMethodGet,
		nil,
		headers,
	)
}

func TestARedirectResponseSendsTheClientToTheUri(t *testing.T) {
	t.Parallel()

	target := newUri(t, "https://example.com/target")

	built := response.NewRedirectResponseToUri(target, 0, nil)

	if built.GetStatusCode() != constant.StatusCodeFound {
		t.Errorf("a redirect must take status 302 where the caller states none, but took: %d",
			built.GetStatusCode())
	}

	if built.GetUri() != target {
		t.Error("the response must carry the URI that it sends the client to, but did not")
	}

	if built.GetHeaders().GetHeaderLine(constant.HeaderNameLocation) == "" {
		t.Error("the response must name the target in its location header, but did not")
	}
}

func TestEachRedirectResponseWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := response.NewRedirectResponseToUri(newUri(t, "https://example.com/first"), 0, nil)
	other := newUri(t, "https://example.com/second")

	if built.WithUri(other).GetUri() != other {
		t.Error("WithUri must hold the new URI, but did not")
	}

	created := built.CreateFromUri(other, constant.StatusCodeMovedPermanently, nil)
	if created.GetStatusCode() != constant.StatusCodeMovedPermanently || created.GetUri() != other {
		t.Error("CreateFromUri must build a response for the URI and the status code, but did not")
	}

	if built.GetUri().GetPath() != "/first" {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestSecureSendsTheClientToThePathOverHttps(t *testing.T) {
	t.Parallel()

	built := response.NewRedirectResponseToUri(newUri(t, "http://example.com/first"), 0, nil)

	secure := built.Secure("/target", newRequestWithReferer(t, ""))

	if secure.GetUri().GetScheme() != constant.SchemeHttps {
		t.Errorf("Secure must send the client over HTTPS, but used: %q", secure.GetUri().GetScheme())
	}

	if secure.GetUri().GetPath() != "/target" {
		t.Errorf("Secure must send the client to the path, but used: %q", secure.GetUri().GetPath())
	}
}

func TestBackSendsTheClientToTheUriThatItCameFrom(t *testing.T) {
	t.Parallel()

	built := response.NewRedirectResponseToUri(newUri(t, "http://example.com/first"), 0, nil)

	back := built.Back(newRequestWithReferer(t, "https://example.com/previous"))

	if back.GetUri().GetPath() != "/previous" {
		t.Errorf("Back must send the client to the referer, but used: %q", back.GetUri().GetPath())
	}
}

func TestBackSendsTheClientToTheRootWhereTheRequestNamesNoReferer(t *testing.T) {
	t.Parallel()

	built := response.NewRedirectResponseToUri(newUri(t, "http://example.com/first"), 0, nil)

	tests := map[string]string{
		"a request that names no referer": "",
		"a referer that no parser reads":  "://example.com",
	}

	for name, referer := range tests {
		back := built.Back(newRequestWithReferer(t, referer))

		if back.GetUri().GetPath() != "/" {
			t.Errorf("%s must send the client to the root, but used: %q", name, back.GetUri().GetPath())
		}
	}
}
