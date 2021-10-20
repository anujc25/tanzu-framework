// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package publish

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	apimachineryjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
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
	return &LocalPublisher{
		StandalonePlugins:     sp,
		ContextAwarePlugins:   cap,
		Version:               version,
		OSArch:                oa,
		LocalDiscoveryPath:    filepath.Join(localDistroPath, "discovery"),
		LocalDistributionPath: filepath.Join(localDistroPath, "distribution"),
		LocalArtifactDir:      localArtifactDir,
	}
}

func (l *LocalPublisher) PublishStandAlonePlugin() error {
	if err := l.ensureLocalDistro(); err != nil {
		return err
	}

	for _, plugin := range l.StandalonePlugins {
		cliPlugin := v1alpha1.CLIPlugin{}
		cliPlugin.SetName(plugin)
		cliPlugin.Spec.Description = plugin
		cliPlugin.Spec.RecommendedVersion = l.Version

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

			artifact := v1alpha1.Artifact{
				Type: "local",
				OS:   os,
				Arch: arch,
				URI:  destPath,
			}
			artifacts = append(artifacts, artifact)
		}
		cliPlugin.Spec.Artifacts[l.Version] = artifacts

		scheme, err := v1alpha1.SchemeBuilder.Build()
		if err != nil {
			return errors.Wrap(err, "failed to create scheme")
		}
		s := apimachineryjson.NewSerializerWithOptions(apimachineryjson.DefaultMetaFactory, scheme, scheme,
			apimachineryjson.SerializerOptions{Yaml: true, Pretty: false, Strict: false})

		discoveryResourceFile := filepath.Join(l.LocalDiscoveryPath, "standalone", plugin+".yaml")
		f, err := os.Create(discoveryResourceFile)
		defer f.Close()

		err = s.Encode(&cliPlugin, f)
		if err != nil {
			return errors.Wrap(err, "could not encode CLIPlugin object")
		}
	}

	return nil
}

func (l *LocalPublisher) PublishContextAwarePlugin() error {
	if err := l.ensureLocalDistro(); err != nil {
		return err
	}
	return nil
}

func (l *LocalPublisher) ensureLocalDistro() error {
	_ = os.RemoveAll(l.LocalDiscoveryPath)
	_ = os.RemoveAll(l.LocalDistributionPath)

	if err := os.MkdirAll(l.LocalDiscoveryPath, 0644); err != nil {
		return errors.Wrapf(err, "unable to create local discovery directory '%v'", l.LocalDiscoveryPath)
	}
	if err := os.MkdirAll(l.LocalDistributionPath, 0644); err != nil {
		return errors.Wrapf(err, "unable to create local distribution directory '%v'", l.LocalDistributionPath)
	}
	return nil
}
