/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package response

import (
	"encoding/json"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
)

// emptyJsonBody is the body of a JSON response that carries no data.
const emptyJsonBody = "{}"

// NewTextResponse builds a response that carries plain text. It sets the content
// type, and it replaces a content type that the headers carry already.
func NewTextResponse(
	text string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *Response {
	return NewResponseFromContent(
		text,
		statusCode,
		withContentType(headers, constant.ContentTypeValueTextPlainUtf8),
	)
}

// NewHtmlResponse builds a response that carries HTML.
func NewHtmlResponse(
	html string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *Response {
	return NewResponseFromContent(
		html,
		statusCode,
		withContentType(headers, constant.ContentTypeValueTextHtmlUtf8),
	)
}

// NewEmptyResponse builds a response that carries no body, at status 204.
func NewEmptyResponse(headers contract.HeaderCollectionContract) *Response {
	return NewResponse(
		stream.NewStream("", constant.ModeRead),
		constant.StatusCodeNoContent,
		headers,
	)
}

// NewJsonResponse builds a response that carries the data as JSON.
func NewJsonResponse(
	data map[string]any,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *Response {
	return NewResponseFromContent(
		encodeJson(data),
		statusCode,
		withContentType(headers, constant.ContentTypeValueApplicationJson),
	)
}

// NewJsonpResponse builds a response that wraps the JSON in the callback, which
// makes it JSONP.
func NewJsonpResponse(
	callback string,
	data map[string]any,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *Response {
	return NewResponseFromContent(
		callback+"("+encodeJson(data)+")",
		statusCode,
		withContentType(headers, constant.ContentTypeValueTextJavascript),
	)
}

// NewRedirectResponse builds a response that sends the client to the URI, at
// status 302.
func NewRedirectResponse(
	redirectUri contract.UriContract,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *Response {
	if statusCode == 0 {
		statusCode = constant.StatusCodeFound
	}

	return NewResponse(
		stream.NewStream("", constant.ModeReadWrite),
		statusCode,
		withLocation(headers, redirectUri),
	)
}

// encodeJson renders the data as JSON, and an empty object where no encoder can
// render it.
func encodeJson(data map[string]any) string {
	if data == nil {
		return emptyJsonBody
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return emptyJsonBody
	}

	return string(encoded)
}

// withContentType returns the headers with the content type in them. It replaces
// a content type that the headers carry already, because a response states one
// content type.
func withContentType(
	headers contract.HeaderCollectionContract,
	contentType string,
) contract.HeaderCollectionContract {
	if headers == nil {
		headers = header.NewHeaderCollection()
	}

	// `NewHeader` reports a failure only for a name that a header cannot carry,
	// and the name here is a constant that always can.
	contentTypeHeader, _ := header.NewHeader(
		constant.HeaderNameContentType,
		value.NewValueFromValue(contentType),
	)

	return message.InjectHeader(contentTypeHeader, headers, true)
}

// withLocation returns the headers with the location of the URI in them.
func withLocation(
	headers contract.HeaderCollectionContract,
	redirectUri contract.UriContract,
) contract.HeaderCollectionContract {
	if headers == nil {
		headers = header.NewHeaderCollection()
	}

	location := ""

	if redirectUri != nil {
		location = redirectUri.String()
	}

	// `NewHeader` reports a failure only for a name that a header cannot carry,
	// and the name here is a constant that always can.
	locationHeader, _ := header.NewHeader(
		constant.HeaderNameLocation,
		value.NewValueFromValue(location),
	)

	return message.InjectHeader(locationHeader, headers, true)
}
