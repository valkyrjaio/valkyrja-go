/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package uri

import (
	"strconv"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

// noPort is the port of a URI that names none.
const noPort = 0

type Uri struct {
	scheme   constant.Scheme
	username string
	password string
	userInfo string
	host     string
	port     int
	path     string
	query    string
	fragment string
}

// NewUri builds a URI from each of its parts. It reports a failure where the
// port, the path, or the query string is not one that a URI carries.
func NewUri(
	scheme constant.Scheme,
	username string,
	password string,
	host string,
	port int,
	path string,
	query string,
	fragment string,
) (*Uri, error) {
	if port == noPort {
		port = getPortFromScheme(scheme)
	} else if !constant.IsValidPort(port) {
		return nil, exception.NewHttpUriInvalidPortError(port)
	}

	filteredPath, err := filterPath(path)
	if err != nil {
		return nil, err
	}

	filteredQuery, err := filterQuery(query)
	if err != nil {
		return nil, err
	}

	return &Uri{
		scheme:   scheme,
		username: username,
		password: password,
		userInfo: filterUserInfo(buildUserInfo(username, password)),
		host:     filterHost(host),
		port:     port,
		path:     filteredPath,
		query:    filteredQuery,
		fragment: filterFragment(fragment),
	}, nil
}

// GetScheme returns the scheme of the URI.
func (u *Uri) GetScheme() constant.Scheme {
	return u.scheme
}

// IsSecure reports whether the scheme is `https`.
func (u *Uri) IsSecure() bool {
	return u.scheme == constant.SchemeHttps
}

// GetAuthority returns the authority of the URI, and an empty string where the
// URI names no host.
func (u *Uri) GetAuthority() string {
	if u.host == "" {
		return ""
	}

	authority := u.host

	if u.userInfo != "" {
		authority = u.userInfo + "@" + authority
	}

	if !isStandardPort(u.scheme, u.host, u.port) {
		authority += ":" + strconv.Itoa(u.port)
	}

	return authority
}

// GetUsername returns the user name of the URI.
func (u *Uri) GetUsername() string {
	return u.username
}

// GetPassword returns the password of the URI.
func (u *Uri) GetPassword() string {
	return u.password
}

// GetUserInfo returns the user name and the password, separated by a colon.
func (u *Uri) GetUserInfo() string {
	return u.userInfo
}

// GetHost returns the host of the URI.
func (u *Uri) GetHost() string {
	return u.host
}

// HasPort reports whether the URI names a port.
func (u *Uri) HasPort() bool {
	return u.port != noPort
}

// GetPort returns the port of the URI, and zero where the port is the standard
// port of the scheme.
func (u *Uri) GetPort() int {
	if isStandardPort(u.scheme, u.host, u.port) {
		return noPort
	}

	return u.port
}

// GetHostPort returns the host and the port, separated by a colon.
func (u *Uri) GetHostPort() string {
	host := u.host
	port := u.GetPort()

	if host != "" && port != noPort {
		host += ":" + strconv.Itoa(port)
	}

	return host
}

// GetSchemeHostPort returns the scheme, the host, and the port.
func (u *Uri) GetSchemeHostPort() string {
	hostPort := u.GetHostPort()

	if hostPort != "" && u.scheme != constant.SchemeEmpty {
		return string(u.scheme) + "://" + hostPort
	}

	return hostPort
}

// GetPath returns the path of the URI.
func (u *Uri) GetPath() string {
	return u.path
}

// GetQuery returns the query string of the URI.
func (u *Uri) GetQuery() string {
	return u.query
}

// GetFragment returns the fragment of the URI.
func (u *Uri) GetFragment() string {
	return u.fragment
}

// WithScheme returns a copy of the URI for another scheme. A URI that names no
// port takes the standard port of the new scheme.
func (u *Uri) WithScheme(scheme constant.Scheme) contract.UriContract {
	copied := *u
	copied.scheme = scheme

	if u.port == noPort {
		copied.port = getPortFromScheme(scheme)
	}

	return &copied
}

// WithUsername returns a copy of the URI for another user name.
func (u *Uri) WithUsername(username string) contract.UriContract {
	return u.WithUserInfo(username, u.password)
}

// WithPassword returns a copy of the URI for another password.
func (u *Uri) WithPassword(password string) contract.UriContract {
	return u.WithUserInfo(u.username, password)
}

// WithUserInfo returns a copy of the URI for another user name and password. A
// URI with no user name carries no password.
func (u *Uri) WithUserInfo(user string, password string) contract.UriContract {
	if user == "" {
		password = ""
	}

	copied := *u
	copied.username = user
	copied.password = password
	copied.userInfo = filterUserInfo(buildUserInfo(user, password))

	return &copied
}

// WithHost returns a copy of the URI for another host.
func (u *Uri) WithHost(host string) contract.UriContract {
	copied := *u
	copied.host = filterHost(host)

	return &copied
}

// WithPort returns a copy of the URI for another port. It keeps the port of the
// receiver where the new one is outside the range that a URI carries.
func (u *Uri) WithPort(port int) contract.UriContract {
	if !constant.IsValidPort(port) {
		return u
	}

	copied := *u
	copied.port = port

	return &copied
}

// WithPath returns a copy of the URI for another path. It keeps the path of the
// receiver where the new one carries a query string.
func (u *Uri) WithPath(path string) contract.UriContract {
	filtered, err := filterPath(path)
	if err != nil {
		return u
	}

	copied := *u
	copied.path = filtered

	return &copied
}

// WithQuery returns a copy of the URI for another query string. It keeps the
// query string of the receiver where the new one carries a fragment.
func (u *Uri) WithQuery(query string) contract.UriContract {
	filtered, err := filterQuery(query)
	if err != nil {
		return u
	}

	copied := *u
	copied.query = filtered

	return &copied
}

// WithFragment returns a copy of the URI for another fragment.
func (u *Uri) WithFragment(fragment string) contract.UriContract {
	copied := *u
	copied.fragment = filterFragment(fragment)

	return &copied
}

// String returns the whole URI as a string.
func (u *Uri) String() string {
	rendered := &strings.Builder{}

	if u.scheme != constant.SchemeEmpty {
		rendered.WriteString(string(u.scheme))
		rendered.WriteByte(':')
	}

	if authority := u.GetAuthority(); authority != "" {
		rendered.WriteString("//")
		rendered.WriteString(authority)
	}

	if u.path != "" {
		if !strings.HasPrefix(u.path, "/") {
			rendered.WriteByte('/')
		}

		rendered.WriteString(u.path)
	}

	if u.query != "" {
		rendered.WriteByte('?')
		rendered.WriteString(u.query)
	}

	if u.fragment != "" {
		rendered.WriteByte('#')
		rendered.WriteString(u.fragment)
	}

	return rendered.String()
}

// buildUserInfo joins the user name and the password with a colon.
func buildUserInfo(username string, password string) string {
	if password == "" {
		return username
	}

	return username + ":" + password
}

// filterUserInfo encodes each character that the user info does not carry.
func filterUserInfo(userInfo string) string {
	return encode(userInfo, userInfoChars)
}

// filterHost lower-cases the host and encodes each character that it does not
// carry.
func filterHost(host string) string {
	host = strings.ToLower(host)

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return host
	}

	return encode(host, hostChars)
}

// filterPath encodes each character that the path does not carry, and collapses
// the leading separators to one. It reports a failure where the path carries a
// query string.
func filterPath(path string) (string, error) {
	if strings.Contains(path, "?") {
		return "", exception.NewHttpUriInvalidPathError(path)
	}

	encoded := encode(path, pathChars)

	if strings.HasPrefix(encoded, "/") {
		return "/" + strings.TrimLeft(encoded, "/"), nil
	}

	return encoded, nil
}

// filterQuery drops the leading question marks and encodes each character that
// the query string does not carry. It reports a failure where the query string
// carries a fragment.
func filterQuery(query string) (string, error) {
	if strings.Contains(query, "#") {
		return "", exception.NewHttpUriInvalidQueryError(query)
	}

	return encode(strings.TrimLeft(query, "?"), queryChars), nil
}

// filterFragment drops the leading number signs and encodes each character that
// the fragment does not carry.
func filterFragment(fragment string) string {
	return encode(strings.TrimLeft(fragment, "#"), queryChars)
}

// getPortFromScheme returns the standard port of the scheme, and zero where the
// scheme has none.
func getPortFromScheme(scheme constant.Scheme) int {
	switch scheme {
	case constant.SchemeHttps:
		return constant.PortHttps
	case constant.SchemeHttp:
		return constant.PortHttp
	case constant.SchemeEmpty:
		return noPort
	default:
		return noPort
	}
}

// isStandardPort reports whether the port is the one that the scheme implies, so
// the authority leaves it out.
func isStandardPort(scheme constant.Scheme, host string, port int) bool {
	if scheme == constant.SchemeEmpty {
		return host != "" && port <= noPort
	}

	if host == "" || port <= noPort {
		return true
	}

	return (scheme == constant.SchemeHttp && port == constant.PortHttp) ||
		(scheme == constant.SchemeHttps && port == constant.PortHttps)
}
