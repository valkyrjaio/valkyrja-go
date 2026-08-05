/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// RequestMethod is the request method of an HTTP request.
type RequestMethod string

// The RequestMethod values that the framework knows.
const (
	RequestMethodGet     RequestMethod = "GET"
	RequestMethodHead    RequestMethod = "HEAD"
	RequestMethodPost    RequestMethod = "POST"
	RequestMethodPut     RequestMethod = "PUT"
	RequestMethodDelete  RequestMethod = "DELETE"
	RequestMethodConnect RequestMethod = "CONNECT"
	RequestMethodOptions RequestMethod = "OPTIONS"
	RequestMethodTrace   RequestMethod = "TRACE"
	RequestMethodPatch   RequestMethod = "PATCH"
	RequestMethodAny     RequestMethod = "ANY"
)

// GetAllRequestMethods returns every request method that a route matches on its
// own. `RequestMethodAny` is absent, because it stands for all of them.
func GetAllRequestMethods() []RequestMethod {
	return []RequestMethod{
		RequestMethodGet,
		RequestMethodHead,
		RequestMethodPost,
		RequestMethodPut,
		RequestMethodDelete,
		RequestMethodConnect,
		RequestMethodOptions,
		RequestMethodTrace,
		RequestMethodPatch,
	}
}
