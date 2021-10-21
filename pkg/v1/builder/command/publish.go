// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/builder/command/publish"
)

var (
	distroType, pluginsString, oa, inputArtifactDir                   string
	localOutputDiscoveryDir, localOutputDistributionDir               string
	ociDiscoveryImageRepository, ociDistributionImageRepositoryPrefix string
)

// PublishCmd publishes plugin resources
var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish operations",
	RunE:  publishPlugins,
}

func init() {
	PublishCmd.Flags().StringVar(&distroType, "type", "", "type of discovery and distribution for publishing plugins. Supported: local")
	PublishCmd.Flags().StringVar(&pluginsString, "plugins", "", "list of plugin names. Example: 'login management-cluster cluster'")
	PublishCmd.Flags().StringVar(&version, "version", "", "version of the plugins to publish")
	PublishCmd.Flags().StringVar(&oa, "os-arch", "", "list of os-arch. Format: 'darwin-amd64 linux-amd64 windows-amd64'")
	PublishCmd.Flags().StringVar(&inputArtifactDir, "input-artifact-dir", "", "artifact directory which is a output of 'tanzu builder cli compile' command")

	PublishCmd.Flags().StringVar(&localOutputDiscoveryDir, "local-output-discovery-dir", "", "local output directory where CLIPlugin resource yamls for discovery will be placed. Applicable to 'local' type")
	PublishCmd.Flags().StringVar(&localOutputDistributionDir, "local-output-distribution-dir", "", "local output directory where plugin binaries will be placed. Applicable to 'local' type")

	PublishCmd.Flags().StringVar(&ociDiscoveryImageRepository, "oci-discovery-image-repository", "", "image path to publish oci image with CLIPlugin resource yamls. Applicable to 'oci' type")
	PublishCmd.Flags().StringVar(&ociDistributionImageRepositoryPrefix, "oci-distribution-image-repository-prefix", "", "image path prefix to publish oci image for plugin binaries. Applicable to 'oci' type")

	PublishCmd.MarkFlagRequired("type")
	PublishCmd.MarkFlagRequired("version")
	PublishCmd.MarkFlagRequired("os-arch")
	PublishCmd.MarkFlagRequired("input-artifact-dir")
}

func publishPlugins(cmd *cobra.Command, args []string) error {
	plugins := strings.Split(pluginsString, " ")
	osArch := strings.Split(oa, " ")

	switch strings.ToLower(distroType) {
	case "local":
		return publish.NewLocalPublisher(plugins, version, osArch, localOutputDiscoveryDir, localOutputDistributionDir, inputArtifactDir).PublishPlugins()
	case "oci":
		return publish.NewOCIPublisher(plugins, version, osArch, ociDiscoveryImageRepository, ociDistributionImageRepositoryPrefix, inputArtifactDir).PublishPlugins()
	default:
		return errors.Errorf("publish plugins with type %s is not yet supported", distroType)
	}
}
