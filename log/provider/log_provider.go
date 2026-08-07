/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package provider holds the log component's providers.
package provider

import (
	"io"
	"os"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	applicationcontract "github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	containerprovider "github.com/valkyrjaio/valkyrja-go/v26/container/provider"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/log/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/log/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/log/data"
	"github.com/valkyrjaio/valkyrja-go/v26/log/logger"
)

// filePermissions is what a new log file is created with. The owner reads and
// writes the file, and every other user reads it.
const filePermissions = 0o644

type LogServiceProvider struct{}

// Publishers returns a publisher for each binding key that the component defers.
func (p *LogServiceProvider) Publishers() map[string]containercontract.PublishFunc {
	return map[string]containercontract.PublishFunc{
		constant.LogConfigContractServiceID: PublishLogConfig,
		constant.LoggerContractServiceID:    PublishLogger,
	}
}

// PublishLogConfig binds the configuration of the component.
func PublishLogConfig(container containercontract.ContainerContract) {
	container.SetSingleton(constant.LogConfigContractServiceID, resolveConfig(container))
}

// PublishLogger binds the logger that the application writes through.
func PublishLogger(container containercontract.ContainerContract) {
	container.SetSingleton(constant.LoggerContractServiceID, logger.NewStreamLogger(resolveWriter(container)))
}

// resolveConfig returns the configuration that the application states, and the
// framework's default where the application states none.
func resolveConfig(container containercontract.ContainerContract) contract.LogConfigContract {
	resolved, err := container.GetSingleton(applicationconstant.ConfigContractServiceID)
	if err != nil {
		return data.NewLogConfig()
	}

	config, isConfig := resolved.(contract.LogConfigContract)
	if !isConfig {
		return data.NewLogConfig()
	}

	return config
}

// resolveWriter returns the stream that the logger writes to.
func resolveWriter(container containercontract.ContainerContract) io.Writer {
	path := resolveStreamPath(container)
	if path == "" {
		return os.Stderr
	}

	// The application names the file, the same way it does in every other port.
	//nolint:gosec // The path is the application's own, not a client's.
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermissions)
	if err != nil {
		return os.Stderr
	}

	return file
}

// resolveStreamPath returns the file that the application names, and an empty
// string where it names none.
func resolveStreamPath(container containercontract.ContainerContract) string {
	resolved, err := container.GetSingleton(applicationconstant.ConfigContractServiceID)
	if err != nil {
		return ""
	}

	config, isConfig := resolved.(contract.LogStreamConfigContract)
	if !isConfig {
		return ""
	}

	return config.GetStreamFilePath()
}

type LogComponentProvider struct{}

// GetComponentProviders returns each component that the log needs.
func (p *LogComponentProvider) GetComponentProviders(
	_ applicationcontract.ApplicationContract,
) []applicationcontract.ComponentProviderContract {
	return []applicationcontract.ComponentProviderContract{
		&containerprovider.ContainerComponentProvider{},
	}
}

// GetContainerProviders returns each service provider of the log.
func (p *LogComponentProvider) GetContainerProviders(
	_ applicationcontract.ApplicationContract,
) []containercontract.ServiceProviderContract {
	return []containercontract.ServiceProviderContract{&LogServiceProvider{}}
}

// GetEventProviders returns each listener provider of the log.
func (p *LogComponentProvider) GetEventProviders(
	_ applicationcontract.ApplicationContract,
) []eventcontract.ListenerProviderContract {
	return []eventcontract.ListenerProviderContract{}
}

// GetCliProviders returns each CLI route provider of the log.
func (p *LogComponentProvider) GetCliProviders(
	_ applicationcontract.ApplicationContract,
) []clicontract.CliRouteProviderContract {
	return []clicontract.CliRouteProviderContract{}
}

// GetHttpProviders returns each HTTP route provider of the log.
func (p *LogComponentProvider) GetHttpProviders(
	_ applicationcontract.ApplicationContract,
) []httpcontract.HttpRouteProviderContract {
	return []httpcontract.HttpRouteProviderContract{}
}
