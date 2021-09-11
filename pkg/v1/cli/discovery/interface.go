// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package discovery

import (
	"github.com/vmware-tanzu/tanzu-framework/apis/config/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
)

// Discovery is an interface to fetch list of available plugins
type Discovery interface {
	// Name of the repository.
	Name() string
	// List available plugins.
	List() ([]common.Plugin, error)
	// Describe a plugin.
	Describe(name string) (common.Plugin, error)
}

func CreateDiscovery(pd v1alpha1.PluginDiscovery) Discovery {
	if pd.GCP != nil {
		return NewGCPDiscovery(pd.GCP.Bucket, pd.GCP.ManifestPath, pd.GCP.Name)
	}
	if pd.OCI != nil {
		return NewOCIDiscovery(pd.OCI.Name, pd.OCI.Registry, pd.OCI.Path, pd.OCI.Tag)
	}
	if pd.Local != nil {
		return NewLocalDiscovery(pd.Local.Name, pd.Local.ManifestPath)
	}
	return nil
}
