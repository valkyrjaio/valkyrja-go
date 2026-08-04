/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// ResponseContract is an HTTP response.
//
// The other ports give each `create` method optional parameters. Go has no
// default parameter, so a caller passes every argument, and passes the zero
// value for one that it leaves to the response.
type ResponseContract interface {
	MessageContract

	// Create returns a response with the content, the status code, and the
	// headers.
	Create(content string, statusCode constant.StatusCode, headers HeaderCollectionContract) ResponseContract

	// GetStatusCode returns the status code of the response.
	GetStatusCode() constant.StatusCode

	// WithStatusCode returns a copy of the response for another status code.
	WithStatusCode(code constant.StatusCode) ResponseContract

	// GetReasonPhrase returns the reason phrase of the response.
	GetReasonPhrase() string

	// WithReasonPhrase returns a copy of the response for another reason
	// phrase.
	WithReasonPhrase(reasonPhrase string) ResponseContract

	// WithCookie returns a copy of the response that sets the cookie.
	WithCookie(cookie CookieContract) ResponseContract

	// WithoutCookie returns a copy of the response that removes the cookie.
	WithoutCookie(cookie CookieContract) ResponseContract
}

// TextResponseContract is a response that carries plain text.
//
// The contract is an alias, because it adds no method. Go compares an
// interface by its method set, so a second interface with the same methods is
// the same type. The TypeScript port declares it as an alias too.
type TextResponseContract = ResponseContract

// HtmlResponseContract is a response that carries HTML.
//
// The contract is an alias, because it adds no method. Go compares an
// interface by its method set, so a second interface with the same methods is
// the same type. The TypeScript port declares it as an alias too.
type HtmlResponseContract = ResponseContract

// EmptyResponseContract is a response that carries no body.
//
// The contract is an alias, because it adds no method. Go compares an
// interface by its method set, so a second interface with the same methods is
// the same type. The TypeScript port declares it as an alias too.
type EmptyResponseContract = ResponseContract

// JsonResponseContract is a response that carries JSON.
type JsonResponseContract interface {
	ResponseContract

	// CreateFromData returns a response that carries the data as JSON.
	CreateFromData(
		data map[string]any,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) JsonResponseContract

	// GetBodyAsJson returns the body of the response as parsed JSON.
	GetBodyAsJson() map[string]any

	// WithJsonAsBody returns a copy of the response that carries the data as
	// JSON.
	WithJsonAsBody(data map[string]any) JsonResponseContract

	// WithCallback returns a copy of the response that wraps the JSON in a
	// callback, which makes it JSONP.
	WithCallback(callback string) JsonResponseContract

	// WithoutCallback returns a copy of the response with no callback.
	WithoutCallback() JsonResponseContract
}

// RedirectResponseContract is a response that sends the client to another URI.
type RedirectResponseContract interface {
	ResponseContract

	// CreateFromUri returns a response that sends the client to the URI.
	CreateFromUri(
		uri UriContract,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) RedirectResponseContract

	// GetUri returns the URI that the response sends the client to.
	GetUri() UriContract

	// WithUri returns a copy of the response for another URI.
	WithUri(uri UriContract) RedirectResponseContract

	// Secure returns a copy of the response that sends the client to the path
	// over HTTPS.
	Secure(path string, request ServerRequestContract) RedirectResponseContract

	// Back returns a copy of the response that sends the client to the URI that
	// it came from.
	Back(request ServerRequestContract) RedirectResponseContract
}

// ResponseFactoryContract builds each kind of response.
type ResponseFactoryContract interface {
	// CreateResponse builds a response that carries the content.
	CreateResponse(
		content string,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) ResponseContract

	// CreateTextResponse builds a response that carries plain text.
	CreateTextResponse(
		content string,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) TextResponseContract

	// CreateJsonResponse builds a response that carries the data as JSON.
	CreateJsonResponse(
		data map[string]any,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) JsonResponseContract

	// CreateJsonpResponse builds a response that wraps the JSON in the
	// callback.
	CreateJsonpResponse(
		callback string,
		data map[string]any,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) JsonResponseContract

	// CreateRedirectResponse builds a response that sends the client to the
	// URI.
	CreateRedirectResponse(
		uri string,
		statusCode constant.StatusCode,
		headers HeaderCollectionContract,
	) RedirectResponseContract
}
