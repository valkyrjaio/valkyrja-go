/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package request holds the HTTP request and the server request.
package request

import (
	"strconv"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

// rootTarget is the request target of a request whose URI names no path.
const rootTarget = "/"

type Request struct {
	message.Message

	requestUri    contract.UriContract
	method        constant.RequestMethod
	requestTarget string
}

// NewRequest builds a request over a URI, a method, a body, and headers. It
// takes the defaults of the other ports where an argument is nil or empty: an
// empty URI, the GET method, an empty body, and no header.
func NewRequest(
	requestUri contract.UriContract,
	method constant.RequestMethod,
	body contract.StreamContract,
	headers contract.HeaderCollectionContract,
) *Request {
	if requestUri == nil {
		requestUri = newEmptyUri()
	}

	if method == "" {
		method = constant.RequestMethodGet
	}

	built := &Request{
		Message:    message.NewMessage("", headers, body),
		requestUri: requestUri,
		method:     method,
	}

	built.SetBody(built.GetBody())
	built.addHostHeaderFromUri()

	return built
}

// GetRequestTarget returns the target of the request. It reads the target from
// the URI where no caller set one.
func (r *Request) GetRequestTarget() string {
	if r.requestTarget != "" {
		return r.requestTarget
	}

	target := r.requestUri.GetPath()

	if query := r.requestUri.GetQuery(); query != "" {
		target += "?" + query
	}

	if target == "" {
		return rootTarget
	}

	return target
}

// WithRequestTarget returns a copy of the request for another target. It keeps
// the target of the receiver where the new one carries whitespace, which no
// request target does.
func (r *Request) WithRequestTarget(requestTarget string) contract.RequestContract {
	if ValidateRequestTarget(requestTarget) != nil {
		return r
	}

	copied := *r
	copied.requestTarget = requestTarget

	return &copied
}

// GetMethod returns the request method.
func (r *Request) GetMethod() constant.RequestMethod {
	return r.method
}

// WithMethod returns a copy of the request for another method.
func (r *Request) WithMethod(method constant.RequestMethod) contract.RequestContract {
	copied := *r
	copied.method = method

	return &copied
}

// GetUri returns the URI of the request.
func (r *Request) GetUri() contract.UriContract {
	return r.requestUri
}

// WithUri returns a copy of the request for another URI.
func (r *Request) WithUri(requestUri contract.UriContract, preserveHost bool) contract.RequestContract {
	copied := *r
	copied.requestUri = requestUri

	if preserveHost && r.GetHeaders().Has(constant.HeaderNameHost) {
		return &copied
	}

	if requestUri.GetHost() == "" {
		return &copied
	}

	copied.SetHeaders(r.GetHeaders().WithHeader(newHostHeader(requestUri)))

	return &copied
}

// WithProtocolVersion returns a copy of the request for another protocol
// version.
func (r *Request) WithProtocolVersion(version constant.ProtocolVersion) contract.MessageContract {
	copied := *r
	copied.SetProtocolVersion(version)

	return &copied
}

// WithHeaders returns a copy of the request with other headers.
func (r *Request) WithHeaders(headers contract.HeaderCollectionContract) contract.MessageContract {
	copied := *r
	copied.SetHeaders(headers)

	return &copied
}

// WithBody returns a copy of the request with another body.
func (r *Request) WithBody(body contract.StreamContract) contract.MessageContract {
	copied := *r
	copied.SetBody(body)

	return &copied
}

// addHostHeaderFromUri adds a `Host` header from the URI, where the headers
// carry none and the URI names a host.
func (r *Request) addHostHeaderFromUri() {
	if r.GetHeaders().Has(constant.HeaderNameHost) {
		return
	}

	if r.requestUri.GetHost() == "" {
		return
	}

	r.SetHeaders(r.GetHeaders().WithHeader(newHostHeader(r.requestUri)))
}

// ValidateRequestTarget reports a failure where the target carries whitespace,
// which no request target does.
func ValidateRequestTarget(requestTarget string) error {
	if strings.ContainsAny(requestTarget, " \t\r\n\v\f") {
		return exception.NewHttpRequestInvalidRequestTargetError(requestTarget)
	}

	return nil
}

// newHostHeader builds the `Host` header of the URI, with the port where the URI
// names one that is not the standard port of its scheme.
func newHostHeader(requestUri contract.UriContract) contract.HeaderContract {
	host := requestUri.GetHost()

	if port := requestUri.GetPort(); port != 0 {
		host += ":" + strconv.Itoa(port)
	}

	// `NewHeader` reports a failure only for a name that a header cannot carry,
	// and the name here is a constant that always can.
	built, _ := header.NewHeader(constant.HeaderNameHost, value.NewValueFromValue(host))

	return built
}

// newEmptyUri builds the URI that a request takes where a caller names none.
func newEmptyUri() contract.UriContract {
	// Every part is empty, so no part can be one that a URI rejects.
	built, _ := uri.NewUri(constant.SchemeEmpty, "", "", "", 0, "", "", "")

	return built
}
