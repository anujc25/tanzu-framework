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

type PublishMetadata struct {
	Plugins            []string
	InputArtifactDir   string
	OSArch             []string
	LocalDiscoveryPath string
	PublisherInterface Publisher
}

func PublishPlugins(g PublishMetadata) error {
	_ = ensureResourceDir(g.LocalDiscoveryPath, true)
	// _ = ensureResourceDir(localDistributionPath, false) // TODO: fix this

	availablePluginInfo := detectAvailablePluginInfo(g.InputArtifactDir, g.Plugins, g.OSArch)

	for plugin, pluginInfo := range availablePluginInfo {
		log.Info("Processing plugin:", plugin)
		mapVersionArtifactList := make(map[string]v1alpha1.ArtifactList)

		// Create version based artifact list
		for version, arrOSArch := range pluginInfo.versions {
			artifacts := make([]v1alpha1.Artifact, 0)
			for _, oa := range arrOSArch {
				sourcePath, digest, err := getPluginPathAndDigestFromMetadata(g.InputArtifactDir, plugin, version, oa.os, oa.arch)
				if err != nil {
					return err
				}

				destPath, err := g.PublisherInterface.PublishPlugin(sourcePath, version, oa.os, oa.arch, plugin)
				if err != nil {
					return err
				}

				artifacts = append(artifacts, newArtifactObject(oa.os, oa.arch, common.DistributionTypeLocal, digest, destPath))
			}
			mapVersionArtifactList[version] = artifacts
		}

		// Create new CLIPlugin resource based on plugin and artifact info
		cliPlugin := newCLIPluginResource(plugin, pluginInfo.description, pluginInfo.recommendedVersion, mapVersionArtifactList)

		err := saveCLIPluginResource(cliPlugin, filepath.Join(g.LocalDiscoveryPath, plugin+".yaml"))
		if err != nil {
			return errors.Wrap(err, "could not write CLIPlugin to file")
		}
	}

	return g.PublisherInterface.PublishDiscovery()
}
