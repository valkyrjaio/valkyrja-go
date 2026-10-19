/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the event component's providers.
package provider

import (
	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/collection"
	"github.com/valkyrjaio/valkyrja-go/v26/event/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/event/data"
	"github.com/valkyrjaio/valkyrja-go/v26/event/dispatcher"
)

type EventServiceProvider struct{}

// Publishers returns a publisher for each binding key that the component defers.
func (p *EventServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.ListenerCollectionContractServiceID: PublishListenerCollection,
		constant.EventDispatcherContractServiceID:    PublishEventDispatcher,
		constant.EventDataServiceID:                  PublishEventData,
	}
}

// PublishListenerCollection binds the collection, and fills it from the
// component's state.
//
// An application that runs in debug mode reads each listener from its providers
// instead, so a change to a provider reaches the run without a new state.
func PublishListenerCollection(container containercontract.ContainerContract) {
	built := collection.NewListenerCollection()

	container.SetSingleton(constant.ListenerCollectionContractServiceID, built)

	if isDebugMode(container) {
		PublishEventData(container)

		return
	}

	built.SetFromData(getEventData(container))
}

// PublishEventDispatcher binds the dispatcher that runs the listeners of an
// event.
func PublishEventDispatcher(container containercontract.ContainerContract) {
	container.SetSingleton(
		constant.EventDispatcherContractServiceID,
		dispatcher.NewEventDispatcher(getCollection(container), container),
	)
}

// PublishEventData files every listener that a listener provider of the
// application registers, and binds the state that holds them.
func PublishEventData(container containercontract.ContainerContract) {
	built := getCollection(container)

	for _, listenerProvider := range getListenerProviders(container) {
		for _, listener := range listenerProvider.GetListeners() {
			built.AddListener(listener)
		}
	}

	container.SetSingleton(constant.EventDataServiceID, built.GetData())
}

// isDebugMode reports whether the application runs in debug mode.
func isDebugMode(container containercontract.ContainerContract) bool {
	app := getApplication(container)

	return app != nil && app.GetDebugMode()
}

// getEventData returns the component's state, and an empty state where the
// container holds none.
func getEventData(container containercontract.ContainerContract) contract.EventDataContract {
	resolved, err := container.GetSingleton(constant.EventDataServiceID)
	if err != nil {
		return data.NewEventData(nil, nil)
	}

	eventData, isData := resolved.(contract.EventDataContract)
	if !isData {
		return data.NewEventData(nil, nil)
	}

	return eventData
}

// getApplication returns the application, and nil where the container holds
// none.
func getApplication(
	container containercontract.ContainerContract,
) applicationcontract.ApplicationContract {
	resolved, err := container.GetSingleton(applicationconstant.ApplicationContractServiceID)
	if err != nil {
		return nil
	}

	app, isApplication := resolved.(applicationcontract.ApplicationContract)
	if !isApplication {
		return nil
	}

	return app
}

// getListenerProviders returns each listener provider that the application
// registers.
func getListenerProviders(
	container containercontract.ContainerContract,
) []contract.ListenerProviderContract {
	app := getApplication(container)
	if app == nil {
		return nil
	}

	return app.GetEventProviders()
}

// getCollection returns the collection that the component published.
func getCollection(container containercontract.ContainerContract) contract.ListenerCollectionContract {
	resolved, err := container.GetSingleton(constant.ListenerCollectionContractServiceID)
	if err != nil {
		return collection.NewListenerCollection()
	}

	built, isCollection := resolved.(contract.ListenerCollectionContract)
	if !isCollection {
		return collection.NewListenerCollection()
	}

	return built
}
