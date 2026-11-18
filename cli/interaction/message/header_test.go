/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message_test

import (
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/format"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/message"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/data"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

const (
	appName    = "Sample"
	appVersion = "1.2.3"
)

// newRoute builds the command that a header test names.
func newRoute() contract.RouteContract {
	return data.NewRoute("cache:clear", "Clear the cache", func(
		_ containercontract.ContainerContract,
		_ contract.RouteContract,
	) contract.OutputContract {
		return nil
	})
}

func TestTheHeaderNamesTheApplicationAndTheCommand(t *testing.T) {
	t.Parallel()

	text := message.NewHeader(appName, appVersion, "/project", newRoute()).GetText()

	for _, part := range []string{appName, "v" + appVersion, "/project", "Clear the cache", "cache:clear"} {
		if !strings.Contains(text, part) {
			t.Errorf("the header must name %q, but printed: %q", part, text)
		}
	}
}

func TestTheHeaderNamesNoCommandWhereItWasGivenNone(t *testing.T) {
	t.Parallel()

	text := message.NewHeader(appName, appVersion, "/project", nil).GetText()

	if !strings.Contains(text, "╰──  · ") {
		t.Errorf("a header with no command must name none, but printed: %q", text)
	}
}

func TestEachHeaderWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := message.NewHeader(appName, appVersion, "/project", newRoute())

	changed := built.
		WithAppName("Other").
		WithAppVersion("2.0.0").
		WithIcon("icon").
		WithValkyrjaVersion("26.9.9").
		WithValkyrjaBuildDate("A date").
		WithGoVersion("go1.99").
		WithProjectRoot("/other").
		WithActionDescription("Another description").
		WithCommandName("other:command")

	text := changed.GetText()

	for _, part := range []string{
		"Other", "v2.0.0", "icon", "26.9.9", "A date", "go1.99",
		"/other", "Another description", "other:command",
	} {
		if !strings.Contains(text, part) {
			t.Errorf("the header must carry %q, but printed: %q", part, text)
		}
	}

	if !strings.Contains(built.GetText(), appName) {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestTheHeaderCarriesNoFormatterOfItsOwn(t *testing.T) {
	t.Parallel()

	built := message.NewHeader(appName, appVersion, "/project", newRoute())

	if built.HasFormatter() || built.GetFormatter() != nil {
		t.Error("a header must carry no formatter, but carried one")
	}

	if built.GetFormattedText() != built.GetText() {
		t.Error("a header must apply no format of its own, but applied one")
	}

	if built.WithoutFormatter() != contract.MessageContract(built) {
		t.Error("WithoutFormatter must return the header as it is, but did not")
	}

	if built.WithFormatter(format.NewErrorFormatter()).GetFormattedText() == built.GetText() {
		t.Error("WithFormatter must return a message that applies the format, but did not")
	}

	if built.WithText("plain").GetText() != "plain" {
		t.Error("WithText must return a message that carries the text, but did not")
	}
}
