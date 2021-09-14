// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package distribution

import (
	"github.com/pkg/errors"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
)

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

func CreateDistribution(dc cliv1alpha1.DistributionConfig) (Distribution, error) {
	if dc.GCP != nil {
		return NewGCPDistribution(dc.GCP.Bucket, dc.GCP.BasePath), nil
	}
	if dc.OCI != nil {
		return NewOCIDistribution(dc.OCI.Registry, dc.OCI.ImagePath), nil
	}
	return nil, errors.New("unknown distribution")
}
