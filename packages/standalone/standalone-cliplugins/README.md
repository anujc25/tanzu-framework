# standalone-cliplugins Package

This standalone-cliplugins package is currently part of management package plugins to streamline the package release process. This package is not (and should not be) installed on the management cluster. Instead, it is used by the Tanzu CLI directly as an OCI bundle. In future, this package should be placed at more appropriate location.

## Components

* Required CLIPlugin resources for standalone plugins like `login`, `management-cluster`, `package`, `secret` etc.

## Details

To learn more about the CLIPlugins and how it is getting used with Tanzu CLI, refer to this [doc](../../docs/design/context-aware-plugin-discovery-design.md)

## NOTE

This standalone-cliplugins package is currently part of management package plugins to streamline the package release process. This package is not directly getting installed on the management cluster but used by tanzu cli directly as an oci bundle. In future, this plugin should be placed at more appropriate location.
