/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/data"
	"github.com/valkyrjaio/valkyrja-go/v26/container/fixtures"
)

const (
	aliasID     = "AliasID"
	serviceID   = "ServiceID"
	singletonID = "SingletonID"
	deferredID  = "DeferredID"
)

func TestNewContainerDataHoldsEachMap(t *testing.T) {
	t.Parallel()

	containerData := data.NewContainerData(
		map[string]string{aliasID: serviceID},
		map[string]contract.PublishFunc{deferredID: func(_ contract.ContainerContract) {}},
		map[string]contract.ServiceFactory{serviceID: fixtures.MakeServiceFixture},
		map[string]string{singletonID: singletonID},
	)

	if containerData.GetAliases()[aliasID] != serviceID {
		t.Errorf("GetAliases must hold the alias, but holds: %v", containerData.GetAliases())
	}

	if _, found := containerData.GetDeferredCallbacks()[deferredID]; !found {
		t.Error("GetDeferredCallbacks must hold the publisher, but does not")
	}

	if _, found := containerData.GetServices()[serviceID]; !found {
		t.Error("GetServices must hold the factory, but does not")
	}

	if containerData.GetSingletons()[singletonID] != singletonID {
		t.Errorf("GetSingletons must hold the binding, but holds: %v", containerData.GetSingletons())
	}
}

func TestNewContainerDataAcceptsNilForEachMap(t *testing.T) {
	t.Parallel()

	containerData := data.NewContainerData(nil, nil, nil, nil)

	if len(containerData.GetAliases()) != 0 {
		t.Errorf("GetAliases must be empty, but holds: %v", containerData.GetAliases())
	}

	if len(containerData.GetDeferredCallbacks()) != 0 {
		t.Errorf("GetDeferredCallbacks must be empty, but holds: %v", containerData.GetDeferredCallbacks())
	}

	if len(containerData.GetServices()) != 0 {
		t.Errorf("GetServices must be empty, but holds: %v", containerData.GetServices())
	}

	if len(containerData.GetSingletons()) != 0 {
		t.Errorf("GetSingletons must be empty, but holds: %v", containerData.GetSingletons())
	}
}

func TestNewContainerDataCopiesEachMapThatItReceives(t *testing.T) {
	t.Parallel()

	aliases := map[string]string{aliasID: serviceID}
	containerData := data.NewContainerData(aliases, nil, nil, nil)

	aliases["SecondAliasID"] = "SecondServiceID"

	if _, found := containerData.GetAliases()["SecondAliasID"]; found {
		t.Error("GetAliases must not follow a later write to the source map, but did")
	}
}

func TestEachGetterReturnsACopy(t *testing.T) {
	t.Parallel()

	containerData := data.NewContainerData(
		map[string]string{aliasID: serviceID},
		map[string]contract.PublishFunc{deferredID: func(_ contract.ContainerContract) {}},
		map[string]contract.ServiceFactory{serviceID: fixtures.MakeServiceFixture},
		map[string]string{singletonID: singletonID},
	)

	delete(containerData.GetAliases(), aliasID)
	delete(containerData.GetDeferredCallbacks(), deferredID)
	delete(containerData.GetServices(), serviceID)
	delete(containerData.GetSingletons(), singletonID)

	if len(containerData.GetAliases()) != 1 {
		t.Error("GetAliases must return a copy, but the delete reached the data")
	}

	if len(containerData.GetDeferredCallbacks()) != 1 {
		t.Error("GetDeferredCallbacks must return a copy, but the delete reached the data")
	}

	if len(containerData.GetServices()) != 1 {
		t.Error("GetServices must return a copy, but the delete reached the data")
	}

	if len(containerData.GetSingletons()) != 1 {
		t.Error("GetSingletons must return a copy, but the delete reached the data")
	}
}
