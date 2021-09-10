// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package discovery

import "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"

// Discovery is an interface to fetch list of available plugins
type Discovery interface {
	// Name of the repository.
	Name() string
	// List available plugins.
	List() ([]common.Plugin, error)
	// Describe a plugin.
	Describe(name string) (common.Plugin, error)
}
