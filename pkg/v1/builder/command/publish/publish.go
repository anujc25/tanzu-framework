// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package publish

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const local = "local"

var minConcurrent = 2
var identifiers = []string{
	string('\U0001F435'),
	string('\U0001F43C'),
	string('\U0001F436'),
	string('\U0001F430'),
	string('\U0001F98A'),
	string('\U0001F431'),
	string('\U0001F981'),
	string('\U0001F42F'),
	string('\U0001F42E'),
	string('\U0001F437'),
	string('\U0001F42D'),
	string('\U0001F428'),
}

var (
	distroType, sp, cp, version, oa           string
	localDistroPath, localArtifactDir         string
	imageRepoDiscovery, imageRepoDistribution string
)

// PublishCmd publishes plugin resources
var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish operations",
	RunE:  publish,
}

func init() {
	PublishCmd.Flags().StringVar(&distroType, "type", "", "")
	PublishCmd.Flags().StringVar(&sp, "standalone-plugins", "", "")
	PublishCmd.Flags().StringVar(&cp, "context-plugins", "", "")
	PublishCmd.Flags().StringVar(&version, "version", "", "")
	PublishCmd.Flags().StringVar(&oa, "os-arch", "", "")
	PublishCmd.Flags().StringVar(&localDistroPath, "local-distro-path", "", "")
	PublishCmd.Flags().StringVar(&localArtifactDir, "local-artifact-dir", "", "")
	PublishCmd.Flags().StringVar(&imageRepoDiscovery, "image-repository-discovery", "", "")
	PublishCmd.Flags().StringVar(&imageRepoDistribution, "image-repository-distribution", "", "")

	PublishCmd.MarkFlagRequired("type")
	PublishCmd.MarkFlagRequired("version")
	PublishCmd.MarkFlagRequired("os-arch")
	PublishCmd.MarkFlagRequired("local-artifact-dir")
}

func publish(cmd *cobra.Command, args []string) error {
	standalonePlugins := strings.Split(sp, " ")
	contextPlugins := strings.Split(cp, " ")
	osArch := strings.Split(oa, " ")

	switch distroType {
	case "local":
		lp := NewLocalPublisher(standalonePlugins, contextPlugins, version, osArch, localDistroPath, localArtifactDir)
		if err := lp.PublishStandAlonePlugin(); err != nil {
			return errors.Wrap(err, "error while publishing standalone plugins")
		}
		if err := lp.PublishContextAwarePlugin(); err != nil {
			return errors.Wrap(err, "error while publishing context-aware plugins")
		}
	case "oci":
		return publishOCI()
	}
	return nil
}

func publishOCI() error {
	return errors.New("publishing to OCI is not yet implemented")
}
