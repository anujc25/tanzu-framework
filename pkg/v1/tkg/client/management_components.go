// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/managementcomponents"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
)

// InstallManagementComponents install management components to the cluster
func (c *TkgClient) InstallManagementComponents(kubeconfig, kubecontext string) error {
	_, err := c.tkgBomClient.GetManagementPackageRepositoryImage()
	if err != nil {
		return errors.Wrap(err, "unable to get management package repository image")
	}

	// TODO: Remove this hardcoded image repo
	managementPackageRepoImage := "gcr.io/eminent-nation-87317/tkg/test/repo/management/packages/management/management@sha256:ab816883ef02a302d43eacf41245e214a1f90281774015102cfcc8de37adf7b8"

	tkgPackageValuesFile, err := c.getTKGPackageConfigValuesFile()
	if err != nil {
		return err
	}

	managementcomponentsInstallOptions := managementcomponents.ManagementComponentsInstallOptions{
		ClusterOptions: managementcomponents.ClusterOptions{
			Kubeconfig:  kubeconfig,
			Kubecontext: kubecontext,
		},
		KappControllerOptions: managementcomponents.KappControllerOptions{
			KappControllerConfigFile:       "https://github.com/vmware-tanzu/carvel-kapp-controller/releases/download/v0.31.0/release.yml",
			KappControllerInstallNamespace: "kapp-controller",
		},
		ManagementPackageRepositoryOptions: managementcomponents.ManagementPackageRepositoryOptions{
			ManagementPackageRepoImage: managementPackageRepoImage,
			TKGPackageValuesFile:       tkgPackageValuesFile,
		},
	}

	return managementcomponents.InstallManagementComponents(managementcomponentsInstallOptions)
}

func (c *TkgClient) getTKGPackageConfigValuesFile() (string, error) {
	userProviderConfigValues, err := c.GetUserConfigVariableValueMap()
	if err != nil {
		return "", err
	}

	tkgPackageConfig, err := managementcomponents.GetTKGPackageConfigValuesFileFromUserConfig(userProviderConfigValues)
	if err != nil {
		return "", err
	}

	configBytes, err := yaml.Marshal(tkgPackageConfig)
	if err != nil {
		return "", err
	}

	valuesFile, err := utils.CreateTempFile("", "")
	if err != nil {
		return "", err
	}

	err = utils.WriteToFile(valuesFile, configBytes)
	if err != nil {
		return "", err
	}
	return valuesFile, nil
}

func (c *TkgClient) GetUserConfigVariableValueMap() (map[string]string, error) {
	path, err := c.tkgConfigPathsClient.GetConfigDefaultsFilePath()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	variables, err := GetConfigVariableListFromYamlData(bytes)
	if err != nil {
		return nil, err
	}

	userProvidedConfigValues := map[string]string{}
	for _, k := range variables {
		if v, e := c.TKGConfigReaderWriter().Get(k); e == nil {
			userProvidedConfigValues[k] = v
		}
	}

	return userProvidedConfigValues, nil
}

func GetConfigVariableListFromYamlData(bytes []byte) ([]string, error) {
	configValues := map[string]interface{}{}
	err := yaml.Unmarshal(bytes, &configValues)
	if err != nil {
		return nil, errors.Wrap(err, "error while unmarshaling")
	}

	keys := make([]string, 0, len(configValues))
	for k := range configValues {
		keys = append(keys, k)
	}

	return keys, nil
}
