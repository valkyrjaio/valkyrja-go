/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package value

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

// deletedLifetime is the number of seconds that a deleted cookie sets its expiry
// into the past by. It is one second over a year, which is what every port uses.
const deletedLifetime = 31536001

// deletedValue is the value that a deleted cookie carries.
const deletedValue = "delete"

// defaultPath is the path that a cookie applies to where it names none.
const defaultPath = "/"

// nowFunc returns the current time. A test replaces it, because a cookie renders
// its expiry and a test cannot assert a value that moves.
var nowFunc = time.Now

type Cookie struct {
	name     string
	value    string
	expire   int64
	path     string
	domain   string
	secure   bool
	httpOnly bool
	raw      bool
	sameSite constant.SameSite
	deleted  bool
}

// NewCookie builds a cookie under a name, with a value. Every other part takes
// the default that the other ports use.
func NewCookie(name string, value string) *Cookie {
	return &Cookie{
		name:     name,
		value:    value,
		path:     defaultPath,
		httpOnly: true,
		sameSite: constant.SameSiteLax,
	}
}

// GetComponents returns each component of the cookie, read from its rendered
// form.
func (c *Cookie) GetComponents() []contract.ComponentContract {
	return NewValueFromValue(c.String()).GetComponents()
}

// WithComponents returns the cookie. A cookie builds its own components from its
// own state, so it takes none.
func (c *Cookie) WithComponents(_ ...contract.ComponentContract) contract.ValueContract {
	return c
}

// WithAddedComponents returns the cookie. A cookie builds its own components
// from its own state, so it takes none.
func (c *Cookie) WithAddedComponents(_ ...contract.ComponentContract) contract.ValueContract {
	return c
}

// Delete returns a copy of the cookie that tells the client to remove it.
func (c *Cookie) Delete() contract.CookieContract {
	copied := *c
	copied.deleted = true

	return &copied
}

// GetMaxAge returns the number of seconds until the cookie expires, and zero
// where the cookie names no expiry.
func (c *Cookie) GetMaxAge() int {
	if c.expire <= 0 {
		return 0
	}

	return int(c.expire - nowFunc().Unix())
}

// GetName returns the name of the cookie.
func (c *Cookie) GetName() string {
	return c.name
}

// WithName returns a copy of the cookie under another name.
func (c *Cookie) WithName(name string) contract.CookieContract {
	copied := *c
	copied.name = name

	return &copied
}

// GetValue returns the value of the cookie.
func (c *Cookie) GetValue() string {
	return c.value
}

// WithValue returns a copy of the cookie with another value.
func (c *Cookie) WithValue(value string) contract.CookieContract {
	copied := *c
	copied.value = value

	return &copied
}

// GetExpire returns the time that the cookie expires at.
func (c *Cookie) GetExpire() int {
	return int(c.expire)
}

// WithExpire returns a copy of the cookie that expires at another time.
func (c *Cookie) WithExpire(expire int) contract.CookieContract {
	copied := *c
	copied.expire = int64(expire)

	return &copied
}

// GetPath returns the path that the cookie applies to.
func (c *Cookie) GetPath() string {
	return c.path
}

// WithPath returns a copy of the cookie for another path.
func (c *Cookie) WithPath(path string) contract.CookieContract {
	copied := *c
	copied.path = path

	return &copied
}

// GetDomain returns the domain that the cookie applies to.
func (c *Cookie) GetDomain() string {
	return c.domain
}

// WithDomain returns a copy of the cookie for another domain.
func (c *Cookie) WithDomain(domain string) contract.CookieContract {
	copied := *c
	copied.domain = domain

	return &copied
}

// IsSecure reports whether the client sends the cookie over HTTPS only.
func (c *Cookie) IsSecure() bool {
	return c.secure
}

// WithSecure returns a copy of the cookie with another secure flag.
func (c *Cookie) WithSecure(secure bool) contract.CookieContract {
	copied := *c
	copied.secure = secure

	return &copied
}

// IsHttpOnly reports whether a script on the client cannot read the cookie.
func (c *Cookie) IsHttpOnly() bool {
	return c.httpOnly
}

// WithHttpOnly returns a copy of the cookie with another HTTP-only flag.
func (c *Cookie) WithHttpOnly(httpOnly bool) contract.CookieContract {
	copied := *c
	copied.httpOnly = httpOnly

	return &copied
}

// IsRaw reports whether the client receives the value without encoding.
func (c *Cookie) IsRaw() bool {
	return c.raw
}

// WithRaw returns a copy of the cookie with another raw flag.
func (c *Cookie) WithRaw(raw bool) contract.CookieContract {
	copied := *c
	copied.raw = raw

	return &copied
}

// GetSameSite returns the SameSite attribute of the cookie.
func (c *Cookie) GetSameSite() constant.SameSite {
	return c.sameSite
}

// WithSameSite returns a copy of the cookie for another SameSite attribute.
func (c *Cookie) WithSameSite(sameSite constant.SameSite) contract.CookieContract {
	copied := *c
	copied.sameSite = sameSite

	return &copied
}

// String returns the whole cookie as a `Set-Cookie` value.
func (c *Cookie) String() string {
	value := c.value
	expire := c.expire
	maxAge := c.GetMaxAge()

	if c.deleted {
		expire = nowFunc().Unix() - deletedLifetime
		maxAge = -deletedLifetime
		value = deletedValue
	}

	parts := []string{c.encodePart(c.name) + "=" + c.encodePart(value)}

	if expire != 0 {
		parts = append(parts,
			"expires="+time.Unix(expire, 0).UTC().Format(time.RFC1123),
			"max-age="+strconv.Itoa(maxAge),
		)
	}

	parts = append(parts, "path="+c.path)

	if c.domain != "" {
		parts = append(parts, "domain="+c.domain)
	}

	if c.secure {
		parts = append(parts, "secure")
	}

	if c.httpOnly {
		parts = append(parts, "httponly")
	}

	if c.sameSite != "" {
		parts = append(parts, "samesite="+string(c.sameSite))
	}

	return strings.Join(parts, "; ")
}

// encodePart returns the part as the client receives it. A raw cookie carries
// its name and its value unencoded.
func (c *Cookie) encodePart(part string) string {
	if c.raw {
		return part
	}

	return url.QueryEscape(part)
}
