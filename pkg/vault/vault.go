package vault

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/api"
)

type VaultClient struct {
	ApiClient *api.Client
	url       *url.URL
}

func NewVaultClient(addr string) (*VaultClient, error) {
	vault := new(VaultClient)

	config := &api.Config{
		Address: addr,
	}

	api, err := api.NewClient(config)
	if err != nil {
		return vault, err
	}

	// Disable namespace usage
	api.SetNamespace("")

	url, err := url.Parse(addr)
	if err != nil {
		return vault, err
	}

	vault.ApiClient = api
	vault.url = url

	return vault, nil
}

func (v *VaultClient) FetchSecret(path string) (map[string]interface{}, error) {
	secret, err := v.ApiClient.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found")
	}

	// Check if this is a KV2 secret
	if data, ok := secret.Data["data"].(map[string]interface{}); ok {
		return data, nil
	}

	// If it's not a KV2 secret, return the data as is
	if secret.Data != nil {
		return secret.Data, nil
	}

	return nil, fmt.Errorf("unexpected secret data format")
}
