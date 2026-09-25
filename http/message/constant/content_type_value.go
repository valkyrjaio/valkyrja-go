/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

// The value of each content type that the framework writes.
const (
	ContentTypeValueApplicationJson       = "application/json"
	ContentTypeValueApplicationJavascript = "application/javascript"
	ContentTypeValueApplicationXml        = "application/xml"
	ContentTypeValueApplicationXmlUtf8    = ContentTypeValueApplicationXml + "; charset=utf-8"
	ContentTypeValueApplicationXWwwForm   = "application/x-www-form-urlencoded"
	ContentTypeValueMultipartFormData     = "multipart/form-data"
	ContentTypeValueTextHtml              = "text/html"
	ContentTypeValueTextHtmlUtf8          = ContentTypeValueTextHtml + "; charset=utf-8"
	ContentTypeValueTextJavascript        = "text/javascript"
	ContentTypeValueTextPlain             = "text/plain"
	ContentTypeValueTextPlainUtf8         = ContentTypeValueTextPlain + "; charset=utf-8"
)
