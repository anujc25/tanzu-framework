// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package distribution

// DistributionConfig contains a specific distribution mechanism. Only one of the
// configs must be set.
type DistributionConfig struct {
	// GCPStorage is set if the plugin is distributed via Google Cloud Storage.
	GCP *GCPStorage `json:"gcp"`
	// GCPStorage is set if the plugin is distributed via an OCI Image Registry.
	OCI *OCIRegistry `json:"oci"`
}

// GCPStorage provides a Google Cloud Storage bucket with an optional base path (or
// object prefix). The object download path name is constructed as
// `{Bucket}/{BasePath}/{Version}/{OS}/{Arch}`.
type GCPStorage struct {
	// Bucket is a Google Cloud Storage bucket.
	// E.g., tanzu-cli
	Bucket string `json:"bucket"`
	// BasePath is a URI path that is prefixed to the object name/path.
	// E.g., plugins/cluster
	BasePath string `json:"basePath"`
}

// OCIRegistry provides an OCI compliant image registry which supports multiple
// architecture image manifests. The fully qualified image name is constructed as
// `{Registry}/{Name}:{Version}`.
type OCIRegistry struct {
	// Registry is an OCI compliant image registry. It MUST be a DNS-compatible name.
	// E.g., harbor.my-domain.local
	Registry string `json:"registry,omitempty"`
	// Name is the unique repository/image name. It MUST be a valid URI path, MAY
	// contain zero or more '/', and SHOULD NOT start or end with '/'.
	// E.g., tanzu/cli/plugins/cluster
	Name string `json:"name"`
}
