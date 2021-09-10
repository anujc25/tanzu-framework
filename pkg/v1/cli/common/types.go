// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
)

// Plugin is an installable CLI plugin.
type Plugin struct {
	// Name is the name of the plugin.
	Name string `json:"name" yaml:"name"`
	// Description is the plugin's description.
	Description string `json:"description" yaml:"description"`
	// Distribution mechanism for the plugin.
	Distribution cliv1alpha1.DistributionConfig `json:"distribution"`
	// VersionConstraints for the plugin describes constraints
	// around using a version of the plugin.
	VersionConstraints cliv1alpha1.VersionConstraints `json:"versionConstraints"`
	// Platforms available for the plugin.
	Platforms []cliv1alpha1.Platform `json:"platforms"`
}
