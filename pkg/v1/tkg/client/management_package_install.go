// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
)

type TKGPackageConfig struct {
	Metadata     Metadata          `yaml:"metadata"`
	ConfigValues map[string]string `yaml:"configvalues"`
}

type Metadata struct {
	InfraProvider string `yaml:"infraProvider"`
}

func (c *TkgClient) InstallManagementPackages(kubeConfig, kubeContext string) error {
	// create package client
	pkgClient, err := tkgpackageclient.NewTKGPackageClient(kubeConfig, kubeContext)
	if err != nil {
		return err
	}

	// install management package repository
	err = c.installManagementPackageRepository(pkgClient)
	if err != nil {
		return err
	}

	// install tkg composite management package
	err = c.installTKGManagementPackage(pkgClient)
	if err != nil {
		return err
	}

	return nil
}

func (c *TkgClient) installManagementPackageRepository(pkgClient tkgpackageclient.TKGPackageClient) error {
	managementPackageRepoImage, err := c.tkgBomClient.GetManagementPackageRepositoryImage()
	if err != nil {
		return errors.Wrap(err, "unable to get management package repository image")
	}

	repositoryOptions := tkgpackagedatamodel.NewRepositoryOptions()
	repositoryOptions.RepositoryName = constants.TKGManagementPackageRepositoryName
	repositoryOptions.RepositoryURL = managementPackageRepoImage

	// TODO(anuj): Remove this hard coded repository url
	repositoryOptions.RepositoryURL = "gcr.io/eminent-nation-87317/tkg/test/repo/management/packages/management/management@sha256:afe1a792c7290e535522b4d3f2bf4f7b9e01ef1ad0cf9720f93a68de8eed539f"
	repositoryOptions.Namespace = constants.TkgNamespace
	repositoryOptions.CreateRepository = true
	repositoryOptions.Wait = true
	repositoryOptions.PollInterval = time.Second * 5
	repositoryOptions.PollTimeout = time.Second * 60

	return pkgClient.UpdateRepository(repositoryOptions, nil, tkgpackagedatamodel.OperationTypeInstall)
}

func (c *TkgClient) installTKGManagementPackage(pkgClient tkgpackageclient.TKGPackageClient) error {
	packageOptions := tkgpackagedatamodel.NewPackageOptions()
	packageOptions.PackageName = constants.TKGManagementPackageName
	packageOptions.PkgInstallName = constants.TKGManagementPackageInstallName
	packageOptions.Namespace = constants.TkgNamespace
	packageOptions.Install = true
	packageOptions.Wait = true
	packageOptions.PollInterval = time.Second * 5
	packageOptions.PollTimeout = time.Minute * 5

	valuesFilepath, err := c.getTKGPackageConfig()
	if err != nil {
		return err
	}
	defer os.Remove(valuesFilepath)

	packageOptions.ValuesFile = valuesFilepath

	return pkgClient.InstallPackage(packageOptions, nil, tkgpackagedatamodel.OperationTypeInstall)
}

func (c *TkgClient) getTKGPackageConfig() (string, error) {
	tkgPackageConfig := TKGPackageConfig{
		Metadata: Metadata{
			InfraProvider: "aws",
		},
		ConfigValues: map[string]string{"AWS_REGION": "us-east-1"},
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
