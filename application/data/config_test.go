/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"os"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/application/data"
	"github.com/valkyrjaio/valkyrja-go/v26/application/fixtures"
)

func TestNewConfigTakesEachDefault(t *testing.T) {
	t.Parallel()

	config := data.NewConfig()

	if config.GetNamespace() != data.DefaultNamespace {
		t.Errorf("GetNamespace must be the default, but is: %s", config.GetNamespace())
	}

	if config.GetVersion() != constant.Version {
		t.Errorf("GetVersion must be the framework version, but is: %s", config.GetVersion())
	}

	if config.GetEnvironment() != data.DefaultEnvironment {
		t.Errorf("GetEnvironment must be the default, but is: %s", config.GetEnvironment())
	}

	if config.GetDebugMode() {
		t.Error("GetDebugMode must be false by default, but is true")
	}

	if config.GetTimezone() != data.DefaultTimezone {
		t.Errorf("GetTimezone must be the default, but is: %s", config.GetTimezone())
	}

	if config.GetKey() != data.DefaultKey {
		t.Errorf("GetKey must be the default, but is: %s", config.GetKey())
	}

	if config.GetDataPath() != data.DefaultDataPath {
		t.Errorf("GetDataPath must be the default, but is: %s", config.GetDataPath())
	}

	if config.GetDataNamespace() != data.DefaultDataNamespace {
		t.Errorf("GetDataNamespace must be the default, but is: %s", config.GetDataNamespace())
	}

	if len(config.GetCallbacks()) != 0 {
		t.Errorf("GetCallbacks must be empty by default, but holds: %d", len(config.GetCallbacks()))
	}
}

func TestNewConfigReadsTheWorkingDirectory(t *testing.T) {
	t.Parallel()

	working, err := os.Getwd()
	if err != nil {
		t.Fatalf("the test must read the working directory, but reported: %v", err)
	}

	if data.NewConfig().GetDir() != working {
		t.Errorf("GetDir must be the working directory, but is: %s", data.NewConfig().GetDir())
	}
}

func TestNewConfigHoldsEachProvider(t *testing.T) {
	t.Parallel()

	config := data.NewConfig(&fixtures.ComponentProviderFixture{}, &fixtures.ComponentProviderFixture{})

	if len(config.GetProviders()) != 2 {
		t.Errorf("GetProviders must hold each provider, but holds: %d", len(config.GetProviders()))
	}
}
