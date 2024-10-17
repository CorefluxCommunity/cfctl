package cmd

import (
	"github.com/CorefluxCommunity/vaultctl/pkg/utils"
	"github.com/CorefluxCommunity/vaultctl/pkg/vault"
	"github.com/spf13/cobra"
)

func init() {
	loginClusterSubCmd.Flags().StringVarP(&method, "method", "m", "", "Authentication method for Vault")
	loginClusterSubCmd.Flags().StringVarP(&user, "user", "u", "", "Username for Vault login")
	loginClusterSubCmd.Flags().StringVarP(&password, "password", "p", "", "Password for Vault login")

	loginCmd.AddCommand(loginClusterSubCmd)

	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Vault",
	Long:  `Login to Vault using the userpass auth method.`,
}

var loginClusterSubCmd = &cobra.Command{
	Use:   "cluster <cluster name>",
	Short: "Login to Vault",
	Long:  `Login to Vault using the userpass auth method.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vaultAddr, err := getVaultAddress(args[0])
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}

		if err := vault.Authenticate(vaultAddr, method, user, password); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}
