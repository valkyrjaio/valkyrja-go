/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds the application's configuration.
package data

import (
	"os"

	"github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
)

// The values that an application takes where it states none of its own.
const (
	DefaultNamespace     = "App"
	DefaultEnvironment   = "production"
	DefaultTimezone      = "UTC"
	DefaultKey           = "some_secret_app_key"
	DefaultDataPath      = "App/Provider/Data"
	DefaultDataNamespace = "App/Provider/Data"
)

// Config is the configuration of the application.
//
// The other ports give the constructor a default for each parameter. Go has no
// default parameter, so `NewConfig` takes the values that an application states
// most often, and a caller changes a field afterwards.
type Config struct {
	Namespace     string
	Dir           string
	Version       string
	Environment   string
	DebugMode     bool
	Timezone      string
	Key           string
	DataPath      string
	DataNamespace string
	Providers     []contract.ComponentProviderContract
	Callbacks     []contract.ApplicationCallbackFunc
}

// NewConfig builds the configuration that every field takes its default value
// in. The directory is the one that the process runs from.
func NewConfig(providers ...contract.ComponentProviderContract) *Config {
	return &Config{
		Namespace:     DefaultNamespace,
		Dir:           getWorkingDirectory(),
		Version:       constant.Version,
		Environment:   DefaultEnvironment,
		DebugMode:     false,
		Timezone:      DefaultTimezone,
		Key:           DefaultKey,
		DataPath:      DefaultDataPath,
		DataNamespace: DefaultDataNamespace,
		Providers:     providers,
		Callbacks:     []contract.ApplicationCallbackFunc{},
	}
}

// getWorkingDirectory returns the directory that the process runs from, and an
// empty string where the process cannot read it.
//
// `os.Getwd` returns an empty directory together with its error, so a guard on
// the error would return the value that the call already gives, and no test
// could reach the guard.
func getWorkingDirectory() string {
	directory, _ := os.Getwd()

	return directory
}

// GetNamespace returns the namespace of the application.
func (c *Config) GetNamespace() string {
	return c.Namespace
}

// GetDir returns the directory that the application runs from.
func (c *Config) GetDir() string {
	return c.Dir
}

// GetVersion returns the version of the application.
func (c *Config) GetVersion() string {
	return c.Version
}

// GetEnvironment returns the environment that the application runs in.
func (c *Config) GetEnvironment() string {
	return c.Environment
}

// GetDebugMode reports whether the application runs in debug mode.
func (c *Config) GetDebugMode() bool {
	return c.DebugMode
}

// GetTimezone returns the timezone that the application runs in.
func (c *Config) GetTimezone() string {
	return c.Timezone
}

// GetKey returns the secret key of the application.
func (c *Config) GetKey() string {
	return c.Key
}

// GetDataPath returns the path that `sindri` writes the generated data to.
func (c *Config) GetDataPath() string {
	return c.DataPath
}

// GetDataNamespace returns the namespace that `sindri` writes the generated data
// under.
func (c *Config) GetDataNamespace() string {
	return c.DataNamespace
}

// GetProviders returns each component provider that the application registers.
func (c *Config) GetProviders() []contract.ComponentProviderContract {
	return c.Providers
}

// GetCallbacks returns each callback that runs once the application publishes
// its providers.
func (c *Config) GetCallbacks() []contract.ApplicationCallbackFunc {
	return c.Callbacks
}
