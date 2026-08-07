/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package factory builds each kind of response that a handler returns.
package factory

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
)

type ResponseFactory struct{}

// NewResponseFactory builds the factory.
func NewResponseFactory() *ResponseFactory {
	return &ResponseFactory{}
}

// CreateResponse builds a response that carries the content.
func (f *ResponseFactory) CreateResponse(
	content string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.ResponseContract {
	return response.NewResponseFromContent(content, statusCode, headers)
}

// CreateTextResponse builds a response that carries plain text.
func (f *ResponseFactory) CreateTextResponse(
	content string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.TextResponseContract {
	return response.NewTextResponse(content, statusCode, headers)
}

// CreateHtmlResponse builds a response that carries HTML.
func (f *ResponseFactory) CreateHtmlResponse(
	content string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.HtmlResponseContract {
	return response.NewHtmlResponse(content, statusCode, headers)
}

// CreateEmptyResponse builds a response that carries nothing.
func (f *ResponseFactory) CreateEmptyResponse(
	headers contract.HeaderCollectionContract,
) contract.EmptyResponseContract {
	return response.NewEmptyResponse(headers)
}

// CreateJsonResponse builds a response that carries the data as JSON.
func (f *ResponseFactory) CreateJsonResponse(
	data map[string]any,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.JsonResponseContract {
	return response.NewJsonResponseFromData(data, statusCode, headers)
}

// CreateJsonpResponse builds a response that wraps the JSON in the callback.
func (f *ResponseFactory) CreateJsonpResponse(
	callback string,
	data map[string]any,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.JsonResponseContract {
	return response.NewJsonpResponseFromData(callback, data, statusCode, headers)
}

// CreateRedirectResponse builds a response that sends the client to the URI.
func (f *ResponseFactory) CreateRedirectResponse(
	target string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.RedirectResponseContract {
	parsed, err := uri.NewUriFromString(target)
	if err != nil {
		return response.NewRedirectResponseToUri(rootUri(), statusCode, headers)
	}

	return response.NewRedirectResponseToUri(parsed, statusCode, headers)
}

// rootUri returns the URI of the root path.
func rootUri() contract.UriContract {
	built, _ := uri.NewUri("", "", "", "", 0, "/", "", "")

	return built
}
