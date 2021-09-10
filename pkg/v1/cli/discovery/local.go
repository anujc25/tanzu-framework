// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package discovery

import "github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"

// LocalDiscovery is a artifact discovery endpoint utilizing a local host os.
type LocalDiscovery struct {
	path string
	name string
}

// NewLocalDiscovery returns a new local repository.
func NewLocalDiscovery(name, localPath string) Discovery {
	return &LocalDiscovery{
		path: localPath,
		name: name,
	}
}

// List available plugins.
func (l *LocalDiscovery) List() (plugins []common.Plugin, err error) {
	return
}

// Describe a plugin.
func (l *LocalDiscovery) Describe(name string) (plugin common.Plugin, err error) {
	return
}

// Name of the repository.
func (l *LocalDiscovery) Name() string {
	return l.name
}
