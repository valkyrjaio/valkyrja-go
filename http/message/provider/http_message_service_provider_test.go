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

	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/provider"
)

func TestTheProviderDefersItsBinding(t *testing.T) {
	t.Parallel()

	publishers := (&provider.HttpMessageServiceProvider{}).Publishers()

	if publishers[constant.ResponseFactoryContractServiceID] == nil {
		t.Error("the provider must defer the response factory, but deferred none")
	}
}

func TestTheProviderBindsTheResponseFactory(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishResponseFactory(container)

	resolved, err := container.GetSingleton(constant.ResponseFactoryContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the response factory, but reported: %v", err)
	}

	if _, isFactory := resolved.(contract.ResponseFactoryContract); !isFactory {
		t.Fatal("the container must hold a response factory, but held another value")
	}
}
