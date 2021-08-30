// Copyright YEAR VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DistributionConfig contains a specific distribution mechanism. Only one of the
// configs must be set.
type DistributionConfig struct {
	// GCPStorage is set if the plugin is distributed via Google Cloud Storage.
	GCP *GCPStorage `json:"gcp"`
	// GCPStorage is set if the plugin is distributed via an OCI Image Registry.
	OCI *OCIRegistry `json:"oci"`
}

// GCPStorage provides a Google Cloud Storage bucket with an optional base path (or
// object prefix). The object download path name is constructed as
// `{Bucket}/{BasePath}/{Version}/{OS}/{Arch}`.
type GCPStorage struct {
	// Bucket is a Google Cloud Storage bucket.
	// E.g., tanzu-cli
	Bucket string `json:"bucket"`
	// BasePath is a URI path that is prefixed to the object name/path.
	// E.g., plugins/cluster
	BasePath string `json:"basePath"`
}

// OCIRegistry provides an OCI compliant image registry which supports multiple
// architecture image manifests. The fully qualified image name is constructed as
// `{Registry}/{Name}:{Version}`.
type OCIRegistry struct {
	// Registry is an OCI compliant image registry. It MUST be a DNS-compatible name.
	// E.g., harbor.my-domain.local
	Registry string `json:"registry,omitempty"`
	// Name is the unique repository/image name. It MUST be a valid URI path, MAY
	// contain zero or more '/', and SHOULD NOT start or end with '/'.
	// E.g., tanzu/cli/plugins/cluster
	Name string `json:"name"`
}

// VersionConstraints
type VersionConstraints struct {
	// RecommendedVersion version that Tanzu CLI should use if available.
	// The value should be a valid semantic version as defined in
	// https://semver.org/.
	RecommendedVersion string `json:"recommendedVersion"`
	// SupportedVersions determines the list of supported CLI plugin versions.
	// The values for each should follow the version constraints format mentioned
	// https://github.com/Masterminds/semver#checking-version-constraints.
	SupportedVersions []string `json:"supportedVersions"`
}

// Platform of the plugin binary.
type Platform struct {
	// OS of the plugin binary in `GOOS` format.
	OS string `json:"os"`
	// Arch of the plugin binary in `GOARCH` format.
	Arch string `json:"arch"`
}

// CLIPluginSpec defines the desired state of CLIPlugin.
type CLIPluginSpec struct {
	// Description is the plugin's description.
	Description string `json:"description"`
	// Distribution mechanism for the plugin.
	Distribution DistributionConfig `json:"distribution"`
	// VersionConstraints for the plugin describes constraints
	// around using a version of the plugin.
	VersionConstraints VersionConstraints `json:"versionConstraints"`
	// Platforms available for the plugin.
	Platforms []Platform `json:"platforms"`
	// Optional specifies whether the plugin is mandatory or optional
	// If optional, the plugin will not get auto-downloaded as part of
	// `tanzu login` or `tanzu plugin sync` command
	// To view the list of plugin, user can use `tanzu plugin list` and
	// to download a specific plugin run, `tanzu plugin install <plugin-name>`
	Optional bool `json:"optional"`
}

//+kubebuilder:object:root=true

// CLIPlugin denotes a Tanzu cli plugin.
type CLIPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CLIPluginSpec `json:"spec"`
}

//+kubebuilder:object:root=true

// CLIPluginList contains a list of CLIPlugin
type CLIPluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CLIPlugin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CLIPlugin{}, &CLIPluginList{})
}
