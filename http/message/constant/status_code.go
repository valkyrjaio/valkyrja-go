/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the HTTP message component's enumerations and its
// binding keys.
//
// Go has no enum, so each enumeration of the other ports is a defined type over
// a `const` block. The taxonomy puts every one of them in the `constant`
// segment for that reason.
package constant

// StatusCode is the status code of an HTTP response.
type StatusCode int

// The status code of each response that the framework builds.
const (
	// StatusCodeContinue is 100 Continue.
	StatusCodeContinue StatusCode = 100

	// StatusCodeSwitchingProtocols is 101 Switching Protocols.
	StatusCodeSwitchingProtocols StatusCode = 101

	// StatusCodeProcessing is 102 Processing.
	StatusCodeProcessing StatusCode = 102

	// StatusCodeEarlyHints is 103 Early Hints.
	StatusCodeEarlyHints StatusCode = 103

	// StatusCodeOk is 200 OK.
	StatusCodeOk StatusCode = 200

	// StatusCodeCreated is 201 Created.
	StatusCodeCreated StatusCode = 201

	// StatusCodeAccepted is 202 Accepted.
	StatusCodeAccepted StatusCode = 202

	// StatusCodeNonAuthoritativeInformation is 203 Non-Authoritative Information.
	StatusCodeNonAuthoritativeInformation StatusCode = 203

	// StatusCodeNoContent is 204 No Content.
	StatusCodeNoContent StatusCode = 204

	// StatusCodeResetContent is 205 Reset Content.
	StatusCodeResetContent StatusCode = 205

	// StatusCodePartialContent is 206 Partial Content.
	StatusCodePartialContent StatusCode = 206

	// StatusCodeMultiStatus is 207 Multi-Status.
	StatusCodeMultiStatus StatusCode = 207

	// StatusCodeAlreadyReported is 208 Already Reported.
	StatusCodeAlreadyReported StatusCode = 208

	// StatusCodeImUsed is 226 IM Used.
	StatusCodeImUsed StatusCode = 226

	// StatusCodeMultipleChoices is 300 Multiple Choices.
	StatusCodeMultipleChoices StatusCode = 300

	// StatusCodeMovedPermanently is 301 Moved Permanently.
	StatusCodeMovedPermanently StatusCode = 301

	// StatusCodeFound is 302 Found.
	StatusCodeFound StatusCode = 302

	// StatusCodeSeeOther is 303 See Other.
	StatusCodeSeeOther StatusCode = 303

	// StatusCodeNotModified is 304 Not Modified.
	StatusCodeNotModified StatusCode = 304

	// StatusCodeUseProxy is 305 Use Proxy.
	StatusCodeUseProxy StatusCode = 305

	// StatusCodeTemporaryRedirect is 307 Temporary Redirect.
	StatusCodeTemporaryRedirect StatusCode = 307

	// StatusCodePermanentRedirect is 308 Permanent Redirect.
	StatusCodePermanentRedirect StatusCode = 308

	// StatusCodeBadRequest is 400 Bad Request.
	StatusCodeBadRequest StatusCode = 400

	// StatusCodeUnauthorized is 401 Unauthorized.
	StatusCodeUnauthorized StatusCode = 401

	// StatusCodePaymentRequired is 402 Payment Required.
	StatusCodePaymentRequired StatusCode = 402

	// StatusCodeForbidden is 403 Forbidden.
	StatusCodeForbidden StatusCode = 403

	// StatusCodeNotFound is 404 Not Found.
	StatusCodeNotFound StatusCode = 404

	// StatusCodeMethodNotAllowed is 405 Method Not Allowed.
	StatusCodeMethodNotAllowed StatusCode = 405

	// StatusCodeNotAcceptable is 406 Not Acceptable.
	StatusCodeNotAcceptable StatusCode = 406

	// StatusCodeProxyAuthenticationRequired is 407 Proxy Authentication Required.
	StatusCodeProxyAuthenticationRequired StatusCode = 407

	// StatusCodeRequestTimeout is 408 Request Timeout.
	StatusCodeRequestTimeout StatusCode = 408

	// StatusCodeConflict is 409 Conflict.
	StatusCodeConflict StatusCode = 409

	// StatusCodeGone is 410 Gone.
	StatusCodeGone StatusCode = 410

	// StatusCodeLengthRequired is 411 Length Required.
	StatusCodeLengthRequired StatusCode = 411

	// StatusCodePreconditionFailed is 412 Precondition Failed.
	StatusCodePreconditionFailed StatusCode = 412

	// StatusCodePayloadTooLarge is 413 Payload Too Large.
	StatusCodePayloadTooLarge StatusCode = 413

	// StatusCodeUriTooLong is 414 URI Too Long.
	StatusCodeUriTooLong StatusCode = 414

	// StatusCodeUnsupportedMediaType is 415 Unsupported Media Type.
	StatusCodeUnsupportedMediaType StatusCode = 415

	// StatusCodeRangeNotSatisfiable is 416 Range Not Satisfiable.
	StatusCodeRangeNotSatisfiable StatusCode = 416

	// StatusCodeExpectationFailed is 417 Expectation Failed.
	StatusCodeExpectationFailed StatusCode = 417

	// StatusCodeIAmATeapot is 418 I Am A Teapot.
	StatusCodeIAmATeapot StatusCode = 418

	// StatusCodeMisdirectedRequest is 421 Misdirected Request.
	StatusCodeMisdirectedRequest StatusCode = 421

	// StatusCodeUnprocessableEntity is 422 Unprocessable Entity.
	StatusCodeUnprocessableEntity StatusCode = 422

	// StatusCodeLocked is 423 Locked.
	StatusCodeLocked StatusCode = 423

	// StatusCodeFailedDependency is 424 Failed Dependency.
	StatusCodeFailedDependency StatusCode = 424

	// StatusCodeTooEarly is 425 Too Early.
	StatusCodeTooEarly StatusCode = 425

	// StatusCodeUpgradeRequired is 426 Upgrade Required.
	StatusCodeUpgradeRequired StatusCode = 426

	// StatusCodePreconditionRequired is 428 Precondition Required.
	StatusCodePreconditionRequired StatusCode = 428

	// StatusCodeTooManyRequests is 429 Too Many Requests.
	StatusCodeTooManyRequests StatusCode = 429

	// StatusCodeRequestHeaderFieldsTooLarge is 431 Request Header Fields Too Large.
	StatusCodeRequestHeaderFieldsTooLarge StatusCode = 431

	// StatusCodeUnavailableForLegalReasons is 451 Unavailable For Legal Reasons.
	StatusCodeUnavailableForLegalReasons StatusCode = 451

	// StatusCodeInternalServerError is 500 Internal Server Error.
	StatusCodeInternalServerError StatusCode = 500

	// StatusCodeNotImplemented is 501 Not Implemented.
	StatusCodeNotImplemented StatusCode = 501

	// StatusCodeBadGateway is 502 Bad Gateway.
	StatusCodeBadGateway StatusCode = 502

	// StatusCodeServiceUnavailable is 503 Service Unavailable.
	StatusCodeServiceUnavailable StatusCode = 503

	// StatusCodeGatewayTimeout is 504 Gateway Timeout.
	StatusCodeGatewayTimeout StatusCode = 504

	// StatusCodeHttpVersionNotSupported is 505 HTTP Version Not Supported.
	StatusCodeHttpVersionNotSupported StatusCode = 505

	// StatusCodeVariantAlsoNegotiates is 506 Variant Also Negotiates.
	StatusCodeVariantAlsoNegotiates StatusCode = 506

	// StatusCodeInsufficientStorage is 507 Insufficient Storage.
	StatusCodeInsufficientStorage StatusCode = 507

	// StatusCodeLoopDetected is 508 Loop Detected.
	StatusCodeLoopDetected StatusCode = 508

	// StatusCodeNotExtendedObsoleted is 510 Not Extended.
	StatusCodeNotExtendedObsoleted StatusCode = 510

	// StatusCodeNetworkAuthenticationRequired is 511 Network Authentication Required.
	StatusCodeNetworkAuthenticationRequired StatusCode = 511
)

// statusTexts holds the reason phrase of each status code.
var statusTexts = map[StatusCode]string{
	StatusCodeContinue:                      "Continue",
	StatusCodeSwitchingProtocols:            "Switching Protocols",
	StatusCodeProcessing:                    "Processing",
	StatusCodeEarlyHints:                    "Early Hints",
	StatusCodeOk:                            "OK",
	StatusCodeCreated:                       "Created",
	StatusCodeAccepted:                      "Accepted",
	StatusCodeNonAuthoritativeInformation:   "Non-Authoritative Information",
	StatusCodeNoContent:                     "No Content",
	StatusCodeResetContent:                  "Reset Content",
	StatusCodePartialContent:                "Partial Content",
	StatusCodeMultiStatus:                   "Multi-Status",
	StatusCodeAlreadyReported:               "Already Reported",
	StatusCodeImUsed:                        "IM Used",
	StatusCodeMultipleChoices:               "Multiple Choices",
	StatusCodeMovedPermanently:              "Moved Permanently",
	StatusCodeFound:                         "Found",
	StatusCodeSeeOther:                      "See Other",
	StatusCodeNotModified:                   "Not Modified",
	StatusCodeUseProxy:                      "Use Proxy",
	StatusCodeTemporaryRedirect:             "Temporary Redirect",
	StatusCodePermanentRedirect:             "Permanent Redirect",
	StatusCodeBadRequest:                    "Bad Request",
	StatusCodeUnauthorized:                  "Unauthorized",
	StatusCodePaymentRequired:               "Payment Required",
	StatusCodeForbidden:                     "Forbidden",
	StatusCodeNotFound:                      "Not Found",
	StatusCodeMethodNotAllowed:              "Method Not Allowed",
	StatusCodeNotAcceptable:                 "Not Acceptable",
	StatusCodeProxyAuthenticationRequired:   "Proxy Authentication Required",
	StatusCodeRequestTimeout:                "Request Timeout",
	StatusCodeConflict:                      "Conflict",
	StatusCodeGone:                          "Gone",
	StatusCodeLengthRequired:                "Length Required",
	StatusCodePreconditionFailed:            "Precondition Failed",
	StatusCodePayloadTooLarge:               "Payload Too Large",
	StatusCodeUriTooLong:                    "URI Too Long",
	StatusCodeUnsupportedMediaType:          "Unsupported Media Type",
	StatusCodeRangeNotSatisfiable:           "Range Not Satisfiable",
	StatusCodeExpectationFailed:             "Expectation Failed",
	StatusCodeIAmATeapot:                    "I Am A Teapot",
	StatusCodeMisdirectedRequest:            "Misdirected Request",
	StatusCodeUnprocessableEntity:           "Unprocessable Entity",
	StatusCodeLocked:                        "Locked",
	StatusCodeFailedDependency:              "Failed Dependency",
	StatusCodeTooEarly:                      "Too Early",
	StatusCodeUpgradeRequired:               "Upgrade Required",
	StatusCodePreconditionRequired:          "Precondition Required",
	StatusCodeTooManyRequests:               "Too Many Requests",
	StatusCodeRequestHeaderFieldsTooLarge:   "Request Header Fields Too Large",
	StatusCodeUnavailableForLegalReasons:    "Unavailable For Legal Reasons",
	StatusCodeInternalServerError:           "Internal Server Error",
	StatusCodeNotImplemented:                "Not Implemented",
	StatusCodeBadGateway:                    "Bad Gateway",
	StatusCodeServiceUnavailable:            "Service Unavailable",
	StatusCodeGatewayTimeout:                "Gateway Timeout",
	StatusCodeHttpVersionNotSupported:       "HTTP Version Not Supported",
	StatusCodeVariantAlsoNegotiates:         "Variant Also Negotiates",
	StatusCodeInsufficientStorage:           "Insufficient Storage",
	StatusCodeLoopDetected:                  "Loop Detected",
	StatusCodeNotExtendedObsoleted:          "Not Extended",
	StatusCodeNetworkAuthenticationRequired: "Network Authentication Required",
}

// GetText returns the reason phrase of the status code, and an empty string
// where the code is not one that the framework knows.
func (s StatusCode) GetText() string {
	return statusTexts[s]
}

// IsValid reports whether the status code is one that the framework knows.
func (s StatusCode) IsValid() bool {
	_, found := statusTexts[s]

	return found
}
