// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package publish

import (
	"path/filepath"

	"github.com/aunum/log"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
)

type LocalPublisher struct {
	Plugins               []string
	InputArtifactDir      string
	OSArch                []string
	LocalDiscoveryPath    string
	LocalDistributionPath string
}

func NewLocalPublisher(plugins []string, version string, arrOSArch []string, localDiscoveryPath, localDistributionPath, inputArtifactDir string) Publisher {
	_ = ensureResourceDir(localDiscoveryPath, true)
	_ = ensureResourceDir(localDistributionPath, false)

	return &LocalPublisher{
		Plugins:               plugins,
		LocalDiscoveryPath:    localDiscoveryPath,
		LocalDistributionPath: localDistributionPath,
		OSArch:                arrOSArch,
		InputArtifactDir:      inputArtifactDir,
	}
}

// PublishPlugins publishes plugins binaries and
// CLIPlugin resource files for discovery to local directories
func (l *LocalPublisher) PublishPlugins() error {

	availablePluginInfo := detectAvailablePluginInfo(l.InputArtifactDir, l.Plugins, l.OSArch)

	for plugin, pluginInfo := range availablePluginInfo {
		log.Info("Processing plugin:", plugin)
		mapVersionArtifactList := make(map[string]v1alpha1.ArtifactList)

		// Create version based artifact list
		for version, arrOSArch := range pluginInfo.versions {
			artifacts := make([]v1alpha1.Artifact, 0)
			for _, oa := range arrOSArch {
				sourcePath, digest, err := getPluginPathAndDigestFromMetadata(l.InputArtifactDir, plugin, version, oa.os, oa.arch)
				if err != nil {
					return err
				}

				destPath, err := l.publishPlugin(sourcePath, version, oa.os, oa.arch, plugin)
				if err != nil {
					return err
				}

				artifacts = append(artifacts, newArtifactObject(oa.os, oa.arch, common.DistributionTypeLocal, digest, destPath))
			}
			mapVersionArtifactList[version] = artifacts
		}

		// Create new CLIPlugin resource based on plugin and artifact info
		cliPlugin := newCLIPluginResource(plugin, pluginInfo.description, pluginInfo.recommendedVersion, mapVersionArtifactList)

		err := saveCLIPluginResource(cliPlugin, filepath.Join(l.LocalDiscoveryPath, plugin+".yaml"))
		if err != nil {
			return errors.Wrap(err, "could not write CLIPlugin to file")
		}
	}

	return nil
}

func (l *LocalPublisher) publishPlugin(sourcePath, version, os, arch, plugin string) (string, error) {
	destPath := filepath.Join(l.LocalDistributionPath, os, arch, "cli", plugin, version, "tanzu-"+plugin+"-"+os+"_"+arch)
	if os == "windows" {
		sourcePath = sourcePath + ".exe"
		destPath = destPath + ".exe"
	}
	err := copy.Copy(sourcePath, destPath)
	if err != nil {
		return "", err
	}
	return destPath, nil
}
