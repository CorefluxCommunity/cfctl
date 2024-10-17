package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/api"
)

type AuthMethod interface {
	Authenticate(client *VaultClient) (*api.Secret, error)
}

type UserPassAuth struct {
	Username string
	Password string
}

func (u *UserPassAuth) Authenticate(client *VaultClient) (*api.Secret, error) {
	return client.ApiClient.Logical().Write("auth/userpass/login/"+u.Username, map[string]interface{}{
		"password": u.Password,
	})
}

func Authenticate(vaultAddr, method string, user string, password string) error {
	vault, err := NewVaultClient(vaultAddr)
	if err != nil {
		return fmt.Errorf("failed to create Vault client: %w", err)
	}

	var authMethod AuthMethod

	switch method {
	case "userpass":
		authMethod = &UserPassAuth{
			Username: user,
			Password: password,
		}
	// Add more cases here for future auth methods
	default:
		return fmt.Errorf("unsupported auth method: %s", method)
	}

	secret, err := authMethod.Authenticate(vault)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if secret == nil || secret.Auth == nil {
		return fmt.Errorf("no auth info returned")
	}

	// Store the token in the CLI's home config directory
	err = storeToken(secret.Auth.ClientToken)
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	return nil
}

func storeToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".vaultctl")
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	tokenFile := filepath.Join(configDir, "token")
	err = os.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}
