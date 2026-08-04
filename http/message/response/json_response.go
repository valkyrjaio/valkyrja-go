/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package response

import (
	"maps"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// JsonResponse is a response that carries JSON.
//
// A callback makes the response JSONP: the body then carries a call that wraps
// the JSON, and the content type states JavaScript rather than JSON.
type JsonResponse struct {
	Response

	data     map[string]any
	callback string
}

// NewJsonResponseFromData builds a response that carries the data as JSON.
func NewJsonResponseFromData(
	data map[string]any,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *JsonResponse {
	return &JsonResponse{
		Response: *NewJsonResponse(data, statusCode, headers),
		data:     data,
	}
}

// NewJsonpResponseFromData builds a response that wraps the JSON in the
// callback.
func NewJsonpResponseFromData(
	callback string,
	data map[string]any,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) *JsonResponse {
	return &JsonResponse{
		Response: *NewJsonpResponse(callback, data, statusCode, headers),
		data:     data,
		callback: callback,
	}
}

// CreateFromData returns a response that carries the data as JSON.
func (r *JsonResponse) CreateFromData(
	data map[string]any,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.JsonResponseContract {
	return NewJsonResponseFromData(data, statusCode, headers)
}

// GetBodyAsJson returns the body of the response as JSON.
//
// The other ports decode the body, because a caller there replaces the body of a
// response. Go's response is closed over the data that built it, so the body
// always carries that data, and no decode can disagree with it.
func (r *JsonResponse) GetBodyAsJson() map[string]any {
	return maps.Clone(r.data)
}

// WithJsonAsBody returns a copy of the response that carries the data as JSON.
func (r *JsonResponse) WithJsonAsBody(data map[string]any) contract.JsonResponseContract {
	return r.rebuild(data, r.callback)
}

// WithCallback returns a copy of the response that wraps the JSON in the
// callback, which makes it JSONP.
func (r *JsonResponse) WithCallback(callback string) contract.JsonResponseContract {
	return r.rebuild(r.data, callback)
}

// WithoutCallback returns a copy of the response with no callback.
func (r *JsonResponse) WithoutCallback() contract.JsonResponseContract {
	return r.rebuild(r.data, "")
}

// rebuild returns a response that carries the data, with the callback where the
// caller states one.
func (r *JsonResponse) rebuild(data map[string]any, callback string) contract.JsonResponseContract {
	if callback == "" {
		return NewJsonResponseFromData(data, r.GetStatusCode(), r.GetHeaders())
	}

	return NewJsonpResponseFromData(callback, data, r.GetStatusCode(), r.GetHeaders())
}
