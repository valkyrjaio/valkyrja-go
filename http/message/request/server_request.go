/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package request

import (
	"encoding/json"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/param"
)

// xmlHttpRequestValue is the value that a script sets on the requested-with
// header, which is how a server tells a script's request from a browser's.
const xmlHttpRequestValue = "XMLHttpRequest"

// ServerRequest is an HTTP request that a server received, with everything that
// the server read from the connection.
type ServerRequest struct {
	Request

	serverParams  contract.ServerParamCollectionContract
	cookieParams  contract.CookieParamCollectionContract
	queryParams   contract.QueryParamCollectionContract
	uploadedFiles contract.UploadedFileCollectionContract
	parsedBody    contract.ParsedBodyParamCollectionContract
	parsedJson    contract.ParsedJsonParamCollectionContract
	attributes    contract.AttributeParamCollectionContract
}

// NewServerRequest builds a server request over a request. Each collection takes
// an empty one where the caller names none.
func NewServerRequest(
	requestUri contract.UriContract,
	method constant.RequestMethod,
	body contract.StreamContract,
	headers contract.HeaderCollectionContract,
) *ServerRequest {
	return &ServerRequest{
		Request:      *NewRequest(requestUri, method, body, headers),
		serverParams: param.NewParamCollection(nil),
		cookieParams: param.NewParamCollection(nil),
		queryParams:  param.NewParamCollection(nil),
		parsedBody:   param.NewParamCollection(nil),
		parsedJson:   param.NewParamCollection(nil),
		attributes:   param.NewParamCollection(nil),
	}
}

// NewJsonServerRequest builds a server request that parses its body as JSON.
//
// The other ports declare a `JsonServerRequest` that extends `ServerRequest`. A
// `With` method promoted from an embedded struct copies only that struct, so an
// inherited `With` would return a plain server request and drop the parsed JSON.
// One struct carries both shapes here, and this constructor is what parses.
//
// The body is parsed as the request is built, where the content type states
// JSON. A body that no decoder reads leaves the parsed JSON empty, because a
// request has no return to carry the failure — whatever reads the body reports
// what is missing instead.
func NewJsonServerRequest(
	requestUri contract.UriContract,
	method constant.RequestMethod,
	body contract.StreamContract,
	headers contract.HeaderCollectionContract,
) *ServerRequest {
	built := NewServerRequest(requestUri, method, body, headers)
	built.parsedJson = param.NewParamCollection(parseJsonBody(built))

	return built
}

// parseJsonBody returns the body of the request as parsed JSON.
//
// It returns nothing where the content type does not state JSON, where the body
// is empty, and where no decoder reads the body.
func parseJsonBody(request *ServerRequest) map[string]any {
	contentType := request.GetHeaders().GetHeaderLine(constant.HeaderNameContentType)
	if !strings.Contains(contentType, constant.ContentTypeValueApplicationJson) {
		return nil
	}

	body := request.GetBody()

	_ = body.Rewind()

	contents, err := body.GetContents()
	if err != nil || contents == "" {
		return nil
	}

	parsed := map[string]any{}

	err = json.Unmarshal([]byte(contents), &parsed)
	if err != nil {
		return nil
	}

	return parsed
}

// GetParsedJson returns the parsed JSON body of the request.
func (r *ServerRequest) GetParsedJson() contract.ParsedJsonParamCollectionContract {
	return r.parsedJson
}

// WithParsedJson returns a copy of the request with another parsed JSON body.
func (r *ServerRequest) WithParsedJson(
	params contract.ParsedJsonParamCollectionContract,
) contract.JsonServerRequestContract {
	copied := *r
	copied.parsedJson = params

	return &copied
}

// GetServerParams returns the server parameters of the request.
func (r *ServerRequest) GetServerParams() contract.ServerParamCollectionContract {
	return r.serverParams
}

// WithServerParams returns a copy of the request with other server parameters.
func (r *ServerRequest) WithServerParams(
	server contract.ServerParamCollectionContract,
) contract.ServerRequestContract {
	copied := *r
	copied.serverParams = server

	return &copied
}

// GetCookieParams returns the cookies of the request.
func (r *ServerRequest) GetCookieParams() contract.CookieParamCollectionContract {
	return r.cookieParams
}

// WithCookieParams returns a copy of the request with other cookies.
func (r *ServerRequest) WithCookieParams(
	cookies contract.CookieParamCollectionContract,
) contract.ServerRequestContract {
	copied := *r
	copied.cookieParams = cookies

	return &copied
}

// GetQueryParams returns the query parameters of the request.
func (r *ServerRequest) GetQueryParams() contract.QueryParamCollectionContract {
	return r.queryParams
}

// WithQueryParams returns a copy of the request with other query parameters.
func (r *ServerRequest) WithQueryParams(
	query contract.QueryParamCollectionContract,
) contract.ServerRequestContract {
	copied := *r
	copied.queryParams = query

	return &copied
}

// GetUploadedFiles returns the files that arrived with the request.
func (r *ServerRequest) GetUploadedFiles() contract.UploadedFileCollectionContract {
	return r.uploadedFiles
}

// WithUploadedFiles returns a copy of the request with other files.
func (r *ServerRequest) WithUploadedFiles(
	uploadedFiles contract.UploadedFileCollectionContract,
) contract.ServerRequestContract {
	copied := *r
	copied.uploadedFiles = uploadedFiles

	return &copied
}

// GetParsedBody returns the parsed body of the request.
func (r *ServerRequest) GetParsedBody() contract.ParsedBodyParamCollectionContract {
	return r.parsedBody
}

// WithParsedBody returns a copy of the request with another parsed body.
func (r *ServerRequest) WithParsedBody(
	params contract.ParsedBodyParamCollectionContract,
) contract.ServerRequestContract {
	copied := *r
	copied.parsedBody = params

	return &copied
}

// GetAttributes returns the attributes that the framework put on the request.
func (r *ServerRequest) GetAttributes() contract.AttributeParamCollectionContract {
	return r.attributes
}

// WithAttributes returns a copy of the request with other attributes.
func (r *ServerRequest) WithAttributes(
	attributes contract.AttributeParamCollectionContract,
) contract.ServerRequestContract {
	copied := *r
	copied.attributes = attributes

	return &copied
}

// IsXmlHttpRequest reports whether the client sent the request from a script.
func (r *ServerRequest) IsXmlHttpRequest() bool {
	line := r.GetHeaders().GetHeaderLine(constant.HeaderNameXRequestedWith)

	return strings.EqualFold(line, xmlHttpRequestValue)
}
