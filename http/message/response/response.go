/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package response holds the HTTP response and each kind of response that the
// framework builds.
package response

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
)

// Response is an HTTP response.
type Response struct {
	message.Message

	statusCode   constant.StatusCode
	statusPhrase string
}

// NewResponse builds a response over a body, a status code, and headers. It
// takes the defaults of the other ports where an argument is nil or zero: an
// empty body, status 200, and no header.
func NewResponse(
	body contract.StreamContract,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *Response {
	if statusCode == 0 {
		statusCode = constant.StatusCodeOk
	}

	built := &Response{
		Message:      message.NewMessage("", headers, body),
		statusCode:   statusCode,
		statusPhrase: statusCode.GetText(),
	}

	built.SetBody(built.GetBody())

	return built
}

// NewResponseFromContent builds a response that carries the content.
func NewResponseFromContent(
	content string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *Response {
	body := stream.NewStream("", constant.ModeReadWrite)

	// The stream opens in a mode that a writer writes, so the write cannot fail.
	_, _ = body.Write(content)

	return NewResponse(body, statusCode, headers)
}

// Create returns a response that carries the content, at the status code, with
// the headers.
func (r *Response) Create(
	content string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.ResponseContract {
	return NewResponseFromContent(content, statusCode, headers)
}

// GetStatusCode returns the status code of the response.
func (r *Response) GetStatusCode() constant.StatusCode {
	return r.statusCode
}

// WithStatusCode returns a copy of the response for another status code. The
// reason phrase follows the code.
func (r *Response) WithStatusCode(code constant.StatusCode) contract.ResponseContract {
	copied := *r
	copied.statusCode = code
	copied.statusPhrase = code.GetText()

	return &copied
}

// GetReasonPhrase returns the reason phrase of the response.
func (r *Response) GetReasonPhrase() string {
	return r.statusPhrase
}

// WithReasonPhrase returns a copy of the response for another reason phrase. An
// empty phrase takes the phrase of the status code.
func (r *Response) WithReasonPhrase(reasonPhrase string) contract.ResponseContract {
	copied := *r

	if reasonPhrase == "" {
		reasonPhrase = r.statusCode.GetText()
	}

	copied.statusPhrase = reasonPhrase

	return &copied
}

// WithProtocolVersion returns a copy of the response for another protocol
// version.
func (r *Response) WithProtocolVersion(version constant.ProtocolVersion) contract.MessageContract {
	copied := *r
	copied.SetProtocolVersion(version)

	return &copied
}

// WithHeaders returns a copy of the response with other headers.
func (r *Response) WithHeaders(headers contract.HeaderCollectionContract) contract.MessageContract {
	copied := *r
	copied.SetHeaders(headers)

	return &copied
}

// WithBody returns a copy of the response with another body.
func (r *Response) WithBody(body contract.StreamContract) contract.MessageContract {
	copied := *r
	copied.SetBody(body)

	return &copied
}

// WithCookie returns a copy of the response that sets the cookie.
func (r *Response) WithCookie(cookie contract.CookieContract) contract.ResponseContract {
	return r.withSetCookie(cookie)
}

// WithoutCookie returns a copy of the response that removes the cookie.
func (r *Response) WithoutCookie(cookie contract.CookieContract) contract.ResponseContract {
	return r.withSetCookie(cookie.Delete())
}

// withSetCookie returns a copy of the response with the cookie added to the
// `Set-Cookie` header.
func (r *Response) withSetCookie(cookie contract.CookieContract) contract.ResponseContract {
	// `NewHeader` reports a failure only for a name that a header cannot carry,
	// and the name here is a constant that always can.
	setCookie, _ := header.NewHeader(constant.HeaderNameSetCookie, cookie)

	copied := *r
	copied.SetHeaders(r.GetHeaders().WithAddedHeaders(setCookie))

	return &copied
}
