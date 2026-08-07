/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package value_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
)

const (
	charsetToken = "charset"
	utf8Text     = "utf-8"
	charsetUtf8  = "charset=utf-8"
	htmlType     = "text/html"
	htmlWithUtf8 = "text/html; charset=utf-8"
)

func TestNewComponentTrimsEachPart(t *testing.T) {
	t.Parallel()

	component := value.NewComponent("  charset  ", "  utf-8  ")

	if component.GetToken() != charsetToken {
		t.Errorf("GetToken must be trimmed, but is: %q", component.GetToken())
	}

	if component.GetText() != utf8Text {
		t.Errorf("GetText must be trimmed, but is: %q", component.GetText())
	}
}

func TestNewComponentFromValueReadsTheTokenAndTheText(t *testing.T) {
	t.Parallel()

	component := value.NewComponentFromValue(" charset = utf-8 ")

	if component.GetToken() != charsetToken {
		t.Errorf("GetToken must be the token, but is: %q", component.GetToken())
	}

	if component.GetText() != utf8Text {
		t.Errorf("GetText must be the text, but is: %q", component.GetText())
	}
}

func TestNewComponentFromValueReadsATokenWithNoText(t *testing.T) {
	t.Parallel()

	component := value.NewComponentFromValue(" text/html ")

	if component.GetToken() != htmlType {
		t.Errorf("GetToken must be the token, but is: %q", component.GetToken())
	}

	if component.GetText() != "" {
		t.Errorf("GetText must be empty, but is: %q", component.GetText())
	}
}

func TestComponentStringJoinsTheTokenAndTheText(t *testing.T) {
	t.Parallel()

	if value.NewComponent(charsetToken, utf8Text).String() != charsetUtf8 {
		t.Errorf("String must join the parts, but is: %q", value.NewComponent(charsetToken, utf8Text).String())
	}
}

func TestComponentStringIsTheTokenWhereThereIsNoText(t *testing.T) {
	t.Parallel()

	if value.NewComponent(htmlType, "").String() != htmlType {
		t.Errorf("String must be the token, but is: %q", value.NewComponent(htmlType, "").String())
	}
}

func TestEachComponentWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	component := value.NewComponent(charsetToken, utf8Text)

	withToken := component.WithToken("  boundary  ")
	withText := component.WithText("  ascii  ")

	if withToken.GetToken() != "boundary" {
		t.Errorf("WithToken must trim and hold the token, but is: %q", withToken.GetToken())
	}

	if withText.GetText() != "ascii" {
		t.Errorf("WithText must trim and hold the text, but is: %q", withText.GetText())
	}

	if component.GetToken() != charsetToken || component.GetText() != utf8Text {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestNewValueDropsAComponentThatRendersToNothing(t *testing.T) {
	t.Parallel()

	headerValue := value.NewValue(
		value.NewComponent(htmlType, ""),
		value.NewComponent("", ""),
		value.NewComponent(charsetToken, utf8Text),
	)

	if len(headerValue.GetComponents()) != 2 {
		t.Errorf("NewValue must drop an empty component, but holds: %d", len(headerValue.GetComponents()))
	}
}

func TestNewValueFromValueReadsEachComponent(t *testing.T) {
	t.Parallel()

	headerValue := value.NewValueFromValue(htmlWithUtf8)

	components := headerValue.GetComponents()

	if len(components) != 2 {
		t.Fatalf("NewValueFromValue must read each component, but read: %d", len(components))
	}

	if components[0].GetToken() != htmlType {
		t.Errorf("the first component must be the media type, but is: %q", components[0].GetToken())
	}

	if components[1].GetText() != utf8Text {
		t.Errorf("the second component must carry the charset, but is: %q", components[1].GetText())
	}
}

func TestValueStringJoinsEachComponent(t *testing.T) {
	t.Parallel()

	headerValue := value.NewValueFromValue(htmlWithUtf8)

	if headerValue.String() != htmlWithUtf8 {
		t.Errorf("String must join each component, but is: %q", headerValue.String())
	}
}

func TestWithComponentsReplacesWhatTheValueHolds(t *testing.T) {
	t.Parallel()

	headerValue := value.NewValueFromValue(htmlWithUtf8)

	replaced := headerValue.WithComponents(value.NewComponent("text/plain", ""))

	if replaced.String() != "text/plain" {
		t.Errorf("WithComponents must replace the components, but is: %q", replaced.String())
	}

	if headerValue.String() != htmlWithUtf8 {
		t.Error("WithComponents must leave the receiver unchanged, but did not")
	}
}

func TestWithAddedComponentsKeepsWhatTheValueHolds(t *testing.T) {
	t.Parallel()

	headerValue := value.NewValueFromValue(htmlType)

	added := headerValue.WithAddedComponents(value.NewComponent(charsetToken, utf8Text))

	if added.String() != htmlWithUtf8 {
		t.Errorf("WithAddedComponents must append the components, but is: %q", added.String())
	}

	if headerValue.String() != htmlType {
		t.Error("WithAddedComponents must leave the receiver unchanged, but did not")
	}
}

func TestTheValueTypesSatisfyTheirContracts(t *testing.T) {
	t.Parallel()

	var component contract.ComponentContract = value.NewComponent(charsetToken, utf8Text)
	var headerValue contract.ValueContract = value.NewValue(component)

	if headerValue.String() != charsetUtf8 {
		t.Errorf("the contracts must render the value, but rendered: %q", headerValue.String())
	}
}
