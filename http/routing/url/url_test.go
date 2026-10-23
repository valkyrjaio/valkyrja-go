/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package url_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/url"
)

const routeName = "users.show"

// newUrl builds a URL generator over a collection that holds the route.
func newUrl(path string) *url.Url {
	built := collection.NewRouteCollection()
	built.Add(data.NewRoute(path, routeName, nil))

	return url.NewUrl(built)
}

func TestGetUrlFillsEachParameter(t *testing.T) {
	t.Parallel()

	got := newUrl("/users/{id}/posts/{slug}").GetUrl(routeName, map[string]string{
		"id":   "42",
		"slug": "hello",
	})

	if got != "/users/42/posts/hello" {
		t.Errorf("GetUrl must fill each parameter, but is: %q", got)
	}
}

func TestGetUrlFillsAnOptionalParameter(t *testing.T) {
	t.Parallel()

	got := newUrl("/users{id?}").GetUrl(routeName, map[string]string{"id": "42"})

	if got != "/users42" {
		t.Errorf("GetUrl must fill an optional parameter, but is: %q", got)
	}
}

func TestGetUrlKeepsAParameterThatTheDataDoesNotName(t *testing.T) {
	t.Parallel()

	got := newUrl("/users/{id}").GetUrl(routeName, map[string]string{})

	if got != "/users/{id}" {
		t.Errorf("GetUrl must keep a parameter that the data does not name, but is: %q", got)
	}
}

func TestGetUrlReturnsThePathOfAStaticRoute(t *testing.T) {
	t.Parallel()

	if newUrl("/users").GetUrl(routeName, nil) != "/users" {
		t.Error("GetUrl must return the path of a static route, but did not")
	}
}

func TestGetUrlIsEmptyForAnUnknownRoute(t *testing.T) {
	t.Parallel()

	if newUrl("/users").GetUrl("missing", nil) != "" {
		t.Error("GetUrl must be empty for an unknown route, but is not")
	}
}

func TestTheUrlSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.UrlContract = newUrl("/users")

	if built.GetUrl(routeName, nil) != "/users" {
		t.Error("the contract must build the URL, but did not")
	}
}
