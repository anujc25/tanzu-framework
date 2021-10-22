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
	"gopkg.in/yaml.v2"

	apimachineryjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type osArch struct {
	os   string
	arch string
}

type pluginInfo struct {
	recommendedVersion string
	description        string
	versions           map[string][]osArch
}

func detectAvailablePluginInfo(artifactDir string, plugins []string, arrOSArch []string) (map[string]*pluginInfo, error) {
	mapPluginInfo := make(map[string]*pluginInfo)

	// For all plugins
	for _, plugin := range plugins {
		// For all supported OS
		for _, osArch := range arrOSArch {
			o, a := splitOSArch(osArch)

			// get all directory under plugin directory
			pluginDir := filepath.Join(artifactDir, o, a, "cli", plugin)
			files, err := os.ReadDir(pluginDir)
			if err != nil {
				return nil, errors.Errorf("unable to find plugin artifact directory for plugin:'%s' os:'%s', arch:'%s' [directory: '%s']", plugin, o, a, pluginDir)
			}

			// Each directory under the plugin directory is considered version directory
			for _, file := range files {
				if file.IsDir() {
					updatePluginInfoMapWithVersionOSArch(mapPluginInfo, plugin, file.Name(), o, a)
				}
			}

			// TODO: How to decide recommanded version of the plugin
			// Currently, using plugin.yaml file to fetch description and recommanded version
			// We need to make sure that `plugin.yaml` points to recommanded version among
			// all available version if there are more than 1 versions
			recommandedVersion, description := getRecommandedVersionDescriptionFromPluginYaml(filepath.Join(artifactDir, o, a, "cli", plugin, "plugin.yaml"))
			// Update recommanded version and Description
			updatePluginInfoMapWithRecommandedVersionDescription(mapPluginInfo, plugin, recommandedVersion, description)
		}
	}

	return mapPluginInfo, nil
}

func updatePluginInfoMapWithRecommandedVersionDescription(mapPluginInfo map[string]*pluginInfo, plugin, recommendedVersion, description string) {
	if mapPluginInfo[plugin] == nil {
		mapPluginInfo[plugin] = &pluginInfo{}
		mapPluginInfo[plugin].versions = make(map[string][]osArch, 0)
	}
	mapPluginInfo[plugin].recommendedVersion = recommendedVersion
	mapPluginInfo[plugin].description = description
}

func updatePluginInfoMapWithVersionOSArch(mapPluginInfo map[string]*pluginInfo, plugin, version, os, arch string) {
	if mapPluginInfo[plugin] == nil {
		mapPluginInfo[plugin] = &pluginInfo{}
		mapPluginInfo[plugin].versions = make(map[string][]osArch, 0)
	}

	if mapPluginInfo[plugin].versions[version] == nil {
		mapPluginInfo[plugin].versions[version] = make([]osArch, 0)
	}

	oa := mapPluginInfo[plugin].versions[version]
	oa = append(oa, osArch{os: os, arch: arch})

	mapPluginInfo[plugin].versions[version] = oa
}

func getRecommandedVersionDescriptionFromPluginYaml(pluginYaml string) (string, string) {
	b, err := os.ReadFile(pluginYaml)
	if err == nil {
		pd := &v1alpha1.PluginDescriptor{}
		err := yaml.Unmarshal(b, pd)
		if err == nil {
			return pd.Version, pd.Description
		}
	}
	return "", ""
}

func newCLIPluginResource(plugin, description, version string, artifacts map[string]v1alpha1.ArtifactList) v1alpha1.CLIPlugin {
	cliPlugin := v1alpha1.CLIPlugin{}
	cliPlugin.SetGroupVersionKind(v1alpha1.GroupVersionKindCLIPlugin)
	cliPlugin.SetName(plugin)
	cliPlugin.Spec.Description = description
	cliPlugin.Spec.RecommendedVersion = version
	cliPlugin.Spec.Artifacts = artifacts
	return cliPlugin
}

func newArtifactObject(os, arch, artifactType, digest, uri string) v1alpha1.Artifact {
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
	if os == "windows" {
		sourcePath = sourcePath + ".exe"
	}
	digest, err := utils.SHA256FromFile(sourcePath)
	if err != nil {
		return "", "", errors.Wrap(err, "error while calculating sha256")
	}
	return sourcePath, digest, nil
}

func saveCLIPluginResource(cliPlugin v1alpha1.CLIPlugin, discoveryResourceFile string) error {
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

func ensureResourceDir(resourceDir string, cleanDir bool) error {
	if cleanDir {
		_ = os.RemoveAll(resourceDir)
	}
	if err := os.MkdirAll(resourceDir, 0755); err != nil {
		return errors.Wrapf(err, "unable to create resource directory '%v'", resourceDir)
	}
	return nil
}

func splitOSArch(osArch string) (string, string) {
	arr := strings.Split(osArch, "-")
	if len(arr) < 2 {
		return "", ""
	}
	return arr[0], arr[1]
}
