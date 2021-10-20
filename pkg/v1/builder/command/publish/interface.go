// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package publish implements plugin and plugin api publishing related function
package publish

// Publisher is an interface to publish plugins and plugin api
type Publisher interface {
	PublishStandAlonePlugin() error

	PublishContextAwarePlugin() error
}
