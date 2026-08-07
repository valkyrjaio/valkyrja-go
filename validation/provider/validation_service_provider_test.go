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

	containerconstant "github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/provider"
	"github.com/valkyrjaio/valkyrja-go/v26/validation/rule"
)

func TestTheProviderDefersItsBinding(t *testing.T) {
	t.Parallel()

	publishers := (&provider.ValidationServiceProvider{}).Publishers()

	if publishers[constant.ValidatorContractServiceID] == nil {
		t.Error("the provider must defer the validator, but deferred none")
	}
}

func TestTheProviderBindsAValidatorThatHoldsNoRule(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishValidator(container)

	built := resolveValidator(t, container)

	if !built.ValidateRules() {
		t.Error("a validator that holds no rule must report that everything passes, but did not")
	}
}

func TestEachResolveBuildsItsOwnValidator(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	provider.PublishValidator(container)

	// A validator holds the rules and the messages of one validation in fields
	// that it writes, so one shared instance across two requests races on them.
	first := resolveValidator(t, container)
	second := resolveValidator(t, container)

	if first == second {
		t.Error("each resolve must build its own validator, but returned the same one")
	}

	first.SetRules(map[string][]contract.RuleContract{"name": {rule.NewRequired("")}})
	first.ValidateRules()

	if second.HasFirstErrorMessage() {
		t.Error("one validator must not reach another, but the messages crossed")
	}
}

// resolveValidator returns a validator that the container built.
func resolveValidator(
	t *testing.T,
	container containercontract.ContainerContract,
) contract.ValidatorContract {
	t.Helper()

	resolved, err := container.Get(
		constant.ValidatorContractServiceID,
		nil,
		containerconstant.NewInstanceOrThrowException,
	)
	if err != nil {
		t.Fatalf("the provider must bind the validator, but reported: %v", err)
	}

	built, isValidator := resolved.(contract.ValidatorContract)
	if !isValidator {
		t.Fatal("the container must hold a validator, but held another value")
	}

	return built
}

func TestTheComponentProviderNamesEveryProviderOfTheComponent(t *testing.T) {
	t.Parallel()

	built := &provider.ValidationComponentProvider{}

	if len(built.GetComponentProviders(nil)) != 1 {
		t.Errorf("validation must name the container, but named: %d", len(built.GetComponentProviders(nil)))
	}

	if len(built.GetContainerProviders(nil)) != 1 {
		t.Errorf("validation must name its service provider, but named: %d",
			len(built.GetContainerProviders(nil)))
	}

	if len(built.GetEventProviders(nil)) != 0 ||
		len(built.GetCliProviders(nil)) != 0 ||
		len(built.GetHttpProviders(nil)) != 0 {
		t.Error("validation must name no listener provider and no route provider, but named one")
	}
}
