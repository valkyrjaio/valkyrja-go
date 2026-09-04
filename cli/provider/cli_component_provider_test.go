/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package provider_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/provider"
)

func TestTheComponentProviderNamesEveryProviderOfTheComponent(t *testing.T) {
	t.Parallel()

	built := &provider.CliComponentProvider{}

	if len(built.GetComponentProviders(nil)) != 1 {
		t.Errorf("the CLI must name the container, but named: %d", len(built.GetComponentProviders(nil)))
	}

	if len(built.GetContainerProviders(nil)) != 3 {
		t.Errorf("the CLI must name each of its service providers, but named: %d",
			len(built.GetContainerProviders(nil)))
	}

	if len(built.GetCliProviders(nil)) != 1 {
		t.Errorf("the CLI must name its route provider, but named: %d", len(built.GetCliProviders(nil)))
	}

	if len(built.GetEventProviders(nil)) != 0 || len(built.GetHttpProviders(nil)) != 0 {
		t.Error("the CLI must name no listener provider and no HTTP route provider, but named one")
	}
}
