/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package uri_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
)

func TestNewUriFromStringReadsEachPart(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUriFromString("https://user:secret@example.com:8443/path?query=1#fragment")
	if err != nil {
		t.Fatalf("the parser must read the string, but reported: %v", err)
	}

	if built.GetScheme() != constant.SchemeHttps {
		t.Errorf("the parser must read the scheme, but read: %q", built.GetScheme())
	}

	if built.GetUsername() != "user" || built.GetPassword() != "secret" {
		t.Error("the parser must read the user information, but did not")
	}

	if built.GetHost() != "example.com" || built.GetPort() != 8443 {
		t.Errorf("the parser must read the host and the port, but read: %s", built.GetHostPort())
	}

	if built.GetPath() != "/path" || built.GetQuery() != "query=1" || built.GetFragment() != "fragment" {
		t.Error("the parser must read the path, the query, and the fragment, but did not")
	}
}

func TestNewUriFromStringTakesTheStandardPortWhereTheStringNamesNone(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUriFromString("https://example.com/path")
	if err != nil {
		t.Fatalf("the parser must read the string, but reported: %v", err)
	}

	// The URI takes the standard port of the scheme, and the authority leaves a
	// standard port out, so the host carries no port.
	if built.GetHostPort() != "example.com" {
		t.Errorf("a URI that names no port must leave it out, but read: %q", built.GetHostPort())
	}
}

func TestNewUriFromStringReadsAPathAlone(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUriFromString("/path")
	if err != nil {
		t.Fatalf("the parser must read a path alone, but reported: %v", err)
	}

	if built.GetPath() != "/path" || built.GetHost() != "" {
		t.Error("a path alone must give a URI with no host, but did not")
	}
}

func TestNewUriFromStringReportsAStringThatIsNotAUri(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"a string that no parser reads": "://example.com",
		"a port that is not a number":   "https://example.com:notaport/path",
		"a port outside the range":      "https://example.com:99999/path",
		"a port that overflows":         "https://example.com:99999999999999999999/path",
		"a path that carries a query":   "https://example.com/pa%3Fth",
	}

	for name, raw := range tests {
		_, err := uri.NewUriFromString(raw)
		if err == nil {
			t.Errorf("%s must report a failure, but reported none", name)
		}
	}
}
