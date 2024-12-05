# vaultctl

``vaultctl`` is a command-line interface (CLI) tool designed to facilitate operations with HashiCorp Vault by allowing users to fetch and export secrets based on predefined configurations. This guide will cover how to use the vaultctl commands, install it as a Nix package, and understand its features.

## Installation

``vaultctl`` is packaged as a Nix Flake. To install and use ``vaultctl``, you need to have Nix installed. Follow the [Nix installation guide](https://nix.dev/manual/nix/2.18/installation/multi-user) if you haven't set up Nix yet.

### Using vaultctl with Nix Flakes

Once Nix is installed, you can enter a shell with ``vaultctl`` ready to use by running the following command:

```bash
nix develop github:CorefluxCommunity/vaultctl
```

Alternatively, you can add ``vaultctl`` to your own Nix Flake as an input. Here's an example of how to include it in your flake:

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    vaultctl.url = "github:CorefluxCommunity/vaultctl";
  };

  outputs = { self, nixpkgs, vaultctl }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      devShell = pkgs.mkShell {
        buildInputs = [ vaultctl.packages."x86_64-linux".default ];
      };
    };
}
```

This configuration sets up a development shell with ``vaultctl`` available.

## Configuration Files

``vaultctl`` uses two main configuration files:

### Vault Configuration File (config.hcl): Defines the Vault clusters and their servers.

**Example configuration (~/.vaultctl/config.hcl):**

```hcl
cluster "prod" {
  address = "https://vault.example.com"
  servers = ["https://vault1.example.com", "https://vault2.example.com"]
}

cluster "dev" {
  address = "https://vault-dev.example.com"
  servers = ["https://vault-dev1.example.com"]
}
```

### Secrets Context File (contexts.hcl): Defines the secrets to be fetched and the keys to be exported.

**Example context file (./contexts.hcl):**

```hcl
context "app-secrets" {
  secret "database" {
    path = "secret/data/prod/database"
    key {
      name = "username"
    }
    key {
      name = "password"
      export_name = "DB_PASSWORD"
    }
  }

  secret "api" {
    path = "secret/data/prod/api"
    key {
      name = "api_key"
      base64_decode = true
    }
  }
}
```

## Command Overview

``vaultctl`` provides the following main commands:

### ``login``

Authenticate to Vault using the userpass authentication method. The token generated after authentication is stored locally and used for subsequent requests.

**Usage:**

```bash
vaultctl login cluster <cluster-name> --method userpass --user <username> --password <password>
```

``cluster-name``: The name of the Vault cluster (from the configuration).
``method``: The authentication method (currently only supports userpass. WIP: client certs).
``user``: The username for Vault login.
``password``: The password for Vault login.

### ``get secrets``

Fetch secrets from Vault based on a predefined context from a secrets context file. The secrets can be exported to a .env file using the --export flag.

**Usage:**

```bash
vaultctl get secrets <cluster-name> <context-name> --context <file-path> [--export]
```

``cluster-name``: The Vault cluster name.
``context-name``: The context name defined in the context file.
``context``: (Optional) Path to the context file. Defaults to ./contexts.hcl.
``export``: (Optional) Export secrets as env file at ``/tmp/secrets.env``.

```bash
vaultctl get secrets prod app-secrets --context ./app-contexts.hcl --export
```

**Important: Exporting Secrets in Shell**

When using the --export flag to export secrets as environment variables, the output needs to be wrapped in eval to actually set the variables in your current shell session.

**Example:**

```bash
eval $(vaultctl get secrets prod app-secrets --context ./app-contexts.hcl --export)
```

This works by generating export commands that need to be evaluated in the current shell context.

### ``list clusters``

Lists the configured Vault clusters from the ``vaultctl`` configuration file.

**Usage:**

```bash
vaultctl list clusters
```

### ``show cluster``

Displays detailed information about a Vault cluster, including its unseal status and the list of servers.

**Usage:**

```bash
vaultctl show cluster <cluster-name>
```
