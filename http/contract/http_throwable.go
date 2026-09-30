/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	throwablecontract "github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
)

type HttpThrowable interface {
	throwablecontract.ValkyrjaThrowable

	// IsHttpThrowable marks the error as one that the HTTP component raised.
	// The mark is what separates this contract from the root contract, which Go
	// otherwise treats as the same type.
	IsHttpThrowable() bool
}

type HttpClientThrowable interface {
	HttpThrowable

	// IsHttpClientThrowable marks the error as one that the HTTP client raised.
	IsHttpClientThrowable() bool
}

type HttpMessageThrowable interface {
	HttpThrowable

	// IsHttpMessageThrowable marks the error as one that the HTTP message sub-component raised.
	IsHttpMessageThrowable() bool
}

type UploadedFileThrowable interface {
	HttpMessageThrowable

	// IsUploadedFileThrowable marks the error as one that an uploaded file raised.
	IsUploadedFileThrowable() bool
}

type HttpHeaderThrowable interface {
	HttpMessageThrowable

	// IsHttpHeaderThrowable marks the error as one that a header raised.
	IsHttpHeaderThrowable() bool
}

type HttpRequestThrowable interface {
	HttpMessageThrowable

	// IsHttpRequestThrowable marks the error as one that a request raised.
	IsHttpRequestThrowable() bool
}

type HttpResponseThrowable interface {
	HttpMessageThrowable

	// IsHttpResponseThrowable marks the error as one that a response raised.
	IsHttpResponseThrowable() bool
}

type HttpStreamThrowable interface {
	HttpMessageThrowable

	// IsHttpStreamThrowable marks the error as one that a stream raised.
	IsHttpStreamThrowable() bool
}

type HttpMiddlewareThrowable interface {
	HttpThrowable

	// IsHttpMiddlewareThrowable marks the error as one that the HTTP middleware raised.
	IsHttpMiddlewareThrowable() bool
}

type HttpRoutingThrowable interface {
	HttpThrowable

	// IsHttpRoutingThrowable marks the error as one that the HTTP routing sub-component raised.
	IsHttpRoutingThrowable() bool
}

type HttpServerThrowable interface {
	HttpThrowable

	// IsHttpServerThrowable marks the error as one that the HTTP server raised.
	IsHttpServerThrowable() bool
}

type HttpStructThrowable interface {
	HttpThrowable

	// IsHttpStructThrowable marks the error as one that a request struct or a response struct raised.
	IsHttpStructThrowable() bool
}
