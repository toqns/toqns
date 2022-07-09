package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toqns/toqns/business/key"
)

var balanceCmd = &cobra.Command{
	Use:   "nodekey",
	Short: "Create a node key",
	Run:   nodeKey,
}

var nodeKeyFile string

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().StringVarP(&nodeKeyFile, "nodekey", "n", "./.node/node.key", "Key file of the node")
}

func nodeKey(cmd *cobra.Command, args []string) {
	fmt.Printf("\nCreating node key...\n")
	k, err := key.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := k.Save(nodeKeyFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Node key file created as:", nodeKeyFile)
}
