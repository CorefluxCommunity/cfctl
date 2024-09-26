package cmd

import (
	"github.com/CorefluxCommunity/vaultctl/pkg/secrets"
	"github.com/CorefluxCommunity/vaultctl/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	getSecretsSubCmd.Flags().StringVarP(&contextFile, "context", "c", "", "secrets context file")
	getSecretsSubCmd.Flags().BoolVarP(&exportSecrets, "export", "e", false, "export secrets as environment variables")

	getCmd.AddCommand(getSecretsSubCmd)

	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Vault secrets",
	Long:  `Get Vault secrets from a secrets context file.`,
}

var getSecretsSubCmd = &cobra.Command{
	Use:   "secrets <cluster name> <context name>",
	Short: "Get Vault secrets",
	Long:  `Get Vault secrets from a secrets context file.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vaultAddr, err := getVaultAddress(args[0])
		if err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
		contextName := args[1]

		if err := secrets.GetSecrets(contextName, contextFile, exportSecrets, vaultAddr); err != nil {
			utils.PrintFatal(err.Error(), 1)
		}
	},
}
