// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package distribution

// OCIRegistry provides an OCI compliant image registry which supports multiple
// architecture image manifests. The fully qualified image name is constructed as
// `{Registry}/{Name}:{Version}`.
type OCIRegistry struct {
	// Registry is an OCI compliant image registry. It MUST be a DNS-compatible name.
	// E.g., harbor.my-domain.local
	Registry string `json:"registry,omitempty"`
	// ImagePath is the unique repository/image name. It MUST be a valid URI path, MAY
	// contain zero or more '/', and SHOULD NOT start or end with '/'.
	// E.g., tanzu/cli/plugins/cluster
	ImagePath string `json:"name"`
}

// NewOCIDistribution returns a new OCI storage distribution.
func NewOCIDistribution(registry, imagePath string) Distribution {
	return &OCIRegistry{
		Registry:  registry,
		ImagePath: imagePath,
	}
}

// Fetch an artifact.
func (g *OCIRegistry) Fetch(name, version string, arch string) ([]byte, error) {
	return nil, nil
}

// FetchTest fetches a test artifact.
func (g *OCIRegistry) FetchTest(name, version string, arch string) ([]byte, error) {
	return nil, nil
}

// Get the relative installation path for a plugin binary
func (g *OCIRegistry) GetInstallationPath() string {
	return ""
}
