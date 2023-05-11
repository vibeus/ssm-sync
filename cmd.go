package main

import (
	"github.com/spf13/cobra"
)

var app App

var rootCmd = &cobra.Command{
	Use:   "ssm-sync",
	Short: "Sync SSM parameters to Terraform",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := app.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&app.AwsRegion, "region", "us-west-2", "AWS region")
	rootCmd.Flags().StringVar(&app.ResFile, "resource", "main.tf", "Resource file name")
	rootCmd.Flags().StringVar(&app.VarFile, "variables", "variables.tf", "Variable file name")
	rootCmd.Flags().StringVar(&app.TfVarFile, "tfvars", "terraform.tfvars", ".tfvars file name")
	rootCmd.Flags().StringVar(&app.DataDir, "data-dir", "data", "Data directory")
	rootCmd.Flags().StringVar(&app.DataFile, "data", "main.tf", "Data file name")
}

func Execute() error {
	return rootCmd.Execute()
}
