/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package directory resolves each directory of an application, from one base
// path.
//
// The other ports hold every path as a mutable static field on a class, so an
// application changes one before it boots. Go has no class, so the paths are
// fields of a `Directory` value that an application builds once. That also
// removes the shared global state, which a Go program reaches from several
// goroutines at once.
package directory

import (
	"strings"
)

// The name of each directory, under the base path.
const (
	DefaultAppPath              = "app"
	DefaultDataPath             = "data"
	DefaultEnvPath              = "env"
	DefaultPublicPath           = "public"
	DefaultResourcesPath        = "resources"
	DefaultSrcPath              = "src"
	DefaultStoragePath          = "storage"
	DefaultFrameworkStoragePath = "framework"
	DefaultCacheStoragePath     = "cache"
	DefaultLogsStoragePath      = "logs"
	DefaultTestsPath            = "tests"
	DefaultVendorPath           = "vendor"
)

type Directory struct {
	BasePath             string
	AppPath              string
	DataPath             string
	EnvPath              string
	PublicPath           string
	ResourcesPath        string
	SrcPath              string
	StoragePath          string
	FrameworkStoragePath string
	CacheStoragePath     string
	LogsStoragePath      string
	TestsPath            string
	VendorPath           string
}

// NewDirectory builds the directories under the base path, each one at its own
// default name.
func NewDirectory(basePath string) *Directory {
	return &Directory{
		BasePath:             basePath,
		AppPath:              DefaultAppPath,
		DataPath:             DefaultDataPath,
		EnvPath:              DefaultEnvPath,
		PublicPath:           DefaultPublicPath,
		ResourcesPath:        DefaultResourcesPath,
		SrcPath:              DefaultSrcPath,
		StoragePath:          DefaultStoragePath,
		FrameworkStoragePath: DefaultFrameworkStoragePath,
		CacheStoragePath:     DefaultCacheStoragePath,
		LogsStoragePath:      DefaultLogsStoragePath,
		TestsPath:            DefaultTestsPath,
		VendorPath:           DefaultVendorPath,
	}
}

// GetPath returns the path with one leading separator, and an empty string where
// the path is empty.
func (d *Directory) GetPath(path string) string {
	if path == "" {
		return ""
	}

	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}

// GetBaseDirectory returns the path under the base path.
func (d *Directory) GetBaseDirectory(path string) string {
	return d.BasePath + d.GetPath(path)
}

// GetAppDirectory returns the path under the application directory.
func (d *Directory) GetAppDirectory(path string) string {
	return d.GetBaseDirectory(d.AppPath + d.GetPath(path))
}

// GetDataDirectory returns the path under the data directory.
func (d *Directory) GetDataDirectory(path string) string {
	return d.GetBaseDirectory(d.DataPath + d.GetPath(path))
}

// GetEnvDirectory returns the path under the environment directory.
func (d *Directory) GetEnvDirectory(path string) string {
	return d.GetBaseDirectory(d.EnvPath + d.GetPath(path))
}

// GetPublicDirectory returns the path under the public directory.
func (d *Directory) GetPublicDirectory(path string) string {
	return d.GetBaseDirectory(d.PublicPath + d.GetPath(path))
}

// GetResourcesDirectory returns the path under the resources directory.
func (d *Directory) GetResourcesDirectory(path string) string {
	return d.GetBaseDirectory(d.ResourcesPath + d.GetPath(path))
}

// GetSrcDirectory returns the path under the source directory.
func (d *Directory) GetSrcDirectory(path string) string {
	return d.GetBaseDirectory(d.SrcPath + d.GetPath(path))
}

// GetStorageDirectory returns the path under the storage directory.
func (d *Directory) GetStorageDirectory(path string) string {
	return d.GetBaseDirectory(d.StoragePath + d.GetPath(path))
}

// GetFrameworkStorageDirectory returns the path under the framework's own
// storage directory.
func (d *Directory) GetFrameworkStorageDirectory(path string) string {
	return d.GetStorageDirectory(d.FrameworkStoragePath + d.GetPath(path))
}

// GetLogsStorageDirectory returns the path under the logs storage directory.
func (d *Directory) GetLogsStorageDirectory(path string) string {
	return d.GetStorageDirectory(d.LogsStoragePath + d.GetPath(path))
}

// GetFrameworkStorageCacheDirectory returns the path under the framework's own
// cache directory.
func (d *Directory) GetFrameworkStorageCacheDirectory(path string) string {
	return d.GetFrameworkStorageDirectory(d.CacheStoragePath + d.GetPath(path))
}

// GetTestsDirectory returns the path under the tests directory.
func (d *Directory) GetTestsDirectory(path string) string {
	return d.GetBaseDirectory(d.TestsPath + d.GetPath(path))
}

// GetVendorDirectory returns the path under the vendor directory.
func (d *Directory) GetVendorDirectory(path string) string {
	return d.GetBaseDirectory(d.VendorPath + d.GetPath(path))
}
