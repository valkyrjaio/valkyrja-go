/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package factory_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	messageconstant "github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	messagefactory "github.com/valkyrjaio/valkyrja-go/v26/http/message/factory"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/factory"
)

// urlFixture returns the URL that the test states, whatever route it is asked
type urlFixture struct {
	url string
}

// GetUrl returns the URL that the fixture holds.
func (u *urlFixture) GetUrl(_ string, _ map[string]string) string {
	return u.url
}

func TestTheFactorySendsTheClientToTheNamedRoute(t *testing.T) {
	t.Parallel()

	var built contract.RoutingResponseFactoryContract = factory.NewRoutingResponseFactory(
		&urlFixture{url: "https://example.com/users/1"},
		messagefactory.NewResponseFactory(),
	)

	redirect := built.CreateRouteRedirectResponse("users.show", map[string]string{"id": "1"}, 0, nil)

	if redirect.GetUri().GetPath() != "/users/1" {
		t.Errorf("the response must send the client to the route, but used: %q", redirect.GetUri().GetPath())
	}

	if redirect.GetStatusCode() != messageconstant.StatusCodeFound {
		t.Errorf("a redirect must take status 302 where the caller states none, but took: %d",
			redirect.GetStatusCode())
	}
}
