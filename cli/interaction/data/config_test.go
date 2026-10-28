/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/data"
)

func TestTheDefaultConfigurationIsInteractiveAndWritesEveryMessage(t *testing.T) {
	t.Parallel()

	built := data.NewCliInteractionConfig()

	if !built.IsInteractive() || built.IsQuiet() || built.IsSilent() {
		t.Error("the default configuration must be interactive and write every message, but was not")
	}
}

func TestTheConfigurationReadsTheValuesThatTheApplicationStates(t *testing.T) {
	t.Parallel()

	built := data.NewCliInteractionConfigFromValues(true, false, true)

	if !built.IsQuiet() || built.IsInteractive() || !built.IsSilent() {
		t.Error("the configuration must read each value that the application states, but did not")
	}
}
