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

	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/application/data"
)

func TestTheCliConfigurationTakesItsDefaults(t *testing.T) {
	t.Parallel()

	var built contract.CliConfigContract = data.NewCliConfig()

	if built.GetApplicationName() != data.DefaultApplicationName {
		t.Errorf("the CLI configuration must name the application, but named: %q", built.GetApplicationName())
	}

	if built.GetDefaultCommandName() != data.DefaultCommandName {
		t.Errorf("the CLI configuration must name a default command, but named: %q", built.GetDefaultCommandName())
	}

	if built.GetNamespace() != data.DefaultNamespace {
		t.Error("the CLI configuration must carry every field of the application configuration, but did not")
	}
}

func TestEachCliMiddlewareStageStartsEmpty(t *testing.T) {
	t.Parallel()

	built := data.NewCliConfig()

	stages := map[string][]string{
		"input received":    built.GetInputReceivedMiddleware(),
		"route matched":     built.GetRouteMatchedMiddleware(),
		"route not matched": built.GetRouteNotMatchedMiddleware(),
		"route dispatched":  built.GetRouteDispatchedMiddleware(),
		"throwable caught":  built.GetThrowableCaughtMiddleware(),
		"process exiting":   built.GetProcessExitingMiddleware(),
	}

	for name, middleware := range stages {
		if len(middleware) != 0 {
			t.Errorf("the %s stage must start empty, but held: %v", name, middleware)
		}
	}
}

func TestEachHttpMiddlewareStageStartsEmpty(t *testing.T) {
	t.Parallel()

	var built contract.HttpConfigContract = data.NewHttpConfig()

	stages := map[string][]string{
		"request received":  built.GetRequestReceivedMiddleware(),
		"route matched":     built.GetRouteMatchedMiddleware(),
		"route not matched": built.GetRouteNotMatchedMiddleware(),
		"route dispatched":  built.GetRouteDispatchedMiddleware(),
		"throwable caught":  built.GetThrowableCaughtMiddleware(),
		"sending response":  built.GetSendingResponseMiddleware(),
		"response sent":     built.GetResponseSentMiddleware(),
	}

	for name, middleware := range stages {
		if len(middleware) != 0 {
			t.Errorf("the %s stage must start empty, but held: %v", name, middleware)
		}
	}

	if built.GetNamespace() != data.DefaultNamespace {
		t.Error("the HTTP configuration must carry every field of the application configuration, but did not")
	}
}
