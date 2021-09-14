// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"time"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
)

// Plugin is an installable CLI plugin.
type Plugin struct {
	// Name is the name of the plugin.
	Name string `json:"name" yaml:"name"`
	// Description is the plugin's description.
	Description string `json:"description" yaml:"description"`
	// Distribution mechanism for the plugin.
	Distribution cliv1alpha1.DistributionConfig `json:"distribution" yaml:"distribution"`
	// VersionConstraints for the plugin describes constraints
	// around using a version of the plugin.
	VersionConstraints cliv1alpha1.VersionConstraints `json:"versionConstraints" yaml:"versionConstraints"`
	// Platforms available for the plugin.
	Platforms []cliv1alpha1.Platform `json:"platforms" yaml:"platforms"`
	// Discovery specificies the name of the discovery from where
	// this plugin is discovered.
	Discovery string
	// Scope specificies the scope of the plugin. Stand-Alone or Context
	Scope string
	// Status specificies the current plugin installation status
	Status string
}

// Manifest is stored in the repository which gives an inventory of the artifacts.
type Manifest struct {
	// Created is the time the manifest was created.
	CreatedTime time.Time `json:"created" yaml:"created"`

	// Plugins is a list of plugin artifacts available.
	Plugins []Plugin `json:"plugins" yaml:"plugins"`
}
