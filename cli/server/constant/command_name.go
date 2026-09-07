/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the CLI server's command names and binding keys.
package constant

// The name of each command that the CLI server ships.
const (
	CommandNameList = "list"

	CommandNameListBash = "list:bash"

	CommandNameHelp = "help"

	CommandNameVersion = "version"
)

// The binding key of each service that the CLI server publishes.
const (
	InputHandlerContractServiceID = "valkyrja.cli.server.handler.InputHandlerContract"
)
