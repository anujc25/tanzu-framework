// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package discovery

import "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"

// GCPDiscovery is a artifact discovery endpoing utilizing a GCP bucket.
type GCPDiscovery struct {
	bucketName   string
	manifestPath string
	name         string
}

// NewGCPDiscovery returns a new GCP bucket repository.
func NewGCPDiscovery(bucket, manifestPath, name string) Discovery {
	return &GCPDiscovery{
		bucketName:   bucket,
		manifestPath: manifestPath,
		name:         name,
	}
}

// List available plugins.
func (g *GCPDiscovery) List() (plugins []common.Plugin, err error) {
	return
}

// Describe a plugin.
func (g *GCPDiscovery) Describe(name string) (plugin common.Plugin, err error) {
	return plugin, err
}

// Name of the repository.
func (g *GCPDiscovery) Name() string {
	return g.name
}
