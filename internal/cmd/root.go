package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "terraform-private-registry",
	Short: "A private registry for Terraform",
	Long:  "A private registry for Terraform",
}

// Execute runs the command, exit on error
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
