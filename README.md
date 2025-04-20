# Peekarole

Peekarole is a CLI tool for analyzing the permissions of a single Kubernetes ClusterRole. It helps administrators and developers understand the structure and permissions of a ClusterRole by summarizing its rules.

## Features

- Load a ClusterRole from a JSON file
- Summarize permissions by API group and resource
- Display allowed verbs and resource names for each resource

## Installation

1. Ensure you have Go 1.22 or later installed.
2. Clone this repository.
3. Run `go install` from the project directory.

## Usage

### Analyze a ClusterRole

First, export a ClusterRole from your Kubernetes cluster as JSON:

```bash
kubectl get clusterrole <role-name> -o json > role.json