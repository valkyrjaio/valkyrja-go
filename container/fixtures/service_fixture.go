/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package fixtures holds the reusable doubles that the container component's
// tests build on. It mirrors the source tree, and each type carries the
// `Fixture` suffix.
package fixtures

import (
	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// ServiceFixture is a service that records the container that built it.
type ServiceFixture struct {
	Container contract.ContainerContract
	Arguments []any
}

// MakeServiceFixture is a service factory that builds a ServiceFixture.
func MakeServiceFixture(container contract.ContainerContract, arguments []any) any {
	return &ServiceFixture{
		Container: container,
		Arguments: arguments,
	}
}

// GetContainer returns the container that built the service.
func (s *ServiceFixture) GetContainer() contract.ContainerContract {
	return s.Container
}

// SingletonFixture is a service that records nothing.
type SingletonFixture struct{}

// MakeSingletonFixture is a service factory that builds a SingletonFixture.
func MakeSingletonFixture(_ contract.ContainerContract, _ []any) any {
	return &SingletonFixture{}
}
