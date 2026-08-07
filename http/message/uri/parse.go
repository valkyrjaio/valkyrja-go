/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package uri

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// ipLiteralSeparator is the character that an IPv6 address separates its groups
// with, and that an authority otherwise reads as the start of a port.
const ipLiteralSeparator = ":"

// ipLiteralOpen opens the brackets that wrap an IPv6 literal in an authority.
const ipLiteralOpen = "["

// ipLiteralClose closes them.
const ipLiteralClose = "]"

// NewUriFromString builds a URI by reading the string.
func NewUriFromString(raw string) (*Uri, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	port, err := readPort(parsed)
	if err != nil {
		return nil, err
	}

	password, _ := parsed.User.Password()

	return NewUri(
		constant.Scheme(parsed.Scheme),
		parsed.User.Username(),
		password,
		readHost(parsed),
		port,
		parsed.Path,
		parsed.RawQuery,
		parsed.Fragment,
	)
}

// readHost returns the host that the authority names.
func readHost(parsed *url.URL) string {
	host := parsed.Hostname()
	if !strings.Contains(host, ipLiteralSeparator) {
		return host
	}

	return ipLiteralOpen + host + ipLiteralClose
}

// readPort returns the port that the authority names, and zero where it names
// none.
//
// Warning: an IPv6 literal carries a colon of its own, and the authority wraps
// it in brackets — `[::1]:8080`. Cutting the authority at its first colon lands
// inside the address, so this reads the port with `url.URL.Port`, which reads
// the brackets.
func readPort(parsed *url.URL) (int, error) {
	port := parsed.Port()
	if port == "" {
		return 0, nil
	}

	return strconv.Atoi(port)
}
