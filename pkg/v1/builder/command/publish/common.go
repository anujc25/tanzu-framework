// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package publish

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
	apimachineryjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func createCLIPluginResource(plugin, description, version string, artifacts []v1alpha1.Artifact) v1alpha1.CLIPlugin {
	cliPlugin := v1alpha1.CLIPlugin{}
	cliPlugin.SetGroupVersionKind(v1alpha1.GroupVersionKindCLIPlugin)
	cliPlugin.SetName(plugin)
	cliPlugin.Spec.Description = description
	cliPlugin.Spec.RecommendedVersion = version
	cliPlugin.Spec.Artifacts = make(map[string]v1alpha1.ArtifactList)
	cliPlugin.Spec.Artifacts[version] = artifacts
	return cliPlugin
}

func createArtifactObject(os, arch, artifactType, digest, uri string) v1alpha1.Artifact {
	artifact := v1alpha1.Artifact{
		Type:   artifactType,
		OS:     os,
		Arch:   arch,
		Digest: digest,
	}

	if artifactType == common.DistributionTypeOCI {
		artifact.Image = uri
	} else {
		artifact.URI = uri
	}
	return artifact
}

func getPluginPathAndDigestFromMetadata(artifactDir, plugin, version, os, arch string) (string, string, error) {
	sourcePath := filepath.Join(artifactDir, os, arch, "cli", plugin, version, "tanzu-"+plugin+"-"+os+"_"+arch)
	digest, err := utils.SHA256FromFile(sourcePath)
	if err != nil {
		return "", "", errors.Wrap(err, "error while calculating sha256")
	}
	return sourcePath, digest, nil
}

func writeCLIPluginToFile(cliPlugin v1alpha1.CLIPlugin, discoveryResourceFile string) error {
	discoveryResourceDir := filepath.Dir(discoveryResourceFile)

	err := os.MkdirAll(discoveryResourceDir, 0755)
	if err != nil {
		return errors.Wrap(err, "could not create dir")
	}

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
	return nil
}

func osArch(osArch string) (string, string) {
	arr := strings.Split(osArch, "-")
	return arr[0], arr[1]
}

func ensureResourceDir(resourceDir string, cleanDir bool) error {
	if cleanDir {
		_ = os.RemoveAll(resourceDir)
	}
	if err := os.MkdirAll(resourceDir, 0755); err != nil {
		return errors.Wrapf(err, "unable to create resource directory '%v'", resourceDir)
	}
	return nil
}
