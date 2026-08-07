/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package value_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
)

const (
	cookieName  = "session"
	cookieValue = "abc123"
)

func TestNewCookieTakesTheDefaultsOfEveryPort(t *testing.T) {
	t.Parallel()

	cookie := value.NewCookie(cookieName, cookieValue)

	if cookie.GetName() != cookieName || cookie.GetValue() != cookieValue {
		t.Error("the cookie must hold its name and its value, but did not")
	}

	if cookie.GetPath() != "/" {
		t.Errorf("the path must default to the root, but is: %q", cookie.GetPath())
	}

	if !cookie.IsHttpOnly() {
		t.Error("the cookie must default to http only, but did not")
	}

	if cookie.GetSameSite() != constant.SameSiteLax {
		t.Errorf("the SameSite attribute must default to lax, but is: %q", cookie.GetSameSite())
	}

	if cookie.IsSecure() || cookie.IsRaw() {
		t.Error("the cookie must default to neither secure nor raw, but did not")
	}

	if cookie.GetExpire() != 0 || cookie.GetMaxAge() != 0 {
		t.Error("a cookie that names no expiry must have no max age, but did")
	}
}

func TestEachCookieWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	cookie := value.NewCookie(cookieName, cookieValue)

	if cookie.WithName("other").GetName() != "other" {
		t.Error("WithName must hold the new name, but did not")
	}

	if cookie.WithValue("other").GetValue() != "other" {
		t.Error("WithValue must hold the new value, but did not")
	}

	if cookie.WithPath("/admin").GetPath() != "/admin" {
		t.Error("WithPath must hold the new path, but did not")
	}

	if cookie.WithDomain("example.com").GetDomain() != "example.com" {
		t.Error("WithDomain must hold the new domain, but did not")
	}

	if !cookie.WithSecure(true).IsSecure() {
		t.Error("WithSecure must hold the new flag, but did not")
	}

	if cookie.WithHttpOnly(false).IsHttpOnly() {
		t.Error("WithHttpOnly must hold the new flag, but did not")
	}

	if cookie.GetName() != cookieName || cookie.GetPath() != "/" || cookie.IsSecure() {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestEachRemainingCookieWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	cookie := value.NewCookie(cookieName, cookieValue)

	if !cookie.WithRaw(true).IsRaw() {
		t.Error("WithRaw must hold the new flag, but did not")
	}

	if cookie.WithSameSite(constant.SameSiteStrict).GetSameSite() != constant.SameSiteStrict {
		t.Error("WithSameSite must hold the new attribute, but did not")
	}

	if cookie.WithExpire(1000).GetExpire() != 1000 {
		t.Error("WithExpire must hold the new time, but did not")
	}

	if cookie.IsRaw() || cookie.GetExpire() != 0 {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestStringRendersEachPartThatTheCookieCarries(t *testing.T) {
	t.Parallel()

	cookie := value.NewCookie(cookieName, cookieValue).
		WithDomain("example.com").
		WithSecure(true).
		WithSameSite(constant.SameSiteStrict)

	rendered := cookie.String()

	for _, part := range []string{
		"session=abc123",
		"path=/",
		"domain=example.com",
		"secure",
		"httponly",
		"samesite=strict",
	} {
		if !strings.Contains(rendered, part) {
			t.Errorf("String must carry %q, but is: %q", part, rendered)
		}
	}
}

func TestStringLeavesOutEachPartThatTheCookieDoesNotCarry(t *testing.T) {
	t.Parallel()

	rendered := value.NewCookie(cookieName, cookieValue).
		WithHttpOnly(false).
		WithSameSite("").
		String()

	for _, part := range []string{"domain=", "secure", "httponly", "samesite=", "expires="} {
		if strings.Contains(rendered, part) {
			t.Errorf("String must leave out %q, but is: %q", part, rendered)
		}
	}
}

func TestStringEncodesTheNameAndTheValue(t *testing.T) {
	t.Parallel()

	cookie := value.NewCookie("a name", "a value")

	if !strings.Contains(cookie.String(), "a+name=a+value") {
		t.Errorf("String must encode the name and the value, but is: %q", cookie.String())
	}

	if !strings.Contains(cookie.WithRaw(true).String(), "a name=a value") {
		t.Errorf("a raw cookie must carry its parts unencoded, but is: %q", cookie.WithRaw(true).String())
	}
}

func TestStringCarriesTheExpiryAndTheMaxAge(t *testing.T) {
	t.Parallel()

	rendered := value.NewCookie(cookieName, cookieValue).WithExpire(2000000000).String()

	if !strings.Contains(rendered, "expires=") || !strings.Contains(rendered, "max-age=") {
		t.Errorf("String must carry the expiry and the max age, but is: %q", rendered)
	}
}

func TestADeletedCookieExpiresInThePast(t *testing.T) {
	t.Parallel()

	rendered := value.NewCookie(cookieName, cookieValue).Delete().String()

	if !strings.Contains(rendered, "session=delete") {
		t.Errorf("a deleted cookie must carry the delete value, but is: %q", rendered)
	}

	if !strings.Contains(rendered, "max-age=-31536001") {
		t.Errorf("a deleted cookie must expire in the past, but is: %q", rendered)
	}
}

func TestTheCookieBuildsItsOwnComponents(t *testing.T) {
	t.Parallel()

	cookie := value.NewCookie(cookieName, cookieValue)

	components := cookie.GetComponents()

	if len(components) == 0 {
		t.Fatal("GetComponents must read the rendered cookie, but read nothing")
	}

	if components[0].GetToken() != cookieName {
		t.Errorf("the first component must be the name, but is: %q", components[0].GetToken())
	}

	if cookie.WithComponents() != contract.ValueContract(cookie) {
		t.Error("WithComponents must return the cookie, but did not")
	}

	if cookie.WithAddedComponents() != contract.ValueContract(cookie) {
		t.Error("WithAddedComponents must return the cookie, but did not")
	}
}

func TestGetMaxAgeCountsFromTheExpiry(t *testing.T) {
	t.Parallel()

	cookie := value.NewCookie(cookieName, cookieValue).WithExpire(2000000000)

	if cookie.GetMaxAge() <= 0 {
		t.Errorf("GetMaxAge must count the seconds until the expiry, but is: %d", cookie.GetMaxAge())
	}
}

func TestTheCookieSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var cookie contract.CookieContract = value.NewCookie(cookieName, cookieValue)

	if cookie.GetName() != cookieName {
		t.Errorf("the contract must read the name, but read: %q", cookie.GetName())
	}
}
