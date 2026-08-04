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

// portSeparator separates the host of an authority from its port.
const portSeparator = ":"

// NewUriFromString builds a URI by reading the string.
//
// The other ports read the string with the platform's own parser, and this port
// reads it with `net/url`. A string that no parser reads reports a failure, and
// so does a port, a path, or a query that a URI does not carry.
func NewUriFromString(raw string) (*Uri, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	port, err := readPort(parsed.Host)
	if err != nil {
		return nil, err
	}

	password, _ := parsed.User.Password()

	return NewUri(
		constant.Scheme(parsed.Scheme),
		parsed.User.Username(),
		password,
		parsed.Hostname(),
		port,
		parsed.Path,
		parsed.RawQuery,
		parsed.Fragment,
	)
}

// readPort returns the port that the authority names, and zero where it names
// none.
func readPort(host string) (int, error) {
	_, port, found := strings.Cut(host, portSeparator)
	if !found || port == "" {
		return 0, nil
	}

	return strconv.Atoi(port)
}
