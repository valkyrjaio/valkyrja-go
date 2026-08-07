/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	throwablecontract "github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
)

type ComponentContract interface {
	// GetToken returns the token of the component.
	GetToken() string

	// WithToken returns a copy of the component for another token.
	WithToken(token string) ComponentContract

	// GetText returns the text of the component.
	GetText() string

	// WithText returns a copy of the component for another text.
	WithText(text string) ComponentContract

	// String returns the whole component as a string.
	String() string
}

type ValueContract interface {
	// GetComponents returns each component of the value.
	GetComponents() []ComponentContract

	// WithComponents returns a copy of the value with other components.
	WithComponents(components ...ComponentContract) ValueContract

	// WithAddedComponents returns a copy of the value with the components
	// added after the ones it holds.
	WithAddedComponents(components ...ComponentContract) ValueContract

	// String returns the whole value as a string.
	String() string
}

//nolint:interfacebloat // Parity with the PHP reference implementation.
type CookieContract interface {
	ValueContract

	// Delete returns a copy of the cookie that tells the client to remove it.
	Delete() CookieContract

	// GetMaxAge returns the number of seconds until the cookie expires.
	GetMaxAge() int

	// GetName returns the name of the cookie.
	GetName() string

	// WithName returns a copy of the cookie under another name.
	WithName(name string) CookieContract

	// GetValue returns the value of the cookie.
	GetValue() string

	// WithValue returns a copy of the cookie with another value.
	WithValue(value string) CookieContract

	// GetExpire returns the time that the cookie expires at.
	GetExpire() int

	// WithExpire returns a copy of the cookie that expires at another time.
	WithExpire(expire int) CookieContract

	// GetPath returns the path that the cookie applies to.
	GetPath() string

	// WithPath returns a copy of the cookie for another path.
	WithPath(path string) CookieContract

	// GetDomain returns the domain that the cookie applies to.
	GetDomain() string

	// WithDomain returns a copy of the cookie for another domain.
	WithDomain(domain string) CookieContract

	// IsSecure reports whether the client sends the cookie over HTTPS only.
	IsSecure() bool

	// WithSecure returns a copy of the cookie with another secure flag.
	WithSecure(secure bool) CookieContract

	// IsHttpOnly reports whether a script on the client cannot read the cookie.
	IsHttpOnly() bool

	// WithHttpOnly returns a copy of the cookie with another HTTP-only flag.
	WithHttpOnly(httpOnly bool) CookieContract

	// IsRaw reports whether the client receives the value without encoding.
	IsRaw() bool

	// WithRaw returns a copy of the cookie with another raw flag.
	WithRaw(raw bool) CookieContract

	// GetSameSite returns the SameSite attribute of the cookie.
	GetSameSite() constant.SameSite

	// WithSameSite returns a copy of the cookie for another SameSite
	// attribute.
	WithSameSite(sameSite constant.SameSite) CookieContract
}

type HeaderContract interface {
	// GetName returns the name of the header, as the sender wrote it.
	GetName() string

	// GetNormalizedName returns the name of the header in lower case, which is
	// the form that the collection keys on.
	GetNormalizedName() string

	// WithName returns a copy of the header under another name.
	WithName(name string) HeaderContract

	// GetValues returns each value of the header.
	GetValues() []ValueContract

	// WithValues returns a copy of the header with other values.
	WithValues(values ...ValueContract) HeaderContract

	// WithAddedValues returns a copy of the header with the values added after
	// the ones it holds.
	WithAddedValues(values ...ValueContract) HeaderContract

	// GetHeaderLine returns every value of the header as one comma-separated
	// string.
	GetHeaderLine() string

	// String returns the whole header as a string.
	String() string
}

type HeaderCollectionContract interface {
	// Has reports whether the collection holds a header under the name.
	Has(name string) bool

	// Get returns the header under the name. It reports a failure where the
	// collection holds no header under the name.
	//
	// The other ports throw here. Go reports a failure with a returned error,
	// which the method naming rules require.
	Get(name string) (HeaderContract, error)

	// GetHeaderLine returns every value of the header as one comma-separated
	// string.
	GetHeaderLine(name string) string

	// GetAll returns every header, keyed by its normalized name.
	GetAll() map[string]HeaderContract

	// GetOnly returns the headers that the names identify.
	GetOnly(names ...string) map[string]HeaderContract

	// GetAllExcept returns every header that the names do not identify.
	GetAllExcept(names ...string) map[string]HeaderContract

	// WithHeader returns a copy of the collection with the header in it.
	WithHeader(header HeaderContract) HeaderCollectionContract

	// WithoutHeader returns a copy of the collection without the header.
	WithoutHeader(name string) HeaderCollectionContract

	// WithHeaders returns a copy of the collection with other headers.
	WithHeaders(headers ...HeaderContract) HeaderCollectionContract

	// WithAddedHeaders returns a copy of the collection with the headers added
	// to the ones it holds.
	WithAddedHeaders(headers ...HeaderContract) HeaderCollectionContract
}

type HttpThrowable interface {
	throwablecontract.ValkyrjaThrowable

	// IsHttpThrowable marks the error as one that the HTTP component raised.
	// The mark is what separates this contract from the root contract, which Go
	// otherwise treats as the same type.
	IsHttpThrowable() bool
}
