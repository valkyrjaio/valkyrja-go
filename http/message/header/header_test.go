/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package header_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

const (
	htmlValue     = "text/html"
	acceptName    = "Accept"
	acceptedTypes = "text/html, application/json"
)

// newHeader builds a header and fails the test where the name is invalid.
func newHeader(t *testing.T, name string, values ...contract.ValueContract) *header.Header {
	t.Helper()

	built, err := header.NewHeader(name, values...)
	if err != nil {
		t.Fatalf("NewHeader must build the header, but reported: %v", err)
	}

	return built
}

func TestIsValidNameAcceptsATokenName(t *testing.T) {
	t.Parallel()

	valid := []string{constant.HeaderNameContentType, "X-Requested-With", "x-custom_header", "a1!#$%&'*+-.^`|~"}

	for _, name := range valid {
		if !header.IsValidName(name) {
			t.Errorf("IsValidName must be true for %q, but is false", name)
		}
	}
}

func TestIsValidNameRejectsAnythingElse(t *testing.T) {
	t.Parallel()

	invalid := []string{"", "Content Type", "Content:Type", "Content\nType", "Content\x00Type", "Ünicode"}

	for _, name := range invalid {
		if header.IsValidName(name) {
			t.Errorf("IsValidName must be false for %q, but is true", name)
		}
	}
}

func TestNewHeaderReportsAnInvalidName(t *testing.T) {
	t.Parallel()

	_, err := header.NewHeader("Content Type")

	target, found := errors.AsType[*exception.HttpHeaderInvalidNameError](err)
	if !found {
		t.Fatalf("NewHeader must report an invalid name, but reported: %v", err)
	}

	if target.GetName() != "Content Type" {
		t.Errorf("the error must name the header, but names: %q", target.GetName())
	}
}

func TestIsValidValueAcceptsAPrintableValue(t *testing.T) {
	t.Parallel()

	valid := []string{htmlValue, "a\tb", "folded\r\n continued", "folded\r\n\tcontinued", ""}

	for _, headerValue := range valid {
		if !header.IsValidValue(headerValue) {
			t.Errorf("IsValidValue must be true for %q, but is false", headerValue)
		}
	}
}

func TestIsValidValueRejectsAHeaderInjection(t *testing.T) {
	t.Parallel()

	invalid := []string{
		"bare\nline feed",
		"bare\rcarriage return",
		"fold\r\nwith no whitespace",
		"trailing\r\n",
		"control\x01character",
		"delete\x7f",
	}

	for _, headerValue := range invalid {
		if header.IsValidValue(headerValue) {
			t.Errorf("IsValidValue must be false for %q, but is true", headerValue)
		}
	}
}

func TestValidateValueReportsAnInvalidValue(t *testing.T) {
	t.Parallel()

	target, found := errors.AsType[*exception.HttpHeaderInvalidValueError](header.ValidateValue("bad\nvalue"))
	if !found {
		t.Fatal("ValidateValue must report an invalid value, but did not")
	}

	if target.GetValue() != "bad\nvalue" {
		t.Errorf("the error must carry the value, but carries: %q", target.GetValue())
	}
}

func TestValidateValueAcceptsAValidValue(t *testing.T) {
	t.Parallel()

	if header.ValidateValue(htmlValue) != nil {
		t.Error("ValidateValue must accept a printable value, but reported a failure")
	}
}

func TestNewHeaderHoldsTheNameAndItsNormalizedForm(t *testing.T) {
	t.Parallel()

	built := newHeader(t, constant.HeaderNameContentType)

	if built.GetName() != constant.HeaderNameContentType {
		t.Errorf("GetName must be the name as written, but is: %q", built.GetName())
	}

	if built.GetNormalizedName() != "content-type" {
		t.Errorf("GetNormalizedName must be lower case, but is: %q", built.GetNormalizedName())
	}
}

func TestNewHeaderFromValueReadsTheNameAndEachValue(t *testing.T) {
	t.Parallel()

	built, err := header.NewHeaderFromValue("Accept: text/html, application/json")
	if err != nil {
		t.Fatalf("NewHeaderFromValue must build the header, but reported: %v", err)
	}

	if built.GetName() != acceptName {
		t.Errorf("GetName must be the name, but is: %q", built.GetName())
	}

	if len(built.GetValues()) != 2 {
		t.Fatalf("the header must hold each value, but holds: %d", len(built.GetValues()))
	}

	if built.GetHeaderLine() != acceptedTypes {
		t.Errorf("GetHeaderLine must join each value, but is: %q", built.GetHeaderLine())
	}
}

func TestNewHeaderFromValueReadsALineWithNoValue(t *testing.T) {
	t.Parallel()

	built, err := header.NewHeaderFromValue("X-Empty")
	if err != nil {
		t.Fatalf("NewHeaderFromValue must build the header, but reported: %v", err)
	}

	if built.GetName() != "X-Empty" {
		t.Errorf("GetName must be the name, but is: %q", built.GetName())
	}

	if built.String() != "" {
		t.Errorf("String must be empty where the header carries no value, but is: %q", built.String())
	}
}

func TestNewHeaderFromValueReportsAnInvalidValue(t *testing.T) {
	t.Parallel()

	_, err := header.NewHeaderFromValue("X-Bad: bad\nvalue")

	if _, found := errors.AsType[*exception.HttpHeaderInvalidValueError](err); !found {
		t.Errorf("NewHeaderFromValue must report an invalid value, but reported: %v", err)
	}
}

func TestNewHeaderFromValueReportsAnInvalidName(t *testing.T) {
	t.Parallel()

	_, err := header.NewHeaderFromValue("Bad Name: text/html")

	if _, found := errors.AsType[*exception.HttpHeaderInvalidNameError](err); !found {
		t.Errorf("NewHeaderFromValue must report an invalid name, but reported: %v", err)
	}
}

func TestHeaderStringJoinsTheNameAndTheValues(t *testing.T) {
	t.Parallel()

	built := newHeader(t, constant.HeaderNameContentType, value.NewValueFromValue(htmlValue))

	if built.String() != "Content-Type: text/html" {
		t.Errorf("String must join the name and the values, but is: %q", built.String())
	}
}

func TestWithNameReturnsACopyUnderTheNewName(t *testing.T) {
	t.Parallel()

	built := newHeader(t, constant.HeaderNameContentType)

	renamed := built.WithName(acceptName)

	if renamed.GetName() != acceptName || renamed.GetNormalizedName() != "accept" {
		t.Errorf("WithName must hold the new name, but is: %q", renamed.GetName())
	}

	if built.GetName() != constant.HeaderNameContentType {
		t.Error("WithName must leave the receiver unchanged, but did not")
	}
}

func TestWithNameKeepsTheNameWhereTheNewOneIsInvalid(t *testing.T) {
	t.Parallel()

	built := newHeader(t, constant.HeaderNameContentType)

	if built.WithName("Bad Name").GetName() != constant.HeaderNameContentType {
		t.Error("WithName must keep the name where the new one is invalid, but did not")
	}
}

func TestWithValuesReplacesWhatTheHeaderHolds(t *testing.T) {
	t.Parallel()

	built := newHeader(t, constant.HeaderNameContentType, value.NewValueFromValue(htmlValue))

	replaced := built.WithValues(value.NewValueFromValue("text/plain"))

	if replaced.GetHeaderLine() != "text/plain" {
		t.Errorf("WithValues must replace the values, but is: %q", replaced.GetHeaderLine())
	}

	if built.GetHeaderLine() != htmlValue {
		t.Error("WithValues must leave the receiver unchanged, but did not")
	}
}

func TestWithAddedValuesKeepsWhatTheHeaderHolds(t *testing.T) {
	t.Parallel()

	built := newHeader(t, acceptName, value.NewValueFromValue(htmlValue))

	added := built.WithAddedValues(value.NewValueFromValue("application/json"))

	if added.GetHeaderLine() != acceptedTypes {
		t.Errorf("WithAddedValues must append the values, but is: %q", added.GetHeaderLine())
	}

	if built.GetHeaderLine() != htmlValue {
		t.Error("WithAddedValues must leave the receiver unchanged, but did not")
	}
}

func TestGetHeaderLineDropsAValueThatRendersToNothing(t *testing.T) {
	t.Parallel()

	built := newHeader(t, acceptName, value.NewValue(), value.NewValueFromValue(htmlValue))

	if built.GetHeaderLine() != htmlValue {
		t.Errorf("GetHeaderLine must drop an empty value, but is: %q", built.GetHeaderLine())
	}
}
