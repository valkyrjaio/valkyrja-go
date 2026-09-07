/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/application/throwable/exception"
	throwableexception "github.com/valkyrjaio/valkyrja-go/v26/throwable/exception"
)

func TestEachBaseErrorSatisfiesTheComponentContract(t *testing.T) {
	t.Parallel()

	runtimeError := &exception.ApplicationRuntimeError{
		ValkyrjaRuntimeError: throwableexception.NewValkyrjaRuntimeError("the runtime failed", nil),
	}
	invalidArgumentError := &exception.ApplicationInvalidArgumentError{
		ValkyrjaInvalidArgumentError: throwableexception.NewValkyrjaInvalidArgumentError(
			"the argument is invalid",
			nil,
		),
	}

	throwables := []contract.ApplicationThrowable{runtimeError, invalidArgumentError}

	for _, throwable := range throwables {
		if !throwable.IsApplicationThrowable() {
			t.Errorf("IsApplicationThrowable must be true, but is false for: %s", throwable.Error())
		}

		if throwable.GetTraceCode() == "" {
			t.Errorf("GetTraceCode must return a code, but is empty for: %s", throwable.Error())
		}
	}
}
