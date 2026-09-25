/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the HTTP server component's binding keys.
package constant

// The HTTP server component's binding keys.
const (
	// RequestHandlerContractServiceID is the binding key for the request
	// handler.
	RequestHandlerContractServiceID = "valkyrja.http.server.handler.RequestHandlerContract"

	// ResponseContractServiceID is the binding key for the response that the
	// handler built. The handler binds it before it returns, so something later
	// in the request reads what the client receives.
	ResponseContractServiceID = "valkyrja.http.message.response.ResponseContract"
)
