package secrets

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/CorefluxCommunity/vaultctl/pkg/vault"
)

type ContextConfig struct {
	Contexts []*Context `hcl:"context,block"`
}

type Context struct {
	Name    string    `hcl:"name,label"`
	Secrets []*Secret `hcl:"secret,block"`
}

type Secret struct {
	Name string `hcl:"name,label"`
	Path string `hcl:"path"`
	Keys []*Key `hcl:"key,block"`
}

type Key struct {
	Name         string `hcl:"name,label"`
	ExportName   string `hcl:"export_name,optional"`
	Base64Decode bool   `hcl:"base64_decode,optional"`
}

func GetSecrets(contextName string, contextFile string, exportSecrets bool, vaultAddr string) error {
	if contextFile == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current directory: %w", err)
		}
		contextFile = filepath.Join(currentDir, "contexts.hcl")
	}

	if _, err := os.Stat(contextFile); os.IsNotExist(err) {
		return fmt.Errorf("contexts file not found: %s", contextFile)
	}

	content, err := os.ReadFile(contextFile)
	if err != nil {
		return fmt.Errorf("error reading contexts file: %w", err)
	}

	var contexts ContextConfig
	err = hclsimple.Decode(contextFile, content, nil, &contexts)
	if err != nil {
		return fmt.Errorf("error decoding contexts file: %w", err)
	}

	var targetContext *Context
	for _, context := range contexts.Contexts {
		if context.Name == contextName {
			targetContext = context
			break
		}
	}

	if targetContext == nil {
		return fmt.Errorf("context '%s' not found in contexts file", contextName)
	}

	// Check for saved token
	token, err := getSavedToken()
	if err != nil || token == "" {
		return fmt.Errorf("no valid token found, please authenticate first")
	}

	// Create Vault client
	vaultClient, err := vault.NewVaultClient(vaultAddr)
	if err != nil {
		return fmt.Errorf("failed to create Vault client: %w", err)
	}
	vaultClient.ApiClient.SetToken(token)

	for _, secret := range targetContext.Secrets {
		secretData, err := vaultClient.FetchSecret(secret.Path)
		if err != nil {
			if strings.Contains(err.Error(), "permission denied") {
				fmt.Printf("Warning: Not authorized to read secret: %s\n", secret.Path)
				continue
			}
			return fmt.Errorf("failed to fetch secret %s: %w", secret.Path, err)
		}

		for _, key := range secret.Keys {
			value, ok := secretData[key.Name].(string)
			if !ok {
				fmt.Printf("Warning: Key %s not found in secret %s\n", key.Name, secret.Path)
				continue
			}

			if key.Base64Decode {
				decodedValue, err := base64.StdEncoding.DecodeString(value)
				if err != nil {
					fmt.Printf("Warning: Failed to base64 decode %s: %v\n", key.Name, err)
				} else {
					value = string(decodedValue)
				}
			}

			if exportSecrets {
				envFile := "/tmp/secrets.env"
				file, err := os.Create(envFile)
				if err != nil {
					return fmt.Errorf("failed to create env file: %w", err)
				}
				defer file.Close()

				for _, secret := range targetContext.Secrets {
					secretData, err := vaultClient.FetchSecret(secret.Path)
					if err != nil {
						fmt.Printf("Warning: Failed to fetch secret %s: %v\n", secret.Path, err)
						continue
					}

					for _, key := range secret.Keys {
						value, ok := secretData[key.Name].(string)
						if !ok {
							fmt.Printf("Warning: Key %s not found in secret %s\n", key.Name, secret.Path)
							continue
						}

						if key.Base64Decode {
							decodedValue, err := base64.StdEncoding.DecodeString(value)
							if err != nil {
								fmt.Printf("Warning: Failed to decode %s: %v\n", key.Name, err)
								continue
							}
							value = string(decodedValue)
						}

						exportName := key.Name
						if key.ExportName != "" {
							exportName = key.ExportName
						}

						// Escape special characters in value
						value = strings.ReplaceAll(value, "'", "'\"'\"'")
						_, err = fmt.Fprintf(file, "%s='%s'\n", exportName, value)
						if err != nil {
							return fmt.Errorf("failed to write to env file: %w", err)
						}
					}
				}

				fmt.Printf("Secrets exported to %s. Source it with 'source %s'\n", envFile, envFile)
			} else {
				fmt.Println("Secrets printed to stdout")
			}
		}
	}

	return nil
}

func getSavedToken() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	tokenFile := filepath.Join(homeDir, ".vaultctl", "token")
	tokenBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(tokenBytes)), nil
}
