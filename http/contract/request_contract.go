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

type RequestContract interface {
	MessageContract

	// GetRequestTarget returns the target of the request.
	GetRequestTarget() string

	// WithRequestTarget returns a copy of the request for another target.
	WithRequestTarget(requestTarget string) RequestContract

	// GetMethod returns the request method.
	GetMethod() constant.RequestMethod

	// WithMethod returns a copy of the request for another method.
	WithMethod(method constant.RequestMethod) RequestContract

	// GetUri returns the URI of the request.
	GetUri() UriContract

	// WithUri returns a copy of the request for another URI. It keeps the host
	// header where preserveHost is true.
	WithUri(uri UriContract, preserveHost bool) RequestContract
}

//nolint:interfacebloat // Parity with the PHP reference implementation.
type ServerRequestContract interface {
	RequestContract

	// GetServerParams returns the server parameters of the request.
	GetServerParams() ServerParamCollectionContract

	// WithServerParams returns a copy of the request with other server
	// parameters.
	WithServerParams(server ServerParamCollectionContract) ServerRequestContract

	// GetCookieParams returns the cookies of the request.
	GetCookieParams() CookieParamCollectionContract

	// WithCookieParams returns a copy of the request with other cookies.
	WithCookieParams(cookies CookieParamCollectionContract) ServerRequestContract

	// GetQueryParams returns the query parameters of the request.
	GetQueryParams() QueryParamCollectionContract

	// WithQueryParams returns a copy of the request with other query
	// parameters.
	WithQueryParams(query QueryParamCollectionContract) ServerRequestContract

	// GetUploadedFiles returns the files that arrived with the request.
	GetUploadedFiles() UploadedFileCollectionContract

	// WithUploadedFiles returns a copy of the request with other files.
	WithUploadedFiles(uploadedFiles UploadedFileCollectionContract) ServerRequestContract

	// GetParsedBody returns the parsed body of the request.
	GetParsedBody() ParsedBodyParamCollectionContract

	// WithParsedBody returns a copy of the request with another parsed body.
	WithParsedBody(params ParsedBodyParamCollectionContract) ServerRequestContract

	// GetAttributes returns the attributes that the framework put on the
	// request.
	GetAttributes() AttributeParamCollectionContract

	// WithAttributes returns a copy of the request with other attributes.
	WithAttributes(attributes AttributeParamCollectionContract) ServerRequestContract

	// IsXmlHttpRequest reports whether the client sent the request from a
	// script.
	IsXmlHttpRequest() bool

	// GetParsedJson returns the parsed JSON body of the request, and an empty
	// collection where the request carries none.
	GetParsedJson() ParsedJsonParamCollectionContract

	// WithParsedJson returns a copy of the request with another parsed JSON
	// body.
	WithParsedJson(params ParsedJsonParamCollectionContract) JsonServerRequestContract
}

type JsonServerRequestContract = ServerRequestContract
