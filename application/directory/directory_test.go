/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package directory_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/application/directory"
)

const firstPath = "/first"

func TestGetPathAddsOneLeadingSeparator(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"":            "",
		"first":       firstPath,
		firstPath:     firstPath,
		"first/other": "/first/other",
	}

	base := directory.NewDirectory("/base")

	for path, expected := range tests {
		if base.GetPath(path) != expected {
			t.Errorf("GetPath for %q must be %q, but is %q", path, expected, base.GetPath(path))
		}
	}
}

func TestEachDirectoryResolvesUnderTheBasePath(t *testing.T) {
	t.Parallel()

	base := directory.NewDirectory("/base")

	tests := map[string]struct {
		resolve  func(path string) string
		expected string
	}{
		"base":      {base.GetBaseDirectory, "/base/first"},
		"app":       {base.GetAppDirectory, "/base/app/first"},
		"data":      {base.GetDataDirectory, "/base/data/first"},
		"env":       {base.GetEnvDirectory, "/base/env/first"},
		"public":    {base.GetPublicDirectory, "/base/public/first"},
		"resources": {base.GetResourcesDirectory, "/base/resources/first"},
		"src":       {base.GetSrcDirectory, "/base/src/first"},
		"storage":   {base.GetStorageDirectory, "/base/storage/first"},
		"tests":     {base.GetTestsDirectory, "/base/tests/first"},
		"vendor":    {base.GetVendorDirectory, "/base/vendor/first"},
		"framework": {base.GetFrameworkStorageDirectory, "/base/storage/framework/first"},
		"logs":      {base.GetLogsStorageDirectory, "/base/storage/logs/first"},
		"cache":     {base.GetFrameworkStorageCacheDirectory, "/base/storage/framework/cache/first"},
	}

	for name, test := range tests {
		if test.resolve("first") != test.expected {
			t.Errorf("the %s directory must be %q, but is %q", name, test.expected, test.resolve("first"))
		}
	}
}

func TestEachDirectoryResolvesWithNoPath(t *testing.T) {
	t.Parallel()

	base := directory.NewDirectory("/base")

	if base.GetAppDirectory("") != "/base/app" {
		t.Errorf("the app directory must be /base/app, but is: %s", base.GetAppDirectory(""))
	}

	if base.GetBaseDirectory("") != "/base" {
		t.Errorf("the base directory must be /base, but is: %s", base.GetBaseDirectory(""))
	}
}

func TestADirectoryTakesAnotherName(t *testing.T) {
	t.Parallel()

	base := directory.NewDirectory("/base")
	base.VendorPath = "node_modules"

	if base.GetVendorDirectory("first") != "/base/node_modules/first" {
		t.Errorf("the vendor directory must follow the field, but is: %s", base.GetVendorDirectory("first"))
	}
}
