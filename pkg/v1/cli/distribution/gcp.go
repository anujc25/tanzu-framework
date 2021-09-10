// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package distribution

import (
	"context"
	"fmt"
	"io"
	"path"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/common"
)

// Fetch an artifact.
func (g *GCPStorage) Fetch(name, version string, arch string) ([]byte, error) {
	ctx := context.Background()

	bkt, err := common.GetGCPBucket(ctx, g.Bucket)
	if err != nil {
		return nil, err
	}

	artifactPath := path.Join(g.BasePath, version, arch)

	return g.fetch(ctx, artifactPath, bkt)
}

// FetchTest fetches a test artifact.
func (g *GCPStorage) FetchTest(name, version string, arch string) ([]byte, error) {
	ctx := context.Background()

	bkt, err := common.GetGCPBucket(ctx, g.Bucket)
	if err != nil {
		return nil, err
	}

	artifactPath := path.Join(g.BasePath, version, "test", arch)
	return g.fetch(ctx, artifactPath, bkt)
}

func (g *GCPStorage) fetch(ctx context.Context, artifactPath string, bkt *storage.BucketHandle) ([]byte, error) {
	obj := bkt.Object(artifactPath)
	if obj == nil {
		return nil, fmt.Errorf("artifact %q not found", artifactPath)
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not read artifact %q", artifactPath))
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch artifact")
	}
	return b, nil
}

// Get the relative installation path for a plugin binary
func (g *GCPStorage) GetInstallationPath() string {
	return ""
}
