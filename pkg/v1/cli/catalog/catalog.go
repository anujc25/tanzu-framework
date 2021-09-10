// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package catalog

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	apimachineryjson "k8s.io/apimachinery/pkg/runtime/serializer/json"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
)

const (
	// catalogCacheDirName is the name of the local directory in which tanzu state is stored.
	catalogCacheDirName = ".cache/tanzu"
	// catalogCacheFileName is the name of the file which holds Catalog cache
	catalogCacheFileName = "catalog_v2.yaml"
)

var (
	// PluginRoot is the plugin root where plugins are installed
	pluginRoot = common.DefaultPluginRoot
)

// getCatalogCacheDir returns the local directory in which tanzu state is stored.
func getCatalogCacheDir() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return path, errors.Wrap(err, "could not locate user home directory")
	}
	path = filepath.Join(home, catalogCacheDirName)
	return
}

// NewCatalog creates an instance of Catalog.
func NewCatalog() (*cliv1alpha1.Catalog, error) {
	c := &cliv1alpha1.Catalog{}

	err := ensureRoot()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// getCatalogCache retrieves the catalog from from the local directory.
func getCatalogCache() (catalog *cliv1alpha1.Catalog, err error) {
	catalogCachePath, err := getCatalogCachePath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(catalogCachePath)
	if err != nil {
		catalog, err = NewCatalog()
		if err != nil {
			return nil, err
		}
		return catalog, nil
	}
	scheme, err := cliv1alpha1.SchemeBuilder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scheme")
	}
	s := apimachineryjson.NewSerializerWithOptions(apimachineryjson.DefaultMetaFactory, scheme, scheme,
		apimachineryjson.SerializerOptions{Yaml: true, Pretty: false, Strict: false})
	var c cliv1alpha1.Catalog
	_, _, err = s.Decode(b, nil, &c)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode catalog file")
	}
	return &c, nil
}

// saveCatalogCache saves the catalog in the local directory.
func saveCatalogCache(catalog *cliv1alpha1.Catalog) error {
	catalogCachePath, err := getCatalogCachePath()
	if err != nil {
		return err
	}
	_, err = os.Stat(catalogCachePath)
	if os.IsNotExist(err) {
		catalogCacheDir, err := getCatalogCacheDir()
		if err != nil {
			return errors.Wrap(err, "could not find tanzu cache dir for OS")
		}
		err = os.MkdirAll(catalogCacheDir, 0755)
		if err != nil {
			return errors.Wrap(err, "could not make tanzu cache directory")
		}
	} else if err != nil {
		return errors.Wrap(err, "could not create catalog cache path")
	}

	scheme, err := cliv1alpha1.SchemeBuilder.Build()
	if err != nil {
		return errors.Wrap(err, "failed to create scheme")
	}

	s := apimachineryjson.NewSerializerWithOptions(apimachineryjson.DefaultMetaFactory, scheme, scheme,
		apimachineryjson.SerializerOptions{Yaml: true, Pretty: false, Strict: false})
	catalog.GetObjectKind().SetGroupVersionKind(cliv1alpha1.GroupVersionKind)
	buf := new(bytes.Buffer)
	if err := s.Encode(catalog, buf); err != nil {
		return errors.Wrap(err, "failed to encode catalog cache file")
	}
	if err = os.WriteFile(catalogCachePath, buf.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "failed to write catalog cache file")
	}
	return nil
}

// savePluginsToCatalogCache saves plugins to catalog cache
func savePluginsToCatalogCache(list []*cliv1alpha1.PluginDescriptor) error {
	catalog, err := getCatalogCache()
	if err != nil {
		catalog, err = NewCatalog()
		if err != nil {
			return err
		}
	}
	catalog.PluginDescriptors = list
	if err := saveCatalogCache(catalog); err != nil {
		return err
	}
	return nil
}

// GetPluginsFromCatalogCache gets plugins from catalog cache
func GetPluginsFromCatalogCache(serverName string) (list []*cliv1alpha1.PluginDescriptor, err error) {
	catalog, err := getCatalogCache()
	if err != nil {
		return nil, err
	}
	if len(catalog.PluginDescriptors) == 0 {
		return nil, errors.New("could not retrieve plugin descriptors from catalog cache")
	}
	return catalog.PluginDescriptors, nil
}

// InsertOrUpdatePluginCacheEntry inserts or updates a plugin entry in catalog cache
func InsertOrUpdatePluginCacheEntry(name string) error {
	list, err := GetPluginsFromCatalogCache()
	if err != nil {
		return err
	}
	list = remove(list, name)
	descriptor, err := DescribePlugin(PluginNameFromBin(name))
	if err != nil {
		return err
	}
	list = append(list, descriptor)
	if err := savePluginsToCatalogCache(list); err != nil {
		return err
	}
	return nil
}

// DeletePluginCacheEntry deletes plugin entry in catalog cache
func DeletePluginCacheEntry(name string) error {
	list, err := GetPluginsFromCatalogCache()
	if err != nil {
		return err
	}
	list = remove(list, name)
	if err := savePluginsToCatalogCache(list); err != nil {
		return err
	}
	return nil
}

// CleanCatalogCache cleans the catalog cache
func CleanCatalogCache() error {
	catalogCachePath, err := getCatalogCachePath()
	if err != nil {
		return err
	}
	if err := os.Remove(catalogCachePath); err != nil {
		return err
	}
	return nil
}

func GetPluginPath(serverName, pluginName string) (string, error) {
	return "", nil
}

// getCatalogCachePath gets the catalog cache path
func getCatalogCachePath() (string, error) {
	catalogCacheDir, err := getCatalogCacheDir()
	if err != nil {
		return "", errors.Wrap(err, "could not locate catalog cache directory")
	}
	return filepath.Join(catalogCacheDir, catalogCacheFileName), nil
}

func remove(list []*cliv1alpha1.PluginDescriptor, name string) []*cliv1alpha1.PluginDescriptor {
	i := 0
	for _, v := range list {
		if v != nil && v.Name != name {
			list[i] = v
			i++
		}
	}
	list = list[:i]
	return list
}

func inExclude(name string, exclude []string) bool {
	for _, e := range exclude {
		if name == e {
			return true
		}
	}
	return false
}

// Ensure the root directory exists.
func ensureRoot() error {
	_, err := os.Stat(testPath())
	if os.IsNotExist(err) {
		err := os.MkdirAll(testPath(), 0755)
		return errors.Wrap(err, "could not make root plugin directory")
	}
	return err
}

// Returns the test path relative to the plugin root
func testPath() string {
	return filepath.Join(pluginRoot, "test")
}
