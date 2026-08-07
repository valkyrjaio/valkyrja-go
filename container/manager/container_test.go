/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package manager_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/data"
	"github.com/valkyrjaio/valkyrja-go/v26/container/fixtures"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/container/throwable/exception"
)

const (
	serviceID   = "ServiceID"
	singletonID = "SingletonID"
	aliasID     = "AliasID"
	missingID   = "MissingID"
)

func TestNewContainerAcceptsNilData(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	if container.Has(serviceID) {
		t.Error("Has must be false for an empty container, but is true")
	}
}

func TestNewContainerLoadsTheData(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(data.NewContainerData(
		map[string]string{aliasID: serviceID},
		map[string]contract.PublishFunc{"DeferredID": func(_ contract.ContainerContract) {}},
		map[string]contract.ServiceFactory{serviceID: fixtures.MakeServiceFixture},
		map[string]string{singletonID: singletonID},
	))

	if !container.IsAlias(aliasID) {
		t.Error("IsAlias must be true for the loaded alias, but is false")
	}

	if !container.IsService(serviceID) {
		t.Error("IsService must be true for the loaded factory, but is false")
	}

	if !container.IsSingletonBinding(singletonID) {
		t.Error("IsSingletonBinding must be true for the loaded binding, but is false")
	}

	if !container.IsDeferred("DeferredID") {
		t.Error("IsDeferred must be true for the loaded publisher, but is false")
	}
}

func TestGetDataReturnsWhatTheContainerHolds(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)
	container.BindAlias(aliasID, serviceID)
	container.BindSingleton(singletonID, fixtures.MakeSingletonFixture)

	containerData := container.GetData()

	if containerData.GetAliases()[aliasID] != serviceID {
		t.Errorf("GetAliases must hold the alias, but holds: %v", containerData.GetAliases())
	}

	if _, found := containerData.GetServices()[serviceID]; !found {
		t.Error("GetServices must hold the factory, but does not")
	}

	if containerData.GetSingletons()[singletonID] != singletonID {
		t.Errorf("GetSingletons must hold the binding, but holds: %v", containerData.GetSingletons())
	}
}

func TestSetFromDataMergesIntoTheContainer(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)

	container.SetFromData(data.NewContainerData(
		nil,
		nil,
		map[string]contract.ServiceFactory{"SecondServiceID": fixtures.MakeServiceFixture},
		nil,
	))

	if !container.IsService(serviceID) {
		t.Error("IsService must stay true for the first factory, but is false")
	}

	if !container.IsService("SecondServiceID") {
		t.Error("IsService must be true for the merged factory, but is false")
	}
}

func TestHasReportsEachKindOfBinding(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		bind func(container contract.ContainerContract)
		id   string
	}{
		"a service": {
			bind: func(container contract.ContainerContract) {
				container.Bind(serviceID, fixtures.MakeServiceFixture)
			},
			id: serviceID,
		},
		"an alias": {
			bind: func(container contract.ContainerContract) {
				container.BindAlias(aliasID, serviceID)
			},
			id: aliasID,
		},
		"a singleton": {
			bind: func(container contract.ContainerContract) {
				container.SetSingleton(singletonID, &fixtures.SingletonFixture{})
			},
			id: singletonID,
		},
		"a deferred publisher": {
			bind: func(container contract.ContainerContract) {
				_ = container.Register(&fixtures.ProviderFixture{})
			},
			id: fixtures.ProvidedID,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			container := manager.NewContainer(nil)
			test.bind(container)

			if !container.Has(test.id) {
				t.Errorf("Has must be true for %s, but is false", name)
			}
		})
	}
}

func TestHasIsFalseForAnUnknownBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	if container.Has(missingID) {
		t.Error("Has must be false for an unknown binding key, but is true")
	}
}

func TestBindReturnsTheContainerForChaining(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	chained := container.
		Bind(serviceID, fixtures.MakeServiceFixture).
		BindAlias(aliasID, serviceID).
		BindSingleton(singletonID, fixtures.MakeSingletonFixture).
		SetSingleton("InstanceID", &fixtures.SingletonFixture{})

	if chained != contract.ContainerContract(container) {
		t.Error("each binding method must return the container itself, but did not")
	}
}

func TestBindMarksTheBindingKeyPublished(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)

	if !container.IsPublished(serviceID) {
		t.Error("IsPublished must be true after Bind, but is false")
	}
}

func TestSetSingletonMarksTheBindingKeyPublished(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(singletonID, &fixtures.SingletonFixture{})

	if !container.IsPublished(singletonID) {
		t.Error("IsPublished must be true after SetSingleton, but is false")
	}
}

func TestEachReportIsFalseForAnUnknownBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	reports := map[string]func(id string) bool{
		"IsAlias":             container.IsAlias,
		"IsService":           container.IsService,
		"IsSingleton":         container.IsSingleton,
		"IsSingletonBinding":  container.IsSingletonBinding,
		"IsSingletonInstance": container.IsSingletonInstance,
		"IsDeferred":          container.IsDeferred,
		"IsPublished":         container.IsPublished,
	}

	for name, report := range reports {
		if report(missingID) {
			t.Errorf("%s must be false for an unknown binding key, but is true", name)
		}
	}
}

func TestIsSingletonReportsABinding(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.BindSingleton(singletonID, fixtures.MakeSingletonFixture)

	if !container.IsSingleton(singletonID) {
		t.Error("IsSingleton must be true for a singleton binding, but is false")
	}
}

func TestIsSingletonReportsAnInstance(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(singletonID, &fixtures.SingletonFixture{})

	if !container.IsSingleton(singletonID) {
		t.Error("IsSingleton must be true for a singleton instance, but is false")
	}
}

func TestGetResolvesAService(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)

	resolved, err := container.Get(serviceID, []any{"first"}, constant.NewInstanceOrThrowException)
	if err != nil {
		t.Fatalf("Get must resolve the service, but reported: %v", err)
	}

	service, isService := resolved.(*fixtures.ServiceFixture)
	if !isService {
		t.Fatalf("Get must return a ServiceFixture, but returned: %T", resolved)
	}

	if service.GetContainer() != contract.ContainerContract(container) {
		t.Error("the factory must receive the container itself, but did not")
	}

	if len(service.Arguments) != 1 || service.Arguments[0] != "first" {
		t.Errorf("the factory must receive the arguments, but received: %v", service.Arguments)
	}
}

func TestGetCallsTheFactoryForEachResolution(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)

	first, _ := container.Get(serviceID, nil, constant.NewInstanceOrThrowException)
	second, _ := container.Get(serviceID, nil, constant.NewInstanceOrThrowException)

	if first == second {
		t.Error("Get must call the factory for each resolution, but returned one instance")
	}
}

func TestGetResolvesASingletonOnce(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.BindSingleton(singletonID, fixtures.MakeSingletonFixture)

	first, err := container.Get(singletonID, nil, constant.NewInstanceOrThrowException)
	if err != nil {
		t.Fatalf("Get must resolve the singleton, but reported: %v", err)
	}

	second, _ := container.Get(singletonID, nil, constant.NewInstanceOrThrowException)

	if first != second {
		t.Error("Get must return one instance for a singleton, but returned two")
	}

	if !container.IsSingletonInstance(singletonID) {
		t.Error("IsSingletonInstance must be true after the first resolution, but is false")
	}
}

func TestGetResolvesAnAlias(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)
	container.BindAlias(aliasID, serviceID)

	resolved, err := container.Get(aliasID, nil, constant.NewInstanceOrThrowException)
	if err != nil {
		t.Fatalf("Get must resolve the alias, but reported: %v", err)
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("Get must return what the alias points to, but returned: %T", resolved)
	}
}

func TestGetReportsAnUnknownBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	_, err := container.Get(missingID, nil, constant.NewInstanceOrThrowException)

	assertInvalidReference(t, err, missingID)
}

func TestGetReportsAnAliasThatPointsAtNothing(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.BindAlias(aliasID, missingID)

	_, err := container.Get(aliasID, nil, constant.NewInstanceOrThrowException)

	assertInvalidReference(t, err, missingID)
}

func TestGetReportsAnUnknownBindingKeyInEachMode(t *testing.T) {
	t.Parallel()

	modes := map[string]constant.InvalidReferenceMode{
		"NewInstanceOrThrowException": constant.NewInstanceOrThrowException,
		"ThrowException":              constant.ThrowException,
	}

	for name, mode := range modes {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			container := manager.NewContainer(nil)

			_, err := container.Get(missingID, nil, mode)

			assertInvalidReference(t, err, missingID)
		})
	}
}

func TestGetAliasedResolvesTheAlias(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)
	container.BindAlias(aliasID, serviceID)

	resolved, err := container.GetAliased(aliasID, nil)
	if err != nil {
		t.Fatalf("GetAliased must resolve the alias, but reported: %v", err)
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("GetAliased must return what the alias points to, but returned: %T", resolved)
	}
}

func TestGetAliasedReportsAKeyThatIsNoAlias(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	_, err := container.GetAliased(missingID, nil)

	assertInvalidReference(t, err, missingID)
}

func TestGetAliasedReportsAnAliasThatPointsAtNothing(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.BindAlias(aliasID, missingID)

	_, err := container.GetAliased(aliasID, nil)

	assertInvalidReference(t, err, missingID)
}

func TestGetServiceResolvesTheFactory(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.Bind(serviceID, fixtures.MakeServiceFixture)

	resolved, err := container.GetService(serviceID, nil)
	if err != nil {
		t.Fatalf("GetService must resolve the factory, but reported: %v", err)
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("GetService must return a ServiceFixture, but returned: %T", resolved)
	}
}

func TestGetServiceReportsAKeyWithNoFactory(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	_, err := container.GetService(missingID, nil)

	assertInvalidReference(t, err, missingID)
}

func TestGetSingletonResolvesTheInstance(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	singleton := &fixtures.SingletonFixture{}
	container.SetSingleton(singletonID, singleton)

	resolved, err := container.GetSingleton(singletonID)
	if err != nil {
		t.Fatalf("GetSingleton must resolve the instance, but reported: %v", err)
	}

	if resolved != any(singleton) {
		t.Errorf("GetSingleton must return the instance that was set, but returned: %v", resolved)
	}
}

func TestGetSingletonReportsAKeyThatIsNoSingleton(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	_, err := container.GetSingleton(missingID)

	assertInvalidReference(t, err, missingID)
}

func TestGetSingletonReportsABindingWithNoFactory(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(data.NewContainerData(
		nil,
		nil,
		nil,
		map[string]string{singletonID: singletonID},
	))

	_, err := container.GetSingleton(singletonID)

	assertInvalidReference(t, err, singletonID)
}

func TestRegisterRecordsEachPublisher(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	err := container.Register(&fixtures.ProviderFixture{})
	if err != nil {
		t.Fatalf("Register must record each publisher, but reported: %v", err)
	}

	if !container.IsDeferred(fixtures.ProvidedID) {
		t.Error("IsDeferred must be true for the first binding key, but is false")
	}

	if !container.IsDeferred(fixtures.ProvidedSecondaryID) {
		t.Error("IsDeferred must be true for the second binding key, but is false")
	}
}

func TestRegisterReportsAProviderWithNoPublisher(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	err := container.Register(&fixtures.InvalidProviderFixture{})

	var target *exception.ContainerInvalidPublishCallbackError

	if !errors.As(err, &target) {
		t.Fatalf("Register must report an invalid publish callback, but reported: %v", err)
	}

	if target.GetID() != fixtures.InvalidProvidedID {
		t.Errorf("the error must name the binding key, but names: %s", target.GetID())
	}
}

func TestPublishRunsThePublisher(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	provider := &fixtures.ProviderFixture{}

	registerErr := container.Register(provider)
	if registerErr != nil {
		t.Fatalf("Register must record each publisher, but reported: %v", registerErr)
	}

	container.Publish(fixtures.ProvidedID)

	if !provider.PublishCalled {
		t.Error("Publish must run the publisher, but did not")
	}

	if !container.IsPublished(fixtures.ProvidedID) {
		t.Error("IsPublished must be true after Publish, but is false")
	}
}

func TestPublishDoesNothingForAKeyWithNoPublisher(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	container.Publish(missingID)

	if container.IsPublished(missingID) {
		t.Error("IsPublished must stay false where no publisher defers the key, but is true")
	}
}

func TestGetPublishesADeferredBindingKey(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	provider := &fixtures.ProviderFixture{}

	registerErr := container.Register(provider)
	if registerErr != nil {
		t.Fatalf("Register must record each publisher, but reported: %v", registerErr)
	}

	resolved, err := container.Get(fixtures.ProvidedID, nil, constant.NewInstanceOrThrowException)
	if err != nil {
		t.Fatalf("Get must publish and resolve the binding key, but reported: %v", err)
	}

	if !provider.PublishCalled {
		t.Error("Get must run the deferred publisher, but did not")
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("Get must return what the publisher set, but returned: %T", resolved)
	}
}

func TestGetPublishesADeferredBindingKeyOnlyOnce(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	provider := &fixtures.ProviderFixture{}

	registerErr := container.Register(provider)
	if registerErr != nil {
		t.Fatalf("Register must record each publisher, but reported: %v", registerErr)
	}

	_, err := container.Get(fixtures.ProvidedID, nil, constant.NewInstanceOrThrowException)
	if err != nil {
		t.Fatalf("Get must publish and resolve the binding key, but reported: %v", err)
	}

	provider.PublishCalled = false

	_, secondErr := container.Get(fixtures.ProvidedID, nil, constant.NewInstanceOrThrowException)
	if secondErr != nil {
		t.Fatalf("Get must resolve the published binding key, but reported: %v", secondErr)
	}

	if provider.PublishCalled {
		t.Error("Get must run the deferred publisher once, but ran it again")
	}
}

// assertInvalidReference fails the test where the error is not an invalid
// reference for the binding key.
func assertInvalidReference(t *testing.T, err error, id string) {
	t.Helper()

	var target *exception.ContainerInvalidReferenceError

	if !errors.As(err, &target) {
		t.Fatalf("the error must be an invalid reference, but is: %v", err)
	}

	if target.GetID() != id {
		t.Errorf("the error must name %s, but names: %s", id, target.GetID())
	}
}

func TestPublishRunsEachPublisherThatTheProviderDefers(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	provider := &fixtures.ProviderFixture{}

	registerErr := container.Register(provider)
	if registerErr != nil {
		t.Fatalf("Register must record each publisher, but reported: %v", registerErr)
	}

	container.Publish(fixtures.ProvidedID)
	container.Publish(fixtures.ProvidedSecondaryID)

	if !provider.PublishCalled || !provider.PublishSecondaryCalled {
		t.Error("Publish must run each publisher that the provider defers, but did not")
	}
}
