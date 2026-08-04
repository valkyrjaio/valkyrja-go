/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the application component's binding keys and the
// metadata of the framework itself.
package constant

// The application component's binding keys.
const (
	// ApplicationContractServiceID is the binding key for the application.
	ApplicationContractServiceID = "Valkyrja.Application.Kernel.ApplicationContract"

	// ConfigContractServiceID is the binding key for the application config.
	ConfigContractServiceID = "Valkyrja.Application.Data.ConfigContract"

	// CliConfigContractServiceID is the binding key for the CLI config.
	CliConfigContractServiceID = "Valkyrja.Application.Data.CliConfigContract"

	// HttpConfigContractServiceID is the binding key for the HTTP config.
	HttpConfigContractServiceID = "Valkyrja.Application.Data.HttpConfigContract"

	// ConfigServiceID is the binding key for the default config.
	ConfigServiceID = "Valkyrja.Application.Data.Config"
)
