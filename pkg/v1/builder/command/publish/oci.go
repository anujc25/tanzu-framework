// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package publish

import (
	"path/filepath"

	"github.com/aunum/log"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
)

type OCIPublisher struct {
	Plugins          []string
	Version          string
	OSArch           []string
	InputArtifactDir string

	OCIDiscoveryImageRepository          string
	OCIDistributionImageRepositoryPrefix string

	localDiscoveryPath string
}

func NewOCIPublisher(plugins []string,
	version string,
	oa []string,
	ociDiscoveryImageRepository,
	ociDistributionImageRepositoryPrefix,
	inputArtifactDir string) Publisher {

	localDiscoveryPath := filepath.Join(common.DefaultLocalPluginDistroDir, "discovery", "oci")
	_ = ensureResourceDir(localDiscoveryPath, true)

	return &OCIPublisher{
		Plugins:                              plugins,
		Version:                              version,
		OSArch:                               oa,
		InputArtifactDir:                     inputArtifactDir,
		OCIDiscoveryImageRepository:          ociDiscoveryImageRepository,
		OCIDistributionImageRepositoryPrefix: ociDistributionImageRepositoryPrefix,
		localDiscoveryPath:                   localDiscoveryPath,
	}
}

// PublishPlugins publishes plugins binaries and
// CLIPlugin resource files for discovery to oci image repository
func (o *OCIPublisher) PublishPlugins() error {

	for _, plugin := range o.Plugins {
		log.Info("Processing plugin:", plugin)

		artifacts := make([]v1alpha1.Artifact, 0)
		for _, oa := range o.OSArch {
			os, arch := osArch(oa)

			sourcePath, digest, err := getPluginPathAndDigestFromMetadata(o.InputArtifactDir, plugin, o.Version, os, arch)
			if err != nil {
				return err
			}

			imagePath, err := o.publishPlugin(sourcePath, os, arch, plugin)
			if err != nil {
				return err
			}

			artifacts = append(artifacts, createArtifactObject(os, arch, common.DistributionTypeOCI, digest, imagePath))
		}

		cliPlugin := createCLIPluginResource(plugin, plugin, o.Version, artifacts)

		err := writeCLIPluginToFile(cliPlugin, filepath.Join(o.localDiscoveryPath, plugin+".yaml"))
		if err != nil {
			return errors.Wrap(err, "could not write CLIPlugin to file")
		}
	}

	return o.publishDiscovery()
}

func (o *OCIPublisher) publishPlugin(sourcePath, os, arch, plugin string) (string, error) {
	return "", nil
}

func (o *OCIPublisher) publishDiscovery() error {
	return nil
}
