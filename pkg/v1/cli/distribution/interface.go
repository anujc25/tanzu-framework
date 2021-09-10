// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package distribution

// Distribution is an interface to download a single plugin binary.
type Distribution interface {
	// Fetch the binary for a plugin version.
	Fetch(version, os, arch string) ([]byte, error)
	// Fetch the test binary for a plugin version.
	FetchTest(version, os, arch string) ([]byte, error)
	// Get the relative installation path for a plugin binary. E.g.,
	// GCP: storage.googleapis.com/{BucketName}/{BasePath}
	// OCI: {Registry}/{Repository}
	GetInstallationPath() string
}
