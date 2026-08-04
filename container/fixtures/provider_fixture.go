/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package fixtures

import (
	"github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// The binding keys that ProviderFixture defers.
const (
	// ProvidedID is the first binding key that ProviderFixture defers.
	ProvidedID = "ProvidedID"

	// ProvidedSecondaryID is the second binding key that ProviderFixture
	// defers.
	ProvidedSecondaryID = "ProvidedSecondaryID"

	// InvalidProvidedID is the binding key that InvalidProviderFixture defers
	// without a publisher.
	InvalidProvidedID = "InvalidProvidedID"
)

// ProviderFixture is a service provider that records each publisher that the
// container runs.
//
// The other ports record the calls in a static field. This port records them on
// the instance, so a test that runs in parallel with another test reads its own
// provider rather than a shared one.
type ProviderFixture struct {
	PublishCalled          bool
	PublishSecondaryCalled bool
}

// Publishers returns a publisher for each binding key that the provider defers.
func (p *ProviderFixture) Publishers() map[string]contract.PublishFunc {
	return map[string]contract.PublishFunc{
		ProvidedID: func(container contract.ContainerContract) {
			p.PublishCalled = true

			container.SetSingleton(ProvidedID, &ServiceFixture{Container: container})
		},
		ProvidedSecondaryID: func(_ contract.ContainerContract) {
			p.PublishSecondaryCalled = true
		},
	}
}

// InvalidProviderFixture is a service provider that defers a binding key
// without a publisher for it.
type InvalidProviderFixture struct{}

// Publishers returns a nil publisher, which is the one value that reaches the
// container without a publisher.
func (p *InvalidProviderFixture) Publishers() map[string]contract.PublishFunc {
	return map[string]contract.PublishFunc{
		InvalidProvidedID: nil,
	}
}
