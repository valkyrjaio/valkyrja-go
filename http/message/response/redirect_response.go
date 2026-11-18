/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package response

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
)

type RedirectResponse struct {
	Response

	uri contract.UriContract
}

// NewRedirectResponseToUri builds a response that sends the client to the URI,
// at status 302 where the caller states no status code.
func NewRedirectResponseToUri(
	redirectUri contract.UriContract,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *RedirectResponse {
	return &RedirectResponse{
		Response: *NewRedirectResponse(redirectUri, statusCode, headers),
		uri:      redirectUri,
	}
}

// CreateFromUri returns a response that sends the client to the URI.
func (r *RedirectResponse) CreateFromUri(
	redirectUri contract.UriContract,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.RedirectResponseContract {
	return NewRedirectResponseToUri(redirectUri, statusCode, headers)
}

// GetUri returns the URI that the response sends the client to.
func (r *RedirectResponse) GetUri() contract.UriContract {
	return r.uri
}

// WithUri returns a copy of the response for another URI.
func (r *RedirectResponse) WithUri(redirectUri contract.UriContract) contract.RedirectResponseContract {
	return NewRedirectResponseToUri(redirectUri, r.GetStatusCode(), r.GetHeaders())
}

// Secure returns a copy of the response that sends the client to the path over
// HTTPS, on the host that the request names.
func (r *RedirectResponse) Secure(
	path string,
	request contract.ServerRequestContract,
) contract.RedirectResponseContract {
	secure := request.GetUri().
		WithScheme(constant.SchemeHttps).
		WithPath(path)

	return r.WithUri(secure)
}

// Back returns a copy of the response that sends the client to the URI that it
// came from.
func (r *RedirectResponse) Back(request contract.ServerRequestContract) contract.RedirectResponseContract {
	referer := request.GetHeaders().GetHeaderLine(constant.HeaderNameReferer)

	if referer == "" {
		return r.WithUri(rootUri())
	}

	back, err := uri.NewUriFromString(referer)
	if err != nil {
		return r.WithUri(rootUri())
	}

	return r.WithUri(back)
}

// rootUri returns the URI of the root path.
func rootUri() contract.UriContract {
	built, _ := uri.NewUri("", "", "", "", 0, "/", "", "")

	return built
}
