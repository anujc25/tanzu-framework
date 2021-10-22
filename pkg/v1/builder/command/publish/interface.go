// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package publish implements plugin and plugin api publishing related function
package publish

// Publisher is an interface to publish plugin and CLIPlugin resource files to discovery
type Publisher interface {
	// PublishPlugin publishes plugin binaries to distribution
	PublishPlugin(version, os, arch, plugin, sourcePath string) (string, error)
	// PublishDiscovery publishes plugin discovery
	PublishDiscovery() error
}
