package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/CorefluxCommunity/vaultctl/pkg/utils"
)

var (
	// CLI config
	config     VaultConfig
	configFile string

	// login
	method   string // TODO: Create enum for auth methods
	user     string
	password string

	// get
	contextFile   string
	exportSecrets bool

	rootCmd = &cobra.Command{
		Use: "vaultctl",
	}
)

type VaultConfig struct {
	Clusters map[string][]*VaultClusterConfig `hcl:"cluster" mapstructure:"cluster"`
}

type VaultClusterConfig struct {
	Address string   `hcl:"address" mapstructure:"address"`
	Servers []string `hcl:"servers" mapstructure:"servers"`
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.vaultctl/config.hcl)")
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in $HOME/vaultctl directory with name "config" (without extension).
		// TODO: Get cli config directory from shell env
		viper.AddConfigPath(home + "/.vaultctl")
		viper.SetConfigName("config")
		viper.SetConfigType("hcl")
	}

	viper.AutomaticEnv()
	viper.ReadInConfig()

	err := viper.Unmarshal(&config)
	if err != nil {
		utils.PrintFatal(fmt.Sprintf("unable to decode into struct, %v", err), 1)
	}

	// get vaultctl config direction and set as cwd
	configDir, err := getConfigDir()
	if err != nil {
		utils.PrintFatal(err.Error(), 1)
	}

	os.Chdir(configDir)
}

func getVaultClusterConfig(clusterName string) (*VaultClusterConfig, error) {
	for name, cluster := range config.Clusters {
		if name == clusterName {
			return cluster[0], nil
		}
	}

	return nil, fmt.Errorf("config for Vault cluster '%s' not found", clusterName)
}

func getVaultAddress(clusterName string) (string, error) {
	// Check if the cluster exists in the config
	cluster, ok := config.Clusters[clusterName]
	if !ok {
		return "", fmt.Errorf("cluster '%s' not found in configuration", clusterName)
	}

	// Use the first config in the slice (assuming there's only one per cluster name)
	clusterConfig := cluster[0]

	// If address is provided, use it
	if clusterConfig.Address != "" {
		return clusterConfig.Address, nil
	}

	// If address is not provided, but servers are, use the first server
	if len(clusterConfig.Servers) > 0 {
		return clusterConfig.Servers[0], nil
	}

	// If neither address nor servers are provided, return an error
	return "", fmt.Errorf("no address or servers found for cluster '%s'", clusterName)
}

func getConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	// TODO: Get cli config directory from shell env
	path := home + "/.vaultctl"

	return path, nil
}
