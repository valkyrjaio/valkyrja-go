/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

type ProtocolVersion string

// The ProtocolVersion values that the framework knows.
const (
	ProtocolVersionV1  ProtocolVersion = "1.0"
	ProtocolVersionV11 ProtocolVersion = "1.1"
	ProtocolVersionV2  ProtocolVersion = "2"
	ProtocolVersionV3  ProtocolVersion = "3"
)
