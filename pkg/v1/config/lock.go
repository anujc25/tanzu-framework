// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"path/filepath"

	"github.com/aunum/log"
	"github.com/juju/fslock"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/constants"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/utils"
)

var tanzuConfigLockFile string
var tanzuConfigLock *fslock.Lock

func AcquireTanzuConfigLock() {
	var err error

	if tanzuConfigLockFile == "" {
		path, err := ClientConfigPath()
		if err != nil {
			log.Warningf("cannot acquire lock for tanzu config file, reason: %v", err)
		}
		tanzuConfigLockFile = filepath.Join(filepath.Dir(path), constants.LocalTanzuFileLock)
	}

	tanzuConfigLock, err = utils.GetFileLockWithTimeOut(tanzuConfigLockFile, utils.DefaultLockTimeout)
	if err != nil {
		log.Warningf("cannot acquire lock for tanzu config file, reason: %v", err)
	}
}

func ReleaseTanzuConfigLock() {
	if tanzuConfigLock == nil {
		return
	}
	if errUnlock := tanzuConfigLock.Unlock(); errUnlock != nil {
		log.Warningf("cannot release lock for tanzu config file, reason: %v", errUnlock)
	}
	return
}
