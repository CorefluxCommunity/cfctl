package cmd

import (
	"fmt"

	"github.com/CorefluxCommunity/vaultctl/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listClustersSubCmd)

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured Vault clusters",
	Long:  `List configured Vault clusters.`,
}

var listClustersSubCmd = &cobra.Command{
	Use:   "clusters",
	Short: "List Vault clusters",
	Long:  `List Vault clusters in vaultctl configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		i := 0
		for name, cluster := range config.Clusters {
			utils.PrintHeader(name)
			utils.PrintKVSlice("Server(s)", cluster[0].Servers)

			if i < len(config.Clusters)-1 {
				fmt.Println()
			}

			i++
		}
	},
}
