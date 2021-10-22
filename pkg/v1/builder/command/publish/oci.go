// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package publish

type OCIPublisher struct {
	OCIDiscoveryImageRepository          string
	OCIDistributionImageRepositoryPrefix string

	LocalDiscoveryPath string
}

func NewOCIPublisher(
	ociDiscoveryImageRepository,
	ociDistributionImageRepositoryPrefix,
	localDiscoveryPath string) Publisher {

	return &OCIPublisher{
		OCIDiscoveryImageRepository:          ociDiscoveryImageRepository,
		OCIDistributionImageRepositoryPrefix: ociDistributionImageRepositoryPrefix,
		LocalDiscoveryPath:                   localDiscoveryPath,
	}
}

func (o *OCIPublisher) PublishPlugin(version, os, arch, plugin, sourcePath string) (string, error) {
	return "", nil
}

func (o *OCIPublisher) PublishDiscovery() error {
	return nil
}
