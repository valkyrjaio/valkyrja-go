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
)

func TestNewChildContainerAcceptsNilData(t *testing.T) {
	t.Parallel()

	child := manager.NewChildContainer(manager.NewContainer(nil), nil)

	if child.Has(serviceID) {
		t.Error("Has must be false for an empty child, but is true")
	}
}

func TestNewChildContainerTakesTheSingletonsAndThePublishers(t *testing.T) {
	t.Parallel()

	child := manager.NewChildContainer(manager.NewContainer(nil), data.NewContainerData(
		map[string]string{aliasID: serviceID},
		map[string]contract.PublishFunc{"DeferredID": func(_ contract.ContainerContract) {}},
		map[string]contract.ServiceFactory{serviceID: fixtures.MakeServiceFixture},
		map[string]string{singletonID: singletonID},
	))

	if !child.IsSingletonBinding(singletonID) {
		t.Error("IsSingletonBinding must be true for the loaded binding, but is false")
	}

	if !child.IsDeferred("DeferredID") {
		t.Error("IsDeferred must be true for the loaded publisher, but is false")
	}

	if child.IsAlias(aliasID) {
		t.Error("the child must not take the aliases from the data, but did")
	}

	if child.IsService(serviceID) {
		t.Error("the child must not take the factories from the data, but did")
	}
}

func TestTheChildReadsEachReportFromItsParent(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		bind   func(container contract.ContainerContract)
		report func(child *manager.ChildContainer, id string) bool
		id     string
	}{
		"IsAlias": {
			bind: func(container contract.ContainerContract) {
				container.BindAlias(aliasID, serviceID)
			},
			report: (*manager.ChildContainer).IsAlias,
			id:     aliasID,
		},
		"IsService": {
			bind: func(container contract.ContainerContract) {
				container.Bind(serviceID, fixtures.MakeServiceFixture)
			},
			report: (*manager.ChildContainer).IsService,
			id:     serviceID,
		},
		"IsSingletonInstance": {
			bind: func(container contract.ContainerContract) {
				container.SetSingleton(singletonID, &fixtures.SingletonFixture{})
			},
			report: (*manager.ChildContainer).IsSingletonInstance,
			id:     singletonID,
		},
		"IsDeferred": {
			bind: func(container contract.ContainerContract) {
				_ = container.Register(&fixtures.ProviderFixture{})
			},
			report: (*manager.ChildContainer).IsDeferred,
			id:     fixtures.ProvidedID,
		},
		"IsPublished": {
			bind: func(container contract.ContainerContract) {
				container.Bind(serviceID, fixtures.MakeServiceFixture)
			},
			report: (*manager.ChildContainer).IsPublished,
			id:     serviceID,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			parent := manager.NewContainer(nil)
			child := manager.NewChildContainer(parent, nil)

			if test.report(child, test.id) {
				t.Errorf("%s must be false before the parent binds, but is true", name)
			}

			test.bind(parent)

			if !test.report(child, test.id) {
				t.Errorf("%s must read the parent's binding, but is false", name)
			}
		})
	}
}

func TestTheChildPrefersItsOwnBinding(t *testing.T) {
	t.Parallel()

	parent := manager.NewContainer(nil)
	parent.Bind(serviceID, fixtures.MakeSingletonFixture)

	child := manager.NewChildContainer(parent, nil)
	child.Bind(serviceID, fixtures.MakeServiceFixture)

	resolved, err := child.Get(serviceID, nil, constant.NewInstanceOrThrowException)
	if err != nil {
		t.Fatalf("Get must resolve the child's factory, but reported: %v", err)
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("Get must use the child's own factory, but returned: %T", resolved)
	}
}

func TestTheChildResolvesTheParentService(t *testing.T) {
	t.Parallel()

	parent := manager.NewContainer(nil)
	parent.Bind(serviceID, fixtures.MakeServiceFixture)

	child := manager.NewChildContainer(parent, nil)

	resolved, err := child.Get(serviceID, nil, constant.NewInstanceOrThrowException)
	if err != nil {
		t.Fatalf("Get must resolve the parent's factory, but reported: %v", err)
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("Get must return what the parent's factory built, but returned: %T", resolved)
	}
}

func TestTheChildResolvesTheParentSingleton(t *testing.T) {
	t.Parallel()

	parent := manager.NewContainer(nil)
	singleton := &fixtures.SingletonFixture{}
	parent.SetSingleton(singletonID, singleton)

	child := manager.NewChildContainer(parent, nil)

	resolved, err := child.GetSingleton(singletonID)
	if err != nil {
		t.Fatalf("GetSingleton must resolve the parent's instance, but reported: %v", err)
	}

	if resolved != any(singleton) {
		t.Errorf("GetSingleton must return the parent's instance, but returned: %v", resolved)
	}
}

func TestTheChildPrefersItsOwnSingletonInstance(t *testing.T) {
	t.Parallel()

	parent := manager.NewContainer(nil)
	parent.SetSingleton(singletonID, &fixtures.SingletonFixture{})

	child := manager.NewChildContainer(parent, nil)
	own := &fixtures.SingletonFixture{}
	child.SetSingleton(singletonID, own)

	resolved, err := child.GetSingleton(singletonID)
	if err != nil {
		t.Fatalf("GetSingleton must resolve the child's instance, but reported: %v", err)
	}

	if resolved != any(own) {
		t.Error("GetSingleton must return the child's own instance, but returned the parent's")
	}
}

func TestTheChildResolvesTheParentAlias(t *testing.T) {
	t.Parallel()

	parent := manager.NewContainer(nil)
	parent.Bind(serviceID, fixtures.MakeServiceFixture)
	parent.BindAlias(aliasID, serviceID)

	child := manager.NewChildContainer(parent, nil)

	resolved, err := child.GetAliased(aliasID, nil)
	if err != nil {
		t.Fatalf("GetAliased must resolve the parent's alias, but reported: %v", err)
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("GetAliased must return what the parent's alias points to, but returned: %T", resolved)
	}
}

func TestTheChildPrefersItsOwnAlias(t *testing.T) {
	t.Parallel()

	parent := manager.NewContainer(nil)
	parent.Bind(serviceID, fixtures.MakeSingletonFixture)
	parent.BindAlias(aliasID, serviceID)

	child := manager.NewChildContainer(parent, nil)
	child.Bind("ChildServiceID", fixtures.MakeServiceFixture)
	child.BindAlias(aliasID, "ChildServiceID")

	resolved, err := child.GetAliased(aliasID, nil)
	if err != nil {
		t.Fatalf("GetAliased must resolve the child's alias, but reported: %v", err)
	}

	if _, isService := resolved.(*fixtures.ServiceFixture); !isService {
		t.Errorf("GetAliased must use the child's own alias, but returned: %T", resolved)
	}
}

func TestTheChildBuildsItsOwnSingletonFromAParentFactory(t *testing.T) {
	t.Parallel()

	parent := manager.NewContainer(nil)
	parent.Bind(singletonID, fixtures.MakeSingletonFixture)

	child := manager.NewChildContainer(parent, data.NewContainerData(
		nil,
		nil,
		nil,
		map[string]string{singletonID: singletonID},
	))

	first, err := child.GetSingleton(singletonID)
	if err != nil {
		t.Fatalf("GetSingleton must build the singleton, but reported: %v", err)
	}

	second, _ := child.GetSingleton(singletonID)

	if first != second {
		t.Error("GetSingleton must return one instance for the child's singleton, but returned two")
	}

	if parent.IsSingletonInstance(singletonID) {
		t.Error("the child must not record its instance on the parent, but did")
	}
}

func TestTheChildReportsAnUnknownBindingKey(t *testing.T) {
	t.Parallel()

	child := manager.NewChildContainer(manager.NewContainer(nil), nil)

	_, err := child.Get(missingID, nil, constant.NewInstanceOrThrowException)

	assertInvalidReference(t, err, missingID)
}

func TestTheChildCarriesEachParentFailure(t *testing.T) {
	t.Parallel()

	tests := map[string]func(child *manager.ChildContainer) error{
		"GetSingleton": func(child *manager.ChildContainer) error {
			_, err := child.GetSingleton(missingID)

			return err
		},
		"GetService": func(child *manager.ChildContainer) error {
			_, err := child.GetService(missingID, nil)

			return err
		},
		"GetAliased": func(child *manager.ChildContainer) error {
			_, err := child.GetAliased(missingID, nil)

			return err
		},
		"Get": func(child *manager.ChildContainer) error {
			_, err := child.Get(missingID, nil, constant.NewInstanceOrThrowException)

			return err
		},
	}

	for name, resolve := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			child := manager.NewChildContainer(fixtures.NewFailingContainerFixture(), nil)

			err := resolve(child)

			if !errors.Is(err, fixtures.ErrFailingContainerFixture) {
				t.Errorf("%s must carry the parent's failure, but reported: %v", name, err)
			}
		})
	}
}

func TestTheChildCarriesTheParentServiceFailureFromGet(t *testing.T) {
	t.Parallel()

	child := manager.NewChildContainer(fixtures.NewFailingServiceContainerFixture(), nil)

	_, err := child.Get(missingID, nil, constant.NewInstanceOrThrowException)

	if !errors.Is(err, fixtures.ErrFailingContainerFixture) {
		t.Errorf("Get must carry the parent's factory failure, but reported: %v", err)
	}
}

func TestTheChildCarriesTheParentServiceFailureFromItsOwnSingletonBinding(t *testing.T) {
	t.Parallel()

	child := manager.NewChildContainer(fixtures.NewFailingServiceContainerFixture(), data.NewContainerData(
		nil,
		nil,
		nil,
		map[string]string{singletonID: singletonID},
	))

	_, err := child.GetSingleton(singletonID)

	if !errors.Is(err, fixtures.ErrFailingContainerFixture) {
		t.Errorf("GetSingleton must carry the parent's factory failure, but reported: %v", err)
	}
}
