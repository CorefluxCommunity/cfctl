package vault

import (
	"fmt"

	"github.com/CorefluxCommunity/vaultctl/pkg/utils"
	"github.com/hashicorp/vault/api"
)

func printSealStatus(resp *api.SealStatusResponse) {
	status := "unsealed"
	if resp.Sealed {
		status = "sealed"
	} else {
		utils.PrintKV("Cluster name", resp.ClusterName)
		utils.PrintKV("Cluster ID", resp.ClusterID)
	}

	utils.PrintKV("Seal status", status)
	utils.PrintKV("Key threshold/shares", fmt.Sprintf("%d/%d", resp.T, resp.N))
	utils.PrintKV("Progress", fmt.Sprintf("%d/%d", resp.Progress, resp.T))
	utils.PrintKV("Version", resp.Version)
}

// ListVaultStatus will output of the status the provided Vault address.
func ListVaultStatus(vaultAddr string) error {
	vault, err := NewVaultClient(vaultAddr)
	if err != nil {
		return err
	}

	resp, err := vault.ApiClient.Sys().SealStatus()
	if err != nil {
		return err
	}

	printSealStatus(resp)

	return nil
}
