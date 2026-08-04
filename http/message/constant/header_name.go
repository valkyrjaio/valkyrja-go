/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// The name of each HTTP header that the framework reads or writes.
const (
	HeaderNameAccept             = "Accept"
	HeaderNameAcceptCharset      = "Accept-Charset"
	HeaderNameAcceptEncoding     = "Accept-Encoding"
	HeaderNameAcceptLanguage     = "Accept-Language"
	HeaderNameAcceptRanges       = "Accept-Ranges"
	HeaderNameAge                = "Age"
	HeaderNameAllow              = "Allow"
	HeaderNameAuthorization      = "Authorization"
	HeaderNameCacheControl       = "Cache-Control"
	HeaderNameConnection         = "Connection"
	HeaderNameContentEncoding    = "Content-Encoding"
	HeaderNameContentLanguage    = "Content-Language"
	HeaderNameContentLength      = "Content-Length"
	HeaderNameContentLocation    = "Content-Location"
	HeaderNameContentMd5         = "Content-MD5"
	HeaderNameContentRange       = "Content-Range"
	HeaderNameContentType        = "Content-Type"
	HeaderNameDate               = "Date"
	HeaderNameETag               = "ETag"
	HeaderNameExpect             = "Expect"
	HeaderNameExpires            = "Expires"
	HeaderNameFrom               = "From"
	HeaderNameHost               = "Host"
	HeaderNameIfMatch            = "If-Match"
	HeaderNameIfModifiedSince    = "If-Modified-Since"
	HeaderNameIfNoneMatch        = "If-None-Match"
	HeaderNameIfRange            = "If-Range"
	HeaderNameIfUnmodifiedSince  = "If-Unmodified-Since"
	HeaderNameLastModified       = "Last-Modified"
	HeaderNameLocation           = "Location"
	HeaderNameMaxForwards        = "Max-Forwards"
	HeaderNamePragma             = "Pragma"
	HeaderNameProxyAuthenticate  = "Proxy-Authenticate"
	HeaderNameProxyAuthorization = "Proxy-Authorization"
	HeaderNameRange              = "Range"
	HeaderNameReferer            = "Referer"
	HeaderNameRetryAfter         = "Retry-After"
	HeaderNameServer             = "Server"
	HeaderNameSetCookie          = "Set-Cookie"
	HeaderNameTe                 = "TE"
	HeaderNameTrailer            = "Trailer"
	HeaderNameTransferEncoding   = "Transfer-Encoding"
	HeaderNameUpgrade            = "Upgrade"
	HeaderNameUserAgent          = "User-Agent"
	HeaderNameVary               = "Vary"
	HeaderNameVia                = "Via"
	HeaderNameWarning            = "Warning"
	HeaderNameWwwAuthenticate    = "WWW-Authenticate"
	HeaderNameXRequestedWith     = "X-Requested-With"
)

// The value of each HTTP header that the framework writes.
const (
	HeaderValueBearer = "Bearer"
)

// The binding key of each service that the HTTP message sub-component publishes.
const (
	ResponseFactoryContractServiceID = "Valkyrja.Http.Message.Factory.ResponseFactoryContract"

	ServerRequestContractServiceID = "Valkyrja.Http.Message.Request.ServerRequestContract"
)
