// Package cmd contains the individual commands for the CLI.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.PersistentFlags().StringVarP(&accountName, "account", "a", "private.ecdsa", "Path to the private key.")
	// rootCmd.PersistentFlags().StringVarP(&accountPath, "account-path", "p", "zblock/accounts/", "Path to the directory with private keys.")
}

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Toqns CLI wallet version " + build,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
