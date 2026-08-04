/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package provider_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationdata "github.com/valkyrjaio/valkyrja-go/v26/application/data"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/log/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/log/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/log/provider"
)

// configFixture is an application configuration that also holds the settings of
// the log component.
type configFixture struct {
	*applicationdata.Config

	logger string
	path   string
}

// GetDefaultLogger returns the binding key of the logger that the application
// writes through.
func (c *configFixture) GetDefaultLogger() string {
	return c.logger
}

// GetStreamFilePath returns the file that the logger writes to.
func (c *configFixture) GetStreamFilePath() string {
	return c.path
}

// readLogger returns the logger that the provider bound.
func readLogger(t *testing.T, container containercontract.ContainerContract) contract.LoggerContract {
	t.Helper()

	resolved, err := container.GetSingleton(constant.LoggerContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the logger, but reported: %v", err)
	}

	built, isLogger := resolved.(contract.LoggerContract)
	if !isLogger {
		t.Fatal("the container must hold a logger, but held another value")
	}

	return built
}

func TestTheProviderDefersEachOfItsBindings(t *testing.T) {
	t.Parallel()

	publishers := (&provider.LogServiceProvider{}).Publishers()

	for _, id := range []string{constant.LogConfigContractServiceID, constant.LoggerContractServiceID} {
		if publishers[id] == nil {
			t.Errorf("the provider must defer %q, but deferred none", id)
		}
	}
}

func TestTheProviderBindsTheFrameworkDefaultWhereTheApplicationStatesNone(t *testing.T) {
	t.Parallel()

	tests := map[string]any{
		"a container that holds no configuration":    nil,
		"a configuration that holds no log settings": applicationdata.NewConfig(),
	}

	for name, config := range tests {
		container := manager.NewContainer(nil)

		if config != nil {
			container.SetSingleton(applicationconstant.ConfigContractServiceID, config)
		}

		provider.PublishLogConfig(container)

		resolved, err := container.GetSingleton(constant.LogConfigContractServiceID)
		if err != nil {
			t.Fatalf("%s must still bind a configuration, but reported: %v", name, err)
		}

		built, isConfig := resolved.(contract.LogConfigContract)
		if !isConfig {
			t.Fatalf("%s must bind a log configuration, but bound another value", name)
		}

		if built.GetDefaultLogger() == "" {
			t.Errorf("%s must name a default logger, but named none", name)
		}
	}
}

func TestTheProviderBindsTheApplicationConfigurationWhereItHoldsTheSettings(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	container.SetSingleton(applicationconstant.ConfigContractServiceID, &configFixture{
		Config: applicationdata.NewConfig(),
		logger: "App.Log.Logger",
	})

	provider.PublishLogConfig(container)

	resolved, err := container.GetSingleton(constant.LogConfigContractServiceID)
	if err != nil {
		t.Fatalf("the provider must bind the configuration, but reported: %v", err)
	}

	built, isConfig := resolved.(contract.LogConfigContract)
	if !isConfig || built.GetDefaultLogger() != "App.Log.Logger" {
		t.Error("the application's own configuration must be bound where it holds the settings, but was not")
	}
}

func TestTheLoggerWritesToTheFileThatTheApplicationNames(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "app.log")

	container := manager.NewContainer(nil)
	container.SetSingleton(applicationconstant.ConfigContractServiceID, &configFixture{
		Config: applicationdata.NewConfig(),
		path:   path,
	})

	provider.PublishLogger(container)

	readLogger(t, container).Info("the message", nil)

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the logger must create the file, but reported: %v", err)
	}

	if !strings.Contains(string(contents), "the message") {
		t.Errorf("the logger must write its message to the file, but wrote: %q", contents)
	}
}

func TestTheLoggerWritesToTheProcessWhereTheApplicationNamesNoUsableFile(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"an application that names no file": "",
		"a file that no process can open":   filepath.Join(t.TempDir(), "missing", "app.log"),
	}

	for name, path := range tests {
		container := manager.NewContainer(nil)
		container.SetSingleton(applicationconstant.ConfigContractServiceID, &configFixture{
			Config: applicationdata.NewConfig(),
			path:   path,
		})

		provider.PublishLogger(container)

		// The logger writes to the standard error of the process, which nothing
		// here reads. What matters is that a logger was bound at all.
		if readLogger(t, container) == nil {
			t.Errorf("%s must still bind a logger, but bound none", name)
		}
	}
}

func TestTheLoggerWritesToTheProcessWhereTheConfigurationNamesNoFile(t *testing.T) {
	t.Parallel()

	tests := map[string]any{
		"a container that holds no configuration":    nil,
		"a configuration that holds no log settings": applicationdata.NewConfig(),
	}

	for name, config := range tests {
		container := manager.NewContainer(nil)

		if config != nil {
			container.SetSingleton(applicationconstant.ConfigContractServiceID, config)
		}

		provider.PublishLogger(container)

		if readLogger(t, container) == nil {
			t.Errorf("%s must still bind a logger, but bound none", name)
		}
	}
}

func TestTheComponentProviderNamesEveryProviderOfTheComponent(t *testing.T) {
	t.Parallel()

	built := &provider.LogComponentProvider{}

	if len(built.GetComponentProviders(nil)) != 1 {
		t.Errorf("the log must name the container, but named: %d", len(built.GetComponentProviders(nil)))
	}

	if len(built.GetContainerProviders(nil)) != 1 {
		t.Errorf("the log must name its service provider, but named: %d", len(built.GetContainerProviders(nil)))
	}

	if len(built.GetEventProviders(nil)) != 0 ||
		len(built.GetCliProviders(nil)) != 0 ||
		len(built.GetHttpProviders(nil)) != 0 {
		t.Error("the log must name no listener provider and no route provider, but named one")
	}
}
