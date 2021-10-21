// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package publish

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aunum/log"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	apimachineryjson "k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
)

type LocalPublisher struct {
	StandalonePlugins     []string
	ContextAwarePlugins   []string
	Version               string
	OSArch                []string
	LocalDiscoveryPath    string
	LocalDistributionPath string
	LocalArtifactDir      string
}

func NewLocalPublisher(sp []string, cap []string, version string, oa []string, localDistroPath string, localArtifactDir string) Publisher {
	localDiscoveryPath := filepath.Join(localDistroPath, "discovery")
	localDistributionPath := filepath.Join(localDistroPath, "distribution")

	_ = os.RemoveAll(localDiscoveryPath)
	_ = os.RemoveAll(localDistributionPath)

	return &LocalPublisher{
		StandalonePlugins:     sp,
		ContextAwarePlugins:   cap,
		Version:               version,
		OSArch:                oa,
		LocalDiscoveryPath:    localDiscoveryPath,
		LocalDistributionPath: localDistributionPath,
		LocalArtifactDir:      localArtifactDir,
	}
}

func (l *LocalPublisher) PublishStandAlonePlugin() error {
	return l.publishPlugins(filepath.Join(l.LocalDiscoveryPath, "standalone"), l.StandalonePlugins)
}

func (l *LocalPublisher) PublishContextAwarePlugin() error {
	return l.publishPlugins(filepath.Join(l.LocalDiscoveryPath, "context"), l.ContextAwarePlugins)
}

func (l *LocalPublisher) publishPlugins(discoveryResourceDir string, plugins []string) error {
	if err := l.ensureResourceDir(discoveryResourceDir); err != nil {
		return err
	}

	for _, plugin := range plugins {
		log.Info("Processing plugin:", plugin)
		cliPlugin := v1alpha1.CLIPlugin{}
		cliPlugin.SetGroupVersionKind(v1alpha1.GroupVersionKindCLIPlugin)
		cliPlugin.SetName(plugin)
		cliPlugin.Spec.Description = plugin
		cliPlugin.Spec.RecommendedVersion = l.Version
		if cliPlugin.Spec.Artifacts == nil {
			cliPlugin.Spec.Artifacts = make(map[string]v1alpha1.ArtifactList)
		}

		artifacts := make([]v1alpha1.Artifact, 0)
		for _, osArch := range l.OSArch {
			arr := strings.Split(osArch, "-")
			os := arr[0]
			arch := arr[1]

			sourcePath := filepath.Join(l.LocalArtifactDir, os, arch, "cli", plugin, l.Version, "tanzu-"+plugin+"-"+os+"_"+arch)
			destPath := filepath.Join(l.LocalDistributionPath, os, arch, "cli", plugin, l.Version, "tanzu-"+plugin+"-"+os+"_"+arch)
			if os == "windows" {
				sourcePath = sourcePath + ".exe"
				destPath = destPath + ".exe"
			}
			copy.Copy(sourcePath, destPath)

			digest, err := utils.SHA256FromFile(sourcePath)
			if err != nil {
				return errors.Wrap(err, "error while calculating sha256")
			}

			artifact := v1alpha1.Artifact{
				Type:   "local",
				OS:     os,
				Arch:   arch,
				URI:    destPath,
				Digest: digest,
			}
			artifacts = append(artifacts, artifact)
		}
		cliPlugin.Spec.Artifacts[l.Version] = artifacts

		err := os.MkdirAll(discoveryResourceDir, 0755)
		if err != nil {
			return errors.Wrap(err, "could not create dir")
		}
		discoveryResourceFile := filepath.Join(discoveryResourceDir, plugin+".yaml")

		fo, err := os.Create(discoveryResourceFile)
		if err != nil {
			return errors.Wrap(err, "could not create resource file")
		}
		defer fo.Close()

		scheme, err := v1alpha1.SchemeBuilder.Build()
		if err != nil {
			return errors.Wrap(err, "failed to create scheme")
		}
		e := apimachineryjson.NewSerializerWithOptions(apimachineryjson.DefaultMetaFactory, scheme, scheme,
			apimachineryjson.SerializerOptions{Yaml: true, Pretty: false, Strict: false})

		err = e.Encode(&cliPlugin, fo)
		if err != nil {
			return errors.Wrap(err, "could not write to CLIPlugin resource file")
		}
	}

	return nil
}

func (l *LocalPublisher) ensureResourceDir(resourceDir string) error {
	_ = os.RemoveAll(resourceDir)
	if err := os.MkdirAll(resourceDir, 0755); err != nil {
		return errors.Wrapf(err, "unable to create resource directory '%v'", l.LocalDiscoveryPath)
	}
	return nil
}
