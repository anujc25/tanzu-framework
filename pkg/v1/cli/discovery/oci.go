// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package discovery

import "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"

// OCIDiscovery is a artifact discovery endpoint utilizing OCI image
type OCIDiscovery struct {
	// name is a name of the discovery
	name string `json:"name"`
	// registry is an OCI compliant image registry. It MUST be a DNS-compatible name.
	// E.g., harbor.my-domain.local
	registry string `json:"registry,omitempty"`
	// path is the unique repository/image name. It MUST be a valid URI path, MAY
	// contain zero or more '/', and SHOULD NOT start or end with '/'.
	// E.g., tanzu/cli/plugins/manifests
	path string `json:"path"`
	// tag is the image tag for the image repository. If not provided `latest` is used
	tag string `json:"tag"`
}

// NewOCIDiscovery returns a new local repository.
func NewOCIDiscovery(name, registry, path, tag string) Discovery {
	return &OCIDiscovery{
		name:     name,
		registry: registry,
		path:     path,
		tag:      tag,
	}
}

// List available plugins.
func (od *OCIDiscovery) List() (plugins []common.Plugin, err error) {
	return
}

// Describe a plugin.
func (od *OCIDiscovery) Describe(name string) (plugin common.Plugin, err error) {
	return
}

// Name of the repository.
func (od *OCIDiscovery) Name() string {
	return od.name
}

// Type of the discovery.
func (od *OCIDiscovery) Type() string {
	return "GCP"
}
