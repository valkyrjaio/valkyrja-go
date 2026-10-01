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

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationdata "github.com/valkyrjaio/valkyrja-go/v26/application/data"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	messageconstant "github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	messagefactory "github.com/valkyrjaio/valkyrja-go/v26/http/message/factory"
	routingconstant "github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
	routingprovider "github.com/valkyrjaio/valkyrja-go/v26/http/routing/provider"
	"github.com/valkyrjaio/valkyrja-go/v26/http/server/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/server/provider"
)

// readHandler returns the request handler that the provider bound.
func readHandler(t *testing.T, container containercontract.ContainerContract) contract.RequestHandlerContract {
	t.Helper()

	resolved, err := container.GetSingleton(constant.RequestHandlerContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the request handler, but reported: %v", err)
	}

	built, isHandler := resolved.(contract.RequestHandlerContract)
	if !isHandler {
		t.Fatal("the container must hold a request handler, but held another value")
	}

	return built
}

func TestTheProviderDefersItsBinding(t *testing.T) {
	t.Parallel()

	publishers := (&provider.HttpServerServiceProvider{}).Publishers()

	if publishers[constant.RequestHandlerContractServiceID] == nil {
		t.Error("the provider must defer the request handler, but deferred none")
	}
}

func TestTheProviderBindsTheRequestHandler(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	config := applicationdata.NewHttpConfig()
	config.DebugMode = true

	container.SetSingleton(applicationconstant.HttpConfigContractServiceID, config)
	container.SetSingleton(
		messageconstant.ResponseFactoryContractServiceID,
		messagefactory.NewResponseFactory(),
	)

	routingprovider.PublishRouter(container)
	provider.PublishRequestHandler(container)

	readHandler(t, container)
}

func TestTheProviderBindsTheRequestHandlerWhereTheApplicationPublishesNothing(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishRequestHandler(container)

	readHandler(t, container)
}

func TestTheProviderIgnoresAServiceOfAnotherTypeUnderItsBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// A binding key that holds a value of another type is the developer's
	// error. The provider builds what it can rather than ending the process.
	container.SetSingleton(applicationconstant.HttpConfigContractServiceID, "not a config")
	container.SetSingleton(routingconstant.RouterContractServiceID, "not a router")
	container.SetSingleton(messageconstant.ResponseFactoryContractServiceID, "not a factory")

	provider.PublishRequestHandler(container)

	readHandler(t, container)
}
