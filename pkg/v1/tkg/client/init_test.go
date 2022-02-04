// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/client"
)

var _ = FDescribe("Unit tests for management package installation", func() {
	var (
		err       error
		tkgClient *TkgClient
	)

	BeforeEach(func() {
		tkgClient, err = CreateTKGClient("../fakes/config/config.yaml", testingDir, defaultTKGBoMFileForTesting, 2*time.Second)
		Expect(err).NotTo(HaveOccurred())

		err = tkgClient.InstallManagementPackages("", "")
	})

	It("should not return an error and all status should be correct", func() {
		Expect(err).NotTo(HaveOccurred())
	})
})
