/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package fixtures

import (
	"errors"

	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
)

// ErrFailingContainerFixture is the failure that FailingContainerFixture
// reports.
var ErrFailingContainerFixture = errors.New("the parent container failed")

// FailingContainerFixture is a container that claims to hold every binding key
// and then fails to resolve it.
//
// A child container reads its parent through the contract, so a real parent
// cannot reach the child's error branches: a parent that reports a binding
// resolves it. The fixture separates the two answers, which is what makes each
// branch reachable.
//
// Each field says which report the fixture answers true for. A child reads the
// three reports in order, so a fixture that claims every binding never reaches
// the second lookup or the third one.
type FailingContainerFixture struct {
	*manager.Container

	ClaimsAlias             bool
	ClaimsService           bool
	ClaimsSingletonInstance bool
}

// NewFailingContainerFixture builds a container that claims every binding key.
func NewFailingContainerFixture() *FailingContainerFixture {
	return &FailingContainerFixture{
		Container:               manager.NewContainer(nil),
		ClaimsAlias:             true,
		ClaimsService:           true,
		ClaimsSingletonInstance: true,
	}
}

// NewFailingServiceContainerFixture builds a container that claims a factory
// for every binding key, and claims no alias and no instance.
func NewFailingServiceContainerFixture() *FailingContainerFixture {
	return &FailingContainerFixture{
		Container:     manager.NewContainer(nil),
		ClaimsService: true,
	}
}

// IsAlias reports whether the fixture claims each binding key as an alias.
func (c *FailingContainerFixture) IsAlias(_ string) bool {
	return c.ClaimsAlias
}

// IsService reports whether the fixture claims a factory for each binding key.
func (c *FailingContainerFixture) IsService(_ string) bool {
	return c.ClaimsService
}

// IsSingletonInstance reports whether the fixture claims an instance for each
// binding key.
func (c *FailingContainerFixture) IsSingletonInstance(_ string) bool {
	return c.ClaimsSingletonInstance
}

// GetAliased reports a failure.
func (c *FailingContainerFixture) GetAliased(_ string, _ []any) (any, error) {
	return nil, ErrFailingContainerFixture
}

// GetService reports a failure.
func (c *FailingContainerFixture) GetService(_ string, _ []any) (any, error) {
	return nil, ErrFailingContainerFixture
}

// GetSingleton reports a failure.
func (c *FailingContainerFixture) GetSingleton(_ string) (any, error) {
	return nil, ErrFailingContainerFixture
}
