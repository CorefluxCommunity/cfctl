package cmd

import (
	"github.com/CorefluxCommunity/vaultctl/pkg/utils"
	"github.com/CorefluxCommunity/vaultctl/pkg/vault"
	"github.com/spf13/cobra"
)

func init() {
	showCmd.AddCommand(showClusterSubCmd)

	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of Vault clusters",
	Long:  `Show details of Vault clusters.`,
}

var showClusterSubCmd = &cobra.Command{
	Use:   "cluster <cluster name>",
	Short: "Show Vault cluster status",
	Long:  `Show overview and unseal status of the specified Vault cluster.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]

		cluster, err := getVaultClusterConfig(clusterName)
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		if len(cluster.Servers) == 0 {
			utils.PrintFatal("no Vault servers in configuration", 1)
		}

		utils.PrintHeader("Vault Cluster Status")
		utils.PrintKVSlice("Server(s)", cluster.Servers)

		if err := vault.ListVaultStatus(cluster.Servers[0]); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}
