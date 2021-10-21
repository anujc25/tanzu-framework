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
	Plugins          []string
	Version          string
	OSArch           []string
	InputArtifactDir string

	LocalDiscoveryPath    string
	LocalDistributionPath string
}

func NewLocalPublisher(plugins []string, version string, oa []string, localDiscoveryPath, localDistributionPath, inputArtifactDir string) Publisher {
	_ = ensureResourceDir(localDiscoveryPath, true)
	_ = ensureResourceDir(localDistributionPath, false)

	return &LocalPublisher{
		Plugins:               plugins,
		Version:               version,
		OSArch:                oa,
		LocalDiscoveryPath:    localDiscoveryPath,
		LocalDistributionPath: localDistributionPath,
		InputArtifactDir:      inputArtifactDir,
	}
}

// PublishPlugins publishes plugins binaries and
// CLIPlugin resource files for discovery to local directories
func (l *LocalPublisher) PublishPlugins() error {
	for _, plugin := range l.Plugins {
		log.Info("Processing plugin:", plugin)

		artifacts := make([]v1alpha1.Artifact, 0)
		for _, oa := range l.OSArch {
			os, arch := osArch(oa)

			sourcePath, digest, err := getPluginPathAndDigestFromMetadata(l.InputArtifactDir, plugin, l.Version, os, arch)
			if err != nil {
				return err
			}

			destPath, err := l.publishPlugin(sourcePath, os, arch, plugin)
			if err != nil {
				return err
			}

			artifacts = append(artifacts, createArtifactObject(os, arch, common.DistributionTypeLocal, digest, destPath))
		}

		cliPlugin := createCLIPluginResource(plugin, plugin, l.Version, artifacts)

		err := writeCLIPluginToFile(cliPlugin, filepath.Join(l.LocalDiscoveryPath, plugin+".yaml"))
		if err != nil {
			return errors.Wrap(err, "could not write CLIPlugin to file")
		}
	}

	return nil
}

func (l *LocalPublisher) publishPlugin(sourcePath, os, arch, plugin string) (string, error) {
	destPath := filepath.Join(l.LocalDistributionPath, os, arch, "cli", plugin, l.Version, "tanzu-"+plugin+"-"+os+"_"+arch)
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
