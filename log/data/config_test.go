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

	"github.com/valkyrjaio/valkyrja-go/v26/log/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/log/data"
)

func TestTheConfigurationNamesTheDefaultLogger(t *testing.T) {
	t.Parallel()

	var built contract.LogConfigContract = data.NewLogConfig()

	if built.GetDefaultLogger() != data.DefaultLogger {
		t.Errorf("the configuration must name the default logger, but named: %q", built.GetDefaultLogger())
	}
}

func TestTheStreamConfigurationNamesNoFileUntilTheApplicationStatesOne(t *testing.T) {
	t.Parallel()

	var built contract.LogStreamConfigContract = data.NewLogStreamConfig()

	if built.GetStreamFilePath() != "" {
		t.Errorf("the configuration must name no file, but named: %q", built.GetStreamFilePath())
	}

	stated := data.NewLogStreamConfig()
	stated.StreamFilePath = "/var/log/app.log"

	if stated.GetStreamFilePath() != "/var/log/app.log" {
		t.Errorf("the configuration must name the file that the application states, but named: %q",
			stated.GetStreamFilePath())
	}
}
