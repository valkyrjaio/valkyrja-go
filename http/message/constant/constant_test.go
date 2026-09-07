/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant_test

import (
	"slices"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

func TestGetTextReturnsTheReasonPhrase(t *testing.T) {
	t.Parallel()

	tests := map[constant.StatusCode]string{
		constant.StatusCodeContinue:            "Continue",
		constant.StatusCodeOk:                  "OK",
		constant.StatusCodeNotFound:            "Not Found",
		constant.StatusCodeIAmATeapot:          "I Am A Teapot",
		constant.StatusCodeInternalServerError: "Internal Server Error",
	}

	for code, text := range tests {
		if code.GetText() != text {
			t.Errorf("GetText for %d must be %q, but is %q", code, text, code.GetText())
		}
	}
}

func TestGetTextIsEmptyForAnUnknownStatusCode(t *testing.T) {
	t.Parallel()

	if constant.StatusCode(999).GetText() != "" {
		t.Errorf("GetText for an unknown code must be empty, but is: %q", constant.StatusCode(999).GetText())
	}
}

func TestIsValidReportsWhetherTheFrameworkKnowsTheStatusCode(t *testing.T) {
	t.Parallel()

	if !constant.StatusCodeOk.IsValid() {
		t.Error("IsValid must be true for a known status code, but is false")
	}

	if constant.StatusCode(999).IsValid() {
		t.Error("IsValid must be false for an unknown status code, but is true")
	}
}

func TestGetAllRequestMethodsLeavesOutTheAnyMethod(t *testing.T) {
	t.Parallel()

	methods := constant.GetAllRequestMethods()

	if len(methods) != 9 {
		t.Errorf("GetAllRequestMethods must return nine methods, but returned: %d", len(methods))
	}

	if slices.Contains(methods, constant.RequestMethodAny) {
		t.Error("GetAllRequestMethods must leave out the any method, but returned it")
	}

	if !slices.Contains(methods, constant.RequestMethodGet) {
		t.Error("GetAllRequestMethods must return the get method, but did not")
	}
}

func TestIsValidPortAcceptsThePortRange(t *testing.T) {
	t.Parallel()

	tests := map[int]bool{
		constant.PortMin - 1: false,
		constant.PortMin:     true,
		constant.PortHttp:    true,
		constant.PortHttps:   true,
		constant.PortMax:     true,
		constant.PortMax + 1: false,
	}

	for port, valid := range tests {
		if constant.IsValidPort(port) != valid {
			t.Errorf("IsValidPort for %d must be %t, but is %t", port, valid, constant.IsValidPort(port))
		}
	}
}

func TestIsReadableReportsEachReadableMode(t *testing.T) {
	t.Parallel()

	readable := []constant.Mode{
		constant.ModeRead,
		constant.ModeReadWrite,
		constant.ModeWriteRead,
		constant.ModeWriteReadEnd,
		constant.ModeCreateWriteRead,
		constant.ModeWriteReadCreate,
	}

	for _, mode := range readable {
		if !mode.IsReadable() {
			t.Errorf("IsReadable for %q must be true, but is false", mode)
		}
	}

	if constant.ModeWrite.IsReadable() {
		t.Error("IsReadable for the write mode must be false, but is true")
	}
}

func TestIsWritableReportsEachWritableMode(t *testing.T) {
	t.Parallel()

	writable := []constant.Mode{
		constant.ModeReadWrite,
		constant.ModeWrite,
		constant.ModeWriteRead,
		constant.ModeWriteEnd,
		constant.ModeWriteReadEnd,
		constant.ModeCreateWrite,
		constant.ModeCreateWriteRead,
		constant.ModeWriteCreate,
		constant.ModeWriteReadCreate,
	}

	for _, mode := range writable {
		if !mode.IsWritable() {
			t.Errorf("IsWritable for %q must be true, but is false", mode)
		}
	}

	if constant.ModeRead.IsWritable() {
		t.Error("IsWritable for the read mode must be false, but is true")
	}
}

func TestTheContentTypeValuesCarryTheirCharset(t *testing.T) {
	t.Parallel()

	if constant.ContentTypeValueTextHtmlUtf8 != "text/html; charset=utf-8" {
		t.Errorf("the html value must carry its charset, but is: %s", constant.ContentTypeValueTextHtmlUtf8)
	}

	if constant.ContentTypeValueApplicationXmlUtf8 != "application/xml; charset=utf-8" {
		t.Errorf("the xml value must carry its charset, but is: %s", constant.ContentTypeValueApplicationXmlUtf8)
	}
}
