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
	ApplicationContractServiceID = "valkyrja.application.kernel.ApplicationContract"

	// ConfigContractServiceID is the binding key for the application config.
	ConfigContractServiceID = "valkyrja.application.data.ConfigContract"

	// CliConfigContractServiceID is the binding key for the CLI config.
	CliConfigContractServiceID = "valkyrja.application.data.CliConfigContract"

	// HttpConfigContractServiceID is the binding key for the HTTP config.
	HttpConfigContractServiceID = "valkyrja.application.data.HttpConfigContract"

	// ConfigServiceID is the binding key for the default config.
	ConfigServiceID = "valkyrja.application.data.Config"
)
