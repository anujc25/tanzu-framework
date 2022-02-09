// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package managementcomponents

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/clusterclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackagedatamodel"
)

type ClusterOptions struct {
	Kubeconfig  string
	Kubecontext string
}

type ManagementPackageRepositoryOptions struct {
	ManagementPackageRepoImage string // "gcr.io/eminent-nation-87317/tkg/test/repo/management/packages/management/management@sha256:ab816883ef02a302d43eacf41245e214a1f90281774015102cfcc8de37adf7b8"
	TKGPackageValuesFile       string
}

type KappControllerOptions struct {
	KappControllerConfigFile       string //  "https://github.com/vmware-tanzu/carvel-kapp-controller/releases/download/v0.31.0/release.yml"
	KappControllerInstallNamespace string //  "kapp-controller"
}

// ManagementComponentsInstallOptions
type ManagementComponentsInstallOptions struct {
	ClusterOptions                     ClusterOptions
	ManagementPackageRepositoryOptions ManagementPackageRepositoryOptions
	KappControllerOptions              KappControllerOptions
}

// InstallManagementComponents installs the management component to cluster
func InstallManagementComponents(mcip ManagementComponentsInstallOptions) error {
	clusterClient, err := clusterclient.NewClient(mcip.ClusterOptions.Kubeconfig, mcip.ClusterOptions.Kubecontext, clusterclient.Options{})
	if err != nil {
		return errors.Wrap(err, "unable to get cluster client")
	}
	if err := InstallKappController(clusterClient, mcip.KappControllerOptions); err != nil {
		return errors.Wrap(err, "unable to install kapp-controller")
	}
	if err = InstallManagementPackages(mcip.ClusterOptions, mcip.ManagementPackageRepositoryOptions); err != nil {
		return errors.Wrap(err, "unable to install management packages")
	}
	return nil
}

// InstallKappController installs kapp-controller to the cluster
func InstallKappController(clusterClient clusterclient.Client, kappControllerOptions KappControllerOptions) error {
	// Apply kapp-controller configuration
	err := clusterClient.ApplyFile(kappControllerOptions.KappControllerConfigFile)
	if err != nil {
		return errors.Errorf("error installing %s", constants.KappControllerDeploymentName)
	}
	// Wait for kapp-controller to be deployed and running
	err = clusterClient.WaitForDeployment(constants.KappControllerDeploymentName, kappControllerOptions.KappControllerInstallNamespace)
	if err != nil {
		return errors.Errorf("failed waiting for deployment %s", constants.KappControllerDeploymentName)
	}
	return nil
}

// InstallManagementPackages installs TKG management packages to the cluster
func InstallManagementPackages(clusterOptions ClusterOptions, mpro ManagementPackageRepositoryOptions) error {
	// create package client
	pkgClient, err := tkgpackageclient.NewTKGPackageClient(clusterOptions.Kubeconfig, clusterOptions.Kubecontext)
	if err != nil {
		return err
	}

	// install management package repository
	err = installManagementPackageRepository(pkgClient, mpro)
	if err != nil {
		return err
	}

	// install tkg composite management package
	err = installTKGManagementPackage(pkgClient, mpro)
	if err != nil {
		return err
	}

	return nil
}

func installManagementPackageRepository(pkgClient tkgpackageclient.TKGPackageClient, mpro ManagementPackageRepositoryOptions) error {
	repositoryOptions := tkgpackagedatamodel.NewRepositoryOptions()
	repositoryOptions.RepositoryName = constants.TKGManagementPackageRepositoryName
	repositoryOptions.RepositoryURL = mpro.ManagementPackageRepoImage
	repositoryOptions.Namespace = constants.TkgNamespace
	repositoryOptions.CreateRepository = true
	repositoryOptions.Wait = true
	repositoryOptions.PollInterval = time.Second * 5
	repositoryOptions.PollTimeout = time.Second * 60

	return pkgClient.UpdateRepository(repositoryOptions, nil, tkgpackagedatamodel.OperationTypeInstall)
}

func installTKGManagementPackage(pkgClient tkgpackageclient.TKGPackageClient, mpro ManagementPackageRepositoryOptions) error {
	packageOptions := tkgpackagedatamodel.NewPackageOptions()
	packageOptions.PackageName = constants.TKGManagementPackageName
	packageOptions.PkgInstallName = constants.TKGManagementPackageInstallName
	packageOptions.Namespace = constants.TkgNamespace
	packageOptions.Install = true
	packageOptions.Wait = true
	packageOptions.PollInterval = time.Second * 5
	packageOptions.PollTimeout = time.Minute * 5
	packageOptions.ValuesFile = mpro.TKGPackageValuesFile
	return pkgClient.InstallPackage(packageOptions, nil, tkgpackagedatamodel.OperationTypeInstall)
}
