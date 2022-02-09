// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package managementcomponents

// TKGPackageConfig
type TKGPackageConfig struct {
	Metadata     Metadata          `yaml:"metadata"`
	ConfigValues map[string]string `yaml:"configvalues"`
}

// Metadata package metadata
type Metadata struct {
	InfraProvider string `yaml:"infraProvider"`
}
