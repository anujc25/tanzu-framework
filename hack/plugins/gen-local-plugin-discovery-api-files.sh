#!/bin/bash
# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

STANDALONE_PLUGINS=$1
CONTEXT_PLUGINS=$2
BUILD_VERSION=$3
GOOS=$4
GOARCH=$5
