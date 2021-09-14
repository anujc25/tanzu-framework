// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package discovery

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
	"gopkg.in/yaml.v2"
)

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
func (l *LocalDiscovery) List() ([]common.Plugin, error) {
	manifest, err := l.Manifest()
	if err != nil {
		return nil, err
	}
	return manifest.Plugins, nil
}

// Describe a plugin.
func (l *LocalDiscovery) Describe(name string) (plugin common.Plugin, err error) {
	manifest, err := l.Manifest()
	if err != nil {
		return
	}

	for _, p := range manifest.Plugins {
		if p.Name == name {
			plugin = p
			return
		}
	}
	err = errors.Errorf("cannot find plugin with name '%v'", name)
	return
}

// Name of the repository.
func (l *LocalDiscovery) Name() string {
	return l.name
}

// Manifest returns the manifest for a local repository.
func (l *LocalDiscovery) Manifest() (manifest common.Manifest, err error) {
	b, err := os.ReadFile(l.path)
	if err != nil {
		err = errors.Wrapf(err, "error while reading manifest file")
		return
	}

	err = yaml.Unmarshal(b, &manifest)
	if err != nil {
		err = fmt.Errorf("could not unmarshal manifest.yaml: %v", err)
	}

	for i := range manifest.Plugins {
		manifest.Plugins[i].Discovery = fmt.Sprintf("%s/%s", l.Type(), l.name)
	}

	return
}

// Type of the repository.
func (l *LocalDiscovery) Type() string {
	return "local"
}
