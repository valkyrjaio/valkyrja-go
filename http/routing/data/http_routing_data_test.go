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

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
)

// newRoutingData builds routing data that holds one static route.
func newRoutingData() *data.HttpRoutingData {
	return data.NewHttpRoutingData(
		map[string]contract.RouteContract{routeName: newRoute()},
		map[constant.RequestMethod]map[string]string{
			constant.RequestMethodGet: {routePath: routeName},
		},
		map[constant.RequestMethod]map[string]string{
			constant.RequestMethodGet: {dynamicRegex: routeName},
		},
	)
}

func TestNewHttpRoutingDataHoldsEachMap(t *testing.T) {
	t.Parallel()

	routingData := newRoutingData()

	if routingData.GetRoutes()[routeName] == nil {
		t.Error("the data must hold the route, but did not")
	}

	if routingData.GetPaths()[constant.RequestMethodGet][routePath] != routeName {
		t.Error("the data must hold the path, but did not")
	}

	if routingData.GetRegexes()[constant.RequestMethodGet][dynamicRegex] != routeName {
		t.Error("the data must hold the regular expression, but did not")
	}
}

func TestNewHttpRoutingDataAcceptsNilForEachMap(t *testing.T) {
	t.Parallel()

	routingData := data.NewHttpRoutingData(nil, nil, nil)

	if len(routingData.GetRoutes()) != 0 || len(routingData.GetPaths()) != 0 || len(routingData.GetRegexes()) != 0 {
		t.Error("the data must be empty, but was not")
	}
}

func TestEachGetterReturnsACopy(t *testing.T) {
	t.Parallel()

	routingData := newRoutingData()

	delete(routingData.GetRoutes(), routeName)
	delete(routingData.GetPaths(), constant.RequestMethodGet)
	delete(routingData.GetPaths()[constant.RequestMethodGet], routePath)
	delete(routingData.GetRegexes(), constant.RequestMethodGet)

	if len(routingData.GetRoutes()) != 1 {
		t.Error("GetRoutes must return a copy, but the delete reached the data")
	}

	if routingData.GetPaths()[constant.RequestMethodGet][routePath] != routeName {
		t.Error("GetPaths must copy the inner map, but the delete reached the data")
	}

	if len(routingData.GetRegexes()) != 1 {
		t.Error("GetRegexes must return a copy, but the delete reached the data")
	}
}

func TestNewHttpRoutingDataCopiesItsSourceMaps(t *testing.T) {
	t.Parallel()

	paths := map[constant.RequestMethod]map[string]string{
		constant.RequestMethodGet: {routePath: routeName},
	}

	routingData := data.NewHttpRoutingData(nil, paths, nil)

	paths[constant.RequestMethodGet]["/other"] = "other"

	if _, found := routingData.GetPaths()[constant.RequestMethodGet]["/other"]; found {
		t.Error("the data must not follow a later write to the source map, but did")
	}
}

func TestTheRoutingDataSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var routingData contract.HttpRoutingDataContract = newRoutingData()

	if len(routingData.GetRoutes()) != 1 {
		t.Error("the contract must read the routes, but did not")
	}
}
