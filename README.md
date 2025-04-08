# RBAC Vision

RBAC Vision is a CLI tool for analyzing and comparing Kubernetes ClusterRoles. It helps administrators and developers understand the differences between ClusterRoles and their permissions.

## Features

- List all ClusterRoles in the cluster
- Filter ClusterRoles by name
- Compare two or more ClusterRoles to see their differences in permissions
- Display unique permissions for each ClusterRole

## Installation

1. Ensure you have Go 1.22 or later installed
2. Clone this repository
3. Run `go install` from the project directory

## Usage

### List ClusterRoles

```bash
# List all ClusterRoles
rbac-vision list

# Filter ClusterRoles by name
rbac-vision list -f admin
```

### Compare ClusterRoles

```bash
# Compare two ClusterRoles
rbac-vision compare cluster-admin edit

# Compare multiple ClusterRoles
rbac-vision compare cluster-admin admin edit view
```

## Building from Source

```bash
go build -o rbac-vision
```

## Requirements

- Go 1.22 or later
- Access to a Kubernetes cluster
- Valid kubeconfig file
