package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version = "v0.1.0"
	rootCmd = &cobra.Command{
		Use:   "rbac-vision",
		Short: "A tool for analyzing and comparing Kubernetes ClusterRoles",
		Long: `rbac-vision is a CLI tool that helps you analyze and compare Kubernetes ClusterRoles.
It provides functionality to list, filter, and compare ClusterRoles to understand their
permission differences.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(compareCmd)
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{printf "rbac-vision %s" .Version}}
`)
}
