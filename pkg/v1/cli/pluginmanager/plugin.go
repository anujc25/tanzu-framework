// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pluginmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aunum/log"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/mod/semver"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/apis/config/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/catalog"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/discovery"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/distribution"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/config"
)

const (
	// exe is an executable file extension
	exe = ".exe"
)

var (
	minConcurrent = 2
	// PluginRoot is the plugin root where plugins are installed
	pluginRoot = common.DefaultPluginRoot
	// Distro is set of plugins that should be included with the CLI.
	distro = common.DefaultDistro
)

// ValidatePlugin validates the plugin descriptor.
func ValidatePlugin(p *cliv1alpha1.PluginDescriptor) (err error) {
	// skip builder plugin for bootstrapping
	if p.Name == "builder" {
		return nil
	}
	if p.Name == "" {
		err = multierr.Append(err, fmt.Errorf("plugin %q name cannot be empty", p.Name))
	}
	if p.Version == "" {
		err = multierr.Append(err, fmt.Errorf("plugin %q version cannot be empty", p.Name))
	}
	if !semver.IsValid(p.Version) && p.Version != "dev" {
		err = multierr.Append(err, fmt.Errorf("version %q %q is not a valid semantic version", p.Name, p.Version))
	}
	if p.Description == "" {
		err = multierr.Append(err, fmt.Errorf("plugin %q description cannot be empty", p.Name))
	}
	if p.Group == "" {
		err = multierr.Append(err, fmt.Errorf("plugin %q group cannot be empty", p.Name))
	}
	return
}

func discoverPlugins(pd []v1alpha1.PluginDiscovery) ([]common.Plugin, error) {
	allPlugins := []common.Plugin{}
	for _, d := range pd {
		discObject := discovery.CreateDiscovery(d)
		plugins, err := discObject.List()
		if err != nil {
			return nil, errors.Wrapf(err, "unable to list plugin from discovery '%v'", discObject.Name())
		}
		allPlugins = append(allPlugins, plugins...)
	}
	return allPlugins, nil
}

// DiscoverStandalonePlugins returns the available standalone plugins
func DiscoverStandalonePlugins() ([]common.Plugin, error) {
	cfg, err := config.GetClientConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get client configuration")
	}

	if cfg == nil || cfg.ClientOptions == nil || cfg.ClientOptions.CLI == nil {
		return []common.Plugin{}, nil
	}

	// TODO: Need to mark these plugins as standalone plugins

	return discoverPlugins(cfg.ClientOptions.CLI.Discoveries)
}

// DiscoverServerPlugins returns the available plugins associated with the given server
func DiscoverServerPlugins(serverName string) ([]common.Plugin, error) {
	server, err := config.GetCurrentServer()
	if err != nil {
		return []common.Plugin{}, nil
	}

	// TODO: Need to mark these plugins as context based plugins

	return discoverPlugins(server.Discoveries)
}

// DiscoverPlugins returns the available plugins that can be used with the given server
func DiscoverPlugins(serverName string) ([]common.Plugin, error) {
	allPlugins := []common.Plugin{}
	serverPlugins, err := DiscoverServerPlugins(serverName)
	if err != nil {
		return allPlugins, errors.Wrapf(err, "unable to discover server plugins")
	}
	standalonePlugins, err := DiscoverStandalonePlugins()
	if err != nil {
		return allPlugins, errors.Wrapf(err, "unable to discover server plugins")
	}
	allPlugins = append(serverPlugins, standalonePlugins...)

	// TODO(anuj): Remove duplicate plugins with server plugins getting higher priority

	return allPlugins, nil
}

// ListPlugins returns the available plugins.
func ListPlugins(serverName string, exclude ...string) ([]*cliv1alpha1.PluginDescriptor, error) {
	pluginDescriptors, err := catalog.GetPluginsFromCatalogCache(serverName)
	if err != nil {
		return nil, errors.Errorf("could not get plugin descriptors %v", err)
	}
	return pluginDescriptors, nil
}

// DescribePlugin describes a plugin.
func DescribePlugin(serverName, pluginName string) (desc *cliv1alpha1.PluginDescriptor, err error) {
	pluginPath, err := catalog.GetPluginPath(serverName, pluginName)
	if err != nil {
		err = fmt.Errorf("could not get plugin path for plugin %q", pluginName)
	}

	b, err := exec.Command(pluginPath, "info").Output()
	if err != nil {
		err = fmt.Errorf("could not describe plugin %q", pluginName)
		return
	}

	var descriptor cliv1alpha1.PluginDescriptor
	err = json.Unmarshal(b, &descriptor)
	if err != nil {
		err = fmt.Errorf("could not unmarshal plugin %q description", pluginName)
	}
	return &descriptor, err
}

// InitializePlugin initializes the plugin configuration
func InitializePlugin(serverName, pluginName string) error {
	pluginPath, err := catalog.GetPluginPath(serverName, pluginName)
	if err != nil {
		err = fmt.Errorf("could not get plugin path for plugin %q", pluginName)
	}

	b, err := exec.Command(pluginPath, "post-install").CombinedOutput()

	// Note: If user is installing old version of plugin than it is possible that
	// the plugin does not implement post-install command. Ignoring the
	// errors if the command does not exist for a particular plugin.
	if err != nil && !strings.Contains(string(b), "unknown command") {
		log.Warningf("Warning: Failed to initialize plugin '%q' after installation. %v", pluginName, string(b))
	}

	return nil
}

// InstallPlugin installs a plugin from the given repository.
func InstallPlugin(serverName, pluginName, version string, distribution distribution.Distribution) error {
	return installOrUpgradePlugin(serverName, pluginName, version, distribution)
}

// UpgradePlugin upgrades a plugin from the given repository.
func UpgradePlugin(serverName, pluginName, version string, distribution distribution.Distribution) error {
	return installOrUpgradePlugin(serverName, pluginName, version, distribution)
}

func installOrUpgradePlugin(serverName, pluginName, version string, distribution distribution.Distribution) error {
	b, err := distribution.Fetch(pluginName, version, string(common.BuildArch()))
	if err != nil {
		return err
	}

	pluginPath := filepath.Join(pluginRoot, distribution.GetInstallationPath(), pluginName, version)

	if common.BuildArch().IsWindows() {
		pluginPath += exe
	}

	err = os.WriteFile(pluginPath, b, 0755)
	if err != nil {
		return errors.Wrap(err, "could not write file")
	}

	b, err = exec.Command(pluginPath, "info").Output()
	if err != nil {
		return fmt.Errorf("could not describe plugin %q", pluginName)
	}
	var descriptor cliv1alpha1.PluginDescriptor
	err = json.Unmarshal(b, &descriptor)
	if err != nil {
		err = fmt.Errorf("could not unmarshal plugin %q description", pluginName)
	}
	descriptor.InstallationPath = pluginPath

	err = catalog.InsertOrUpdatePluginCacheEntry(serverName, pluginName, descriptor)
	if err != nil {
		log.Debug("Plugin descriptor could not be updated in cache")
	}
	err = InitializePlugin(serverName, pluginName)
	if err != nil {
		log.Infof("could not initialize plugin after installing: %v", err.Error())
	}
	return nil
}

// DeletePlugin deletes a plugin.
func DeletePlugin(serverName, pluginName string) error {
	pluginPath, err := catalog.GetPluginPath(serverName, pluginName)
	if err != nil {
		err = fmt.Errorf("could not get plugin path for plugin %q", pluginName)
	}

	err = catalog.DeletePluginCacheEntry(serverName, pluginName)
	if err != nil {
		log.Debugf("Plugin descriptor could not be deleted from cache %v", err)
	}

	return os.Remove(pluginPath)
}

// Clean deletes all plugins and tests.
func Clean() error {
	if err := catalog.CleanCatalogCache(); err != nil {
		return errors.Errorf("Failed to clean the catalog cache %v", err)
	}
	return os.RemoveAll(pluginRoot)
}
